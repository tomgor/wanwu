package v1

import (
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/service"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	mp "github.com/UnicomAI/wanwu/pkg/model-provider"
	"github.com/gin-gonic/gin"
)

// ListLlmModelsByWorkflow
//
//	@Tags		workflow
//	@Summary	llm模型列表（用于workflow）
//	@Description
//	@Security	JWT
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	response.Response{data=response.ListResult{list=response.CozeWorkflowModelInfo}}
//	@Router		/appspace/workflow/model/select/llm [get]
func ListLlmModelsByWorkflow(ctx *gin.Context) {
	resp, err := service.ListLlmModelsByWorkflow(ctx, getUserID(ctx), getOrgID(ctx), mp.ModelTypeLLM)
	gin_util.Response(ctx, resp, err)
}

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
	resp, err := service.CreateWorkflow(ctx, getOrgID(ctx), req.Name, req.Desc, req.Avatar.Key)
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
