package service

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	app_service "github.com/UnicomAI/wanwu/api/proto/app-service"
	assistant_service "github.com/UnicomAI/wanwu/api/proto/assistant-service"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/pkg/ahocorasick"
	"github.com/UnicomAI/wanwu/pkg/constant"
	"github.com/UnicomAI/wanwu/pkg/log"
	sse_util "github.com/UnicomAI/wanwu/pkg/sse-util"
	"github.com/UnicomAI/wanwu/pkg/util"
	"github.com/gin-gonic/gin"
)

type AgentChatService struct{}

// buildContent implements ChatService.
func (a *AgentChatService) buildContent(text string) (contentList []string, id string) {
	text = strings.TrimPrefix(text, "data:")
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, ""
	}
	resp := struct {
		Code           int      `json:"code"`
		Message        string   `json:"message"`
		Response       string   `json:"response"`
		MsgID          string   `json:"msg_id"`
		GenFileURLList []string `json:"gen_file_url_list"`
		History        []struct {
			Response string `json:"response"`
		} `json:"history"`
		Finish int `json:"finish"`
		Usage  struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
		SearchList []string `json:"search_list"`
		QAType     []int    `json:"qa_type"`
	}{}
	if err := json.Unmarshal([]byte(text), &resp); err != nil {
		return nil, ""
	}
	var contents []string
	if resp.Response != "" {
		contents = append(contents, resp.Response)
	}
	return contents, resp.MsgID
}

// buildChatType implements ChatService.
func (a *AgentChatService) buildChatType() string {
	return constant.AppTypeAgent
}

// buildSensitiveResp implements ChatService.
func (a *AgentChatService) buildSensitiveResp(id string, content string) string {
	resp := map[string]interface{}{
		"code":              0,
		"message":           "success",
		"response":          content,
		"gen_file_url_list": []interface{}{},
		"history":           []interface{}{},
		"finish":            1, // Note: The original JSON has "finish" misspelled as "finish"
		"usage": map[string]interface{}{
			"prompt_tokens":     0,
			"completion_tokens": 0,
			"total_tokens":      0,
		},
		"search_list": []interface{}{},
		"qa_type":     []int{1},
	}

	marshal, _ := json.Marshal(resp)
	return "data:" + string(marshal)
}
func CallAssistantConversationStream(ctx *gin.Context, userId, orgId string, req request.ConversionStreamRequest) (<-chan string, error) {
	// 根据agentID获取敏感词配置
	agentInfo, err := assistant.GetAssistantInfo(ctx, &assistant_service.GetAssistantInfoReq{
		AssistantId: req.AssistantId,
	})
	if err != nil {
		return nil, err
	}
	var matchDicts []ahocorasick.DictConfig
	// 如果Enable为true,则处理敏感词
	if agentInfo.SafetyConfig.Enable {
		var ids []string
		for _, idx := range agentInfo.SafetyConfig.SensitiveTable {
			ids = append(ids, idx.TableId)
		}
		matchDicts, err = BuildSensitiveDict(ctx, ids)
		if err != nil {
			return nil, err
		}
		ret, err := ahocorasick.ContentMatch(req.Prompt, matchDicts, true)
		if err != nil {
			return nil, err
		}
		if len(ret) > 0 {
			return nil, fmt.Errorf("您的提问中含有敏感词:%v", ret[0].Reply)
		}
	}
	appList, err := app.GetAppListByIds(ctx.Request.Context(), &app_service.GetAppListByIdsReq{
		AppIdsList: []string{req.AssistantId},
	})
	if err != nil {
		return nil, err
	}
	var publishType string
	if len(appList.Infos) > 0 {
		publishType = appList.Infos[0].PublishType
	}
	stream, err := assistant.AssistantConversionStream(ctx.Request.Context(), &assistant_service.AssistantConversionStreamReq{
		AssistantId:    req.AssistantId,
		ConversationId: req.ConversationId,
		FileInfo: &assistant_service.ConversionStreamFile{
			FileName: req.FileInfo.FileName,
			FileSize: req.FileInfo.FileSize,
			FileUrl:  req.FileInfo.FileUrl,
		},
		Trial:          req.Trial,
		Prompt:         req.Prompt,
		AppPublishType: publishType,
		Identity: &assistant_service.Identity{
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
		log.Infof("[Agent] %v conversation %v user %v org %v start, query: %s", req.AssistantId, req.ConversationId, userId, orgId, req.Prompt)
		for {
			s, err := stream.Recv()
			if err == io.EOF {
				log.Infof("[Agent] %v conversation %v user %v org %v stop", req.AssistantId, req.ConversationId, userId, orgId)
				break
			}
			if err != nil {
				log.Errorf("[Agent] %v conversation %v user %v org %v recv err: %v", req.AssistantId, req.ConversationId, userId, orgId, err)
				break
			}
			ret <- s.Content
		}
	}()
	if !agentInfo.SafetyConfig.Enable {
		return ret, nil
	}
	// 敏感词过滤
	filteredCh := ProcessSensitiveWords(ctx, ret, matchDicts, constant.AppTypeAgent)
	return filteredCh, nil
}

func AssistantConversionStream(ctx *gin.Context, userId, orgId string, req request.ConversionStreamRequest) error {

	// 1. CallAssistantConversationStream
	chatCh, err := CallAssistantConversationStream(ctx, userId, orgId, req)
	if err != nil {
		return err
	}
	// 2. 流式返回结果
	_ = sse_util.NewSSEWriter(ctx, fmt.Sprintf("[Agent] %v conversation %v user %v org %v recv", req.AssistantId, req.ConversationId, userId, orgId), sse_util.DONE_MSG).
		WriteStream(chatCh, nil, buildAgentChatRespLineProcessor(), nil)
	return nil
}

// buildAgentChatRespLineProcessor 构造agent对话结果行处理器
func buildAgentChatRespLineProcessor() func(*gin.Context, string, interface{}) (string, bool, error) {
	return func(c *gin.Context, lineText string, params interface{}) (string, bool, error) {
		if strings.HasPrefix(lineText, "error:") {
			errorText := fmt.Sprintf("data: {\"code\": \"-1\", \"message\": \"%s\"}\n\n", strings.TrimPrefix(lineText, "error:"))
			return errorText, false, nil
		}
		if strings.HasPrefix(lineText, "data:") {
			return lineText + "\n\n", false, nil
		}
		return "data:" + lineText + "\n\n", false, nil
	}
}
