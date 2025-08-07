package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/config"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/http"
	http_client "github.com/UnicomAI/wanwu/pkg/http-client"
	"github.com/UnicomAI/wanwu/pkg/log"
	"time"
)

const (
	keywordsSuccessCode = 200
)

type RagOperateKeywordsParams struct {
	Id               uint32   `json:"id"`
	UserId           string   `json:"user_id"`
	Action           string   `json:"action"`
	Name             string   `json:"name"`
	Alias            []string `json:"alias"`
	KnowledgeBaseIds []string `json:"knowledge_base_list"`
}

// RagOperateKeywords rag添加关键词
func RagOperateKeywords(ctx context.Context, ragOperateKeywordsParams *RagOperateKeywordsParams) error {
	ragServer := config.GetConfig().RagServer
	url := ragServer.Endpoint + ragServer.KeywordsUri
	paramsByte, err := json.Marshal(ragOperateKeywordsParams)
	if err != nil {
		return err
	}
	result, err := http.GetClient().PostJson(ctx, &http_client.HttpRequestParams{
		Url:        url,
		Body:       paramsByte,
		Timeout:    time.Duration(ragServer.Timeout) * time.Second,
		MonitorKey: fmt.Sprintf("rag_keywords_%s", ragOperateKeywordsParams.Action),
		LogLevel:   http_client.LogAll,
	})
	if err != nil {
		return err
	}
	var resp RagCommonResp
	if err := json.Unmarshal(result, &resp); err != nil {
		log.Errorf(err.Error())
		return err
	}
	if resp.Code != keywordsSuccessCode {
		return errors.New(resp.Message)
	}
	return nil
}
