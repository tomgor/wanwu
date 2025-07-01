package v1

import (
	"net/http"

	mid "github.com/UnicomAI/wanwu/internal/bff-service/pkg/gin-util/mid-wrap"
	v1 "github.com/UnicomAI/wanwu/internal/bff-service/server/http/handler/v1"
	"github.com/UnicomAI/wanwu/internal/bff-service/server/http/middleware"
	"github.com/UnicomAI/wanwu/pkg/constant"
	"github.com/gin-gonic/gin"
)

func registerRag(apiV1 *gin.RouterGroup) {
	mid.Sub("rag").Reg(apiV1, "/appspace/rag", http.MethodPost, v1.CreateRag, "创建rag")
	mid.Sub("rag").Reg(apiV1, "/appspace/rag", http.MethodPut, v1.UpdateRag, "修改rag基本信息")
	mid.Sub("rag").Reg(apiV1, "/appspace/rag/config", http.MethodPut, v1.UpdateRagConfig, "修改rag配置信息")
	mid.Sub("rag").Reg(apiV1, "/appspace/rag", http.MethodDelete, v1.DeleteRag, "删除rag")
	mid.Sub("rag").Reg(apiV1, "/appspace/rag", http.MethodGet, v1.GetRag, "获取rag详情")

	mid.Sub("rag").Reg(apiV1, "/rag/chat", http.MethodPost, v1.ChatRag, "rag流式接口", middleware.AppHistoryRecord("ragId", constant.AppTypeRag))
}
