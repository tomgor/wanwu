package service

import (
	"fmt"
	net_url "net/url"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/bff-service/config"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	"github.com/UnicomAI/wanwu/pkg/constant"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	"github.com/UnicomAI/wanwu/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func ListWorkflow(ctx *gin.Context, userID, orgID, name string) (*response.CozeWorkflowListData, error) {
	url, _ := net_url.JoinPath(config.Cfg().Workflow.Endpoint, config.Cfg().Workflow.ListUri)
	ret := &response.CozeWorkflowListResp{}
	if resp, err := resty.New().
		R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeaders(workflowHttpReqHeader(ctx, userID, orgID)).
		SetQueryParams(map[string]string{
			"space_id": orgID,
			"name":     name,
			"page":     "1",
			"size":     "100",
		}).
		SetResult(ret).
		Post(url); err != nil {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_apps_list", err.Error())
	} else if resp.StatusCode() >= 300 {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_apps_list", fmt.Sprintf("[%v] %v", resp.StatusCode(), resp.String()))
	} else if ret.Code != 0 {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_apps_list", fmt.Sprintf("code %v msg %v", ret.Code, ret.Msg))
	}
	return ret.Data, nil
}

func CreateWorkflow(ctx *gin.Context, userID, orgID, name, desc string) (*response.CozeWorkflowCreateData, error) {
	url, _ := net_url.JoinPath(config.Cfg().Workflow.Endpoint, config.Cfg().Workflow.CreateUri)
	ret := &response.CozeWorkflowCreateResp{}
	if resp, err := resty.New().
		R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeaders(workflowHttpReqHeader(ctx, userID, orgID)).
		SetQueryParams(map[string]string{
			"space_id": orgID,
			"name":     name,
			"desc":     desc,
		}).
		SetResult(ret).
		Post(url); err != nil {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_app_create", err.Error())
	} else if resp.StatusCode() >= 300 {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_app_create", fmt.Sprintf("[%v] %v", resp.StatusCode(), resp.String()))
	} else if ret.Code != 0 {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_app_create", fmt.Sprintf("code %v msg %v", ret.Code, ret.Msg))
	}
	return ret.Data, nil
}

func DeleteWorkflow(ctx *gin.Context, userID, orgID, workflowID string) error {
	url, _ := net_url.JoinPath(config.Cfg().Workflow.Endpoint, config.Cfg().Workflow.DeleteUri)
	ret := &response.CozeWorkflowDeleteResp{}
	if resp, err := resty.New().
		R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeaders(workflowHttpReqHeader(ctx, userID, orgID)).
		SetQueryParams(map[string]string{
			"workflow_id": workflowID,
			"space_id":    orgID,
		}).
		SetResult(ret).
		Post(url); err != nil {
		return grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_app_delete", err.Error())
	} else if resp.StatusCode() >= 300 {
		return grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_app_delete", fmt.Sprintf("[%v] %v", resp.StatusCode(), resp.String()))
	} else if ret.Code != 0 || (ret.Data != nil && ret.Data.Status != 0) {
		return grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_app_delete", fmt.Sprintf("code %v msg %v status %v", ret.Code, ret.Msg, ret.Data.GetStatus()))
	}
	return nil
}

// --- internal ---

func workflowHttpReqHeader(ctx *gin.Context, userID, orgID string) map[string]string {
	return map[string]string{
		"Authorization": ctx.GetHeader("Authorization"),
		"Content-Type":  "application/json",
		"X-Org-Id":      orgID,
		"X-User-Id":     userID,
	}
}

func cozeWorkflowInfo2Model(workflowInfo response.CozeWorkflowListDataWorkflow) response.AppBriefInfo {
	return response.AppBriefInfo{
		AppId:     workflowInfo.WorkflowId,
		AppType:   constant.AppTypeWorkflow,
		Name:      workflowInfo.Name,
		Desc:      workflowInfo.Desc,
		CreatedAt: util.Time2Str(workflowInfo.CreateTime * 1000),
		UpdatedAt: util.Time2Str(workflowInfo.UpdateTime * 1000),
	}
}
