package service

import (
	"commodity-rpc/pb"
	"context"
	"github.com/golang/protobuf/ptypes/empty"
)

type CommodityServer struct {
	pb.UnimplementedCommodityServer
}

// 商品列表
func (c *CommodityServer) CommodityList(context.Context, *pb.CommodityFilterReq) (*pb.CommodityListRes, error) {

}

// 商品增删改查
func (c *CommodityServer) CreateCommodity(context.Context, *pb.CreateCommodityReq) (*pb.CommodityInfoRes, error) {

}

// 删除商品
func (c *CommodityServer) DeleteCommodity(context.Context, *pb.DeleteCommodityInfo) (*empty.Empty, error) {

}

// 更新商品
func (c *CommodityServer) UpdateCommodity(context.Context, *pb.CreateCommodityReq) (*empty.Empty, error) {

}

//获取商品
func (c *CommodityServer) GetCommodity(context.Context, *pb.CommodityInfoReq) (*pb.CommodityInfoRes, error) {

}

// 批量获取商品信息(比如一个订单中有多个商品,多个商品的id)
func (c *CommodityServer) GetBatchCommodity(context.Context, *pb.BatchCommodityIdInfo) (*pb.CommodityListRes, error) {

}
