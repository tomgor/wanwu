package orm

import (
	"context"

	"github.com/UnicomAI/wanwu/internal/knowledge-service/client/model"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/client/orm/sqlopt"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/db"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/service"
	"gorm.io/gorm"
)

// SelectDocMetaList 查询知识库列表
func SelectDocMetaList(ctx context.Context, userId, orgId, docId string) ([]*model.KnowledgeDocMeta, error) {
	var docMetaList []*model.KnowledgeDocMeta
	err := sqlopt.SQLOptions(sqlopt.WithDocID(docId), sqlopt.WithPermit(orgId, userId)).
		Apply(db.GetHandle(ctx), &model.KnowledgeDocMeta{}).
		Order("create_at desc").
		Find(&docMetaList).
		Error
	if err != nil {
		return nil, err
	}
	return docMetaList, nil
}

// UpdateDocStatusDocMeta 更新文档tag
func UpdateDocStatusDocMeta(ctx context.Context, docId string, metaDataList []*model.KnowledgeDocMeta, ragDocMetaParams *service.RagDocMetaParams) error {
	return db.GetHandle(ctx).Transaction(func(tx *gorm.DB) error {
		//todo 文档元数据应该不会特别多，所以先这么做，如果比较多，后续优化
		//删除所有知识库标签
		err := tx.Unscoped().Model(&model.KnowledgeDocMeta{}).Where("doc_id = ?", docId).Delete(&model.KnowledgeDocMeta{}).Error
		if err != nil {
			return err
		}
		//插入数据
		err = tx.Model(&model.KnowledgeDocMeta{}).CreateInBatches(metaDataList, len(metaDataList)).Error
		if err != nil {
			return err
		}
		//调用rag
		return service.RagDocMeta(ctx, ragDocMetaParams)
	})
}
