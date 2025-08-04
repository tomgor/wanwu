package orm

import (
	"context"
	"errors"
	"fmt"
	"time"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/app-service/client/model"
	"github.com/UnicomAI/wanwu/internal/app-service/client/orm/sqlopt"
	"github.com/UnicomAI/wanwu/internal/app-service/pkg"
	"github.com/UnicomAI/wanwu/pkg/minio"
	"github.com/UnicomAI/wanwu/pkg/util"
	"gorm.io/gorm"
)

const (
	AppSafetySensitiveUploadSingle       = "single"
	AppSafetySensitiveUploadFile         = "file"
	MaxSensitiveUploadSize         int64 = 100
)

func (c *Client) CreateSensitiveWordTable(ctx context.Context, userId, orgId, tableName, remark string) (string, *errs.Status) {
	err := sqlopt.SQLOptions(
		sqlopt.WithOrgID(orgId),
		sqlopt.WithUserID(userId),
		sqlopt.WithName(tableName),
	).Apply(c.db.WithContext(ctx)).First(&model.SensitiveWordTable{}).Error
	if err == nil {
		return "", toErrStatus("app_safety_sensitive_table_exist")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", toErrStatus("app_safety_sensitive_table_get", tableName)
	}
	table := &model.SensitiveWordTable{
		Name:    tableName,
		Remark:  remark,
		Version: getSensitiveTableVersion(),
		UserID:  userId,
		OrgID:   orgId,
	}
	if err := c.db.WithContext(ctx).Create(table).Error; err != nil {
		return "", toErrStatus("app_safety_sensitive_table_create", tableName, err.Error())
	}
	return util.Int2Str(table.ID), nil
}

func (c *Client) UpdateSensitiveWordTable(ctx context.Context, tableId, tableName, remark string) *errs.Status {
	var existingTable model.SensitiveWordTable
	err := sqlopt.SQLOptions(
		sqlopt.WithName(tableName),
	).Apply(c.db.WithContext(ctx)).First(&existingTable).Error
	if err == nil && util.Int2Str(existingTable.ID) != tableId {
		return toErrStatus("app_safety_sensitive_table_exist")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return toErrStatus("app_safety_sensitive_table_get", tableName)
	}
	updates := map[string]interface{}{
		"name":   tableName,
		"remark": remark,
	}
	updateErr := sqlopt.WithID(tableId).
		Apply(c.db.WithContext(ctx)).
		Model(&model.SensitiveWordTable{}).
		Updates(updates).Error

	if updateErr != nil {
		return toErrStatus("app_safety_sensitive_table_update", tableId, updateErr.Error())
	}
	return nil
}

func (c *Client) UpdateSensitiveWordTableReply(ctx context.Context, tableId, reply string) *errs.Status {
	var table model.SensitiveWordTable
	if err := sqlopt.WithID(tableId).Apply(c.db.WithContext(ctx)).Model(&table).
		Updates(map[string]interface{}{
			"reply":   reply,
			"version": getSensitiveTableVersion(),
		}).Update("reply", reply).Error; err != nil {
		return toErrStatus("app_safety_sensitive_table_reply_update", tableId, err.Error())
	}
	return nil
}

func (c *Client) DeleteSensitiveWordTable(ctx context.Context, tableId string) *errs.Status {
	err := c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := sqlopt.SQLOptions(
			sqlopt.WithTableID(tableId),
		).Apply(tx).Delete(&model.SensitiveWordVocabulary{}).Error; err != nil {
			return fmt.Errorf("failed to delete sensitiveWordVocabulary: %v", err)
		}
		if err := sqlopt.SQLOptions(
			sqlopt.WithID(tableId),
		).Apply(tx).Delete(&model.SensitiveWordTable{}).Error; err != nil {
			return fmt.Errorf("failed to delete sensitiveWordTable: %v", err)
		}
		return nil
	})
	if err != nil {
		return toErrStatus("app_safety_sensitive_table_delete", tableId, err.Error())
	}
	return nil
}

func (c *Client) GetSensitiveWordTableList(ctx context.Context, userId, orgId string) ([]*model.SensitiveWordTable, *errs.Status) {
	var tables []*model.SensitiveWordTable
	if err := sqlopt.SQLOptions(
		sqlopt.WithOrgID(orgId),
		sqlopt.WithUserID(userId),
	).Apply(c.db.WithContext(ctx)).Find(&tables).Error; err != nil {
		return nil, toErrStatus("app_safety_sensitive_table_list_get", err.Error())
	}
	return tables, nil
}

func (c *Client) GetSensitiveVocabularyList(ctx context.Context, tableId string, offset, limit int32) ([]*model.SensitiveWordVocabulary, int64, *errs.Status) {
	var vocabularies []*model.SensitiveWordVocabulary
	var count int64
	// 查询分页数据
	if err := sqlopt.SQLOptions(
		sqlopt.WithTableID(tableId),
	).Apply(c.db.WithContext(ctx)).Offset(int(offset)).Limit(int(limit)).Order("id DESC").Find(&vocabularies).
		Offset(-1).Limit(-1).Count(&count).Error; err != nil {
		return nil, 0, toErrStatus("app_safety_sensitive_vocabulary_list_get", tableId, err.Error())
	}
	return vocabularies, count, nil
}

func (c *Client) UploadSensitiveVocabulary(ctx context.Context, userId, orgId, tableId, importType, word, sensitiveType, filePath string) *errs.Status {
	var words []*model.SensitiveWordVocabulary
	var count int64
	if err := sqlopt.WithTableID(tableId).Apply(c.db.WithContext(ctx)).Find(&words).Count(&count).Error; err != nil {
		return toErrStatus("app_safety_sensitive_vocabulary_list_get", tableId, err.Error())
	}
	if count >= MaxSensitiveUploadSize {
		return toErrStatus("app_safety_sensitive_table_full", util.Int2Str(MaxSensitiveUploadSize))
	}
	// single上传
	if importType == AppSafetySensitiveUploadSingle {
		var existingRecord model.SensitiveWordVocabulary
		err := sqlopt.SQLOptions(
			sqlopt.WithTableID(tableId),
			sqlopt.WithContent(word),
		).Apply(c.db.WithContext(ctx)).First(&existingRecord).Error
		if err == nil {
			return toErrStatus("app_safety_sensitive_vocabulary_exist", word)
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newRecord := &model.SensitiveWordVocabulary{
				OrgID:         orgId,
				UserID:        userId,
				Content:       word,
				SensitiveType: sensitiveType,
				TableID:       tableId,
			}
			err = c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
				if err := tx.Create(newRecord).Error; err != nil {
					return fmt.Errorf("create sensitive word failed: %w", err)
				}
				if err := sqlopt.WithID(tableId).Apply(tx).Model(&model.SensitiveWordTable{}).
					Update("version", getSensitiveTableVersion()).Error; err != nil {
					return fmt.Errorf("update table version failed: %w", err)
				}
				return nil
			})
			if err != nil {
				return toErrStatus("app_safety_sensitive_vocabulary_create", word, err.Error())
			}
			return nil
		}
		return toErrStatus("app_safety_sensitive_vocabulary_create", word, err.Error())
	}
	// 1. 从MinIO下载文件到内存
	fileData, err := minio.DownloadFileToMemory(ctx, filePath)
	if err != nil {
		return toErrStatus("app_safety_sensitive_download_fail", err.Error())
	}
	// 2. 解析Excel文件
	sensitiveWords, parseErr := pkg.ParseSensitiveExcel(fileData)
	if parseErr != nil {
		return toErrStatus("app_safety_sensitive_download_fail", parseErr.Error())
	}
	// 3. 构造完整敏感词数据表
	allWords := make([]*model.SensitiveWordVocabulary, len(sensitiveWords))
	for i, raw := range sensitiveWords {
		allWords[i] = &model.SensitiveWordVocabulary{
			TableID:       tableId,
			SensitiveType: raw.SensitiveType,
			Content:       raw.Content,
			UserID:        userId,
			OrgID:         orgId,
		}
	}
	wordContents := make([]string, 0, len(allWords))
	for _, word := range allWords {
		wordContents = append(wordContents, word.Content)
	}
	// 4. 查询已存在的词条
	var existingWords []*model.SensitiveWordVocabulary
	query := sqlopt.SQLOptions(
		sqlopt.WithTableID(tableId),
		sqlopt.WithContents(wordContents),
	)
	if err := query.Apply(c.db.WithContext(ctx)).Find(&existingWords).Error; err != nil {
		return toErrStatus("app_safety_sensitive_vocabulary_list_get", tableId, err.Error())
	}
	// 5. 构建重复词条映射
	existingMap := make(map[string]bool)
	for _, word := range existingWords {
		existingMap[word.Content] = true
	}
	// 6. 过滤掉重复词条
	filteredList := make([]*model.SensitiveWordVocabulary, 0, len(allWords))
	for _, word := range allWords {
		if !existingMap[word.Content] {
			filteredList = append(filteredList, word)
		}
	}
	if len(filteredList) == 0 {
		return nil
	}
	// 7. 计算有效数据量
	remaining := int(MaxSensitiveUploadSize - count)
	if remaining < len(filteredList) {
		return toErrStatus("app_safety_sensitive_table_full", util.Int2Str(MaxSensitiveUploadSize))
	}
	// 8. 批量插入数据
	err = c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.SensitiveWordVocabulary{}).
			Create(filteredList).Error; err != nil {
			return fmt.Errorf("batch create failed: %w", err)
		}
		if err := sqlopt.WithID(tableId).Apply(tx).Model(&model.SensitiveWordTable{}).
			Update("version", getSensitiveTableVersion()).Error; err != nil {
			return fmt.Errorf("update table version failed: %w", err)
		}
		return nil
	})
	if err != nil {
		return toErrStatus("app_safety_sensitive_word_file_create_err", tableId, err.Error())
	}
	return nil
}

func (c *Client) DeleteSensitiveVocabulary(ctx context.Context, tableId, wordId string) *errs.Status {
	err := c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := sqlopt.SQLOptions(
			sqlopt.WithID(tableId),
		).Apply(tx).Model(&model.SensitiveWordTable{}).
			Update("version", getSensitiveTableVersion()).Error; err != nil {
			return fmt.Errorf("update table version failed: %w", err)
		}
		if err := sqlopt.SQLOptions(
			sqlopt.WithTableID(tableId),
			sqlopt.WithID(wordId),
		).Apply(tx).Delete(&model.SensitiveWordVocabulary{}).Error; err != nil {
			return fmt.Errorf("failed to delete sensitiveWordVocabulary: %v", err)
		}
		return nil
	})
	if err != nil {
		return toErrStatus("app_safety_sensitive_vocabulary_delete", wordId, err.Error())
	}
	return nil
}

func (c *Client) GetSensitiveWordTableListWithWordsByIDs(ctx context.Context, tableIds []string) ([]*SensitiveWordTableWithWord, *errs.Status) {
	var vocabularies []*model.SensitiveWordVocabulary
	if err := sqlopt.WithTableIDs(tableIds).Apply(c.db.WithContext(ctx)).
		Find(&vocabularies).Error; err != nil {
		return nil, toErrStatus("app_safety_sensitive_vocabulary_list_get_by_ids", err.Error())
	}
	var tables []*model.SensitiveWordTable
	if err := sqlopt.WithIDs(tableIds).Apply(c.db.WithContext(ctx)).
		Find(&tables).Error; err != nil {
		return nil, toErrStatus("app_safety_sensitive_table_list_get", err.Error())
	}
	result := make([]*SensitiveWordTableWithWord, 0, len(tables))

	for _, t := range tables {
		tableID := util.Int2Str(t.ID)
		item := &SensitiveWordTableWithWord{
			SensitiveWordTable: *t,
			SensitiveWords:     make([]string, 0),
		}
		for _, v := range vocabularies {
			if tableID == v.TableID {
				item.SensitiveWords = append(item.SensitiveWords, v.Content)
			}
		}
		result = append(result, item)
	}
	return result, nil
}

func (c *Client) GetSensitiveWordTableListByIDs(ctx context.Context, tableIds []string) ([]*model.SensitiveWordTable, *errs.Status) {
	var tables []*model.SensitiveWordTable
	if err := sqlopt.SQLOptions(
		sqlopt.WithIDs(tableIds),
	).Apply(c.db.WithContext(ctx)).Find(&tables).Error; err != nil {
		return nil, toErrStatus("app_safety_sensitive_table_list_get", err.Error())
	}
	return tables, nil
}

func (c *Client) GetSensitiveWordTableByID(ctx context.Context, tableId string) (*model.SensitiveWordTable, *errs.Status) {
	var table model.SensitiveWordTable
	if err := sqlopt.SQLOptions(
		sqlopt.WithID(tableId),
	).Apply(c.db.WithContext(ctx)).First(&table).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, toErrStatus("app_safety_sensitive_table_not_found", tableId)
		}
		return nil, toErrStatus("app_safety_sensitive_table_get", tableId, err.Error())
	}
	return &table, nil
}
func getSensitiveTableVersion() string {
	return util.Int2Str(time.Now().UnixMilli())
}
