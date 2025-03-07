// Copyright 2025 fsyyft-ai
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

// package task 包含 MCP 时间服务相关的任务实现。
package task

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/fsyyft-ai/mcp-time/internal/config"
	"github.com/fsyyft-go/kit/log"
)

type (
	// MCPTimeTask 定义了MCP时间任务的接口。
	MCPTimeTask interface {
		// Run 执行MCP时间任务。
		Run(ctx context.Context) error
	}

	// mcpTimeTask 实现了 MCPTimeTask 接口，提供获取不同时区当前时间的功能。
	mcpTimeTask struct {
		// logger 用于记录任务执行过程中的日志信息。
		logger log.Logger
		// cfg 存储应用配置信息。
		cfg *config.Config
	}
)

// NewMCPTimeTask 创建并返回一个新的 MCPTimeTask 接口实例。
// 参数：
//   - logger: 用于记录日志的 Logger 实例。
//   - cfg: 应用配置信息。
//
// 返回：
//   - MCPTimeTask 接口的实现实例。
func NewMCPTimeTask(logger log.Logger, cfg *config.Config) MCPTimeTask {
	return &mcpTimeTask{
		logger: logger.WithField("ddd", "task"),
		cfg:    cfg,
	}
}

// Run 实现了 MCPTimeTask 接口的 Run 方法，启动 MCP 时间服务。
// 参数：
//   - ctx: 上下文，用于控制任务的生命周期。
//
// 返回：
//   - error: 如果启动过程中出现错误，返回相应的错误信息。
func (t *mcpTimeTask) Run(ctx context.Context) error {
	// 创建 MCP 服务器。
	s := server.NewMCPServer(
		"MCP Time", // 服务名称。
		"0.0.1",    // 服务版本号。
	)
	// 添加工具。
	tool := mcp.NewTool("current time",
		// 设置工具描述，说明默认时区为 Asia/Shanghai。
		mcp.WithDescription("Get current time with timezone, Asia/Shanghai is default"),
		// 添加名为 timezone 的字符串参数。
		mcp.WithString("timezone",
			mcp.Required(),                           // 标记参数为必填项。
			mcp.Description("current time timezone"), // 参数描述。
		),
	)
	// 为工具添加处理器。
	s.AddTool(tool, t.currentTimeHandler)
	// 启动标准输入输出服务。
	if err := server.ServeStdio(s); err != nil {
		// 如果服务启动出错，记录错误日志。
		t.logger.Error("Server error: %v\n", err)
	}

	return nil
}

// currentTimeHandler 处理获取当前时间的请求。
// 参数：
//   - ctx: 上下文对象。
//   - request: MCP 工具调用请求，包含时区参数。
//
// 返回：
//   - *mcp.CallToolResult: 工具调用结果，包含当前时间信息或错误信息。
//   - error: 处理过程中的错误，正常情况下为 nil。
func (t *mcpTimeTask) currentTimeHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 从请求参数中获取时区字符串。
	timezone, ok := request.Params.Arguments["timezone"].(string)
	if !ok {
		// 如果时区参数不是字符串类型，返回错误。
		return mcp.NewToolResultError("timezone must be a string"), nil
	}

	// 加载指定的时区。
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		// 如果时区加载失败，返回错误信息。
		return mcp.NewToolResultError(fmt.Sprintf("parse timezone with error: %v", err)), nil
	}
	// 返回指定时区的当前时间。
	return mcp.NewToolResultText(fmt.Sprintf(`current time is %s`, time.Now().In(loc))), nil
}
