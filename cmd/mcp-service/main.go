package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/UnicomAI/wanwu/internal/mcp-service/client/orm"
	"github.com/UnicomAI/wanwu/internal/mcp-service/config"
	"github.com/UnicomAI/wanwu/internal/mcp-service/server/grpc"
	"github.com/UnicomAI/wanwu/pkg/db"
	"github.com/UnicomAI/wanwu/pkg/log"
)

var (
	configFile   = flag.String("config", "configs/microservice/mcp-service/configs/config.yaml", "mcp-service config")
	isVersion    bool
	buildTime    string //编译时间
	buildVersion string //编译版本
	gitCommitID  string //git的commit id
	gitBranch    string //git branch
	builder      string //构建者

)

func main() {
	//打印编译信息
	flag.BoolVar(&isVersion, "v", false, "编译信息")
	flag.Parse()

	if len(os.Args) > 1 && !isVersion {
		flag.Usage()
		return
	}

	if isVersion {
		versionPrint()
		return
	}

	ctx := context.Background()

	flag.Parse()
	if err := config.LoadConfig(*configFile); err != nil {
		log.Fatalf("init cfg err: %v", err)
	}

	if err := log.InitLog(config.Cfg().Log.Std, config.Cfg().Log.Level, config.Cfg().Log.Logs...); err != nil {
		log.Fatalf("init log err: %v", err)
	}

	db, err := db.New(config.Cfg().DB)
	if err != nil {
		log.Fatalf("init db failed, err: %v", err)
	}

	c, err := orm.NewClient(ctx, db)
	if err != nil {
		log.Fatalf("init client failed, err: %v", err)
	}

	s, err := grpc.NewServer(config.Cfg(), c)
	if err != nil {
		log.Fatalf("init server failed, err: %v", err)
	}
	if err := s.Start(ctx); err != nil {
		log.Fatalf("start grpc server failed, err: %s", err)
	}
	//detail, err := s.Mcp.GetSquareMCP(ctx, &mcp_service.GetSquareMCPReq{
	//	McpSquareId: "converttime",
	//})
	//if err != nil {
	//	return
	//}
	//if detail == nil {
	//	log.Errorf("MCP detail not found")
	//	return
	//}
	//
	////打印基本信息
	//printSection := func(title string, items map[string]string) {
	//	fmt.Printf("\n=== %s ===\n", title)
	//	for k, v := range items {
	//		fmt.Printf("%-15s: %s\n", k, v)
	//	}
	//}
	//
	//// 打印Info部分
	//info := map[string]string{
	//	"McpSquareId": detail.Info.McpSquareId,
	//	"Name":        detail.Info.Name,
	//	"AvatarPath":  detail.Info.AvatarPath,
	//	"Desc":        detail.Info.Desc,
	//	"Category":    detail.Info.Category,
	//	"From":        detail.Info.From,
	//}
	//printSection("基本信息", info)
	//
	//// 打印Intro部分
	//intro := map[string]string{
	//	"Detail":   detail.Intro.Detail,
	//	"Summary":  detail.Intro.Summary,
	//	"Feature":  detail.Intro.Feature,
	//	"Scenario": detail.Intro.Scenario,
	//	"Manual":   detail.Intro.Manual,
	//}
	//printSection("介绍信息", intro)
	//
	//// 打印Tool部分
	//fmt.Printf("\n=== 工具信息 ===\n")
	//fmt.Printf("SSE URL: %s\n", detail.Tool.SseUrl)
	//fmt.Printf("HasCustom: %v\n", detail.Tool.HasCustom)
	//
	//// 打印Tools列表
	//for _, tool := range detail.Tool.Tools {
	//	fmt.Printf("工具名称: %s\n", tool.Name)
	//	fmt.Printf("描述: %s\n", tool.Description)
	//	fmt.Printf("InputSchema: \n")
	//	fmt.Printf("输入类型: %s\n", tool.InputSchema.Type)
	//	fmt.Printf("必填字段: %v\n", tool.InputSchema.Required)
	//
	//	fmt.Println("properties:")
	//	for name, field := range tool.InputSchema.Properties {
	//		fmt.Printf("属性名: %s\n", name)
	//		fmt.Printf("- %s: %s\n", field.Type, field.Description)
	//	}
	//}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGTERM)
	<-sc
	s.Stop(ctx)
}

func versionPrint() {
	fmt.Printf("build_time: %s\n", buildTime)
	fmt.Printf("build_version: %s\n", buildVersion)
	fmt.Printf("git_commit_id: %s\n", gitCommitID)
	fmt.Printf("git branch: %s\n", gitBranch)
	fmt.Printf("runtime version: %s\n", runtime.Version())
	fmt.Printf("builder: %s\n", builder)
}
