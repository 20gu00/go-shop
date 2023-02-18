package service

import (
	"commodity-rpc/dao"
	"commodity-rpc/model"
	"commodity-rpc/pb"
	"context"
	"encoding/json"
	"github.com/golang/protobuf/ptypes/empty"
)

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

// 商品分类列表
func (c *CommodityServer) AllCategoryList(context.Context, *empty.Empty) (*pb.CategoryListRes, error) {
	// 三级分类样式
	/*
		[
			{
				"id":xxx,
				"name":"",
				"level":1,
				"is_tab":false,
				"parent":13xxx,
				"sub_category":[
					"id":xxx,
					"name":"",
					"level":1,
					"is_tab":false,
					"sub_category":[]
				]
			}
		]
	*/
	var categorys []model.Category
	// 预加载方式先查询多表数据,在查询一表数据
	// 制定预加载的字段
	// SubCategory只拿到一级目录的子目录也就是两层,SubCategory.SubCategory这就是三层
	// 默认是将全部分类目录都拿出来,这里做个过滤,只拿一级目录,但子目录下的其实也获取到,但是不展开.不然整个表的分类不管几级的都拿出来
	dao.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categorys)
	//sql:先查询处一级的分类,那拿到外键也就父id,从而拿到对应的分类表的数据,以及对应的分类表下的父id,在拿到对应的分类表
	//自动过滤了deleted_at is null

	// api层访问直接返回这个json字符串,然后直接给前端这个字符串
	// 也就是层级全部做好,并且数据按层级加载好的字符串
	b, _ := json.Marshal(&categorys)
	return &pb.CategoryListRes{JsonData: string(b)}, nil
}
