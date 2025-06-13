package v1

import (
	"net/http"

	mid "github.com/UnicomAI/wanwu/internal/bff-service/pkg/gin-util/mid-wrap"
	v1 "github.com/UnicomAI/wanwu/internal/bff-service/server/http/handler/v1"
	"github.com/gin-gonic/gin"
)

func registerExploration(apiV1 *gin.RouterGroup) {
	mid.Sub("exploration").Reg(apiV1, "/exploration/app/history", http.MethodGet, v1.GetExplorationHistoryApp, "获取历史应用")
	mid.Sub("exploration").Reg(apiV1, "/exploration/app/list", http.MethodGet, v1.GetExplorationAppList, "获取探索广场应用")
	mid.Sub("exploration").Reg(apiV1, "/exploration/app/favorite", http.MethodPost, v1.ChangeExplorationAppFavorite, "更改App收藏状态")

	// app 相关接口
	mid.Sub("exploration").Reg(apiV1, "/appspace/app/url", http.MethodGet, v1.GetApiBaseUrl, "获取Api根地址")
	mid.Sub("exploration").Reg(apiV1, "/appspace/app/key", http.MethodPost, v1.GenApiKey, "生成ApiKey")
	mid.Sub("exploration").Reg(apiV1, "/appspace/app/key", http.MethodDelete, v1.DelApiKey, "删除ApiKey")
	mid.Sub("exploration").Reg(apiV1, "/appspace/app/key/list", http.MethodGet, v1.GetApiKeyList, "获取ApiKey列表")
	// rag 相关接口
	mid.Sub("exploration").Reg(apiV1, "/appspace/rag", http.MethodGet, v1.GetRag, "获取rag详情")
	// agent 相关接口
	mid.Sub("exploration").Reg(apiV1, "/assistant", http.MethodGet, v1.GetAssistantInfo, "查看智能体详情")
	mid.Sub("exploration").Reg(apiV1, "/assistant/conversation", http.MethodPost, v1.ConversationCreate, "创建智能体对话")
	mid.Sub("exploration").Reg(apiV1, "/assistant/conversation", http.MethodDelete, v1.ConversationDelete, "删除智能体对话")
	mid.Sub("exploration").Reg(apiV1, "/assistant/conversation/list", http.MethodGet, v1.GetConversationList, "智能体对话列表")
	mid.Sub("exploration").Reg(apiV1, "/assistant/conversation/detail", http.MethodGet, v1.GetConversationDetailList, "智能体对话详情历史列表")
}
