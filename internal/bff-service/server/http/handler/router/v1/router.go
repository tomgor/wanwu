package v1

import (
	"github.com/gin-gonic/gin"
)

func Register(apiV1 *gin.RouterGroup) {
	// guest
	registerGuest(apiV1)

	// common
	registerCommon(apiV1)

	// permission
	registerPermission(apiV1)

	// model
	registerModel(apiV1)

	// workspace.appspace
	registerAppSpace(apiV1)

	// exploration
	registerExploration(apiV1)

	// rag
	registerRag(apiV1)

	// knowledge
	registerKnowledge(apiV1)

	// assistant
	registerAssistant(apiV1)

	//callback
	registerV1Callback(apiV1)
}
