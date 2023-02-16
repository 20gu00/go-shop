package service

import (
	"commodity-rpc/pb"
	"context"
	"github.com/golang/protobuf/ptypes/empty"
)

// 品牌列表
func (c *CommodityServer) BrandList(context.Context, *pb.BrandFilterReq) (*pb.BrandListRes, error) {}

// 创建品牌
func (c *CommodityServer) CreateBrand(context.Context, *pb.BrandReq) (*pb.BrandInfoRes, error) {}

// 删除品牌
func (c *CommodityServer) DeleteBrand(context.Context, *pb.BrandReq) (*empty.Empty, error) {}

// 更新品牌
func (c *CommodityServer) UpdateBrand(context.Context, *pb.BrandReq) (*empty.Empty, error) {}
