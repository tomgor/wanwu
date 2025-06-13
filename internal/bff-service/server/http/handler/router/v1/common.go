package v1

import (
	"net/http"

	mid "github.com/UnicomAI/wanwu/internal/bff-service/pkg/gin-util/mid-wrap"
	v1 "github.com/UnicomAI/wanwu/internal/bff-service/server/http/handler/v1"
	"github.com/UnicomAI/wanwu/internal/bff-service/server/http/middleware"
	"github.com/UnicomAI/wanwu/pkg/constant"
	"github.com/gin-gonic/gin"
)

func registerCommon(apiV1 *gin.RouterGroup) {
	mid.Sub("common").Reg(apiV1, "/user/permission", http.MethodGet, v1.GetUserPermission, "获取用户权限")
	mid.Sub("common").Reg(apiV1, "/user/info", http.MethodGet, v1.GetUserInfo, "获取用户信息")
	mid.Sub("common").Reg(apiV1, "/org/select", http.MethodGet, v1.GetOrgSelect, "获取用户组织列表")
	mid.Sub("common").Reg(apiV1, "/user/password", http.MethodPut, v1.ChangeUserPassword, "修改用户密码（by 个人）")
	mid.Sub("common").Reg(apiV1, "/avatar", http.MethodPost, v1.UploadAvatar, "上传自定义图标")

	mid.Sub("common").Reg(apiV1, "/file/check", http.MethodGet, v1.CheckFile, "校验文件")
	mid.Sub("common").Reg(apiV1, "/file/check/list", http.MethodGet, v1.CheckFileList, "校验文件列表")
	mid.Sub("common").Reg(apiV1, "/file/upload", http.MethodPost, v1.UploadFile, "上传文件")
	mid.Sub("common").Reg(apiV1, "/file/merge", http.MethodPost, v1.MergeFile, "合并文件")
	mid.Sub("common").Reg(apiV1, "/file/clean", http.MethodPost, v1.CleanFile, "清除文件")
	mid.Sub("common").Reg(apiV1, "/file/delete", http.MethodDelete, v1.DeleteFile, "刪除文件")

	mid.Sub("common").Reg(apiV1, "/doc_center", http.MethodGet, v1.GetDocCenter, "获取文档中心路径")

	mid.Sub("common").Reg(apiV1, "/rag/chat", http.MethodPost, v1.ChatRag, "rag流式接口", middleware.AppHistoryRecord("ragId", constant.AppTypeRag))
	mid.Sub("common").Reg(apiV1, "/assistant/stream", http.MethodPost, v1.AssistantConversionStream, "智能体流式问答", middleware.AppHistoryRecord("assistantId", constant.AppTypeAgent))
}
