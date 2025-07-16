package service

import (
	"fmt"
	"net/http"

	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	model_service "github.com/UnicomAI/wanwu/api/proto/model-service"
	"github.com/UnicomAI/wanwu/internal/bff-service/config"
	gin_util "github.com/UnicomAI/wanwu/internal/bff-service/pkg/gin-util"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	mp "github.com/UnicomAI/wanwu/pkg/model-provider"
	mp_common "github.com/UnicomAI/wanwu/pkg/model-provider/mp-common"
	"github.com/gin-gonic/gin"
)

func ModelEmbeddings(ctx *gin.Context, modelID string, req *mp_common.EmbeddingReq) {
	// modelInfo by modelID
	modelInfo, err := model.GetModelById(ctx.Request.Context(), &model_service.GetModelByIdReq{ModelId: modelID})
	if err != nil {
		gin_util.Response(ctx, nil, err)
		return
	}

	// 校验model字段
	if req != nil {
		if req.Model != "" && req.Model != modelInfo.Model {
			gin_util.Response(ctx, nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("model %v embeddings err: model mismatch!", modelInfo.ModelId)))
			return
		}
	}

	// embedding config
	embedding, err := mp.ToModelConfig(modelInfo.Provider, modelInfo.ModelType, modelInfo.ProviderConfig)
	if err != nil {
		gin_util.Response(ctx, nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("model %v embeddings err: %v", modelInfo.ModelId, err)))
		return
	}
	iEmbedding, ok := embedding.(mp.IEmbedding)
	if !ok {
		gin_util.Response(ctx, nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("model %v embeddings err: invalid provider", modelInfo.ModelId)))
		return
	}
	// embeddings
	embeddingReq, err := iEmbedding.NewReq(req)
	if err != nil {
		gin_util.Response(ctx, nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("model %v embeddings NewReq err: %v", modelInfo.ModelId, err)))
		return
	}
	resp, err := iEmbedding.Embeddings(ctx.Request.Context(), embeddingReq)
	if err != nil {
		gin_util.Response(ctx, nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("model %v embeddings err: %v", modelInfo.ModelId, err)))
		return
	}
	if data, ok := resp.ConvertResp(); ok {
		status := http.StatusOK
		ctx.Set(config.STATUS, status)
		//ctx.Set(config.RESULT, resp.String())
		ctx.JSON(status, data)
		return
	}
	gin_util.Response(ctx, nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("model %v embeddings err: invalid resp", modelInfo.ModelId)))
}
