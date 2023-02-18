package service

import (
	paginate_list "commodity-rpc/common/paginate-list"
	"commodity-rpc/dao"
	"commodity-rpc/model"
	"commodity-rpc/pb"
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

//通过category获取brands
func (c *CommodityServer) GetCategoryBrandList(ctx context.Context, req *pb.CategoryInfoReq) (*pb.BrandListRes, error) {
	brandListResponse := pb.BrandListRes{}

	var category model.Category
	if result := dao.DB.Find(&category, req.Id).First(&category); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var categoryBrands []model.CategoryBrand
	// 预加载brand表的数据,也就是加载外键数据
	if result := dao.DB.Preload("Brands").Where(&model.CategoryBrand{CategoryId: req.Id}).Find(&categoryBrands); result.RowsAffected > 0 {
		brandListResponse.Total = int32(result.RowsAffected)
	}

	var brandInfoResponses []*pb.BrandInfoRes
	for _, categoryBrand := range categoryBrands {
		brandInfoResponses = append(brandInfoResponses, &pb.BrandInfoRes{
			Id:   categoryBrand.Brand.ID,
			Name: categoryBrand.Brand.Name,
			Logo: categoryBrand.Brand.Logo,
		})
	}
	brandListResponse.Data = brandInfoResponses
	return &brandListResponse, nil
}

// 创建分类品牌
func (c *CommodityServer) CreateCategoryBrand(ctx context.Context, req *pb.CategoryBrandReq) (*pb.CategoryBrandRes, error) {
	var category model.Category
	if result := dao.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brand
	if result := dao.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	//新增这两个外键(最好应该是关联append)
	categoryBrand := model.CategoryBrand{
		CategoryId: req.CategoryId,
		BrandId:    req.BrandId,
	}

	dao.DB.Save(&categoryBrand)
	return &pb.CategoryBrandRes{Id: categoryBrand.ID}, nil
}

// 删除分类品牌
func (c *CommodityServer) DeleteCategoryBrand(ctx context.Context, req *pb.CategoryBrandReq) (*empty.Empty, error) {
	if result := dao.DB.Delete(&model.CategoryBrand{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌分类不存在")
	}
	return &emptypb.Empty{}, nil
}

// 更新分类品牌
func (c *CommodityServer) UpdateCategoryBrand(ctx context.Context, req *pb.CategoryBrandReq) (*empty.Empty, error) {
	var categoryBrand model.CategoryBrand

	if result := dao.DB.First(&categoryBrand, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌分类不存在")
	}

	var category model.Category
	if result := dao.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brand
	if result := dao.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	categoryBrand.CategoryId = req.CategoryId
	categoryBrand.BrandId = req.BrandId

	dao.DB.Save(&categoryBrand)

	return &emptypb.Empty{}, nil
}

// 分类品牌列表
func (c *CommodityServer) CategoryBrandList(ctx context.Context, req *pb.CategoryBrandFilterReq) (*pb.CategoryBrandListRes, error) {
	var categoryBrands []model.CategoryBrand
	categoryBrandListResponse := pb.CategoryBrandListRes{}

	var total int64
	dao.DB.Model(&model.CategoryBrand{}).Count(&total)
	categoryBrandListResponse.Total = int32(total)

	// 预加载这两个表
	dao.DB.Preload("Category").Preload("Brands").Scopes(paginate_list.Paginate(int(req.Pages), int(req.PagePerNums))).Find(&categoryBrands)

	var categoryResponses []*pb.CategoryBrandRes
	for _, categoryBrand := range categoryBrands {
		categoryResponses = append(categoryResponses, &pb.CategoryBrandRes{
			Category: &pb.CategoryInfoRes{
				Id:             categoryBrand.Category.ID,
				Name:           categoryBrand.Category.Name,
				Level:          categoryBrand.Category.Level,
				IsTab:          categoryBrand.Category.IsTop,
				ParentCategory: categoryBrand.Category.ParCategoryId,
			},
			Brand: &pb.BrandInfoRes{
				Id:   categoryBrand.Brand.ID,
				Name: categoryBrand.Brand.Name,
				Logo: categoryBrand.Brand.Logo,
			},
		})
	}

	categoryBrandListResponse.Data = categoryResponses
	return &categoryBrandListResponse, nil
}
