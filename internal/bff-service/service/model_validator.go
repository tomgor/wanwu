package service

import (
	"fmt"
	model_service "github.com/UnicomAI/wanwu/api/proto/model-service"
	mp "github.com/UnicomAI/wanwu/pkg/model-provider"
	mp_common "github.com/UnicomAI/wanwu/pkg/model-provider/mp-common"
	"github.com/gin-gonic/gin"
	"strings"
	"sync"
)

// 定义校验函数类型
type ModelValidator func(ctx *gin.Context, modelInfo *model_service.ModelInfo) error

// 校验器注册表
var validators = sync.OnceValue(func() map[string]ModelValidator {
	return map[string]ModelValidator{
		"llm":       ValidateLLMModel,
		"rerank":    ValidateRerankModel,
		"embedding": ValidateEmbeddingModel,
	}
})

// 统一校验入口
func ValidateModel(ctx *gin.Context, modelInfo *model_service.ModelInfo) error {
	validator, exists := validators()[strings.ToLower(modelInfo.ModelType)]
	if !exists {
		return fmt.Errorf("unsupported model type: %s", modelInfo.ModelType)
	}
	return validator(ctx, modelInfo)
}

func ValidateLLMModel(ctx *gin.Context, modelInfo *model_service.ModelInfo) error {
	llm, err := mp.ToModelConfig(modelInfo.Provider, modelInfo.ModelType, modelInfo.ProviderConfig)
	if err != nil {
		return err
	}
	iLLM, ok := llm.(mp.ILLM)
	if !ok {
		return fmt.Errorf("invalid provider")
	}
	// mock  request
	req := map[string]interface{}{
		"model": modelInfo.Model,
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": "hello",
			},
		},
		"stream": false,
	}
	llmReq := mp_common.NewLLMReq(req)
	_, _, err = iLLM.ChatCompletions(ctx.Request.Context(), llmReq)
	if err != nil {
		return fmt.Errorf("invalid resp: %v", err)
	}
	return nil
}

func ValidateEmbeddingModel(ctx *gin.Context, modelInfo *model_service.ModelInfo) error {
	embedding, err := mp.ToModelConfig(modelInfo.Provider, modelInfo.ModelType, modelInfo.ProviderConfig)
	if err != nil {
		return err
	}
	iEmbedding, ok := embedding.(mp.IEmbedding)
	if !ok {
		return fmt.Errorf("invalid provider")
	}
	// mock  request
	req := &mp_common.EmbeddingReq{
		Model: modelInfo.Model,
		Input: []string{"你好"},
	}
	embeddingReq, err := iEmbedding.NewReq(req)
	if err != nil {
		return err
	}
	_, err = iEmbedding.Embeddings(ctx.Request.Context(), embeddingReq)
	if err != nil {
		{
			return fmt.Errorf("invalid resp: %v", err)
		}
	}
	return nil
}

func ValidateRerankModel(ctx *gin.Context, modelInfo *model_service.ModelInfo) error {
	rerank, err := mp.ToModelConfig(modelInfo.Provider, modelInfo.ModelType, modelInfo.ProviderConfig)
	if err != nil {
		return err
	}
	iRerank, ok := rerank.(mp.IRerank)
	if !ok {
		return fmt.Errorf("invalid provider")
	}
	// mock  request
	req := &mp_common.RerankReq{
		Model: modelInfo.Model,
		Query: "乌萨奇",
		Documents: []string{
			"乌萨奇",
			"尖尖我噶奶～",
		},
	}
	rerankReq, err := iRerank.NewReq(req)
	if err != nil {
		return err
	}
	_, err = iRerank.Rerank(ctx.Request.Context(), rerankReq)
	if err != nil {
		return fmt.Errorf("invalid resp: %v", err)
	}
	return nil
}
