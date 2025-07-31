package orm

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/app-service/client/model"
	"github.com/UnicomAI/wanwu/internal/app-service/client/orm/sqlopt"
	"github.com/UnicomAI/wanwu/internal/app-service/pkg"
	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/UnicomAI/wanwu/pkg/minio"
	"github.com/UnicomAI/wanwu/pkg/util"
	"gorm.io/gorm"
)

const (
	SensitiveTypeAll SensitiveType = iota
	SensitiveTypePolitical
	SensitiveTypeRevile
	SensitiveTypePornography
	SensitiveTypeViolentTerror
	SensitiveTypeIllegal
	SensitiveTypeInformationSecurity
	SensitiveTypeOther
	AppSafetySensitiveUploadSingle       = "single"
	AppSafetySensitiveUploadFile         = "file"
	MaxSensitiveUploadSize         int64 = 100
	BatchUploadSize                      = 30
)

type SensitiveType int

var SensitiveTypeToString = map[SensitiveType]string{
	SensitiveTypePolitical:           "Political",
	SensitiveTypeRevile:              "Revile",
	SensitiveTypePornography:         "Pornography",
	SensitiveTypeViolentTerror:       "ViolentTerror",
	SensitiveTypeIllegal:             "Illegal",
	SensitiveTypeInformationSecurity: "InformationSecurity",
	SensitiveTypeOther:               "Other",
}

func (c *Client) CreateSensitiveWordTable(ctx context.Context, userId, orgId, tableName, remark string) (string, *errs.Status) {
	err := sqlopt.SQLOptions(
		sqlopt.WithOrgID(orgId),
		sqlopt.WithUserID(userId),
		sqlopt.WithName(tableName),
	).Apply(c.db.WithContext(ctx)).First(&model.SensitiveWordTable{}).Error
	if err == nil {
		return "", toErrStatus("app_safety_sensitive_table_exist")
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		table := &model.SensitiveWordTable{
			Name:   tableName,
			Remark: remark,
			UserID: userId,
			OrgID:  orgId,
		}
		if err := c.db.WithContext(ctx).Create(table).Error; err != nil {
			return "", toErrStatus("app_safety_sensitive_table_create", tableName, err.Error())
		}
		return util.Int2Str(table.ID), nil
	}
	return "", toErrStatus("app_safety_sensitive_table_get", tableName)
}

func (c *Client) UpdateSensitiveWordTable(ctx context.Context, tableId, tableName, remark string) *errs.Status {
	var table model.SensitiveWordTable
	updates := map[string]interface{}{
		"name":   tableName,
		"remark": remark,
	}
	if err := sqlopt.WithID(tableId).Apply(c.db.WithContext(ctx)).Model(&table).Updates(updates).Error; err != nil {
		return toErrStatus("app_safety_sensitive_table_update", tableId, err.Error())
	}
	return nil
}

func (c *Client) UpdateSensitiveWordTableReply(ctx context.Context, tableId, reply string) *errs.Status {
	var table model.SensitiveWordTable
	if err := sqlopt.WithID(tableId).Apply(c.db.WithContext(ctx)).Model(&table).Update("reply", reply).Error; err != nil {
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
			_, exceeded, err := c.checkSensitiveWordCount(ctx, tableId)
			if err != nil {
				return toErrStatus("app_safety_sensitive_vocabulary_count_get", tableId, err.Error())
			}
			if exceeded {
				return toErrStatus("app_safety_sensitive_table_full", util.Int2Str(MaxSensitiveUploadSize))
			}
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
	// 2. 内存处理Excel
	var dataList []*model.SensitiveWordVocabulary
	dataMap := make(map[string]bool)

	_, err = pkg.ReadExcelFromMemory(fileData, func(lineCount int64, lineText []string) bool {
		if lineCount == 0 { // 跳过表头
			return true
		}
		_, _, err1 := c.checkSensitiveWordCount(ctx, tableId)
		if err1 != nil {
			dataList = nil
			log.Errorf("checkTableWordCount error %s", err1)
			return false
		}
		for index, text := range lineText {
			text = strings.TrimSpace(text)
			if text == "" || dataMap[text] {
				continue
			}
			dataMap[text] = true
			dataList = append(dataList, &model.SensitiveWordVocabulary{
				TableID:       tableId,
				SensitiveType: SensitiveTypeToString[SensitiveType(index+1)],
				Content:       text,
				CreatedAt:     time.Now().UnixMilli(),
				UserID:        userId,
				OrgID:         orgId,
			})
			if len(dataList) >= BatchUploadSize {
				if err := c.CreateSensitiveWordsBatch(ctx, dataList); err != nil {
					log.Errorf("batch insert failed at line %d: %v", lineCount, err)
					return false
				}
				dataList = nil // 清空但保留底层数组
			}
		}
		return true
	})

	if len(dataList) > 0 {
		if err := c.CreateSensitiveWordsBatch(ctx, dataList); err != nil {
			return err
		}
	}

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

func (c *Client) checkSensitiveWordCount(ctx context.Context, tableId string) (int64, bool, error) {
	var count int64
	if err := sqlopt.WithTableID(tableId).Apply(c.db.WithContext(ctx)).Model(&model.SensitiveWordVocabulary{}).Count(&count).Error; err != nil {
		log.Errorf("GetSensitiveWordCount failed for table %d: %v", tableId, err)
		return 0, false, err
	}
	if count >= MaxSensitiveUploadSize {
		return count, true, nil
	}
	return count, false, nil
}

func (c *Client) FilterDuplicateWords(ctx context.Context, tableId string, wordDataList []*model.SensitiveWordVocabulary) ([]*model.SensitiveWordVocabulary, error) {
	// 1. 提取待检查的词条内容列表
	wordContents := make([]string, 0, len(wordDataList))
	for _, word := range wordDataList {
		wordContents = append(wordContents, word.Content)
	}
	// 2. 查询已存在的词条
	var existingWords []*model.SensitiveWordVocabulary
	query := sqlopt.SQLOptions(
		sqlopt.WithTableID(tableId),
		sqlopt.WithContents(wordContents),
	)
	if err := query.Apply(c.db.WithContext(ctx)).Find(&existingWords).Error; err != nil {
		log.Errorf("FilterDuplicateWords failed for table %s: %v", tableId, err)
		return nil, err
	}
	// 3. 如果没有重复词条，直接返回原始列表
	if len(existingWords) == 0 {
		return wordDataList, nil
	}
	// 4. 构建重复词条映射
	existingMap := make(map[string]bool)
	for _, word := range existingWords {
		existingMap[word.Content] = true
	}
	// 5. 过滤掉重复词条
	filteredList := make([]*model.SensitiveWordVocabulary, 0, len(wordDataList))
	for _, word := range wordDataList {
		if !existingMap[word.Content] {
			filteredList = append(filteredList, word)
		}
	}
	return filteredList, nil
}

func (c *Client) CreateSensitiveWordsBatch(ctx context.Context, dataList []*model.SensitiveWordVocabulary) *errs.Status {
	if len(dataList) == 0 {
		return nil
	}
	tableId := dataList[0].TableID
	filteredWords, err := c.FilterDuplicateWords(ctx, tableId, dataList)
	if err != nil {
		log.Errorf("Error filtering duplicate word list: %v", err)
	}
	if len(filteredWords) == 0 {
		return nil
	}
	count, exceeded, err := c.checkSensitiveWordCount(ctx, tableId)
	if err != nil {
		return toErrStatus("app_safety_sensitive_vocabulary_count_get", tableId, err.Error())
	}
	if exceeded {
		return toErrStatus("app_safety_sensitive_table_full", util.Int2Str(MaxSensitiveUploadSize))
	}
	remaining := int(MaxSensitiveUploadSize - count)
	if remaining <= 0 {
		return nil
	}
	if remaining < len(filteredWords) {
		filteredWords = filteredWords[:remaining]
	}
	err = c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.SensitiveWordVocabulary{}).
			CreateInBatches(filteredWords, len(filteredWords)).Error; err != nil {
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

func getSensitiveTableVersion() string {
	return util.Int2Str(time.Now().UnixMilli())
}
