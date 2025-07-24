package v1

import (
	"net/http"

	v1 "github.com/UnicomAI/wanwu/internal/bff-service/server/http/handler/v1"
	"github.com/UnicomAI/wanwu/internal/bff-service/server/http/middleware"
	"github.com/UnicomAI/wanwu/pkg/constant"
	mid "github.com/UnicomAI/wanwu/pkg/gin-util/mid-wrap"
	"github.com/gin-gonic/gin"
)

func registerAssistant(apiV1 *gin.RouterGroup) {
	mid.Sub("agent").Reg(apiV1, "/assistant/template/list", http.MethodGet, v1.GetAssistantTemplateList, "智能体模板列表")
	mid.Sub("agent").Reg(apiV1, "/assistant/template", http.MethodGet, v1.GetAssistantTemplate, "获取智能体模板")
	mid.Sub("agent").Reg(apiV1, "/assistant/template", http.MethodPost, v1.AssistantTemplateCreate, "复制智能体模板")

	mid.Sub("agent").Reg(apiV1, "/assistant", http.MethodPost, v1.AssistantCreate, "创建智能体")
	mid.Sub("agent").Reg(apiV1, "/assistant", http.MethodPut, v1.AssistantUpdate, "修改智能体基本信息")
	mid.Sub("agent").Reg(apiV1, "/assistant/config", http.MethodPut, v1.AssistantConfigUpdate, "修改智能体配置信息")
	mid.Sub("agent").Reg(apiV1, "/assistant", http.MethodGet, v1.GetAssistantInfo, "查看智能体详情")

	mid.Sub("agent").Reg(apiV1, "/assistant/workflow", http.MethodPost, v1.AssistantWorkFlowCreate, "添加工作流")
	mid.Sub("agent").Reg(apiV1, "/assistant/workflow", http.MethodDelete, v1.AssistantWorkFlowDelete, "删除工作流")
	mid.Sub("agent").Reg(apiV1, "/assistant/workflow/enable", http.MethodPut, v1.AssistantWorkFlowEnableSwitch, "启用/停用工作流")

	mid.Sub("agent").Reg(apiV1, "/assistant/action", http.MethodPost, v1.AssistantActionCreate, "添加action")
	mid.Sub("agent").Reg(apiV1, "/assistant/action", http.MethodDelete, v1.AssistantActionDelete, "删除action")
	mid.Sub("agent").Reg(apiV1, "/assistant/action", http.MethodPut, v1.AssistantActionUpdate, "编辑action")
	mid.Sub("agent").Reg(apiV1, "/assistant/action", http.MethodGet, v1.GetAssistantActionInfo, "查看智能体action详情")
	mid.Sub("agent").Reg(apiV1, "/assistant/action/enable", http.MethodPut, v1.AssistantActionEnableSwitch, "启用/停用action")

	mid.Sub("agent").Reg(apiV1, "/assistant/conversation", http.MethodPost, v1.ConversationCreate, "创建智能体对话")
	mid.Sub("agent").Reg(apiV1, "/assistant/conversation", http.MethodDelete, v1.ConversationDelete, "删除智能体对话")
	mid.Sub("agent").Reg(apiV1, "/assistant/conversation/list", http.MethodGet, v1.GetConversationList, "智能体对话列表")
	mid.Sub("agent").Reg(apiV1, "/assistant/conversation/detail", http.MethodGet, v1.GetConversationDetailList, "智能体对话详情历史列表")

	mid.Sub("agent").Reg(apiV1, "/assistant/stream", http.MethodPost, v1.AssistantConversionStream, "智能体流式问答", middleware.AppHistoryRecord("assistantId", constant.AppTypeAgent))
}
