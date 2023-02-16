package service

import (
	"commodity-rpc/dao"
	"commodity-rpc/model"
	"commodity-rpc/pb"
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// 轮播图列表
func (c *CommodityServer) BannerList(context.Context, *empty.Empty) (*pb.BannerListRes, error) {
	bannerListRes := pb.BannerListRes{}

	var banners []model.Banner
	result := dao.DB.Find(&banners)
	bannerListRes.Total = int32(result.RowsAffected)

	var bannerRes []*pb.BannerRes
	for _, banner := range banners {
		bannerRes = append(bannerRes, &pb.BannerRes{
			Id:    banner.ID,
			Image: banner.Picture,
			Index: banner.Index,
			Url:   banner.Url,
		})
	}

	bannerListRes.Data = bannerRes
	return &bannerListRes, nil
}

// 创建轮播图
func (c *CommodityServer) CreateBanner(ctx context.Context, req *pb.BannerReq) (*pb.BannerRes, error) {
	banner := model.Banner{}

	banner.Picture = req.Image
	banner.Index = req.Index
	banner.Url = req.Url

	dao.DB.Save(&banner)

	return &pb.BannerRes{Id: banner.ID}, nil
}

// 删除轮播图
func (c *CommodityServer) DeleteBanner(ctx context.Context, req *pb.BannerReq) (*empty.Empty, error) {
	if result := dao.DB.Delete(&model.Banner{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "此轮播图不存在")
	}
	return &emptypb.Empty{}, nil
}

// 更新轮播图
func (c *CommodityServer) UpdateBanner(ctx context.Context, req *pb.BannerReq) (*empty.Empty, error) {
	var banner model.Banner

	if result := dao.DB.First(&banner, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "轮播图不存在")
	}

	if req.Url != "" {
		banner.Url = req.Url
	}
	if req.Image != "" {
		banner.Picture = req.Image
	}
	if req.Index != 0 {
		banner.Index = req.Index
	}

	dao.DB.Save(&banner)

	return &emptypb.Empty{}, nil
}
