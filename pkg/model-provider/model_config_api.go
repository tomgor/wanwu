package mp

import (
	"context"
	"encoding/json"
	"fmt"
	mp_huoshan "github.com/UnicomAI/wanwu/pkg/model-provider/mp-huoshan"
	mp_ollama "github.com/UnicomAI/wanwu/pkg/model-provider/mp-ollama"
	mp_qwen "github.com/UnicomAI/wanwu/pkg/model-provider/mp-qwen"

	mp_common "github.com/UnicomAI/wanwu/pkg/model-provider/mp-common"
	mp_openai_compatible "github.com/UnicomAI/wanwu/pkg/model-provider/mp-openai-compatible"
	mp_yuanjing "github.com/UnicomAI/wanwu/pkg/model-provider/mp-yuanjing"
)

type ILLM interface {
	NewReq(req *mp_common.LLMReq) (mp_common.ILLMReq, error)
	ChatCompletions(ctx context.Context, req mp_common.ILLMReq, headers ...mp_common.Header) (mp_common.ILLMResp, <-chan mp_common.ILLMResp, error)
}

type IEmbedding interface {
	NewReq(req *mp_common.EmbeddingReq) (mp_common.IEmbeddingReq, error)
	Embeddings(ctx context.Context, req mp_common.IEmbeddingReq, headers ...mp_common.Header) (mp_common.IEmbeddingResp, error)
}

type IRerank interface {
	NewReq(req *mp_common.RerankReq) (mp_common.IRerankReq, error)
	Rerank(ctx context.Context, req mp_common.IRerankReq, headers ...mp_common.Header) (mp_common.IRerankResp, error)
}

// ToModelConfig 返回ILLM、IEmbedding或IRerank
func ToModelConfig(provider, modelType, cfg string) (interface{}, error) {
	if cfg == "" {
		return nil, nil
	}
	var ret interface{} // 前端需要的结构体
	switch provider {
	case ProviderOpenAICompatible:
		switch modelType {
		case ModelTypeLLM:
			ret = &mp_openai_compatible.LLM{}
		case ModelTypeRerank:
			ret = &mp_openai_compatible.Rerank{}
		case ModelTypeEmbedding:
			ret = &mp_openai_compatible.Embedding{}
		default:
			return nil, fmt.Errorf("invalid model type: %v", modelType)
		}
	case ProviderYuanJing:
		switch modelType {
		case ModelTypeLLM:
			ret = &mp_yuanjing.LLM{}
		case ModelTypeRerank:
			ret = &mp_yuanjing.Rerank{}
		case ModelTypeEmbedding:
			ret = &mp_yuanjing.Embedding{}
		default:
			return nil, fmt.Errorf("invalid model type: %v", modelType)
		}
	case ProviderHuoshan:
		switch modelType {
		case ModelTypeLLM:
			ret = &mp_huoshan.LLM{}
		case ModelTypeEmbedding:
			ret = &mp_huoshan.Embedding{}
		}
	case ProviderQwen:
		switch modelType {
		case ModelTypeLLM:
			ret = &mp_qwen.LLM{}
		case ModelTypeRerank:
			ret = &mp_qwen.Rerank{}
		case ModelTypeEmbedding:
			ret = &mp_qwen.Embedding{}
		}
	case ProviderOllama:
		switch modelType {
		case ModelTypeLLM:
			ret = &mp_ollama.LLM{}
		case ModelTypeEmbedding:
			ret = &mp_ollama.Embedding{}
		}
	default:
		return nil, fmt.Errorf("invalid provider: %v", modelType)
	}

	if err := json.Unmarshal([]byte(cfg), ret); err != nil {
		return nil, fmt.Errorf("unmarshal model config err: %v", err)
	}
	return ret, nil
}

type ProviderModelConfig struct {
	ProviderYuanJing         ProviderModelByYuanjing         `json:"providerYuanJing"`
	ProviderOpenAICompatible ProviderModelByOpenAICompatible `json:"providerOpenAICompatible"`
	ProviderHuoshan          ProviderModelByHuoshan          `json:"providerHuoshan"`
	ProviderQwen             ProviderModelByQwen             `json:"providerQwen"`
	ProviderOllama           ProviderModelByOllama           `json:"providerOllama"`
}

type ProviderModelByOpenAICompatible struct {
	Llm       mp_openai_compatible.LLM       `json:"llm"`
	Rerank    mp_openai_compatible.Rerank    `json:"rerank"`
	Embedding mp_openai_compatible.Embedding `json:"embedding"`
}

type ProviderModelByYuanjing struct {
	Llm       mp_yuanjing.LLM       `json:"llm"`
	Rerank    mp_yuanjing.Rerank    `json:"rerank"`
	Embedding mp_yuanjing.Embedding `json:"embedding"`
}

type ProviderModelByHuoshan struct {
	Llm       mp_huoshan.LLM       `json:"llm"`
	Embedding mp_huoshan.Embedding `json:"embedding"`
}

type ProviderModelByQwen struct {
	Llm       mp_qwen.LLM       `json:"llm"`
	Rerank    mp_qwen.Rerank    `json:"rerank"`
	Embedding mp_qwen.Embedding `json:"embedding"`
}

type ProviderModelByOllama struct {
	Llm       mp_ollama.LLM       `json:"llm"`
	Embedding mp_ollama.Embedding `json:"embedding"`
}
