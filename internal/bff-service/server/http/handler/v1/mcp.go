package v1

import (
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/service"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	"github.com/gin-gonic/gin"
)

// GetMCPSquareDetail
//
//	@Tags			mcp.square
//	@Summary		获取广场MCP详情
//	@Description	获取广场MCP详情
//	@Accept			json
//	@Produce		json
//	@Param			mcpSquareId	query		string	true	"mcpSquareId"
//	@Success		200			{object}	response.Response{data=response.MCPSquareDetail}
//	@Router			/mcp/square [get]
func GetMCPSquareDetail(ctx *gin.Context) {
	resp, err := service.GetMCPSquareDetail(ctx, getUserID(ctx), getOrgID(ctx), ctx.Query("mcpSquareId"))
	gin_util.Response(ctx, resp, err)
}

// GetMCPSquareList
//
//	@Tags			mcp.square
//	@Summary		获取广场MCP列表
//	@Description	获取广场MCP列表
//	@Accept			json
//	@Produce		json
//	@Param			category	query		string	false	"mcp类型"	Enums(all,data,create,search)
//	@Param			name		query		string	false	"mcp名称"
//	@Success		200			{object}	response.Response{data=response.ListResult{list=[]response.MCPSquareInfo}}
//	@Router			/mcp/square/list [get]
func GetMCPSquareList(ctx *gin.Context) {
	resp, err := service.GetMCPSquareList(ctx, getUserID(ctx), getOrgID(ctx), ctx.Query("category"), ctx.Query("name"))
	gin_util.Response(ctx, resp, err)
}

// GetMCPSquareRecommends
//
//	@Tags			mcp.square
//	@Summary		获取广场MCP推荐列表
//	@Description	获取广场MCP推荐列表
//	@Accept			json
//	@Produce		json
//	@Param			mcpId		query		string	false	"mcpId"
//	@Param			mcpSquareId	query		string	false	"mcpSquareId"
//	@Success		200			{object}	response.Response{data=response.ListResult{list=[]response.MCPSquareInfo}}
//	@Router			/mcp/square/recommend [get]
func GetMCPSquareRecommends(ctx *gin.Context) {
	resp, err := service.GetMCPSquareList(ctx, getUserID(ctx), getOrgID(ctx), "", "")
	gin_util.Response(ctx, resp, err)
}

// CreateMCP
//
//	@Tags			mcp
//	@Summary		创建自定义MCP
//	@Description	创建自定义MCP
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.MCPCreate	true	"自定义MCP信息"
//	@Success		200		{object}	response.Response{}
//	@Router			/mcp [post]
func CreateMCP(ctx *gin.Context) {
	var req request.MCPCreate
	if !gin_util.Bind(ctx, &req) {
		return
	}
	gin_util.Response(ctx, nil, service.CreateMCP(ctx, getUserID(ctx), getOrgID(ctx), req))
}

// GetMCP
//
//	@Tags			mcp
//	@Summary		获取自定义MCP详情
//	@Description	获取自定义MCP详情
//	@Accept			json
//	@Produce		json
//	@Param			mcpId	query		string	true	"mcpId"
//	@Success		200		{object}	response.Response{data=response.MCPDetail}
//	@Router			/mcp [get]
func GetMCP(ctx *gin.Context) {
	resp, err := service.GetMCP(ctx, ctx.Query("mcpId"))
	gin_util.Response(ctx, resp, err)
}

// DeleteMCP
//
//	@Tags			mcp
//	@Summary		删除自定义MCP
//	@Description	删除自定义MCP
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.MCPIDReq	true	"mcpId"
//	@Success		200		{object}	response.Response{}
//	@Router			/mcp [delete]
func DeleteMCP(ctx *gin.Context) {
	var req request.MCPIDReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	gin_util.Response(ctx, nil, service.DeleteMCP(ctx, req.MCPID))
}

// GetMCPList
//
//	@Tags			mcp
//	@Summary		获取自定义MCP列表
//	@Description	获取自定义MCP列表
//	@Accept			json
//	@Produce		json
//	@Param			name	query		string	false	"mcp名称"
//	@Success		200		{object}	response.Response{data=response.ListResult{list=[]response.MCPInfo}}
//	@Router			/mcp/list [get]
func GetMCPList(ctx *gin.Context) {
	resp, err := service.GetMCPList(ctx, getUserID(ctx), getOrgID(ctx), ctx.Query("name"))
	gin_util.Response(ctx, resp, err)
}

// GetMCPSelect
//
//	@Tags			mcp
//	@Summary		获取自定义MCP列表（用于下拉选择）
//	@Description	获取自定义MCP列表（用于下拉选择）
//	@Accept			json
//	@Produce		json
//	@Param			name	query		string	false	"mcp名称"
//	@Success		200		{object}	response.Response{data=response.ListResult{list=[]response.MCPSelect}}
//	@Router			/mcp/select [get]
func GetMCPSelect(ctx *gin.Context) {
	resp, err := service.GetMCPSelect(ctx, getUserID(ctx), getOrgID(ctx), ctx.Query("name"))
	gin_util.Response(ctx, resp, err)
}

// GetMCPTools
//
//	@Tags			mcp
//	@Summary		获取MCP Tool列表
//	@Description	获取MCP Tool列表
//	@Accept			json
//	@Produce		json
//	@Param			mcpId		query		string	false	"mcpId(和serverUrl传一个)"
//	@Param			serverUrl	query		string	false	"serverUrl,就是sseUrl(和mcpId传一个)"
//	@Success		200			{object}	response.Response{data=response.MCPToolList}
//	@Router			/mcp/tool/list [get]
func GetMCPTools(ctx *gin.Context) {
	resp, err := service.GetMCPToolList(ctx, ctx.Query("mcpId"), ctx.Query("serverUrl"))
	gin_util.Response(ctx, resp, err)
}

// CreateCustomTool
//
//	@Tags			tool
//	@Summary		创建自定义工具
//	@Description	创建自定义工具
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.CustomToolCreate	true	"自定义工具信息"
//	@Success		200		{object}	response.Response{}
//	@Router			/tool/custom [post]
func CreateCustomTool(ctx *gin.Context) {
	var req request.CustomToolCreate
	if !gin_util.Bind(ctx, &req) {
		return
	}
	gin_util.Response(ctx, nil, service.CreateCustomTool(ctx, getUserID(ctx), getOrgID(ctx), req))
}

// GetCustomTool
//
//	@Tags			tool
//	@Summary		获取自定义工具详情
//	@Description	获取自定义工具详情
//	@Accept			json
//	@Produce		json
//	@Param			customToolId	query		string	true	"customToolId"
//	@Success		200				{object}	response.Response{data=response.CustomToolDetail}
//	@Router			/tool/custom [get]
func GetCustomTool(ctx *gin.Context) {
	resp, err := service.GetCustomToolInfo(ctx, getUserID(ctx), getOrgID(ctx), ctx.Query("customToolId"))
	gin_util.Response(ctx, resp, err)
}

// DeleteCustomTool
//
//	@Tags			tool
//	@Summary		删除自定义工具
//	@Description	删除自定义工具
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.CustomToolIDReq	true	"自定义工具ID"
//	@Success		200		{object}	response.Response{}
//	@Router			/tool/custom [delete]
func DeleteCustomTool(ctx *gin.Context) {
	var req request.CustomToolIDReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	gin_util.Response(ctx, nil, service.DeleteCustomTool(ctx, getUserID(ctx), getOrgID(ctx), req))
}

// UpdateCustomTool
//
//	@Tags			tool
//	@Summary		修改自定义工具
//	@Description	修改自定义工具
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.CustomToolUpdateReq	true	"自定义工具信息"
//	@Success		200		{object}	response.Response{}
//	@Router			/tool/custom [put]
func UpdateCustomTool(ctx *gin.Context) {
	var req request.CustomToolUpdateReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	gin_util.Response(ctx, nil, service.UpdateCustomTool(ctx, getUserID(ctx), getOrgID(ctx), req))
}

// GetCustomToolList
//
//	@Tags			tool
//	@Summary		获取自定义工具列表
//	@Description	获取自定义工具列表
//	@Accept			json
//	@Produce		json
//	@Param			name	query		string	false	"CustomTool名称"
//	@Success		200		{object}	response.Response{data=response.ListResult{list=[]response.CustomToolCell}}
//	@Router			/tool/custom/list [get]
func GetCustomToolList(ctx *gin.Context) {
	resp, err := service.GetCustomToolList(ctx, getUserID(ctx), getOrgID(ctx), ctx.Query("name"))
	gin_util.Response(ctx, resp, err)
}

// GetCustomToolSelect
//
//	@Tags			tool
//	@Summary		获取自定义工具列表（用于下拉选择）
//	@Description	获取自定义工具列表（用于下拉选择）
//	@Accept			json
//	@Produce		json
//	@Param			name	query		string	false	"CustomTool名称"
//	@Success		200		{object}	response.Response{data=response.ListResult{list=[]response.CustomToolSelect}}
//	@Router			/tool/custom/select [get]
func GetCustomToolSelect(ctx *gin.Context) {
	resp, err := service.GetCustomToolSelect(ctx, getUserID(ctx), getOrgID(ctx), ctx.Query("name"))
	gin_util.Response(ctx, resp, err)
}

// GetCustomToolSchemaAPI
//
//	@Tags			tool
//	@Summary		获取可用API列表（根据Schema）
//	@Description	解析自定义工具的Schema转换为API相关数据
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.CustomToolSchemaReq	true	"Schema格式数据"
//	@Success		200		{object}	response.Response{data=response.ListResult{list=[]response.CustomToolApiResponse}}
//	@Router			/tool/custom/schema [post]
func GetCustomToolSchemaAPI(ctx *gin.Context) {
	var req request.CustomToolSchemaReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	resp, err := service.GetCustomToolSchemaAPI(ctx, getUserID(ctx), getOrgID(ctx), req)
	gin_util.Response(ctx, resp, err)
}
