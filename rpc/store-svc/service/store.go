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
