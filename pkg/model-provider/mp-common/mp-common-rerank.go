package mp_common

import (
	"encoding/json"
	"fmt"

	"github.com/UnicomAI/wanwu/pkg/log"
)

// --- openapi request ---

type RerankReq struct {
	Documents       []string `json:"documents" validate:"required"`
	Model           string   `json:"model" validate:"required"`
	Query           string   `json:"query" validate:"required"`
	ReturnDocuments *bool    `json:"return_documents,omitempty"`
	TopN            *int     `json:"top_n,omitempty"`
}

func (req *RerankReq) Check() error {
	if req.TopN != nil && *req.TopN < 0 {
		return fmt.Errorf("top_n must greater than 0")
	}
	return nil
}

// --- openapi response ---

type RerankResp struct {
	Results   []Result `json:"results"`
	Model     string   `json:"model"`
	Object    *string  `json:"object,omitempty"`
	Usage     Usage    `json:"usage"`
	RequestId *string  `json:"request_id,omitempty"`
}

type Result struct {
	Index          int       `json:"index"`
	Document       *Document `json:"document,omitempty"`
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

// --- request ---

type IRerankReq interface {
	Data() map[string]interface{}
}

// rerankReq implementation of IRerankReq
type rerankReq struct {
	data map[string]interface{}
}

func NewRerankReq(data map[string]interface{}) IRerankReq {
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
