package service

import (
	paginate_list "commodity-rpc/common/paginate-list"
	"commodity-rpc/dao"
	"commodity-rpc/model"
	"commodity-rpc/pb"
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CommodityServer struct {
	pb.UnimplementedCommodityServer
}

// 商品列表
func (c *CommodityServer) CommodityList(ctx context.Context, req *pb.CommodityFilterReq) (*pb.CommodityListRes, error) {
	//列表过滤:关键词搜索、查询新品、查询热门商品、通过价格区间筛选， 通过商品分类筛选
	CommodityListRes := &pb.CommodityListRes{}
	commoditys := make([]model.Commodity, 0, 0)
	//query := map[string]interface{}{}

	// map查询,弊端就是key=value这种形式,如果是价格区间,那么就是>或者<了
	// 更灵活的方式是在前一个sql基础上进行sql而不是基于
	db := dao.DB.Model(&model.Commodity{})

	// 模糊查询 where("name LIKE ?","%"+a+"%")  虽然sql大小写不敏感,但关键字最好还是大写
	if req.Key != "" {
		//dao.DB.Where("name LIKE ?","%"+req.Key+"%").Find(CommodityListRes)
		//这样的搜索如果有多个条件,比如关键字同时还要是热门商品那就不符合了,要使用map查询
		//query["name"] = req.Key

		//不要修改全局的DB
		db = db.Where("name Like ?", "%"+req.Key+"%")

	}
	if req.Hot {
		//query["hot"] = req.Hot

		//db = db.Where("hot =true")
		//model也行
		db = db.Where(model.Commodity{IsHot: true})
	}
	if req.New {
		db = db.Where(model.Commodity{IsNew: true})
	}

	// 可能是proto的默认值
	if req.MinPri > 0 {
		db = db.Where("min_pri >= ?", req.MinPri)
	}
	if req.MaxPri > 0 {
		db = db.Where("max_pri <= ?", req.MaxPri)
	}

	// 品牌的id
	if req.Brand > 0 {
		db = db.Where("brand_id=?", req.Brand)
	}

	//通过category去查询商品
	// 子查询可以实现
	// select * from comodity where category_id in (select id from category where parent_category_id in (select id from category where parent_category_id=99999))

	var subQuery string

	// 有上层分类,上层分类的id
	if req.TopCategory > 0 {
		var category model.Category
		// ID
		if result := dao.DB.First(&category, req.TopCategory); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "商品分类不存在")
		}

		// 1 2 3级分类
		if category.Level == 1 {
			// 那就是要往下两层,拿到3级目录的商品分类
			subQuery = fmt.Sprintf("select id from category where parent_category_id in (select id from category WHERE parent_category_id=%d)", req.TopCategory)
		} else if category.Level == 2 {
			subQuery = fmt.Sprintf("select id from category WHERE parent_category_id=%d", req.TopCategory)
		} else if category.Level == 3 {
			// 直接分类id
			subQuery = fmt.Sprintf("select id from category WHERE id=%d", req.TopCategory)
		}

		// 商品表有个category_id
		db = db.Where(fmt.Sprintf("category_id in (%s)", subQuery))

	}

	// count要在scope之前做,不然就是pagenum
	var num int64
	db.Count(&num)
	CommodityListRes.Total = int32(num)

	if result := db.Preload("Category").Preload("Brand").Scopes(paginate_list.Paginate(int(req.Pages), int(req.PagePerNum))).Find(&commoditys); result.Error != nil {
		return nil, result.Error
	}

	for _, commodity := range commoditys {
		goodsInfoResponse := ModelToResponse(commodity)
		CommodityListRes.Data = append(CommodityListRes.Data, &goodsInfoResponse)
	}

	return CommodityListRes, nil
}

//获取商品详情
func (c *CommodityServer) GetCommodity(ctx context.Context, req *pb.CommodityInfoReq) (*pb.CommodityInfoRes, error) {
	var goods model.Commodity

	if result := dao.DB.Preload("Category").Preload("Brands").First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	goodsInfoResponse := ModelToResponse(goods)
	return &goodsInfoResponse, nil
}

// 批量获取商品信息(比如一个订单中有多个商品,多个商品的id)
func (c *CommodityServer) GetBatchCommodity(ctx context.Context, req *pb.BatchCommodityIdInfo) (*pb.CommodityListRes, error) {
	CommodityListResponse := &pb.CommodityListRes{}
	var commodity []model.Commodity

	//First(&commodity,[]{1,2,3}) 批量查询
	//Where([]{1,2,3})
	//实际上sql就是主键id in (1,2,3)
	result := dao.DB.Where(req.Id).Find(&commodity)
	for _, good := range commodity {
		goodsInfoResponse := ModelToResponse(good)
		CommodityListResponse.Data = append(CommodityListResponse.Data, &goodsInfoResponse)
	}
	CommodityListResponse.Total = int32(result.RowsAffected)
	return CommodityListResponse, nil
}

func ModelToResponse(commodity model.Commodity) pb.CommodityInfoRes {
	return pb.CommodityInfoRes{
		Id:          commodity.ID,
		Categoryid:  commodity.CategoryId,
		Name:        commodity.Name,
		CommoditySn: commodity.Sn,
		ReadNum:     commodity.ReadNum,
		SaleNum:     commodity.SaleNum,
		FavNum:      commodity.FavNum,
		CommonPri:   commodity.CommonPri,
		ShopPri:     commodity.LocalPri,
		EasyDesc:    commodity.EasyDesc,
		FreeShip:    commodity.IsFreeShip,
		FrontImage:  commodity.FrontImage,
		New:         commodity.IsNew,
		Hot:         commodity.IsHot,
		Sale:        commodity.IsSale,
		DescImages:  commodity.DescImages,
		Images:      commodity.Images,

		// 外键
		Category: &pb.CategoryEasyInfoRes{
			Id:   commodity.Category.ID,
			Name: commodity.Category.Name,
		},
		Brand: &pb.BrandInfoRes{
			Id:   commodity.Brand.ID,
			Name: commodity.Brand.Name,
			Logo: commodity.Brand.Logo,
		},
	}
}

// 商品增删改查
func (c *CommodityServer) CreateCommodity(ctx context.Context, req *pb.CreateCommodityReq) (*pb.CommodityInfoRes, error) {
	var category model.Category
	if result := dao.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brand
	if result := dao.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	//这里没有看到图片文件是如何上传， 在微服务中(grpc) 普通的文件上传已经不再使用
	goods := model.Commodity{
		// 外键部分
		Brand:      brand,
		BrandId:    brand.ID,
		Category:   category,
		CategoryId: category.ID,

		Name:       req.Name,
		Sn:         req.GoodsSn,
		CommonPri:  req.MarketPrice,
		LocalPri:   req.ShopPrice,
		EasyDesc:   req.GoodsBrief,
		IsFreeShip: req.ShipFree,
		Images:     req.Images,
		DescImages: req.DescImages,
		FrontImage: req.GoodsFrontImage,
		IsNew:      req.IsNew,
		IsHot:      req.IsHot,
		IsSale:     req.OnSale,
	}

	//srv之间互相调用了
	// 写入的操作都要小心,要做事务,可能多个服务同时调用你这个服务
	tx := dao.DB.Begin()
	result := tx.Save(&goods)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	return &pb.CommodityInfoRes{
		Id: goods.ID,
	}, nil
}

// 删除商品
func (c *CommodityServer) DeleteCommodity(ctx context.Context, req *pb.DeleteCommodityInfo) (*empty.Empty, error) {
	if result := dao.DB.Delete(&model.Commodity{Base: model.Base{ID: req.Id}}, req.Id); result.Error != nil {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	return &emptypb.Empty{}, nil
}

// 更新商品
func (c *CommodityServer) UpdateCommodity(ctx context.Context, req *pb.CreateCommodityReq) (*empty.Empty, error) {
	var goods model.Commodity

	if result := dao.DB.First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}

	var category model.Category
	if result := dao.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brand
	if result := dao.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	goods.Brand = brand
	goods.BrandId = brand.ID
	goods.Category = category
	goods.CategoryId = category.ID
	goods.Name = req.Name
	goods.Sn = req.GoodsSn
	goods.CommonPri = req.MarketPrice
	goods.LocalPri = req.ShopPrice
	goods.EasyDesc = req.GoodsBrief
	goods.IsFreeShip = req.ShipFree
	goods.Images = req.Images
	goods.DescImages = req.DescImages
	goods.FrontImage = req.GoodsFrontImage
	goods.IsNew = req.IsNew
	goods.IsHot = req.IsHot
	goods.IsSale = req.OnSale

	tx := dao.DB.Begin()
	result := tx.Save(&goods)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}
