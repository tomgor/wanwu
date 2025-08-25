package v1

import (
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/service"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	"github.com/gin-gonic/gin"
)

// CreateWorkflow
//
//	@Tags		workflow
//	@Summary	创建Workflow
//	@Description
//	@Security	JWT
//	@Accept		json
//	@Produce	json
//	@Param		data	body		request.AppBriefConfig	true	"创建Workflow的请求参数"
//	@Success	200		{object}	response.Response{data=response.CozeWorkflowIDData}
//	@Router		/appspace/workflow [post]
func CreateWorkflow(ctx *gin.Context) {
	var req request.AppBriefConfig
	if !gin_util.Bind(ctx, &req) {
		return
	}
	resp, err := service.CreateWorkflow(ctx, getOrgID(ctx), req.Name, req.Desc)
	gin_util.Response(ctx, resp, err)
}

// CopyWorkflow
//
//	@Tags		workflow
//	@Summary	拷贝Workflow
//	@Description
//	@Security	JWT
//	@Accept		json
//	@Produce	json
//	@Param		data	body		request.WorkflowIDReq	true	"创建Workflow的请求参数"
//	@Success	200		{object}	response.Response{data=response.CozeWorkflowIDData}
//	@Router		/appspace/workflow/copy [post]
func CopyWorkflow(ctx *gin.Context) {
	var req request.WorkflowIDReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	resp, err := service.CopyWorkflow(ctx, getOrgID(ctx), req.WorkflowID)
	gin_util.Response(ctx, resp, err)
}
