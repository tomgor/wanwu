package orm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	knowledgebase_keywords_service "github.com/UnicomAI/wanwu/api/proto/knowledgebase-keywords-service"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/client/model"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/client/orm/sqlopt"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/config"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/db"
	http_client "github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/http-client"
	"github.com/UnicomAI/wanwu/pkg/log"
	"gorm.io/gorm"
	"io"
	"net/http"
	"time"
)

const (
	KeywordsAdd    = "add"    // 添加关键词
	KeywordsDelete = "delete" // 删除关键词
	KeywordsUpdate = "update" // 更新关键词
)

// GetKeywordsList 根据 userId 和 orgId 查询关键词列表
func GetKeywordsList(ctx context.Context, req *knowledgebase_keywords_service.GetKnowledgeKeywordsListReq, isSpecify bool) ([]*model.KnowledgeKeywords, int64, error) {
	var keywordsList []*model.KnowledgeKeywords
	var sqlOptions []sqlopt.SQLOption
	if req.Name != "" {
		// 精确查询（创建/更新关键词）
		if isSpecify {
			sqlOptions = append(sqlOptions, sqlopt.WithOrgID(req.Identity.OrgId), sqlopt.WithUserID(req.Identity.UserId), sqlopt.WithName(req.Name))
		} else {
			// 模糊查询
			sqlOptions = append(sqlOptions, sqlopt.WithOrgID(req.Identity.OrgId), sqlopt.WithUserID(req.Identity.UserId), sqlopt.WithNameOrAliasLike(req.Name))
		}
	} else {
		// 不查询关键词
		sqlOptions = append(sqlOptions, sqlopt.WithOrgID(req.Identity.OrgId), sqlopt.WithUserID(req.Identity.UserId))
	}
	tx := sqlopt.SQLOptions(sqlOptions...).Apply(db.GetHandle(ctx), &model.KnowledgeKeywords{})
	var total int64
	err := tx.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	limit := req.PageSize
	offset := req.PageSize * (req.PageNum - 1)
	err = tx.Order("updated_at desc").Limit(int(limit)).Offset(int(offset)).Find(&keywordsList).Error
	if err != nil {
		return nil, 0, err
	}
	return keywordsList, total, nil
}

// GetKeywordsById 根据关键词id查询关键词
func GetKeywordsById(ctx context.Context, id uint32) (*model.KnowledgeKeywords, error) {
	var keywords model.KnowledgeKeywords
	err := sqlopt.SQLOptions(sqlopt.WithID(id)).
		Apply(db.GetHandle(ctx), &model.KnowledgeKeywords{}).First(&keywords).Error
	if err != nil {
		return nil, err
	}
	return &keywords, nil
}

// CreateKeywords 创建关键词
func CreateKeywords(ctx context.Context, keywords *model.KnowledgeKeywords) error {
	// 创建关键词
	err := db.GetHandle(ctx).Create(keywords).Error
	if err != nil {
		return err
	}
	// 同步RAG
	params, err := buildHttpParams(keywords, KeywordsAdd)
	if err != nil {
		return err
	}
	err = SyncRAG(ctx, params)
	if err != nil {
		return err
	}
	return nil
}

func SyncRAG(ctx context.Context, params *http_client.HttpRequestParams) error {
	// 发送请求
	resp, err := http_client.GetClient().PostJsonOriResp(ctx, params)
	if err != nil {
		log.Errorf("request %+v keywords err: %v", params, err)
		return err
	}
	// 判断状态
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		errBody := string(bodyBytes)
		if errBody == "" {
			errBody = "(empty response)"
		} else if len(errBody) > 4096 {
			errBody = errBody[:4096] + "...(truncated)"
		}

		// 详细错误日志（包含完整的请求和响应信息）
		log.Errorf("===== HTTP ERROR DETAILS =====")
		log.Errorf("URL: %s", params.Url)
		log.Errorf("Status: %d %s", resp.StatusCode, resp.Status)
		log.Errorf("Response Body: %s", errBody)
		log.Errorf("Request Headers: %v", params.Headers)
		log.Errorf("Query Params: %v", params.Params)
		log.Errorf("request keywords returned non-OK status: %d", resp.StatusCode)
		return errors.New(resp.Status)
	}
	log.Infof("resp: %v", resp)
	return nil
}

func buildHttpParams(keywords *model.KnowledgeKeywords, action string) (*http_client.HttpRequestParams, error) {
	// 反序列化id列表
	knowledgeIds, err := JsonToList(keywords.KnowledgeBaseIds)
	if err != nil {
		return nil, err
	}
	params := model.AddKeywords{
		Id:               keywords.Id,
		UserId:           keywords.UserId,
		Action:           action,
		Name:             keywords.Name,
		Alias:            []string{keywords.Alias},
		KnowledgeBaseIds: knowledgeIds,
	}
	// 打印结构化参数
	paramsJSON, _ := json.MarshalIndent(params, "", "  ")
	log.Infof("构建请求参数:\n%s", string(paramsJSON))
	url := fmt.Sprintf("%s%s", config.GetConfig().RagServer.Endpoint, config.GetConfig().RagServer.KeywordsUri)
	//url := fmt.Sprintf("%s%s", config.GetConfig().RagServer.Endpoint, "/rag/proper_noun")

	log.Infof("构建完整URL: %s", url)
	body, err := json.Marshal(params)
	if err != nil {
		log.Errorf("参数序列化错误: %v", err)
		return nil, err
	}
	return &http_client.HttpRequestParams{
		Url:        url,
		Body:       body,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Timeout:    time.Minute * 3,
		MonitorKey: "keyword_service",
		LogLevel:   http_client.LogAll,
	}, nil
}

func JsonToList(knowledgeBaseIdStr string) ([]string, error) {
	var knowledgeIds []string
	err := json.Unmarshal([]byte(knowledgeBaseIdStr), &knowledgeIds)
	if err != nil {
		log.Errorf("反序列化错误: %v, 原始字符串: %s", err, knowledgeBaseIdStr)
		return nil, err
	}
	return knowledgeIds, nil
}

// DeleteKeywords 删除知识库关键词
func DeleteKeywords(ctx context.Context, id uint32) error {
	return db.GetHandle(ctx).Transaction(func(tx *gorm.DB) error {
		// 获取关键词
		keywords, err := GetKeywordsById(ctx, id)
		if err != nil {
			return err
		}
		// 同步RAG
		params, err := buildHttpParams(keywords, KeywordsDelete)
		if err != nil {
			return err
		}
		err = SyncRAG(ctx, params)
		if err != nil {
			return err
		}
		// 删除关键词
		err = tx.Unscoped().Model(&model.KnowledgeKeywords{}).Where("id = ?", id).Delete(&model.KnowledgeKeywords{}).Error
		if err != nil {
			log.Errorf("DeleteKnowledgeKeywords err: %v", err)
			return err
		}
		return nil
	})
}

// UpdateKeywords 更新知识库关键词
func UpdateKeywords(ctx context.Context, keywords *model.KnowledgeKeywords) error {
	// 打印结构化参数
	paramsJSON, _ := json.MarshalIndent(keywords, "", "  ")
	log.Infof("更新关键词参数:\n%s", string(paramsJSON))
	err := db.GetHandle(ctx).Model(&model.KnowledgeKeywords{}).Where("id = ?", keywords.Id).Updates(keywords).Debug().Error
	if err != nil {
		return err
	}
	params, err := buildHttpParams(keywords, KeywordsUpdate)
	if err != nil {
		return err
	}
	err = SyncRAG(ctx, params)
	if err != nil {
		return err
	}
	return nil
}
