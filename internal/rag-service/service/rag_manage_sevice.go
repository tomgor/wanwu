package service

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/UnicomAI/wanwu/internal/rag-service/config"
	"net/http"
	"strconv"
	"time"

	knowledgeBase_service "github.com/UnicomAI/wanwu/api/proto/knowledgebase-service"
	rag_service "github.com/UnicomAI/wanwu/api/proto/rag-service"
	"github.com/UnicomAI/wanwu/internal/rag-service/client/model"
	http_client "github.com/UnicomAI/wanwu/internal/rag-service/pkg/http-client"
	"github.com/UnicomAI/wanwu/pkg/log"
)

const (
	DefaultTemperature      = 0.14
	DefaultTopP             = 0.85
	DefaultFrequencyPenalty = 1.1
)

type RagChatParams struct {
	KnowledgeBase     []string         `json:"knowledgeBase"`
	Question          string           `json:"question"`
	Threshold         float32          `json:"threshold"`
	TopK              int32            `json:"topK"`
	Stream            bool             `json:"stream"`
	Chichat           bool             `json:"chichat"` // 当知识库召回结果为空时是否使用默认话术（兜底），默认为true
	RerankModelId     string           `json:"rerank_model_id"`
	CustomModelInfo   *CustomModelInfo `json:"custom_model_info"`
	History           []*HistoryItem   `json:"history"`
	MaxHistory        int32            `json:"max_history"`
	RewriteQuery      bool             `json:"rewrite_query"`   // 是否query改写
	RerankMod         string           `json:"rerank_mod"`      // rerank_model:重排序模式，weighted_score：权重搜索
	RetrieveMethod    string           `json:"retrieve_method"` // hybrid_search:混合搜索， semantic_search:向量搜索， full_text_search：文本搜索
	Weight            *WeightParams    `json:"weights"`         // 权重搜索下的权重配置
	Temperature       float32          `json:"temperature"`
	TopP              float32          `json:"top_p"`              // 多样性
	RepetitionPenalty float32          `json:"repetition_penalty"` // 重复惩罚/频率惩罚
	ReturnMeta        bool             `json:"return_meta"`        // 是否返回元数据
	AutoCitation      bool             `json:"auto_citation"`      // 是否自动角标
}

type WeightParams struct {
	VectorWeight float32 `json:"vector_weight"` //语义权重
	TextWeight   float32 `json:"text_weight"`   //关键字权重
}

type CustomModelInfo struct {
	LlmModelID string `json:"llm_model_id"`
}

type HistoryItem struct {
	Query       string `json:"query"`
	Response    string `json:"response"`
	NeedHistory bool   `json:"needHistory"`
}

type ModelConfig struct {
	Temperature            float32 `json:"temperature"`
	TemperatureEnable      bool    `json:"temperatureEnable"`
	TopP                   float32 `json:"topP"`
	TopPEnable             bool    `json:"topPEnable"`
	FrequencyPenalty       float32 `json:"frequencyPenalty"`
	FrequencyPenaltyEnable bool    `json:"frequencyPenaltyEnable"`
	PresencePenalty        float32 `json:"presencePenalty"`
	PresencePenaltyEnable  bool    `json:"presencePenaltyEnable"`
}

func RagStreamChat(ctx context.Context, userId string, req *RagChatParams) (<-chan string, error) {
	log.Infof("ragStreamChat")
	params, err := buildHttpParams(userId, req)
	if err != nil {
		log.Errorf("build http params fail", "err", err)
		return nil, err
	}
	ret := make(chan string, 1024)
	go func() {
		// 确保通道最终被关闭
		defer close(ret)

		// 捕获 panic 并记录日志（不重新抛出，避免崩溃）
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("RagStreamChat panic: %v", r)
			}
		}()

		//1.开启超时监控
		if params.Timeout == 0 {
			params.Timeout = time.Minute * 10
		}
		ctx, cancel := context.WithTimeout(ctx, params.Timeout)
		defer cancel()

		resp, err := http_client.GetClient().PostJsonOriResp(ctx, params)
		if err != nil {
			log.Errorf("request %+v rag stream err: %v", params, err)
			ret <- fmt.Sprintf("error: 调用下游服务失败: %v", err)
			return
		}
		defer resp.Body.Close() // 确保响应体关闭

		if resp.StatusCode != http.StatusOK {
			log.Errorf("request %+v rag stream returned non-OK status: %d", params, resp.StatusCode)
			ret <- fmt.Sprintf("error: 调用下游服务失败: %s", strconv.Itoa(resp.StatusCode))
			return
		}
		log.Infof("resp: %v", resp)
		scan := bufio.NewScanner(resp.Body)
		for scan.Scan() {
			ret <- scan.Text()
		}
		if err := scan.Err(); err != nil {
			log.Errorf("error reading stream from %v: %v", params, err)
			ret <- fmt.Sprintf("error: 调用下游服务失败: %s", strconv.Itoa(resp.StatusCode))
		}
	}()

	return ret, nil
}

func buildHttpParams(userId string, req *RagChatParams) (*http_client.HttpRequestParams, error) {
	log.Infof("build http param")
	url := fmt.Sprintf("%s%s", config.Cfg().RagServer.ChatEndpoint, config.Cfg().RagServer.ChatUrl)
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return &http_client.HttpRequestParams{
		Url:        url,
		Body:       body,
		Headers:    map[string]string{"X-uid": userId},
		Timeout:    time.Minute * 10,
		MonitorKey: "rag_search_service",
		LogLevel:   http_client.LogAll,
	}, nil
}

// BuildChatConsultParams 构造rag 会话参数
func BuildChatConsultParams(req *rag_service.ChatRagReq, rag *model.RagInfo, knowledgeInfoList *knowledgeBase_service.KnowledgeDetailSelectListResp) *RagChatParams {
	// 知识库参数
	ragChatParams := &RagChatParams{}
	knowledgeConfig := rag.KnowledgeBaseConfig
	ragChatParams.MaxHistory = int32(knowledgeConfig.MaxHistory)
	ragChatParams.Threshold = float32(knowledgeConfig.Threshold)
	ragChatParams.TopK = int32(knowledgeConfig.TopK)
	ragChatParams.RetrieveMethod = buildRetrieveMethod(knowledgeConfig.MatchType)
	ragChatParams.RerankMod = buildRerankMod(knowledgeConfig.PriorityMatch)
	ragChatParams.Weight = buildWeight(knowledgeConfig)
	var kbNameList []string
	for _, v := range knowledgeInfoList.List {
		kbNameList = append(kbNameList, v.Name)
	}
	ragChatParams.KnowledgeBase = kbNameList
	ragChatParams.RerankModelId = buildRerankId(knowledgeConfig.PriorityMatch, rag.RerankConfig.ModelId)

	// RAG属性参数
	ragChatParams.Question = req.Question
	ragChatParams.Stream = true
	ragChatParams.Chichat = true
	ragChatParams.History = []*HistoryItem{}
	ragChatParams.RewriteQuery = true
	ragChatParams.ReturnMeta = true
	//自动角标
	ragChatParams.AutoCitation = true

	// 模型参数
	ragChatParams.CustomModelInfo = &CustomModelInfo{LlmModelID: rag.ModelConfig.ModelId}
	modelConfigStr := rag.ModelConfig.Config
	modelConfig := ModelConfig{}
	err := json.Unmarshal([]byte(modelConfigStr), &modelConfig)
	if err != nil {
		log.Errorf("model config unmarshal fail: %s", modelConfigStr)
		ragChatParams.Temperature = DefaultTemperature
		ragChatParams.TopP = DefaultTopP
		ragChatParams.RepetitionPenalty = DefaultFrequencyPenalty
		return ragChatParams
	}
	if modelConfig.TemperatureEnable {
		ragChatParams.Temperature = modelConfig.Temperature
	} else {
		ragChatParams.Temperature = DefaultTemperature
	}
	if modelConfig.TopPEnable {
		ragChatParams.TopP = modelConfig.TopP
	} else {
		ragChatParams.TopP = DefaultTopP
	}
	if modelConfig.FrequencyPenaltyEnable {
		ragChatParams.RepetitionPenalty = modelConfig.FrequencyPenalty
	} else {
		ragChatParams.RepetitionPenalty = DefaultFrequencyPenalty
	}

	log.Infof("ragparams = %+v", ragChatParams)
	return ragChatParams
}

// buildRerankId 构造重排序模型id
func buildRerankId(priorityType int32, rerankId string) string {
	if priorityType == 1 {
		return ""
	}
	return rerankId
}

// buildRetrieveMethod 构造检索方式
func buildRetrieveMethod(matchType string) string {
	switch matchType {
	case "vector":
		return "semantic_search"
	case "text":
		return "full_text_search"
	case "mix":
		return "hybrid_search"
	}
	return ""
}

// buildRerankMod 构造重排序模式
func buildRerankMod(priorityType int32) string {
	if priorityType == 1 {
		return "weighted_score"
	}
	return "rerank_model"
}

// buildWeight 构造权重信息
func buildWeight(knowConfig model.KnowledgeBaseConfig) *WeightParams {
	if knowConfig.PriorityMatch != 1 {
		return nil
	}
	return &WeightParams{
		VectorWeight: float32(knowConfig.SemanticsPriority),
		TextWeight:   float32(knowConfig.KeywordPriority),
	}
}
