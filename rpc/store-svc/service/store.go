package service

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"store-rpc/dao"
	"store-rpc/model"
	"store-rpc/pb"
)

type StoreServer struct {
	pb.UnimplementedInventoryServer
}

// 设置库存
func (s *StoreServer) SetInv(ctx context.Context, req *pb.GoodsInvInfo) (*empty.Empty, error) {
	//有则更新,没有则创建

	var inv model.Inventory
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

// 扣减库存
func (s *StoreServer) Sell(ctx context.Context, req *pb.SellInfo) (*emptypb.Empty, error) {
	// 涉及本地事务和分布式事务
	//本地事务,如果一个订单a商品库存扣减成功,b商品扣减失败比如库存不足
	//这时候就不能成功下单支付,应该要么全部成功要么全部失败
	for _, commodityInfo := range req.GoodsInfo {
		var inv model.Inventory
		if result := dao.DB.First(&inv, commodityInfo.GoodsId); result.RowsAffected {
			return nil, status.Errorf(codes.NotFound, "没有该商品的库存库存信息")
		}
		// 库存不足
		if inv.Stocks < commodityInfo.Num {
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}

		inv.Stocks -= commodityInfo.Num
		dao.DB.Save(&inv)
	}
	return &emptypb.Empty{}, nil
}
