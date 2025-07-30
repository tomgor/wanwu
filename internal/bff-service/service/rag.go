package service

import (
	knowledgebase_service "github.com/UnicomAI/wanwu/api/proto/knowledgebase-service"
	model_service "github.com/UnicomAI/wanwu/api/proto/model-service"
	rag_service "github.com/UnicomAI/wanwu/api/proto/rag-service"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	"github.com/gin-gonic/gin"
)

func CreateRag(ctx *gin.Context, userId, orgId string, req request.AppBriefConfig) (*request.RagReq, error) {
	resp, err := rag.CreateRag(ctx.Request.Context(), &rag_service.CreateRagReq{
		AppBrief: appBriefConfigModel2Proto(req),
		Identity: &rag_service.Identity{
			UserId: userId,
			OrgId:  orgId,
		},
	})
	if err != nil {
		return nil, err
	}
	return &request.RagReq{
		RagID: resp.RagId,
	}, err
}

func UpdateRag(ctx *gin.Context, req request.RagBrief) error {
	_, err := rag.UpdateRag(ctx.Request.Context(), &rag_service.UpdateRagReq{
		RagId:    req.RagID,
		AppBrief: appBriefConfigModel2Proto(req.AppBriefConfig),
	})
	return err
}

func UpdateRagConfig(ctx *gin.Context, req request.RagConfig) error {
	modelConfig, err := appModelConfigModel2Proto(req.ModelConfig)
	if err != nil {
		return err
	}
	rerankConfig, err := appModelConfigModel2Proto(req.RerankConfig)
	if err != nil {
		return err
	}

	configParams := req.KnowledgeBaseConfig.Config
	var sensitiveTableIds []string
	for _, v := range req.SafetyConfig.Tables {
		sensitiveTableIds = append(sensitiveTableIds, v.TableId)
	}
	_, err = rag.UpdateRagConfig(ctx.Request.Context(), &rag_service.UpdateRagConfigReq{
		RagId:        req.RagID,
		ModelConfig:  modelConfig,
		RerankConfig: rerankConfig,
		KnowledgeBaseConfig: &rag_service.RagKnowledgeBaseConfig{
			KnowledgeBaseId:   req.KnowledgeBaseConfig.Knowledgebases[0].ID,
			MaxHistory:        configParams.MaxHistory,
			MaxHistoryEnable:  configParams.MaxHistoryEnable,
			Threshold:         configParams.Threshold,
			ThresholdEnable:   configParams.ThresholdEnable,
			TopK:              configParams.TopK,
			TopKEnable:        configParams.TopKEnable,
			MatchType:         configParams.MatchType,
			PriorityMatch:     configParams.PriorityMatch,
			SemanticsPriority: configParams.SemanticsPriority,
			KeywordPriority:   configParams.KeywordPriority,
		},
		SensitiveConfig: &rag_service.RagSensitiveConfig{
			Enable:   req.SafetyConfig.Enable,
			TableIds: sensitiveTableIds,
		},
	})
	return err
}

func DeleteRag(ctx *gin.Context, req request.RagReq) error {
	_, err := rag.DeleteRag(ctx.Request.Context(), &rag_service.RagDeleteReq{
		RagId: req.RagID,
	})
	return err
}

func GetRag(ctx *gin.Context, req request.RagReq) (*response.RagInfo, error) {
	resp, err := rag.GetRagDetail(ctx.Request.Context(), &rag_service.RagDetailReq{RagId: req.RagID})
	if err != nil {
		return nil, err
	}
	var modelInfo, rerankInfo *model_service.ModelInfo
	var modelConfig, rerankConfig request.AppModelConfig
	var knowledgeInfo *knowledgebase_service.KnowledgeInfo
	var ragInfo = &response.RagInfo{}
	if resp.ModelConfig.ModelId != "" {
		modelInfo, err = model.GetModelById(ctx.Request.Context(), &model_service.GetModelByIdReq{ModelId: resp.ModelConfig.ModelId})
		if err != nil {
			return nil, err
		}
		modelConfig, err = appModelConfigProto2Model(resp.ModelConfig, modelInfo.DisplayName)
		if err != nil {
			return nil, err
		}
	}
	if resp.RerankConfig.ModelId != "" {
		rerankInfo, err = model.GetModelById(ctx.Request.Context(), &model_service.GetModelByIdReq{ModelId: resp.RerankConfig.ModelId})
		if err != nil {
			return nil, err
		}
		rerankConfig, err = appModelConfigProto2Model(resp.RerankConfig, rerankInfo.DisplayName)
		if err != nil {
			return nil, err
		}
	}
	if resp.KnowledgeBaseConfig.KnowledgeBaseId != "" {
		knowledgeInfo, _ = knowledgeBase.SelectKnowledgeDetailById(ctx, &knowledgebase_service.KnowledgeDetailSelectReq{
			KnowledgeId: resp.KnowledgeBaseConfig.KnowledgeBaseId,
		})
	}

	var sensitiveTableList []request.SensitiveTable
	if resp.SensitiveConfig.TableIds != nil {
		for _, tableId := range resp.SensitiveConfig.TableIds {
			sensitiveTableList = append(sensitiveTableList, request.SensitiveTable{
				TableId: tableId,
			})
		}
	}

	knowledgeConfig := resp.KnowledgeBaseConfig
	ragInfo = &response.RagInfo{
		RagID:          resp.RagId,
		AppBriefConfig: appBriefConfigProto2Model(ctx, resp.BriefConfig),
		ModelConfig:    modelConfig,
		RerankConfig:   rerankConfig,
		KnowledgeBaseConfig: request.AppKnowledgebaseConfig{
			Config: request.AppKnowledgebaseParams{
				MaxHistory:        knowledgeConfig.MaxHistory,
				MaxHistoryEnable:  knowledgeConfig.MaxHistoryEnable,
				Threshold:         knowledgeConfig.Threshold,
				ThresholdEnable:   knowledgeConfig.ThresholdEnable,
				TopK:              knowledgeConfig.TopK,
				TopKEnable:        knowledgeConfig.TopKEnable,
				MatchType:         knowledgeConfig.MatchType,
				KeywordPriority:   knowledgeConfig.KeywordPriority,
				PriorityMatch:     knowledgeConfig.PriorityMatch,
				SemanticsPriority: knowledgeConfig.SemanticsPriority,
			},
		},
		SafetyConfig: request.AppSafetyConfig{
			Enable: resp.SensitiveConfig.Enable,
			Tables: sensitiveTableList,
		},
	}
	if knowledgeInfo != nil {
		ragInfo.KnowledgeBaseConfig.Knowledgebases = []request.AppKnowledgeBase{
			{
				ID:   resp.KnowledgeBaseConfig.KnowledgeBaseId,
				Name: knowledgeInfo.Name,
			},
		}
	}
	return ragInfo, nil
}
