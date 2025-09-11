package request

type OpenAPIAgentCreateConversationRequest struct {
	Title string `json:"title"`
}

func (req *OpenAPIAgentCreateConversationRequest) Check() error {
	return nil
}

type OpenAPIAgentChatRequest struct {
	ConversationID string `json:"conversation_id" validate:"required"`
	Query          string `json:"query" validate:"required"`
	Stream         bool   `json:"stream"`
}

func (req *OpenAPIAgentChatRequest) Check() error {
	return nil
}

type OpenAPIRagChatRequest struct {
	Query  string `json:"query" validate:"required"`
	Stream bool   `json:"stream"`
}

func (req *OpenAPIRagChatRequest) Check() error {
	return nil
}

type OpenAPIWorkflowRunRequest struct {
	Input map[string]any `json:"input"`
}

func (req *OpenAPIWorkflowRunRequest) Check() error {
	return nil
}
