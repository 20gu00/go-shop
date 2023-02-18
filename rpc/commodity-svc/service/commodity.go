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
)

type CommodityServer struct {
	pb.UnimplementedCommodityServer
}

// 商品列表
func (c *CommodityServer) CommodityList(ctx context.Context, req *pb.CommodityFilterReq) (*pb.CommodityListRes, error) {
	//列表过滤:关键词搜索、查询新品、查询热门商品、通过价格区间筛选， 通过商品分类筛选
	CommodityListRes := &pb.CommodityListRes{}
	commoditys := make([]model.Commodity, 0, 0)
	query := map[string]interface{}{}

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
	categoryIds := make([]interface{}, 0)

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

		type Result struct {
			ID int32
		}
		var results []Result
		global.DB.Model(model.Category{}).Raw(subQuery).Scan(&results)
		for _, re := range results {
			categoryIds = append(categoryIds, re.ID)
		}

		//生成terms查询
		q = q.Filter(elastic.NewTermsQuery("category_id", categoryIds...))
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
