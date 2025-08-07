package assistant

import (
	"fmt"

	knowledgeBase_service "github.com/UnicomAI/wanwu/api/proto/knowledgebase-service"
	mcp_service "github.com/UnicomAI/wanwu/api/proto/mcp-service"
	"github.com/UnicomAI/wanwu/internal/assistant-service/config"
	"github.com/UnicomAI/wanwu/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	maxMsgSize            = 1024 * 1024 * 4 // 4M
	headlessServiceSchema = "dns:///"
)

var (
	Knowledge knowledgeBase_service.KnowledgeBaseServiceClient
	MCP       mcp_service.MCPServiceClient
)

func StartService() error {
	// grpc connections
	knowledgeConn, err := newConn(config.Cfg().Knowledge.Host)
	if err != nil {
		return fmt.Errorf("init knowledgebase-service connection err: %v", err)
	}
	Knowledge = knowledgeBase_service.NewKnowledgeBaseServiceClient(knowledgeConn)
	log.Infof("Knowledge init success")
	log.Infof("Knowledge: %s", config.Cfg().Knowledge.Host)

	MCPConn, err := newConn(config.Cfg().MCP.Host)
	if err != nil {
		return fmt.Errorf("init mcp-service connection err: %v", err)
	}
	MCP = mcp_service.NewMCPServiceClient(MCPConn)
	log.Infof("MCP init success")
	log.Infof("MCP: %s", config.Cfg().MCP.Host)
	return nil
}

func newConn(host string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(headlessServiceSchema+host,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(maxMsgSize),
			grpc.MaxCallSendMsgSize(maxMsgSize)),
	)
	if err != nil {
		return nil, err
	}
	return conn, err
}
