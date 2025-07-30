package service

import (
	queue_util "github.com/UnicomAI/wanwu/internal/bff-service/pkg/queue-util"

	safety_service "github.com/UnicomAI/wanwu/api/proto/safety-service"
	"github.com/UnicomAI/wanwu/internal/bff-service/pkg/ahocorasick"
	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/gin-gonic/gin"
)

func InitAddChatContainer() {
	AddChatContainer(&RagChatService{})
	AddChatContainer(&AgentChatService{})
}

var serviceMap = make(map[string]ChatService)

func AddChatContainer(service ChatService) {
	serviceMap[service.buildChatType()] = service
}

type ChatService interface {
	buildChatType() string
	buildSensitiveResp(id, content string) string
	buildContent(text string) (contentList []string, id string)
}

// 构建敏感词字典
func BuildSensitiveDict(ctx *gin.Context, ids []string) ([]ahocorasick.DictConfig, error) {
	var dicts []ahocorasick.DictConfig
	sensitiveTableResp, err := safety.GetSensitiveWordTableWithVersion(ctx, &safety_service.GetSensitiveWordTableWithVersionReq{
		TableIds: ids,
	})
	if err != nil {
		return nil, err
	}
	for _, index := range sensitiveTableResp.Details {
		dicts = append(dicts, ahocorasick.DictConfig{
			DictID:  index.TableId,
			Version: index.Version,
		})
	}
	// 检测内存中的敏感词表
	status, err := ahocorasick.CheckDictStatus(dicts)
	if err != nil {
		return nil, err
	}
	// 拼接id,version与内存不匹配的tableID
	var needLoadWords []string
	var ret []ahocorasick.DictConfig
	for _, stat := range status {
		if !stat.Status {
			needLoadWords = append(needLoadWords, stat.DictCfg.DictID)
		}
		ret = append(ret, ahocorasick.DictConfig{
			DictID:  stat.DictCfg.DictID,
			Version: stat.DictCfg.Version,
		})
	}
	// 访问safey 更新词表信息
	tableWithWords, err := safety.GetSensitiveWordTableWithWord(ctx, &safety_service.GetSensitiveWordTableWithWordReq{
		TableIds: needLoadWords,
	})
	if err != nil {
		return nil, err
	}
	// 重新构建version不匹配的词表
	for _, detail := range tableWithWords.Details {
		dict := ahocorasick.DictConfig{
			DictID:  detail.TableId,
			Version: detail.Version,
		}
		err := ahocorasick.BuildDict(dict, detail.Reply, detail.SensitiveWords)
		if err != nil {
			log.Errorf("build dict id %v & dict version %v err:%v", dict.DictID, dict.Version, err)
		}
		ret = append(ret, ahocorasick.DictConfig{
			DictID:  detail.TableId,
			Version: detail.Version,
		})
	}
	return ret, nil
}

// processSensitiveWords 中间处理函数，负责敏感词检测并返回处理后的通道
func ProcessSensitiveWords(ctx *gin.Context, inputCh <-chan string, matchDicts []ahocorasick.DictConfig, serviceType string) <-chan string {
	appService := serviceMap[serviceType]
	outputCh := make(chan string, defaultChannelSize)
	go func() {
		defer close(outputCh)
		// 初始化队列
		var id string
		queue := queue_util.NewOverridableQueue(defaultCheckWindowSize)
		outputQueue := queue_util.NewBoundedQueue(defaultOutputWindowSize)
		for text := range inputCh {
			convertTextList, id := appService.buildContent(text)
			if len(convertTextList) > 0 {
				for _, s := range convertTextList {
					queue.EnQueue(s)
				}
			}
			if outputQueue.IsFull() {
				// 校验敏感词
				if !checkStreamResponse(outputCh, queue, matchDicts, appService, id) {
					return
				}
				// 输出队列内容
				for !outputQueue.IsEmpty() {
					if dequeue, ok := outputQueue.Dequeue(); ok {
						outputCh <- dequeue
					}
				}
			}
			outputQueue.Enqueue(text)
		}
		// 处理剩余内容
		if !checkStreamResponse(outputCh, queue, matchDicts, appService, id) {
			return
		}
		valueList := outputQueue.AllValue()
		if len(valueList) > 0 {
			for _, value := range valueList {
				outputCh <- value
			}
		}
	}()
	return outputCh
}

func checkStreamResponse(retChannel chan string, queue *queue_util.OverridableCircularQueue, matchDicts []ahocorasick.DictConfig, appService ChatService, id string) bool {
	//准备出队的时候再校验 并且将进行出队
	ret, err := ahocorasick.ContentMatch(queue.AllValue(), matchDicts, true)
	if err != nil {
		return false
	}
	if len(ret) > 0 {
		messageList := appService.buildSensitiveResp(id, ret[0].Reply)
		if len(messageList) > 0 {
			for _, message := range []string{messageList} {
				retChannel <- "\n" // 流式返回 两次返回之间一定要加空行，否则前端或者智能体解析有问题
				retChannel <- message
				retChannel <- "\n"
			}
		}
		return false
	}
	return true
}
