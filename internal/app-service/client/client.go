package client

import (
	"context"

	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/app-service/client/model"
	"github.com/UnicomAI/wanwu/internal/app-service/client/orm"
)

type IClient interface {
	// --- api key ---
	GetApiKeyList(ctx context.Context, userId, orgId, appId, appType string) ([]*model.ApiKey, *err_code.Status)
	DelApiKey(ctx context.Context, apiId string) *err_code.Status
	GenApiKey(ctx context.Context, userId, orgId, appId, appType, apiKey string) (*model.ApiKey, *err_code.Status)
	GetApiKeyByKey(ctx context.Context, apiKey string) (*model.ApiKey, *err_code.Status)

	// --- explore ---
	GetExplorationAppList(ctx context.Context, userId, name, appType, searchType string) ([]*orm.ExplorationAppInfo, *err_code.Status)
	ChangeExplorationAppFavorite(ctx context.Context, userId, orgId, appId, appType string, isFavorite bool) *err_code.Status

	// --- app ---
	PublishApp(ctx context.Context, userId, orgId, appId, appType, publishType string) *err_code.Status
	UnPublishApp(ctx context.Context, appId, appType, userId string) *err_code.Status
	GetAppList(ctx context.Context, userId, orgId, appType string) ([]*model.App, *err_code.Status)
	DeleteApp(ctx context.Context, appId, appType string) *err_code.Status
	RecordAppHistory(ctx context.Context, userId, appId, appType string) *err_code.Status
	GetAppListByIds(ctx context.Context, ids []string) ([]*model.App, *err_code.Status)

	// --- safety ---
	CreateSensitiveWordTable(ctx context.Context, userId, orgId, tableName, remark string) (string, *err_code.Status)
	UpdateSensitiveWordTable(ctx context.Context, tableId, tableName, remark string) *err_code.Status
	UpdateSensitiveWordTableReply(ctx context.Context, tableId, reply string) *err_code.Status
	DeleteSensitiveWordTable(ctx context.Context, tableId string) *err_code.Status
	GetSensitiveWordTableList(ctx context.Context, userId, orgId string) ([]*model.SensitiveWordTable, *err_code.Status)
	GetSensitiveVocabularyList(ctx context.Context, tableId string, offset, limit int32) ([]*model.SensitiveWordVocabulary, int64, *err_code.Status)
	UploadSensitiveVocabulary(ctx context.Context, userId, orgId, tableId, importType, word, sensitiveType, filePath string) *err_code.Status
	DeleteSensitiveVocabulary(ctx context.Context, tableId, wordId string) *err_code.Status
	GetSensitiveWordTableWithWord(ctx context.Context, tableIds []string) ([]*orm.SensitiveWordTableWithWord, *err_code.Status)
	GetSensitiveWordTableByIds(ctx context.Context, tableIds []string) ([]*model.SensitiveWordTable, *err_code.Status)
}
