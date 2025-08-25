package v1

import (
	"net/http"

	v1 "github.com/UnicomAI/wanwu/internal/bff-service/server/http/handler/v1"
	mid "github.com/UnicomAI/wanwu/pkg/gin-util/mid-wrap"
	"github.com/gin-gonic/gin"
)

func registerWorkflow(apiV1 *gin.RouterGroup) {
	mid.Sub("workflow").Reg(apiV1, "/appspace/workflow", http.MethodPost, v1.CreateWorkflow, "创建workflow")
	mid.Sub("workflow").Reg(apiV1, "/appspace/workflow/copy", http.MethodPost, v1.CopyWorkflow, "拷贝workflow")
	mid.Sub("workflow").Reg(apiV1, "/appspace/workflow/model/select/llm", http.MethodGet, v1.ListLlmModelsByWorkflow, "lm模型列表（用于workflow）")
}
