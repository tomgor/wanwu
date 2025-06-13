package v1

import (
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	gin_util "github.com/UnicomAI/wanwu/internal/bff-service/pkg/gin-util"
	"github.com/UnicomAI/wanwu/internal/bff-service/service"
	"github.com/gin-gonic/gin"
)

// GetExplorationHistoryApp
//
//	@Tags			exploration
//	@Summary		获取历史应用
//	@Description	获取历史应用
//	@Produce		json
//	@Success		200	{object}	response.Response{data=response.ListResult{list=[]response.ExplorationAppInfo}}
//	@Router			/exploration/app/history [get]
func GetExplorationHistoryApp(ctx *gin.Context) {
	resp, err := service.GetExplorationHistoryApp(ctx, getUserID(ctx))
	gin_util.Response(ctx, resp, err)
}

// GetExplorationAppList
//
//	@Tags			exploration
//	@Summary		获取探索广场应用
//	@Description	获取探索广场应用
//	@Accept			json
//	@Produce		json
//	@Param			data	query		request.GetExplorationAppListRequest	true	"获取探索广场应用参数"
//	@Success		200		{object}	response.Response{data=response.ListResult{list=[]response.ExplorationAppInfo}}
//	@Router			/exploration/app/list [get]
func GetExplorationAppList(ctx *gin.Context) {
	var req request.GetExplorationAppListRequest
	if !gin_util.BindQuery(ctx, &req) {
		return
	}
	resp, err := service.GetExplorationAppList(ctx, getUserID(ctx), req)
	gin_util.Response(ctx, resp, err)
}

// ChangeExplorationAppFavorite
//
//	@Tags			exploration
//	@Summary		更改App收藏状态
//	@Description	更改App收藏状态
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.ChangeExplorationAppFavoriteRequest	true	"更改App收藏状态参数"
//	@Success		200		{object}	response.Response
//	@Router			/exploration/app/favorite [post]
func ChangeExplorationAppFavorite(ctx *gin.Context) {
	var req request.ChangeExplorationAppFavoriteRequest
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.ChangeExplorationAppFavorite(ctx, getUserID(ctx), getOrgID(ctx), req)
	gin_util.Response(ctx, nil, err)
}
