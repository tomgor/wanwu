package knowledge

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/UnicomAI/wanwu/internal/knowledge-service/client/model"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/client/orm"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/generator"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/util"
	rag_service "github.com/UnicomAI/wanwu/internal/knowledge-service/service"
	"github.com/UnicomAI/wanwu/pkg/log"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	knowledgebase_service "github.com/UnicomAI/wanwu/api/proto/knowledgebase-service"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	HitTopK              = 3
	HitThreshold float64 = 0.4
)

func (s *Service) SelectKnowledgeList(ctx context.Context, req *knowledgebase_service.KnowledgeSelectReq) (*knowledgebase_service.KnowledgeSelectListResp, error) {
	list, err := orm.SelectKnowledgeList(ctx, req.UserId, req.OrgId, req.Name, req.TagIdList)
	if err != nil {
		log.Errorf(fmt.Sprintf("获取知识库列表失败(%v)  参数(%v)", err, req))
		return nil, util.ErrCode(errs.Code_KnowledgeBaseSelectFailed)
	}

	var tagMap = make(map[string][]*orm.TagRelationDetail)
	if len(list) > 0 {
		var knowledgeList []string
		for _, k := range list {
			knowledgeList = append(knowledgeList, k.KnowledgeId)
		}
		relation := orm.SelectKnowledgeTagListWithRelation(ctx, req.UserId, req.OrgId, "", knowledgeList)
		tagMap = buildKnowledgeTagMap(relation)
	}

	return buildKnowledgeListResp(list, tagMap), nil
}

func (s *Service) SelectKnowledgeDetailById(ctx context.Context, req *knowledgebase_service.KnowledgeDetailSelectReq) (*knowledgebase_service.KnowledgeInfo, error) {
	knowledgeInfo, err := orm.SelectKnowledgeById(ctx, req.KnowledgeId, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf(fmt.Sprintf("获取知识库详情(%v)  参数(%v)", err, req))
		return nil, err
	}
	return buildKnowledgeInfo(knowledgeInfo), nil
}

func (s *Service) SelectKnowledgeDetailByName(ctx context.Context, req *knowledgebase_service.KnowledgeDetailSelectReq) (*knowledgebase_service.KnowledgeInfo, error) {
	knowledgeInfo, err := orm.SelectKnowledgeByName(ctx, req.KnowledgeName, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf(fmt.Sprintf("根据名称获取知识库详情失败(%v)  参数(%v)", err, req))
		return nil, err
	}
	return buildKnowledgeInfo(knowledgeInfo), nil
}

func (s *Service) SelectKnowledgeDetailByIdList(ctx context.Context, req *knowledgebase_service.KnowledgeDetailSelectListReq) (*knowledgebase_service.KnowledgeDetailSelectListResp, error) {
	knowledgeInfoList, err := orm.SelectKnowledgeByIdList(ctx, req.KnowledgeIds, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf(fmt.Sprintf("根据id列表获取知识库详情列表失败(%v)  参数(%v)", err, req))
		return nil, err
	}
	return buildKnowledgeInfoList(knowledgeInfoList), nil
}

func (s *Service) CreateKnowledge(ctx context.Context, req *knowledgebase_service.CreateKnowledgeReq) (*knowledgebase_service.CreateKnowledgeResp, error) {
	//1.重名校验
	err := orm.CheckSameKnowledgeName(ctx, req.UserId, req.OrgId, req.Name)
	if err != nil {
		return nil, err
	}
	//2.创建创建知识库
	knowledgeModel, err := buildKnowledgeBaseModel(req)
	if err != nil {
		log.Errorf("buildKnowledgeBaseModel error %s", err)
		return nil, util.ErrCode(errs.Code_KnowledgeBaseCreateFailed)
	}
	err = orm.CreateKnowledge(ctx, knowledgeModel, req.EmbeddingModelInfo.ModelId)
	if err != nil {
		log.Errorf("CreateKnowledge error %v params %v", err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeBaseCreateFailed)
	}
	//3.返回结果
	return &knowledgebase_service.CreateKnowledgeResp{
		KnowledgeId: knowledgeModel.KnowledgeId,
	}, nil
}

func (s *Service) UpdateKnowledge(ctx context.Context, req *knowledgebase_service.UpdateKnowledgeReq) (*emptypb.Empty, error) {
	//1.查询知识库详情
	knowledge, err := orm.SelectKnowledgeById(ctx, req.KnowledgeId, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf(fmt.Sprintf("没有操作该知识库的权限 参数(%v)", req))
		return nil, err
	}
	//2.重名校验
	err = orm.CheckSameKnowledgeName(ctx, req.UserId, req.OrgId, req.Name)
	if err != nil {
		return nil, err
	}
	//3.更新知识库
	err = orm.UpdateKnowledge(ctx, req.Name, req.Description, knowledge)
	if err != nil {
		log.Errorf("知识库更新失败(%v)  参数(%v)", err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeBaseUpdateFailed)
	}
	return &emptypb.Empty{}, nil
}

// DeleteKnowledge 删除知识库
func (s *Service) DeleteKnowledge(ctx context.Context, req *knowledgebase_service.DeleteKnowledgeReq) (*emptypb.Empty, error) {
	//1.查询知识库详情
	knowledge, err := orm.SelectKnowledgeById(ctx, req.KnowledgeId, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf(fmt.Sprintf("没有操作该知识库的权限 参数(%v)", req))
		return nil, err
	}
	//2.校验导入状态
	err = orm.SelectKnowledgeRunningImportTask(ctx, knowledge.KnowledgeId)
	if err != nil {
		return nil, err
	}
	//3.先删除知识库，异步删除资源数据
	err = orm.DeleteKnowledge(ctx, knowledge)
	if err != nil {
		log.Errorf("删除知识库失败 error %v params %v", err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeBaseDeleteFailed)
	}
	return &emptypb.Empty{}, nil
}

// KnowledgeHit 知识库命中测试
func (s *Service) KnowledgeHit(ctx context.Context, req *knowledgebase_service.KnowledgeHitReq) (*knowledgebase_service.KnowledgeHitResp, error) {
	list, err := orm.SelectKnowledgeByIdList(ctx, req.KnowledgeIdList, req.UserId, req.OrgId)
	if err != nil {
		return nil, err
	}
	matchParams := req.KnowledgeMatchParams
	priorityMatch := matchParams.PriorityMatch
	hitResp, err := rag_service.RagKnowledgeHit(ctx, &rag_service.KnowledgeHitParams{
		UserId:         req.UserId,
		Question:       req.Question,
		KnowledgeBase:  buildKnowledgeNameList(list),
		TopK:           matchParams.TopK,
		Threshold:      float64(matchParams.Score),
		RerankModelId:  buildRerankId(priorityMatch, matchParams.RerankModelId),
		RetrieveMethod: buildRetrieveMethod(matchParams.MatchType),
		RerankMod:      buildRerankMod(priorityMatch),
		Weight:         buildWeight(priorityMatch, matchParams.SemanticsPriority, matchParams.KeywordPriority),
	})
	if err != nil {
		log.Errorf("RagKnowledgeHit error %s", err)
		return nil, util.ErrCode(errs.Code_KnowledgeBaseHitFailed)
	}
	return buildKnowledgeBaseHitResp(hitResp), nil
}

// buildKnowledgeListResp 构造知识库列表返回结果
func buildKnowledgeListResp(knowledgeList []*model.KnowledgeBase, knowledgeTagMap map[string][]*orm.TagRelationDetail) *knowledgebase_service.KnowledgeSelectListResp {
	if len(knowledgeList) == 0 {
		return &knowledgebase_service.KnowledgeSelectListResp{}
	}
	var retList []*knowledgebase_service.KnowledgeInfo
	for _, knowledge := range knowledgeList {
		knowledgeInfo := buildKnowledgeInfo(knowledge)
		knowledgeInfo.KnowledgeTagInfoList = buildKnowledgeTagList(knowledge.KnowledgeId, knowledgeTagMap)
		retList = append(retList, knowledgeInfo)
	}
	return &knowledgebase_service.KnowledgeSelectListResp{
		KnowledgeList: retList,
	}
}

func buildKnowledgeTagMap(tagRelation *orm.TagRelation) map[string][]*orm.TagRelationDetail {
	if tagRelation.RelationErr != nil || tagRelation.TagErr != nil {
		return make(map[string][]*orm.TagRelationDetail)
	}
	var knowledgeTagMap = make(map[string][]*orm.TagRelationDetail)
	for _, relation := range tagRelation.RelationList {
		details := knowledgeTagMap[relation.KnowledgeId]
		if details == nil {
			details = make([]*orm.TagRelationDetail, 0)
		}
		for _, tag := range tagRelation.TagList {
			if tag.TagId == relation.TagId {
				details = append(details, &orm.TagRelationDetail{
					TagId:   tag.TagId,
					TagName: tag.Name,
				})
			}
		}
		knowledgeTagMap[relation.KnowledgeId] = details
	}
	return knowledgeTagMap
}

func buildKnowledgeTagList(knowledgeId string, knowledgeTagMap map[string][]*orm.TagRelationDetail) []*knowledgebase_service.KnowledgeTagInfo {
	if len(knowledgeTagMap) == 0 {
		return []*knowledgebase_service.KnowledgeTagInfo{}
	}
	tagList := knowledgeTagMap[knowledgeId]
	if len(tagList) == 0 {
		return []*knowledgebase_service.KnowledgeTagInfo{}
	}
	var retList []*knowledgebase_service.KnowledgeTagInfo
	for _, tag := range tagList {
		retList = append(retList, &knowledgebase_service.KnowledgeTagInfo{
			TagId:   tag.TagId,
			TagName: tag.TagName,
		})
	}
	return retList
}

// buildKnowledgeInfo 构造知识库信息
func buildKnowledgeInfo(knowledge *model.KnowledgeBase) *knowledgebase_service.KnowledgeInfo {
	embeddingModelInfo := &knowledgebase_service.EmbeddingModelInfo{}
	_ = json.Unmarshal([]byte(knowledge.EmbeddingModel), embeddingModelInfo)
	return &knowledgebase_service.KnowledgeInfo{
		KnowledgeId:        knowledge.KnowledgeId,
		Name:               knowledge.Name,
		Description:        knowledge.Description,
		DocCount:           int32(knowledge.DocCount),
		EmbeddingModelInfo: embeddingModelInfo,
	}
}

// buildKnowledgeInfoList 构造知识库信息列表
func buildKnowledgeInfoList(knowledgeList []*model.KnowledgeBase) *knowledgebase_service.KnowledgeDetailSelectListResp {
	var retList []*knowledgebase_service.KnowledgeInfo
	for _, v := range knowledgeList {
		info := buildKnowledgeInfo(v)
		retList = append(retList, info)
	}
	return &knowledgebase_service.KnowledgeDetailSelectListResp{
		List:  retList,
		Total: int32(len(retList)),
	}
}

// buildKnowledgeBaseModel 构造知识库模型
func buildKnowledgeBaseModel(req *knowledgebase_service.CreateKnowledgeReq) (*model.KnowledgeBase, error) {
	embeddingModelInfo, err := json.Marshal(req.EmbeddingModelInfo)
	if err != nil {
		return nil, err
	}
	return &model.KnowledgeBase{
		KnowledgeId:    generator.GetGenerator().NewID(),
		Name:           req.Name,
		Description:    req.Description,
		OrgId:          req.OrgId,
		UserId:         req.UserId,
		EmbeddingModel: string(embeddingModelInfo),
		CreatedAt:      time.Now().UnixMilli(),
		UpdatedAt:      time.Now().UnixMilli(),
	}, nil
}

// buildKnowledgeNameList 构造知识库名称
func buildKnowledgeNameList(knowledgeList []*model.KnowledgeBase) []string {
	if len(knowledgeList) == 0 {
		return make([]string, 0)
	}
	var knowledgeNameList []string
	for _, knowledge := range knowledgeList {
		knowledgeNameList = append(knowledgeNameList, knowledge.Name)
	}
	return knowledgeNameList
}

// buildKnowledgeBaseHitResp 构造知识库命中返回
func buildKnowledgeBaseHitResp(ragKnowledgeHitResp *rag_service.RagKnowledgeHitResp) *knowledgebase_service.KnowledgeHitResp {
	knowledgeHitData := ragKnowledgeHitResp.Data
	var searchList = make([]*knowledgebase_service.KnowledgeSearchInfo, 0)
	list := knowledgeHitData.SearchList
	if len(list) > 0 {
		for _, search := range list {
			searchList = append(searchList, &knowledgebase_service.KnowledgeSearchInfo{
				Title:         search.Title,
				Snippet:       search.Snippet,
				KnowledgeName: search.KbName,
			})
		}
	}
	return &knowledgebase_service.KnowledgeHitResp{
		Prompt:     knowledgeHitData.Prompt,
		Score:      knowledgeHitData.Score,
		SearchList: searchList,
	}
}

// buildRerankId 构造重排序模型id
func buildRerankId(priorityType int32, rerankId string) string {
	if priorityType == 1 {
		return ""
	}
	return rerankId
}

// buildRetrieveMethod 构造检索方式
func buildRetrieveMethod(matchType string) string {
	switch matchType {
	case "vector":
		return "semantic_search"
	case "text":
		return "full_text_search"
	case "mix":
		return "hybrid_search"
	}
	return ""
}

// buildRerankMod 构造重排序模式
func buildRerankMod(priorityType int32) string {
	if priorityType == 1 {
		return "weighted_score"
	}
	return "rerank_model"
}

// buildWeight 构造权重信息
func buildWeight(priorityType int32, semanticsPriority float32, keywordPriority float32) *rag_service.WeightParams {
	if priorityType != 1 {
		return nil
	}
	return &rag_service.WeightParams{
		VectorWeight: semanticsPriority,
		TextWeight:   keywordPriority,
	}
}
