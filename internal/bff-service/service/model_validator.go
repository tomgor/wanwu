package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	model_service "github.com/UnicomAI/wanwu/api/proto/model-service"
	"github.com/UnicomAI/wanwu/pkg/log"
	mp "github.com/UnicomAI/wanwu/pkg/model-provider"
	mp_common "github.com/UnicomAI/wanwu/pkg/model-provider/mp-common"
	"github.com/gin-gonic/gin"
)

// 定义校验函数类型
type ModelValidator func(ctx *gin.Context, modelInfo *model_service.ModelInfo) error

// 校验器注册表
var validators = sync.OnceValue(func() map[string]ModelValidator {
	return map[string]ModelValidator{
		mp.ModelTypeLLM:       ValidateLLMModel,
		mp.ModelTypeRerank:    ValidateRerankModel,
		mp.ModelTypeEmbedding: ValidateEmbeddingModel,
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
	var stream bool = false
	req := &mp_common.LLMReq{
		Model: modelInfo.Model,
		Messages: []mp_common.OpenAIMsg{
			{
				Role:    mp_common.MsgRoleUser,
				Content: "几点了",
			},
		},
		Stream: &stream,
	}
	// ToolCall 校验
	var result map[string]interface{}
	err = json.Unmarshal([]byte(modelInfo.ProviderConfig), &result)
	if err != nil {
		return err
	}
	fc, ok := result["functionCalling"].(mp_common.FCType)
	if ok && fc == mp_common.FCTypeToolCall {
		tools := []mp_common.OpenAITool{
			{
				Type: mp_common.ToolTypeFunction,
				Function: &mp_common.OpenAIFunction{
					Name:        "get_current_time",
					Description: "当你想知道现在的时间时非常有用。",
					Parameters: &mp_common.OpenAIFunctionParameters{
						Type:       "object",
						Properties: map[string]mp_common.OpenAIFunctionParametersProperty{},
					},
				},
			},
		}
		req.Tools = tools
		llmReq, err := iLLM.NewReq(req)
		if err != nil {
			return err
		}
		resp, _, err := iLLM.ChatCompletions(context.Background(), llmReq)
		if err != nil {
			return err
		}
		openAIResp, ok := resp.ConvertResp()
		if !ok {
			return fmt.Errorf("invalid resp: %v", err)
		}
		if len(openAIResp.Choices) == 0 || openAIResp.Choices[0].Message.ToolCalls == nil {
			return fmt.Errorf("model does not support toolcall functionality")
		} else {
			data, _ := json.MarshalIndent(openAIResp.Choices[0].Message.ToolCalls, "", "  ")
			log.Debugf("tool call: %v", string(data))
		}
		return nil
	}
	llmReq, err := iLLM.NewReq(req)
	if err != nil {
		return err
	}
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
