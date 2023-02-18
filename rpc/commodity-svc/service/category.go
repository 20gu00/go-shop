package service

import (
	"commodity-rpc/dao"
	"commodity-rpc/model"
	"commodity-rpc/pb"
	"context"
	"encoding/json"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// 获取商品的子分类
func (c *CommodityServer) SubCategoryList(ctx context.Context, req *pb.CategoryListReq) (*pb.SubCategoryListRes, error) {
	categoryListRes := pb.SubCategoryListRes{}
	var category model.Category
	// 商品不存在
	// where 或者 category{ID:req.Id}
	if result := dao.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "该商品分类不存在")
	}

	categoryListRes.Info = &pb.CategoryInfoRes{
		Id:             category.ID,
		Name:           category.Name,
		Level:          category.Level,
		IsTab:          category.IsTop,
		ParentCategory: category.ParCategoryId,
	}

	//判断是几级目录
	//这里查询的时候不用preload也行,因为就两层,直接拿到两层子分类
	//preloads := "SubCategory"
	//if category.Level == 1 {
	//	preloads = "SubCategory.SubCategory"
	//}

	var subCategorys []model.Category
	var subCategoryResponse []*pb.CategoryInfoRes

	//根据外键查询到所有的子分类
	dao.DB.Where(&model.Category{ParCategoryId: req.Id}).Find(&subCategorys)
	for _, subCategory := range subCategorys {
		subCategoryResponse = append(subCategoryResponse, &pb.CategoryInfoRes{
			Id:             subCategory.ID,
			Name:           subCategory.Name,
			Level:          subCategory.Level,
			IsTab:          subCategory.IsTop,
			ParentCategory: subCategory.ParCategoryId,
		})
	}

	categoryListRes.SubCategorys = subCategoryResponse
	return &categoryListRes, nil
}

// 增加商品

func (c *CommodityServer) CreateCategory(ctx context.Context, req *pb.CategoryInfoReq) (*pb.CategoryInfoRes, error) {
	category := model.Category{}
	cMap := map[string]interface{}{} // 初始化
	cMap["name"] = req.Name
	cMap["level"] = req.Level
	cMap["is_tab"] = req.IsTab

	// 不是一级分类
	// 得判断父分类存在,这里就不判断了选择信任传递进来的数据
	// 由获取到父分类的方法,通过api层调用,前端处理即可,也就是信任前端传递进来的数据
	if req.Level != 1 {
		cMap["parent_category_id"] = req.ParentCategory
	}
	dao.DB.Model(&model.Category{}).Create(cMap)
	return &pb.CategoryInfoRes{Id: category.ID}, nil
}

// 删除商品
func (c *CommodityServer) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryReq) (*empty.Empty, error) {
	if result := dao.DB.Delete(&model.Category{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	return &emptypb.Empty{}, nil
}

// 更新商品
func (c *CommodityServer) UpdateCategory(ctx context.Context, req *pb.CategoryInfoReq) (*empty.Empty, error) {
	var category model.Category
	if result := dao.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}

	if req.Name != "" {
		category.Name = req.Name
	}

	// proto由零值,主要判断,不然填入零值,不符合我们的预期会报错
	if req.ParentCategory != 0 {
		category.ParCategoryId = req.ParentCategory
	}
	if req.Level != 0 {
		category.Level = req.Level
	}
	if req.IsTab {
		category.IsTop = req.IsTab
	}

	dao.DB.Save(&category)
	return &emptypb.Empty{}, nil
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
	// SubCategory只拿到一级目录的子目录也就是两层,SubCategory.SubCategory这就是三层(二层也加载)
	// 默认是将全部分类目录都拿出来(因为是自身嵌套),这里做个过滤,只拿一级目录,但子目录下的其实也获取到,但是不展开.不然整个表的分类不管几级的都拿出来
	dao.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categorys)
	//sql:先查询处一级的分类,那拿到外键也就父id,从而拿到对应的分类表的数据,以及对应的分类表下的父id,在拿到对应的分类表
	//自动过滤了deleted_at is null

	// api层访问直接返回这个json字符串,然后直接给前端这个字符串
	// 也就是层级全部做好,并且数据按层级加载好的字符串
	b, _ := json.Marshal(&categorys)
	return &pb.CategoryListRes{JsonData: string(b)}, nil
}
