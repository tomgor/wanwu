package service

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	rag_service "github.com/UnicomAI/wanwu/api/proto/rag-service"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/pkg/ahocorasick"
	"github.com/UnicomAI/wanwu/pkg/constant"
	"github.com/UnicomAI/wanwu/pkg/log"
	sse_util "github.com/UnicomAI/wanwu/pkg/sse-util"
	"github.com/UnicomAI/wanwu/pkg/util"
	"github.com/gin-gonic/gin"
)

type RagChatService struct{}

func (dp *RagChatService) buildContent(text string) (contentList []string, id string) {
	// 1. 清理数据前缀
	text = strings.TrimPrefix(text, "data:")
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, ""
	}
	// 2. 解析JSON
	resp := struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		MsgID   string `json:"msg_id"`
		Data    struct {
			Output string `json:"output"`
		} `json:"data"`
		History []struct {
			Response string `json:"response"`
		} `json:"history"`
		Finish int `json:"finish"`
	}{}

	if err := json.Unmarshal([]byte(text), &resp); err != nil {
		return nil, ""
	}
	// 3. 构建返回内容
	var contents []string
	// 主输出内容
	if resp.Data.Output != "" {
		contents = append(contents, resp.Data.Output)
	}
	return contents, resp.MsgID
}

func (dp *RagChatService) buildSensitiveResp(id string, content string) string {
	resp := map[string]interface{}{
		"code":    0,
		"message": "success",
		"msg_id":  id,
		"data": map[string]interface{}{
			"output":     content,
			"searchList": []interface{}{},
		},
		"history": []interface{}{},
		"finish":  0,
	}

	marshal, _ := json.Marshal(resp)
	return "data:" + string(marshal)
}

func (dp *RagChatService) buildChatType() string {
	return constant.AppTypeRag
}

// ChatRagStream rag私域问答
func ChatRagStream(ctx *gin.Context, userId, orgId string, req request.ChatRagRequest) error {
	// 1.CallRagChatStream
	chatCh, err := CallRagChatStream(ctx, userId, orgId, req)
	if err != nil {
		return err
	}
	// 2.流式返回结果
	_ = sse_util.NewSSEWriter(ctx, fmt.Sprintf("[RAG] %v user %v org %v", req.RagID, userId, orgId), sse_util.DONE_MSG).
		WriteStream(chatCh, nil, buildRagChatRespLineProcessor(), nil)
	return nil
}

// CallRagChatStream 调用Rag对话
func CallRagChatStream(ctx *gin.Context, userId, orgId string, req request.ChatRagRequest) (<-chan string, error) {
	// 根据ragID获取敏感词配置
	ragInfo, err := rag.GetRagDetail(ctx, &rag_service.RagDetailReq{
		RagId: req.RagID,
	})
	if err != nil {
		return nil, err
	}
	var matchDicts []ahocorasick.DictConfig
	// 如果Enable为true,则处理敏感词
	if ragInfo.SensitiveConfig.Enable {
		matchDicts, err = BuildSensitiveDict(ctx, ragInfo.SensitiveConfig.TableIds)
		if err != nil {
			return nil, err
		}
		ret, err := ahocorasick.ContentMatch(req.Question, matchDicts, true)
		if err != nil {
			return nil, err
		}
		if len(ret) > 0 {
			return nil, fmt.Errorf("您的提问中含有敏感词:%v", ret[0].Reply)
		}
	}
	stream, err := rag.ChatRag(ctx, &rag_service.ChatRagReq{
		RagId:    req.RagID,
		Question: req.Question,
		Identity: &rag_service.Identity{
			UserId: userId,
			OrgId:  orgId,
		},
	})
	if err != nil {
		return nil, err
	}

	ret := make(chan string, 128)
	go func() {
		defer util.PrintPanicStack()
		defer close(ret)
		log.Infof("[RAG] %v user %v org %v start, query: %s", req.RagID, userId, orgId, req.Question)
		for {
			s, err := stream.Recv()
			if err == io.EOF {
				log.Infof("[RAG] %v user %v org %v stop", req.RagID, userId, orgId)
				break
			}
			if err != nil {
				log.Errorf("[RAG] %v user %v org %v recv err: %v", req.RagID, userId, orgId, err)
				break
			}
			ret <- s.Content
		}
	}()
	if !ragInfo.SensitiveConfig.Enable {
		return ret, nil
	}
	// 敏感词过滤
	filteredCh := ProcessSensitiveWords(ctx, ret, matchDicts, constant.AppTypeRag)
	return filteredCh, nil
}

// buildRagChatRespLineProcessor 构造rag对话结果行处理器
func buildRagChatRespLineProcessor() func(*gin.Context, string, interface{}) (string, bool, error) {
	return func(c *gin.Context, lineText string, params interface{}) (string, bool, error) {
		if strings.HasPrefix(lineText, "error:") {
			errorText := fmt.Sprintf("data: {\"code\": \"-1\", \"msg\": \"%s\"}\n\n", strings.TrimPrefix(lineText, "error:"))
			return errorText, false, nil
		}
		if strings.HasPrefix(lineText, "data:") {
			return lineText + "\n\n", false, nil
		}
		return lineText + "\n\n", false, nil
	}
}
