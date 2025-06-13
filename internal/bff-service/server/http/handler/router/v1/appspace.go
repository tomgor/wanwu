package v1

import (
	"net/http"

	mid "github.com/UnicomAI/wanwu/internal/bff-service/pkg/gin-util/mid-wrap"
	v1 "github.com/UnicomAI/wanwu/internal/bff-service/server/http/handler/v1"
	"github.com/gin-gonic/gin"
)

func registerAppSpace(apiV1 *gin.RouterGroup) {
	mid.Sub("workspace.appspace").Reg(apiV1, "/appspace/app", http.MethodDelete, v1.DeleteAppSapceApp, "刪除应用")
	mid.Sub("workspace.appspace").Reg(apiV1, "/appspace/app/list", http.MethodGet, v1.GetAppSpaceAppList, "获取应用列表")
	mid.Sub("workspace.appspace").Reg(apiV1, "/appspace/app/publish", http.MethodPost, v1.PublishApp, "发布应用")

	mid.Sub("workspace.appspace").Reg(apiV1, "/appspace/app/url", http.MethodGet, v1.GetApiBaseUrl, "获取Api根地址")
	mid.Sub("workspace.appspace").Reg(apiV1, "/appspace/app/key", http.MethodPost, v1.GenApiKey, "生成ApiKey")
	mid.Sub("workspace.appspace").Reg(apiV1, "/appspace/app/key", http.MethodDelete, v1.DelApiKey, "删除ApiKey")
	mid.Sub("workspace.appspace").Reg(apiV1, "/appspace/app/key/list", http.MethodGet, v1.GetApiKeyList, "获取ApiKey列表")
}
