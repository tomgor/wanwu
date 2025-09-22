package service

import (
	"fmt"
	net_url "net/url"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/bff-service/config"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	"github.com/UnicomAI/wanwu/pkg/constant"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	mp "github.com/UnicomAI/wanwu/pkg/model-provider"
	"github.com/UnicomAI/wanwu/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func ListLlmModelsByWorkflow(ctx *gin.Context, userId, orgId, modelT string) (*response.ListResult, error) {
	modelResp, err := ListTypeModels(ctx, userId, orgId, &request.ListTypeModelsRequest{ModelType: mp.ModelTypeLLM})
	if err != nil {
		return nil, err
	}
	var rets []response.CozeWorkflowModelInfo
	for _, modelBrief := range modelResp.List.([]*response.ModelBrief) {
		rets = append(rets, toModelBriefByWorkflow(modelBrief))
	}
	return &response.ListResult{
		List:  rets,
		Total: modelResp.Total,
	}, nil
}

func toModelBriefByWorkflow(modelBrief *response.ModelBrief) response.CozeWorkflowModelInfo {
	return response.CozeWorkflowModelInfo{
		ModelBrief:  *modelBrief,
		ModelParams: config.Cfg().Workflow.ModelParams,
	}
}

// ListWorkflow userID/orgID数据隔离，用于【工作流】
func ListWorkflow(ctx *gin.Context, orgID, name string) (*response.CozeWorkflowListData, error) {
	url, _ := net_url.JoinPath(config.Cfg().Workflow.Endpoint, config.Cfg().Workflow.ListUri)
	ret := &response.CozeWorkflowListResp{}
	if resp, err := resty.New().
		R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeaders(workflowHttpReqHeader(ctx)).
		SetQueryParams(map[string]string{
			"login_user_create": "true",
			"space_id":          orgID,
			"name":              name,
			"page":              "1",
			"size":              "99999",
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

// ListWorkflowByIDs 无userID或orgID隔离，用于【智能体选工作流】【应用广场】业务流程中
func ListWorkflowByIDs(ctx *gin.Context, name string, workflowIDs []string) (*response.CozeWorkflowListData, error) {
	var ids []string
	for _, workflowID := range workflowIDs {
		if _, err := util.I64(workflowID); err == nil {
			// AgentScope Workflow ID为uuid，将这部分脏数据过滤掉；Coze Workflow ID可转为数字
			ids = append(ids, workflowID)
		}
	}
	url, _ := net_url.JoinPath(config.Cfg().Workflow.Endpoint, config.Cfg().Workflow.ListUri)
	ret := &response.CozeWorkflowListResp{}
	request := resty.New().
		R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeaders(workflowHttpReqHeader(ctx)).
		SetQueryParams(map[string]string{
			"name": name,
			"page": "1",
			"size": "99999",
		})
	if len(ids) > 0 {
		request = request.SetBody(map[string]interface{}{
			"workflow_ids": ids,
		})
	}
	if resp, err := request.SetResult(ret).Post(url); err != nil {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_apps_list", err.Error())
	} else if resp.StatusCode() >= 300 {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_apps_list", fmt.Sprintf("[%v] %v", resp.StatusCode(), resp.String()))
	} else if ret.Code != 0 {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_apps_list", fmt.Sprintf("code %v msg %v", ret.Code, ret.Msg))
	}
	return ret.Data, nil
}

func CreateWorkflow(ctx *gin.Context, orgID, name, desc, iconUri string) (*response.CozeWorkflowIDData, error) {
	url, _ := net_url.JoinPath(config.Cfg().Workflow.Endpoint, config.Cfg().Workflow.CreateUri)
	ret := &response.CozeWorkflowIDResp{}
	if resp, err := resty.New().
		R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeaders(workflowHttpReqHeader(ctx)).
		SetQueryParams(map[string]string{
			"space_id": orgID,
			"name":     name,
			"desc":     desc,
			"icon_uri": iconUri,
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

func CopyWorkflow(ctx *gin.Context, orgID, workflowID string) (*response.CozeWorkflowIDData, error) {
	url, _ := net_url.JoinPath(config.Cfg().Workflow.Endpoint, config.Cfg().Workflow.CopyUri)
	ret := &response.CozeWorkflowIDResp{}
	if resp, err := resty.New().
		R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeaders(workflowHttpReqHeader(ctx)).
		SetQueryParams(map[string]string{
			"space_id":    orgID,
			"workflow_id": workflowID,
		}).
		SetResult(ret).
		Post(url); err != nil {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_app_copy", err.Error())
	} else if resp.StatusCode() >= 300 {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_app_copy", fmt.Sprintf("[%v] %v", resp.StatusCode(), resp.String()))
	} else if ret.Code != 0 {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_app_copy", fmt.Sprintf("code %v msg %v", ret.Code, ret.Msg))
	}
	return ret.Data, nil
}

func DeleteWorkflow(ctx *gin.Context, orgID, workflowID string) error {
	url, _ := net_url.JoinPath(config.Cfg().Workflow.Endpoint, config.Cfg().Workflow.DeleteUri)
	ret := &response.CozeWorkflowDeleteResp{}
	if resp, err := resty.New().
		R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeaders(workflowHttpReqHeader(ctx)).
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

func workflowHttpReqHeader(ctx *gin.Context) map[string]string {
	return map[string]string{
		"Authorization": ctx.GetHeader("Authorization"),
		"X-Org-Id":      ctx.GetHeader(gin_util.X_ORG_ID),
		"X-User-Id":     ctx.GetString(gin_util.USER_ID),
		"Content-Type":  "application/json",
	}
}

func cozeWorkflowInfo2Model(workflowInfo *response.CozeWorkflowListDataWorkflow) response.AppBriefInfo {
	return response.AppBriefInfo{
		AppId:     workflowInfo.WorkflowId,
		AppType:   constant.AppTypeWorkflow,
		Name:      workflowInfo.Name,
		Desc:      workflowInfo.Desc,
		Avatar:    cacheWorkflowAvatar(workflowInfo.URL),
		CreatedAt: util.Time2Str(workflowInfo.CreateTime * 1000),
		UpdatedAt: util.Time2Str(workflowInfo.UpdateTime * 1000),
	}
}
