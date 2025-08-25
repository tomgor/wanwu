package config

type AgentScopeWorkFlowSchemaResp struct {
	Code int `json:"code"`
	Data struct {
		Base64OpenAPISchema string `json:"base64OpenAPISchema"`
	} `json:"data"`
}
