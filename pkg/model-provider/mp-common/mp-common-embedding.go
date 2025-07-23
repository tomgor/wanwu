package mp_common

import (
	"encoding/json"

	"github.com/UnicomAI/wanwu/pkg/log"
)

// --- openapi request ---

type EmbeddingReq struct {
	Model          string   `json:"model" validate:"required"`
	Input          []string `json:"input" validate:"required"`
	EncodingFormat *string  `json:"encoding_format,omitempty"`
}

func (req *EmbeddingReq) Check() error { return nil }

// --- openapi response ---

type EmbeddingResp struct {
	Model  string          `json:"model"`
	Object string          `json:"object"`
	Data   []EmbeddingData `json:"data"`
	Usage  Usage           `json:"usage"`
}

type EmbeddingData struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

// --- request ---

type IEmbeddingReq interface {
	Data() map[string]interface{}
}

// embeddingReq implementation of IEmbeddingReq
type embeddingReq struct {
	data map[string]interface{}
}

func NewEmbeddingReq(data map[string]interface{}) IEmbeddingReq {
	return &embeddingReq{data: data}
}

func (req *embeddingReq) Data() map[string]interface{} {
	return req.data
}

// --- response ---

type IEmbeddingResp interface {
	String() string
	Data() (map[string]interface{}, bool)
	ConvertResp() (*EmbeddingResp, bool)
}

// embeddingResp implementation of IEmbeddingResp
type embeddingResp struct {
	raw string
}

func NewEmbeddingResp(raw string) IEmbeddingResp {
	return &embeddingResp{raw: raw}
}

func (resp *embeddingResp) String() string {
	return resp.raw
}

func (resp *embeddingResp) Data() (map[string]interface{}, bool) {
	ret := make(map[string]interface{})
	if err := json.Unmarshal([]byte(resp.raw), &ret); err != nil {
		log.Errorf("embedding resp (%v) convert to data err: %v", resp.raw, err)
		return nil, false
	}
	return ret, true
}

func (resp *embeddingResp) ConvertResp() (*EmbeddingResp, bool) {
	var ret *EmbeddingResp
	if err := json.Unmarshal([]byte(resp.raw), &ret); err != nil {
		log.Errorf("embedding resp (%v) convert to data err: %v", resp.raw, err)
		return nil, false
	}
	return ret, true
}
