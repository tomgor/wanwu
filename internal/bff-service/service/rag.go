package service

import (
	knowledgebase_service "github.com/UnicomAI/wanwu/api/proto/knowledgebase-service"
	model_service "github.com/UnicomAI/wanwu/api/proto/model-service"
	rag_service "github.com/UnicomAI/wanwu/api/proto/rag-service"
	safety_service "github.com/UnicomAI/wanwu/api/proto/safety-service"
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
	_, err = rag.UpdateRagConfig(ctx.Request.Context(), &rag_service.UpdateRagConfigReq{
		RagId:               req.RagID,
		ModelConfig:         modelConfig,
		RerankConfig:        rerankConfig,
		KnowledgeBaseConfig: ragKBConfigToProto(req.KnowledgeBaseConfig),
		SensitiveConfig:     ragSensitiveConfigToProto(req.SafetyConfig),
	})
	return err
}

func ragSensitiveConfigToProto(req request.AppSafetyConfig) *rag_service.RagSensitiveConfig {
	var sensitiveTableIds []string
	for _, v := range req.Tables {
		sensitiveTableIds = append(sensitiveTableIds, v.TableId)
	}
	sensitiveConfig := &rag_service.RagSensitiveConfig{
		Enable:   req.Enable,
		TableIds: sensitiveTableIds,
	}
	return sensitiveConfig
}

func ragKBConfigToProto(knowledgeConfig request.AppKnowledgebaseConfig) *rag_service.RagKnowledgeBaseConfig {
	var knowledgeBaseIds []string
	for _, v := range knowledgeConfig.Knowledgebases {
		knowledgeBaseIds = append(knowledgeBaseIds, v.ID)
	}
	configParams := knowledgeConfig.Config
	knowledgeBaseConfig := &rag_service.RagKnowledgeBaseConfig{
		KnowledgeBaseIds:  knowledgeBaseIds,
		MaxHistory:        configParams.MaxHistory,
		Threshold:         configParams.Threshold,
		TopK:              configParams.TopK,
		MatchType:         configParams.MatchType,
		PriorityMatch:     configParams.PriorityMatch,
		SemanticsPriority: configParams.SemanticsPriority,
		KeywordPriority:   configParams.KeywordPriority,
	}
	return knowledgeBaseConfig
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
	modelConfig, rerankConfig, err := appModelRerankProto2Model(ctx, resp)
	if err != nil {
		return &response.RagInfo{}, err
	}
	ragInfo := &response.RagInfo{
		RagID:               resp.RagId,
		AppBriefConfig:      appBriefConfigProto2Model(ctx, resp.BriefConfig),
		ModelConfig:         modelConfig,
		RerankConfig:        rerankConfig,
		KnowledgeBaseConfig: ragKBConfigProto2Model(ctx, resp.KnowledgeBaseConfig),
		SafetyConfig:        ragSafetyConfigProto2Model(ctx, resp.SensitiveConfig),
	}

	return ragInfo, nil
}

func appModelRerankProto2Model(ctx *gin.Context, resp *rag_service.RagInfo) (request.AppModelConfig, request.AppModelConfig, error) {
	var modelConfig, rerankConfig request.AppModelConfig
	if resp.ModelConfig.ModelId != "" {
		modelInfo, err := model.GetModelById(ctx.Request.Context(), &model_service.GetModelByIdReq{ModelId: resp.ModelConfig.ModelId})
		if err != nil {
			return request.AppModelConfig{}, request.AppModelConfig{}, err
		}
		modelConfig, err = appModelConfigProto2Model(resp.ModelConfig, modelInfo.DisplayName)
		if err != nil {
			return request.AppModelConfig{}, request.AppModelConfig{}, err
		}
	}
	if resp.RerankConfig.ModelId != "" {
		rerankInfo, err := model.GetModelById(ctx.Request.Context(), &model_service.GetModelByIdReq{ModelId: resp.RerankConfig.ModelId})
		if err != nil {
			return request.AppModelConfig{}, request.AppModelConfig{}, err
		}
		rerankConfig, err = appModelConfigProto2Model(resp.RerankConfig, rerankInfo.DisplayName)
		if err != nil {
			return request.AppModelConfig{}, request.AppModelConfig{}, err
		}
	}
	return modelConfig, rerankConfig, nil
}

func ragSafetyConfigProto2Model(ctx *gin.Context, sensitiveCfg *rag_service.RagSensitiveConfig) request.AppSafetyConfig {
	var sensitiveTableList []request.SensitiveTable
	if len(sensitiveCfg.GetTableIds()) != 0 {
		sensitiveWordTable, _ := safety.GetSensitiveWordTableListByIDs(ctx, &safety_service.GetSensitiveWordTableListByIDsReq{TableIds: sensitiveCfg.GetTableIds()})
		if sensitiveWordTable != nil {
			for _, table := range sensitiveWordTable.List {
				sensitiveTableList = append(sensitiveTableList, request.SensitiveTable{
					TableId:   table.TableId,
					TableName: table.TableName,
				})
			}
		}
	}
	safetyConfig := request.AppSafetyConfig{
		Enable: sensitiveCfg.Enable,
		Tables: sensitiveTableList,
	}
	return safetyConfig
}

func ragKBConfigProto2Model(ctx *gin.Context, kbConfig *rag_service.RagKnowledgeBaseConfig) request.AppKnowledgebaseConfig {
	var knowledgeBases []request.AppKnowledgeBase
	if len(kbConfig.KnowledgeBaseIds) > 0 {
		knowledgeInfoList, _ := knowledgeBase.SelectKnowledgeDetailByIdList(ctx, &knowledgebase_service.KnowledgeDetailSelectListReq{
			KnowledgeIds: kbConfig.KnowledgeBaseIds,
		})
		if knowledgeInfoList != nil {
			for _, v := range knowledgeInfoList.List {
				kb := request.AppKnowledgeBase{
					ID:   v.KnowledgeId,
					Name: v.Name,
				}
				knowledgeBases = append(knowledgeBases, kb)
			}
		}
	}
	knowledgeBaseConfig := request.AppKnowledgebaseConfig{
		Knowledgebases: knowledgeBases,
		Config: request.AppKnowledgebaseParams{
			MaxHistory:        kbConfig.MaxHistory,
			Threshold:         kbConfig.Threshold,
			TopK:              kbConfig.TopK,
			MatchType:         kbConfig.MatchType,
			KeywordPriority:   kbConfig.KeywordPriority,
			PriorityMatch:     kbConfig.PriorityMatch,
			SemanticsPriority: kbConfig.SemanticsPriority,
		},
	}
	return knowledgeBaseConfig
}
