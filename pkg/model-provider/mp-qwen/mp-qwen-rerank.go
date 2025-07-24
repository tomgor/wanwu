package mp_qwen

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/UnicomAI/wanwu/pkg/log"
	mp_common "github.com/UnicomAI/wanwu/pkg/model-provider/mp-common"
	"github.com/go-resty/resty/v2"
	"io"
	"net/url"
)

type Rerank struct {
	ApiKey      string `json:"apiKey"`      // ApiKey
	EndpointUrl string `json:"endpointUrl"` // 推理url
}

func (cfg *Rerank) NewReq(req *mp_common.RerankReq) (mp_common.IRerankReq, error) {
	m := map[string]interface{}{
		"model": req.Model,
		"input": map[string]interface{}{
			"documents": req.Documents,
			"query":     req.Query,
		},
	}
	parameters := make(map[string]interface{})
	if req.TopN != nil || req.ReturnDocuments != nil {
		if req.TopN != nil {
			parameters["top_n"] = req.TopN
		}
		if req.ReturnDocuments != nil {
			parameters["return_documents"] = req.ReturnDocuments
		}
		m["parameters"] = parameters
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
		SetTimeout(0). // 关闭请求超时
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
		return nil, fmt.Errorf("request %v huoshan rerank err: %v", url, err)
	} else if resp.StatusCode() >= 300 {
		return nil, fmt.Errorf("request %v huoshan rerank http status %v msg: %v", url, resp.StatusCode(), resp.String())
	}
	b, err := io.ReadAll(resp.RawResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("request %v huoshan rerank read response body err: %v", url, err)
	}
	return &rerankResp{raw: string(b)}, nil
}

func (cfg *Rerank) rerankUrl() string {
	ret, _ := url.JoinPath(cfg.EndpointUrl, "/services/rerank/text-rerank/text-rerank")
	return ret
}

// --- rerankResp ---

type rerankResp struct {
	raw string

	Output    RerankRespOutput `json:"output"`
	Usage     mp_common.Usage  `json:"usage"`
	RequestId string           `json:"request_id"`
}

type RerankRespOutput struct {
	Results []mp_common.Result `json:"results"`
}

func (resp *rerankResp) String() string {
	return resp.raw
}

func (resp *rerankResp) Data() (interface{}, bool) {
	ret := []map[string]interface{}{}
	if err := json.Unmarshal([]byte(resp.raw), &ret); err != nil {
		log.Errorf("huoshan rerank resp (%v) convert to data err: %v", resp.raw, err)
		return nil, false
	}
	return ret, true
}
func (resp *rerankResp) ConvertResp() (*mp_common.RerankResp, bool) {
	data := &rerankResp{}
	if err := json.Unmarshal([]byte(resp.raw), &data); err != nil {
		log.Errorf("huoshan rerank resp (%v) convert to data err: %v", resp.raw, err)
		return nil, false
	}
	res := &mp_common.RerankResp{
		Results:   data.Output.Results,
		Usage:     data.Usage,
		RequestId: data.RequestId,
	}
	return res, true
}
