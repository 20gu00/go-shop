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

// 创建品牌
func (c *CommodityServer) CreateBrand(ctx context.Context, req *pb.BrandReq) (*pb.BrandInfoRes, error) {
	// 查询品牌已经存在了
	if result := dao.DB.Where("name=?", req.Name).First(&model.Brand{}); result.RowsAffected == 1 {
		return nil, status.Errorf(codes.InvalidArgument, "该品牌已存在")
	}

	brand := &model.Brand{
		Name: req.Name,
		Logo: req.Logo,
	}

	dao.DB.Save(brand)

	return &pb.BrandInfoRes{Id: brand.ID}, nil
}

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

// 删除品牌
func (c *CommodityServer) DeleteBrand(ctx context.Context, req *pb.BrandReq) (*empty.Empty, error) {
	// 删除本身不论记录是否存在,不会报错
	if result := dao.DB.Delete(&model.Brand{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "该品牌不存在")
	}
	return &emptypb.Empty{}, nil
}

// 更新品牌
func (c *CommodityServer) UpdateBrand(ctx context.Context, req *pb.BrandReq) (*empty.Empty, error) {
	brands := model.Brand{}
	if result := dao.DB.First(&brands); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	// 可以更新品牌名称或者logo,避免用户输入空
	if req.Name != "" {
		brands.Name = req.Name
	}
	if req.Logo != "" {
		brands.Logo = req.Logo
	}

	dao.DB.Save(&brands)

	return &emptypb.Empty{}, nil
}
