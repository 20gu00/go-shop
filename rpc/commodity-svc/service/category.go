package service

import (
	"commodity-rpc/pb"
	"context"
	"github.com/golang/protobuf/ptypes/empty"
)

// 商品分类列表
func (c *CommodityServer) AllCategoryList(context.Context, *empty.Empty) (*pb.CategoryListRes, error) {

}

// 获取商品的子分类
func (c *CommodityServer) SubCategoryList(context.Context, *pb.CategoryListReq) (*pb.SubCategoryListRes, error) {

}

// 增加商品
func (c *CommodityServer) CreateCategory(context.Context, *pb.CreateCommodityReq) (*pb.CategoryInfoRes, error) {

}

// 删除商品
func (c *CommodityServer) DeleteCategory(context.Context, *pb.DeleteCategoryReq) (*empty.Empty, error) {

}

// 更新商品
func (c *CommodityServer) UpdateCategory(context.Context, *pb.CategoryInfoReq) (*empty.Empty, error) {

}
