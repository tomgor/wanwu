package v1

import (
	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/bff-service/service"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	"github.com/UnicomAI/wanwu/pkg/util"
	"github.com/gin-gonic/gin"
)

// GetUserPermission
//
//	@Tags		common
//	@Summary	获取用户权限
//	@Security	JWT
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	response.Response{data=response.UserPermission}
//	@Router		/user/permission [get]
func GetUserPermission(ctx *gin.Context) {
	resp, err := service.GetUserPermission(ctx, getUserID(ctx), getOrgID(ctx))
	gin_util.Response(ctx, resp, err)
}

// GetUserInfo
//
//	@Tags		common
//	@Summary	获取用户信息
//	@Security	JWT
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	response.Response{data=response.UserInfo}
//	@Router		/user/info [get]
func GetUserInfo(ctx *gin.Context) {
	resp, err := service.GetUserInfo(ctx, getUserID(ctx), getOrgID(ctx))
	gin_util.Response(ctx, resp, err)
}

// GetOrgSelect
//
//	@Tags		common
//	@Summary	获取用户组织列表（用于下拉选择）
//	@Security	JWT
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	response.Response{data=response.Select}
//	@Router		/org/select [get]
func GetOrgSelect(ctx *gin.Context) {
	resp, err := service.GetOrgSelect(ctx, getUserID(ctx))
	gin_util.Response(ctx, resp, err)
}

// UploadAvatar
//
//	@Tags		common
//	@Summary	上传自定义图标
//	@Security	JWT
//	@Accept		multipart/form-data
//	@Produce	json
//	@Param		avatar	formData	file	true	"自定义图标（JPG/JPEG/PNG）"
//	@Success	200		{object}	response.Response{data=request.Avatar}
//	@Router		/avatar [post]
func UploadAvatar(ctx *gin.Context) {
	avatarFile, err := ctx.FormFile("avatar")
	if err != nil {
		gin_util.ResponseErrCodeKey(ctx, err_code.Code_BFFGeneral, "bff_avatar_upload_error", err.Error())
		return
	}
	avatarObjectPath, err := service.UploadAvatar(ctx, avatarFile)
	if err != nil {
		gin_util.Response(ctx, nil, err)
		return
	}
	avatar := service.CacheAvatar(ctx, avatarObjectPath)
	gin_util.Response(ctx, avatar, nil)
}

// SearchDocCenter
//
//	@Tags		common
//	@Summary	查找文档中心内容
//	@Security	JWT
//	@Accept		json
//	@Produce	json
//	@Param		content	query		string	true	"搜索关键字"
//	@Success	200		{object}	response.Response{data=[]response.DocSearchResp}
//	@Router		/doc_center/search [get]
func SearchDocCenter(ctx *gin.Context) {
}

// GetDocCenterMenu
//
//	@Tags		common
//	@Summary	获取文档中心目录
//	@Security	JWT
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	response.Response{data=[]response.DocMenu}
//	@Router		/doc_center/menu [get]
func GetDocCenterMenu(ctx *gin.Context) {
}

// GetDocCenterMarkdown
//
//	@Tags		common
//	@Summary	获取文档中心Markdown文件内容
//	@Security	JWT
//	@Accept		json
//	@Produce	json
//	@Param		path	query		string	true	"目录path"
//	@Success	200		{object}	response.Response{data=string}
//	@Router		/doc_center/markdown [get]
func GetDocCenterMarkdown(ctx *gin.Context) {
}

// --- internal ---

// 获取当前用户ID
func getUserID(ctx *gin.Context) string {
	return ctx.GetString(gin_util.USER_ID)
}

// 获取当前组织ID
func getOrgID(ctx *gin.Context) string {
	return ctx.GetHeader(gin_util.X_ORG_ID)
}

// 获取当前系统语言
func getLanguage(ctx *gin.Context) string {
	return ctx.GetHeader(gin_util.X_LANGUAGE)
}

// 当前用户是否是当前组织内置管理员角色
func isAdmin(ctx *gin.Context) bool {
	return ctx.GetBool(gin_util.IS_ADMIN)
}

// 当前组织是否是内置顶级【系统】组织
func isSystem(ctx *gin.Context) bool {
	return ctx.GetBool(gin_util.IS_SYSTEM)
}

func getPageNo(ctx *gin.Context) int32 {
	return util.MustI32(ctx.Query(gin_util.PageNo))
}

func getPageSize(ctx *gin.Context) int32 {
	return util.MustI32(ctx.Query(gin_util.PageSize))
}
