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
func UpdateDocStatusDocMeta(ctx context.Context, docId string, addList []*model.KnowledgeDocMeta,
	updateList []*model.KnowledgeDocMeta, deleteDataIdList []string, ragDocMetaParams *service.RagDocMetaParams) error {
	return db.GetHandle(ctx).Transaction(func(tx *gorm.DB) error {
		//todo 文档元数据应该不会特别多，所以先这么做，如果比较多，后续优化
		if len(deleteDataIdList) > 0 {
			err := tx.Unscoped().Model(&model.KnowledgeDocMeta{}).Where("meta_id IN ?", deleteDataIdList).Delete(&model.KnowledgeDocMeta{}).Error
			if err != nil {
				return err
			}
		}
		if len(addList) > 0 {
			//插入数据
			err := tx.Model(&model.KnowledgeDocMeta{}).CreateInBatches(addList, len(addList)).Error
			if err != nil {
				return err
			}
		}
		if len(updateList) > 0 {
			for _, meta := range updateList {
				//更新数据
				updateMap := map[string]interface{}{
					"value": meta.Value,
				}
				err := tx.Model(&model.KnowledgeDocMeta{}).Where("meta_id = ?", meta.MetaId).Updates(updateMap).Error
				if err != nil {
					return err
				}
			}
		}
		//调用rag
		//return service.RagDocMeta(ctx, ragDocMetaParams)
		return nil
	})
}
