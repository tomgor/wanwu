package mp_common

import (
	"encoding/json"

	"github.com/UnicomAI/wanwu/pkg/log"
)

// --- request ---

type IRerankReq interface {
	Data() map[string]interface{}
}
type RerankReq struct {
	Documents       []string `json:"documents" validate:"required"`
	Model           string   `json:"model" validate:"required"`
	Query           string   `json:"query" validate:"required"`
	ReturnDocuments bool     `json:"return_documents"`
	TopN            int      `json:"top_n" validate:"gte=0"`
}

type RerankResp struct {
	Results []Result `json:"results"`
	Model   string   `json:"model"`
	Object  string   `json:"object"`
	Usage   Usage    `json:"usage"`
}
type Result struct {
	Index          int       `json:"index"`
	Document       *Document `json:"document"`
	RelevanceScore float64   `json:"relevance_score"`
}

type Document struct {
	Text string `json:"text"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
	CompletionTokens int `json:"completion_tokens"`
}

// rerankReq implementation of IRerankReq
type rerankReq struct {
	data map[string]interface{}
}

func NewRerankReq(data map[string]interface{}) IRerankReq {
	if val, exists := data["top_n"]; exists {
		if num, ok := val.(int); ok && num <= 0 {
			delete(data, "top_n")
		}
	}
	return &rerankReq{data: data}
}

func (req *rerankReq) Data() map[string]interface{} {
	return req.data
}

// --- response ---

type IRerankResp interface {
	String() string
	Data() (interface{}, bool)
	ConvertResp() (*RerankResp, bool)
}

// rerankResp implementation of IRerankResp
type rerankResp struct {
	raw string
}

func NewRerankResp(raw string) IRerankResp {
	return &rerankResp{raw: raw}
}

func (resp *rerankResp) String() string {
	return resp.raw
}

func (resp *rerankResp) Data() (interface{}, bool) {
	ret := make(map[string]interface{})
	if err := json.Unmarshal([]byte(resp.raw), &ret); err != nil {
		log.Errorf("rerank resp (%v) convert to data err: %v", resp.raw, err)
		return nil, false
	}
	return ret, true
}

func (resp *rerankResp) ConvertResp() (*RerankResp, bool) {
	var ret *RerankResp
	if err := json.Unmarshal([]byte(resp.raw), &ret); err != nil {
		log.Errorf("rerank resp (%v) convert to data err: %v", resp.raw, err)
		return nil, false
	}
	return ret, true
}
