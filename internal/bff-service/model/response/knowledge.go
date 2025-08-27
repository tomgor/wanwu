package response

type KnowledgeListResp struct {
	KnowledgeList []*KnowledgeInfo `json:"knowledgeList"`
}

type CreateKnowledgeResp struct {
	KnowledgeId string `json:"knowledgeId"`
}

type KnowledgeHitResp struct {
	Prompt     string             `json:"prompt"`     //提示词列表
	SearchList []*ChunkSearchList `json:"searchList"` //种种结果
	Score      []float64          `json:"score"`      //打分信息
}

type EmbeddingModelInfo struct {
	ModelId string `json:"modelId"`
}

type KnowledgeInfo struct {
	KnowledgeId        string              `json:"knowledgeId"`        //知识库id
	Name               string              `json:"name"`               //知识库名称
	Description        string              `json:"description"`        //知识库描述
	DocCount           int                 `json:"docCount"`           //文档数量
	EmbeddingModelInfo *EmbeddingModelInfo `json:"embeddingModelInfo"` //embedding模型信息
	KnowledgeTagList   []*KnowledgeTag     `json:"knowledgeTagList"`   //知识库标签列表
	CreateAt           string              `json:"createAt"`           //创建时间
}

type ChunkSearchList struct {
	Title         string `json:"title"`
	Snippet       string `json:"snippet"`
	KnowledgeName string `json:"knowledgeName"`
}
