package mp_yuanjing

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/UnicomAI/wanwu/pkg/log"
	mp_common "github.com/UnicomAI/wanwu/pkg/model-provider/mp-common"
	"github.com/go-resty/resty/v2"
)

type Rerank struct {
	ApiKey      string `json:"apiKey"`      // ApiKey
	EndpointUrl string `json:"endpointUrl"` // 推理url
}

func (cfg *Rerank) NewReq(req *mp_common.RerankReq) (mp_common.IRerankReq, error) {
	m := map[string]interface{}{
		"texts": req.Documents,
		"query": req.Query,
	}
	return mp_common.NewRerankReq(m), nil
}

func (cfg *Rerank) Rerank(ctx context.Context, req mp_common.IRerankReq, headers ...mp_common.Header) (mp_common.IRerankResp, error) {
	if cfg.ApiKey != "" {
		headers = append(headers, mp_common.Header{
			Key:   "Authorization",
			Value: "Bearer " + cfg.ApiKey,
		})
	}

	request := resty.New().
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}). // 关闭证书校验
		SetTimeout(0).                                             // 关闭请求超时
		R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(req.Data()).
		SetDoNotParseResponse(true)
	for _, header := range headers {
		request.SetHeader(header.Key, header.Value)
	}

	url := cfg.rerankUrl()
	resp, err := request.Post(url)
	if err != nil {
		return nil, fmt.Errorf("request %v yuanjing rerank err: %v", url, err)
	} else if resp.StatusCode() >= 300 {
		return nil, fmt.Errorf("request %v yuanjing rerank http status %v msg: %v", url, resp.StatusCode(), resp.String())
	}
	b, err := io.ReadAll(resp.RawResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("request %v yuanjing rerank read response body err: %v", url, err)
	}
	return &rerankResp{raw: string(b)}, nil
}

func (cfg *Rerank) rerankUrl() string {
	ret, _ := url.JoinPath(cfg.EndpointUrl, "/rerank")
	return ret
}

// --- rerankResp ---

type rerankResp struct {
	raw string

	Index    int     `json:"index"`
	Score    float64 `json:"score"`
	Document string  `json:"document"`
}

func (resp *rerankResp) String() string {
	return resp.raw
}

func (resp *rerankResp) Data() (interface{}, bool) {
	ret := []map[string]interface{}{}
	if err := json.Unmarshal([]byte(resp.raw), &ret); err != nil {
		log.Errorf("yuanjing rerank resp (%v) convert to data err: %v", resp.raw, err)
		return nil, false
	}
	return ret, true
}
func (resp *rerankResp) ConvertResp() (*mp_common.RerankResp, bool) {
	var data []map[string]interface{}
	if err := json.Unmarshal([]byte(resp.raw), &data); err != nil {
		log.Errorf("yuanjing rerank resp (%v) convert to data err: %v", resp.raw, err)
		return nil, false
	}

	var results []mp_common.Result
	for _, item := range data {
		b, err := json.Marshal(item)
		if err != nil {
			log.Errorf("yuanjing rerank resp (%v) item (%v) convert err: %v", resp.raw, item, err)
			return nil, false
		}
		if err = json.Unmarshal(b, resp); err != nil {
			log.Errorf("yuanjing rerank resp (%v) item (%v) unmarshal err: %v", resp.raw, item, err)
			return nil, false
		}
		results = append(results, mp_common.Result{
			Index:          resp.Index,
			RelevanceScore: resp.Score,
			Document: &mp_common.Document{
				Text: resp.Document,
			},
		})
	}

	return &mp_common.RerankResp{
		Results: results,
	}, true
}
