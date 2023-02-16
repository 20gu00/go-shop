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

// 品牌列表
func (c *CommodityServer) BrandList(ctx context.Context, req *pb.BrandFilterReq) (*pb.BrandListRes, error) {
	brandListRes := pb.BrandListRes{}

	var brands []model.Brand
	// 使用gorm的scope来分页查询,也可以使用offset limit来分页查询
	result := dao.DB.Scopes(paginate_list.Paginate(int(req.Pages), int(req.PagePerNums))).Find(&brands)
	if result.Error != nil {
		return nil, result.Error
	}

	var total int64
	// Count获取记录数目
	dao.DB.Model(&model.Brand{}).Count(&total)
	brandListRes.Total = int32(total)

	var brandRes []*pb.BrandInfoRes
	for _, brand := range brands {
		brandRes = append(brandRes, &pb.BrandInfoRes{
			Id:   brand.ID,
			Name: brand.Name,
			Logo: brand.Logo,
		})
	}
	brandListRes.Data = brandRes
	return &brandListRes, nil
}

// 创建品牌
func (c *CommodityServer) CreateBrand(context.Context, *pb.BrandReq) (*pb.BrandInfoRes, error) {
	if result := global.DB.Where("name=?", req.Name).First(&model.Brands{}); result.RowsAffected == 1 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌已存在")
	}

	brand := &model.Brands{
		Name: req.Name,
		Logo: req.Logo,
	}
	global.DB.Save(brand)

	return &proto.BrandInfoResponse{Id: brand.ID}, nil
}

// 删除品牌
func (c *CommodityServer) DeleteBrand(context.Context, *pb.BrandReq) (*empty.Empty, error) {
	if result := global.DB.Delete(&model.Brands{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}
	return &emptypb.Empty{}, nil
}

// 更新品牌
func (c *CommodityServer) UpdateBrand(context.Context, *pb.BrandReq) (*empty.Empty, error) {
	brands := model.Brands{}
	if result := global.DB.First(&brands); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	if req.Name != "" {
		brands.Name = req.Name
	}
	if req.Logo != "" {
		brands.Logo = req.Logo
	}

	global.DB.Save(&brands)

	return &emptypb.Empty{}, nil
}
