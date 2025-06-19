package v1

import "github.com/gin-gonic/gin"

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

}

// DeleteMCP
//
//	@Tags			mcp
//	@Summary		删除自定义MCP
//	@Description	删除自定义MCP
//	@Accept			json
//	@Produce		json
//	@Param			mcpId	query		string	true	"mcpId"
//	@Success		200		{object}	response.Response{}
//	@Router			/mcp [delete]
func DeleteMCP(ctx *gin.Context) {

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

}

// GetMCPSelect
//
//	@Tags			mcp
//	@Summary		获取自定义MCP列表（用于下拉选择）
//	@Description	获取自定义MCP列表（用于下拉选择）
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.Response{data=response.MCPSelectList}
//	@Router			/mcp/select [get]
func GetMCPSelect(ctx *gin.Context) {

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
//	@Success		200			{object}	response.Response{data=response.ListResult{list=[]response.MCPTool}}
//	@Router			/mcp/tool/list [get]
func GetMCPTools(ctx *gin.Context) {

}
