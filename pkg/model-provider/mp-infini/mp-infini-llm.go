package mp_infini

import (
	"context"
	mp_common "github.com/UnicomAI/wanwu/pkg/model-provider/mp-common"
	"net/url"
)

type LLM struct {
	ApiKey          string `json:"apiKey"`                                                           // ApiKey
	EndpointUrl     string `json:"endpointUrl"`                                                      // 推理url
	FunctionCalling string `json:"functionCalling" validate:"oneof=noSupport toolCall functionCall"` // 函数调用是否支持
}

func (cfg *LLM) NewReq(req *mp_common.LLMReq) (mp_common.ILLMReq, error) {
	m, err := req.Data()
	if err != nil {
		return nil, err
	}
	return mp_common.NewLLMReq(m), nil
}

func (cfg *LLM) ChatCompletions(ctx context.Context, req mp_common.ILLMReq, headers ...mp_common.Header) (mp_common.ILLMResp, <-chan mp_common.ILLMResp, error) {
	return mp_common.ChatCompletions(ctx, "qwen", cfg.ApiKey, cfg.chatCompletionsUrl(), req, mp_common.NewLLMResp, headers...)
}

func (cfg *LLM) chatCompletionsUrl() string {
	ret, _ := url.JoinPath(cfg.EndpointUrl, "/chat/completions")
	return ret
}
