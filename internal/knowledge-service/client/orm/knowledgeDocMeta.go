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

// UpdateDocMetaData 更新文档元数据
func UpdateDocMetaData(ctx context.Context, name, description string, knowledgeBase *model.KnowledgeBase) error {
	return db.GetHandle(ctx).Transaction(func(tx *gorm.DB) error {
		//1.更新数据
		err := updateKnowledge(tx, knowledgeBase.Id, name, description)
		if err != nil {
			return err
		}
		//2.通知rag更新知识库
		return service.RagKnowledgeUpdate(ctx, &service.RagUpdateParams{
			UserId:          knowledgeBase.UserId,
			KnowledgeBaseId: knowledgeBase.KnowledgeId,
			OldKbName:       knowledgeBase.Name,
			NewKbName:       name,
		})
	})
}
