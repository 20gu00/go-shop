package service

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"store-rpc/dao"
	"store-rpc/model"
	"store-rpc/pb"
	"sync"
)

type StoreServer struct {
	pb.UnimplementedInventoryServer
}

// 设置库存
func (s *StoreServer) SetInv(ctx context.Context, req *pb.GoodsInvInfo) (*empty.Empty, error) {
	//有则更新,没有则创建
	var inv model.Inventory
	// First(&inv,req.GoodsId)  这就要设置Inventory的GoodsId为primarykey
	dao.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv)
	inv.Goods = req.GoodsId
	inv.Stocks = req.Num

	dao.DB.Save(&inv)
	return &emptypb.Empty{}, nil
}

// 获取库存详情
func (s *StoreServer) InvDetail(ctx context.Context, req *pb.GoodsInvInfo) (*pb.GoodsInvInfo, error) {
	var inv model.Inventory
	if result := dao.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有该商品的库存库存信息")
	}
	return &pb.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     inv.Stocks,
	}, nil
}

// 全局锁
var mu sync.Mutex

// 扣减库存
func (s *StoreServer) Sell(ctx context.Context, req *pb.SellInfo) (*emptypb.Empty, error) {
	// 涉及本地事务和分布式事务
	//本地事务,如果一个订单a商品库存扣减成功,b商品扣减失败比如库存不足
	//这时候就不能成功下单支付,应该要么全部成功要么全部失败
	//并发情况下可能会出现超卖,其实就是多个进程都去处理stock这个变量,哪怕是本地事务,本地事务只是保证了自身这个进程的完整执行,这就需要分布式事务

	//可以启动一个服务,里边并发调用这个扣减服务,也就是设置wg,wg.Add()和wg.Wait(),注意给并发逻辑的函数传递指针*wg,wg.Donw()

	tx := dao.DB.Begin()
	mu.Lock() //在查询和更新的逻辑之前上锁
	for _, commodityInfo := range req.GoodsInfo {
		var inv model.Inventory
		//if result := dao.DB.First(&inv, commodityInfo.GoodsId); result.RowsAffected {

		//更改成使用分布式锁(悲观锁)
		//这样这条记录就会行锁了甚至表锁
		//if result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&inv, commodityInfo.GoodsId); result.RowsAffected == 0 {
		//	return nil, status.Errorf(codes.NotFound, "没有该商品的库存库存信息")
		//}

		//操作扣减库存失败就重试
		for {
			//乐观锁
			if result := dao.DB.First(&inv, commodityInfo.GoodsId); result.RowsAffected == 0 {
				return nil, status.Errorf(codes.NotFound, "没有该商品的库存库存信息")
			}
			// 库存不足(库存0的时候就会扣减失败)
			if inv.Stocks < commodityInfo.Num {
				return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
			}
			tx.Rollback()

			inv.Stocks -= commodityInfo.Num

			//更改乐观锁字段
			//tx.Model(&inv) inv赋值了,使用的话等于添加了查询条件
			//version已经读取出来并解析到version中
			//updates更新多个字段(结构体) update更新单个列字段key:value

			//gorm的坑,数据类型的零值,这里是int也就是0,会忽略(比如库存为1,扣减了这里设置0,忽略)
			//解决方法就是强制更新字段,也就是指定要强制更新的字段
			//加上select()
			//大写,结构体字段
			if result := tx.Model(&model.Inventory{}).Select("Stocks", "Version").Where("goods=? and version = ?", commodityInfo.GoodsId, inv.Version).Updates(model.Inventory{Stocks: inv.Stocks, Version: inv.Version + 1}); result.RowsAffected == 0 {
				zap.L().Info("库存扣减失败")
				//重试
			} else {
				//退出循环
				break
			}
		}
		dao.DB.Save(&inv)
	}
	tx.Commit()
	mu.Unlock() //事务执行之后才释放锁
	return &emptypb.Empty{}, nil
}

func (s *StoreServer) Reback(ctx context.Context, req *pb.SellInfo) (*emptypb.Empty, error) {
	//库存归还：
	//1：订单超时归还
	//2. 订单创建失败(创建订单的时候会先扣减库存,然后生成实际的订单存入数据库表中,如果扣减了库存,但是订单创建失败,那么也会归还库存)
	//3. 手动归还,用户取消
	//批量归还,也就是一个订单中多个商品,不然一个一个归还容易涉及到分布式事务的问题
	tx := dao.DB.Begin()
	mu.Lock()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		if result := dao.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}

		inv.Stocks += goodInfo.Num
		tx.Save(&inv)
	}
	tx.Commit()
	mu.Unlock()
	return &emptypb.Empty{}, nil
}
