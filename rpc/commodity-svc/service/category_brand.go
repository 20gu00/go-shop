package service

import (
	"commodity-rpc/pb"
	"context"
	"github.com/golang/protobuf/ptypes/empty"
)

// 分类品牌列表
func (c *CommodityServer) CategoryBrandList(context.Context, *pb.CategoryBrandFilterReq) (*pb.CategoryBrandListRes, error) {
}

//通过category获取brands
func (c *CommodityServer) GetCategoryBrandList(context.Context, *pb.CategoryInfoReq) (*pb.BrandListRes, error) {
}

// 创建分类品牌
func (c *CommodityServer) CreateCategoryBrand(context.Context, *pb.CategoryBrandReq) (*pb.CategoryBrandRes, error) {
}

// 删除分类品牌
func (c *CommodityServer) DeleteCategoryBrand(context.Context, *pb.CategoryBrandReq) (*empty.Empty, error) {
}

// 更新分类品牌
func (c *CommodityServer) UpdateCategoryBrand(context.Context, *pb.CategoryBrandReq) (*empty.Empty, error) {
}
