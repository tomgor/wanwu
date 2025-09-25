package service

import (
	"encoding/json"
	"strings"

	assistant_service "github.com/UnicomAI/wanwu/api/proto/assistant-service"
	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	mcp_service "github.com/UnicomAI/wanwu/api/proto/mcp-service"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

func CreateCustomTool(ctx *gin.Context, userID, orgID string, req request.CustomToolCreate) error {

	_, err := GetCustomToolSchemaAPI(ctx, userID, orgID, request.CustomToolSchemaReq{Schema: req.Schema})
	if err != nil {
		return err
	}

	_, err = mcp.CreateCustomTool(ctx.Request.Context(), &mcp_service.CreateCustomToolReq{
		Schema:        req.Schema,
		Name:          req.Name,
		Description:   req.Description,
		PrivacyPolicy: req.PrivacyPolicy,
		ApiAuth: &mcp_service.ApiAuthWebRequest{
			Type:             req.ApiAuth.Type,
			ApiKey:           req.ApiAuth.APIKey,
			CustomHeaderName: req.ApiAuth.CustomHeaderName,
			AuthType:         req.ApiAuth.AuthType,
		},
		Identity: &mcp_service.Identity{
			UserId: userID,
			OrgId:  orgID,
		},
	})
	return err
}

func GetCustomToolInfo(ctx *gin.Context, userID, orgID string, customToolId string) (*response.CustomToolDetail, error) {
	info, err := mcp.GetCustomToolInfo(ctx.Request.Context(), &mcp_service.GetCustomToolInfoReq{
		CustomToolId: customToolId,
	})
	if err != nil {
		return nil, err
	}
	apiList, err := GetCustomToolSchemaAPI(ctx, userID, orgID, request.CustomToolSchemaReq{Schema: info.Schema})
	if err != nil {
		return nil, err
	}
	return &response.CustomToolDetail{
		CustomToolId:  info.CustomToolId,
		ToolSquareID:  info.ToolSquareId,
		Schema:        info.Schema,
		Name:          info.Name,
		Description:   info.Description,
		PrivacyPolicy: info.PrivacyPolicy,
		ApiAuth: response.CustomToolApiAuthWebRequest{
			Type:             info.ApiAuth.Type,
			APIKey:           info.ApiAuth.ApiKey,
			CustomHeaderName: info.ApiAuth.CustomHeaderName,
			AuthType:         info.ApiAuth.AuthType,
		},
		ApiList: apiList.List.([]response.CustomToolApiResponse),
	}, nil
}

func DeleteCustomTool(ctx *gin.Context, userID, orgID string, req request.CustomToolIDReq) error {
	// 删除智能体AssistantCustom中记录
	_, err := assistant.AssistantCustomToolDeleteByCustomToolId(ctx.Request.Context(), &assistant_service.AssistantCustomToolDeleteByCustomToolIdReq{
		CustomToolId: req.CustomToolID,
	})
	if err != nil {
		return err
	}

	_, err = mcp.DeleteCustomTool(ctx.Request.Context(), &mcp_service.DeleteCustomToolReq{
		CustomToolId: req.CustomToolID,
	})
	return err
}

func UpdateCustomTool(ctx *gin.Context, userID, orgID string, req request.CustomToolUpdateReq) error {

	_, err := GetCustomToolSchemaAPI(ctx, userID, orgID, request.CustomToolSchemaReq{Schema: req.Schema})
	if err != nil {
		return err
	}

	_, err = mcp.UpdateCustomTool(ctx.Request.Context(), &mcp_service.UpdateCustomToolReq{
		CustomToolId: req.CustomToolID,
		Name:         req.Name,
		Description:  req.Description,
		ApiAuth: &mcp_service.ApiAuthWebRequest{
			Type:             req.ApiAuth.Type,
			ApiKey:           req.ApiAuth.APIKey,
			CustomHeaderName: req.ApiAuth.CustomHeaderName,
			AuthType:         req.ApiAuth.AuthType,
		},
		Schema:        req.Schema,
		PrivacyPolicy: req.PrivacyPolicy,
	})
	return err
}

func GetCustomToolList(ctx *gin.Context, userID, orgID, name string) (*response.ListResult, error) {
	resp, err := mcp.GetCustomToolList(ctx.Request.Context(), &mcp_service.GetCustomToolListReq{
		Identity: &mcp_service.Identity{
			UserId: userID,
			OrgId:  orgID,
		},
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	var list []response.CustomToolCell
	for _, item := range resp.List {
		list = append(list, response.CustomToolCell{
			CustomToolId: item.CustomToolId,
			Name:         item.Name,
			Description:  item.Description,
		})
	}
	return &response.ListResult{
		List:  list,
		Total: int64(len(list)),
	}, nil
}

func GetCustomToolSelect(ctx *gin.Context, userID, orgID, name string) (*response.ListResult, error) {
	resp, err := mcp.GetCustomToolList(ctx.Request.Context(), &mcp_service.GetCustomToolListReq{
		Identity: &mcp_service.Identity{
			UserId: userID,
			OrgId:  orgID,
		},
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	var list []response.CustomToolSelect
	for _, item := range resp.List {
		list = append(list, response.CustomToolSelect{
			UniqueId:     "tool-" + item.CustomToolId,
			CustomToolId: item.CustomToolId,
			Name:         item.Name,
			Description:  item.Description,
		})
	}
	return &response.ListResult{
		List:  list,
		Total: int64(len(list)),
	}, nil
}

func GetCustomToolSchemaAPI(ctx *gin.Context, userID, orgID string, req request.CustomToolSchemaReq) (*response.ListResult, error) {
	var list []response.CustomToolApiResponse
	infos, err := ParseOpenAPI(req.Schema)
	if err != nil {
		return nil, err
	} else {
		for _, item := range infos {
			list = append(list, response.CustomToolApiResponse{
				Name:   item.Name,
				Method: item.Method,
				Path:   item.Path,
			})
		}
	}
	return &response.ListResult{
		List:  list,
		Total: int64(len(list)),
	}, nil
}

// ParseOpenAPI 同时支持JSON和YAML格式的OpenAPI规范解析
func ParseOpenAPI(spec string) ([]response.CustomToolApiResponse, error) {
	var openAPI response.OpenAPI
	var err error

	// 先尝试JSON解析
	if err = json.Unmarshal([]byte(spec), &openAPI); err != nil {
		log.Errorf("JSON解析错误: %v\n", err)
		// JSON解析失败，尝试YAML解析
		if err = yaml.Unmarshal([]byte(spec), &openAPI); err != nil {
			log.Errorf("YAML解析错误: %v\n", err)
			return nil, grpc_util.ErrorStatus(errs.Code_BFFInvalidArg, "Schema is not valid OpenAPI format")
		}
	}

	var result []response.CustomToolApiResponse

	// 遍历所有路径和方法
	for path, item := range openAPI.Paths {
		addOperation(&result, item.Get, "GET", path)
		addOperation(&result, item.Post, "POST", path)
		addOperation(&result, item.Put, "PUT", path)
		addOperation(&result, item.Delete, "DELETE", path)
		addOperation(&result, item.Patch, "PATCH", path)
		addOperation(&result, item.Head, "HEAD", path)
		addOperation(&result, item.Options, "OPTIONS", path)
	}

	return result, nil
}

// 添加操作到结果列表
func addOperation(list *[]response.CustomToolApiResponse, op *response.Operation, method, path string) {
	if op == nil {
		return
	}

	name := op.OperationID
	if name == "" {
		name = op.Summary
	}
	if name == "" {
		name = "unknown"
	}

	*list = append(*list, response.CustomToolApiResponse{
		Name:   name,
		Method: method,
		Path:   path,
	})
}

func GetToolSquareDetail(ctx *gin.Context, userID, orgID, toolSquareID string) (*response.ToolSquareDetail, error) {
	resp, err := mcp.GetSquareTool(ctx.Request.Context(), &mcp_service.GetSquareToolReq{
		ToolSquareId: toolSquareID,
		Identity: &mcp_service.Identity{
			UserId: userID,
			OrgId:  orgID,
		},
	})
	if err != nil {
		return nil, err
	}
	return toToolSquareDetail(ctx, resp), nil
}

func toToolSquareDetail(ctx *gin.Context, toolSquare *mcp_service.SquareToolDetail) *response.ToolSquareDetail {
	ret := &response.ToolSquareDetail{
		ToolSquareInfo: toToolSquareInfo(ctx, toolSquare.Info),
		BuiltInTools: response.BuiltInTools{
			NeedApiKeyInput: toolSquare.BuiltInTools.NeedApiKeyInput,
			APIKey:          toolSquare.BuiltInTools.ApiKey,
			Detail:          toolSquare.BuiltInTools.Detail,
			ActionSum:       int64(toolSquare.BuiltInTools.ActionSum),
		},
	}
	for _, tool := range toolSquare.BuiltInTools.Tools {
		ret.BuiltInTools.Tools = append(ret.BuiltInTools.Tools, toMCPTool(tool))
	}
	return ret
}

func toToolSquareInfo(ctx *gin.Context, toolSquareInfo *mcp_service.ToolSquareInfo) response.ToolSquareInfo {
	return response.ToolSquareInfo{
		ToolSquareID: toolSquareInfo.ToolSquareId,
		Avatar:       cacheMCPAvatar(ctx, toolSquareInfo.AvatarPath),
		Name:         toolSquareInfo.Name,
		Desc:         toolSquareInfo.Desc,
		Tags:         getToolTags(toolSquareInfo.Tags),
	}
}

func GetToolSquareList(ctx *gin.Context, userID, orgID, name string) (*response.ListResult, error) {
	resp, err := mcp.GetSquareToolList(ctx.Request.Context(), &mcp_service.GetSquareToolListReq{
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	var list []response.ToolSquareInfo
	for _, item := range resp.Infos {
		list = append(list, response.ToolSquareInfo{
			ToolSquareID: item.ToolSquareId,
			Avatar:       cacheMCPAvatar(ctx, item.AvatarPath),
			Name:         item.Name,
			Desc:         item.Desc,
			Tags:         getToolTags(item.Tags),
		})
	}
	return &response.ListResult{
		List:  list,
		Total: int64(len(list)),
	}, nil
}

func UpdateBuiltInTool(ctx *gin.Context, userID, orgID string, req request.BuiltInToolReq) error {
	toolInfo, _ := mcp.GetCustomToolInfo(ctx.Request.Context(), &mcp_service.GetCustomToolInfoReq{
		ToolSquareId: req.ToolSquareID,
		Identity: &mcp_service.Identity{
			UserId: userID,
			OrgId:  orgID,
		},
	})
	if toolInfo == nil {
		return grpc_util.ErrorStatus(errs.Code_MCPGetCustomToolInfoErr, "tool not found")
	}
	if toolInfo.ApiAuth.ApiKey == "" {
		_, _ = mcp.CreateCustomTool(ctx.Request.Context(), &mcp_service.CreateCustomToolReq{
			ToolSquareId: req.ToolSquareID,
			ApiAuth: &mcp_service.ApiAuthWebRequest{
				ApiKey: req.APIKey,
			},
			Identity: &mcp_service.Identity{
				UserId: userID,
				OrgId:  orgID,
			},
		})
	}
	if toolInfo.CustomToolId != "" {
		_, err := mcp.UpdateCustomTool(ctx.Request.Context(), &mcp_service.UpdateCustomToolReq{
			CustomToolId: toolInfo.CustomToolId,
			ApiAuth: &mcp_service.ApiAuthWebRequest{
				ApiKey: req.APIKey,
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func getToolTags(tagString string) []string {
	if tagString == "" {
		return []string{}
	}
	return strings.Split(tagString, ",")
}
