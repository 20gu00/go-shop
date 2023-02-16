package serivce

import (
	"commodity-rpc/pb"
	"context"
	"github.com/golang/protobuf/ptypes/empty"
)

// 轮播图列表
func (c *CommodityServer) BannerList(context.Context, *empty.Empty) (*pb.BannerListRes, error) {

}

// 创建轮播图
func (c *CommodityServer) CreateBanner(context.Context, *pb.BannerReq) (*pb.BannerRes, error) {

}

// 删除轮播图
func (c *CommodityServer) DeleteBanner(context.Context, *pb.BannerReq) (*empty.Empty, error) {

}

// 更新轮播图
func (c *CommodityServer) UpdateBanner(context.Context, *pb.BannerReq) (*empty.Empty, error) {

}
