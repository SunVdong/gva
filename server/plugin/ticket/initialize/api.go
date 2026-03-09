package initialize

import (
	"context"
	model "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
)

func Api(ctx context.Context) {
	entities := []model.SysApi{
		{Path: "/ticket/scenic/createScenic", Description: "创建景区", ApiGroup: "景点门票", Method: "POST"},
		{Path: "/ticket/scenic/deleteScenic", Description: "删除景区", ApiGroup: "景点门票", Method: "DELETE"},
		{Path: "/ticket/scenic/deleteScenicByIds", Description: "批量删除景区", ApiGroup: "景点门票", Method: "DELETE"},
		{Path: "/ticket/scenic/updateScenic", Description: "更新景区", ApiGroup: "景点门票", Method: "PUT"},
		{Path: "/ticket/scenic/findScenic", Description: "查询景区", ApiGroup: "景点门票", Method: "GET"},
		{Path: "/ticket/scenic/getScenicList", Description: "景区列表", ApiGroup: "景点门票", Method: "GET"},
		{Path: "/ticket/scenicOpenTime/getByScenic", Description: "景区开放时间", ApiGroup: "景点门票", Method: "GET"},
		{Path: "/ticket/scenicOpenTime/save", Description: "保存景区开放时间", ApiGroup: "景点门票", Method: "POST"},
		{Path: "/ticket/product/createProduct", Description: "创建门票商品", ApiGroup: "景点门票", Method: "POST"},
		{Path: "/ticket/product/deleteProduct", Description: "删除门票商品", ApiGroup: "景点门票", Method: "DELETE"},
		{Path: "/ticket/product/deleteProductByIds", Description: "批量删除门票商品", ApiGroup: "景点门票", Method: "DELETE"},
		{Path: "/ticket/product/updateProduct", Description: "更新门票商品", ApiGroup: "景点门票", Method: "PUT"},
		{Path: "/ticket/product/findProduct", Description: "查询门票商品", ApiGroup: "景点门票", Method: "GET"},
		{Path: "/ticket/product/getProductList", Description: "门票商品列表", ApiGroup: "景点门票", Method: "GET"},
		{Path: "/ticket/sku/createSku", Description: "创建门票SKU", ApiGroup: "景点门票", Method: "POST"},
		{Path: "/ticket/sku/deleteSku", Description: "删除门票SKU", ApiGroup: "景点门票", Method: "DELETE"},
		{Path: "/ticket/sku/deleteSkuByIds", Description: "批量删除门票SKU", ApiGroup: "景点门票", Method: "DELETE"},
		{Path: "/ticket/sku/updateSku", Description: "更新门票SKU", ApiGroup: "景点门票", Method: "PUT"},
		{Path: "/ticket/sku/findSku", Description: "查询门票SKU", ApiGroup: "景点门票", Method: "GET"},
		{Path: "/ticket/sku/getSkuList", Description: "门票SKU列表", ApiGroup: "景点门票", Method: "GET"},
		{Path: "/ticket/rule/getByProduct", Description: "门票规则列表", ApiGroup: "景点门票", Method: "GET"},
		{Path: "/ticket/rule/save", Description: "保存门票规则", ApiGroup: "景点门票", Method: "POST"},
		{Path: "/ticket/audience/getBySku", Description: "适用人群列表", ApiGroup: "景点门票", Method: "GET"},
		{Path: "/ticket/audience/save", Description: "保存适用人群", ApiGroup: "景点门票", Method: "POST"},
		{Path: "/ticket/calendar/getBySku", Description: "日历库存", ApiGroup: "景点门票", Method: "GET"},
		{Path: "/ticket/calendar/set", Description: "设置日历库存", ApiGroup: "景点门票", Method: "POST"},
	}
	utils.RegisterApis(entities...)
}
