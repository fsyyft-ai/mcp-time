// Copyright 2025 fsyyft-ai
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

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

	mcpTimeTask struct {
		logger log.Logger
		cfg    *config.Config
	}
)

func NewMCPTimeTask(logger log.Logger, cfg *config.Config) MCPTimeTask {
	return &mcpTimeTask{
		logger: logger.WithField("ddd", "task"),
		cfg:    cfg,
	}
}

func (t *mcpTimeTask) Run(ctx context.Context) error {
	// Create MCP server
	s := server.NewMCPServer(
		"MCP Time",
		"0.0.1",
	)
	// Add tool
	tool := mcp.NewTool("current time",
		mcp.WithDescription("Get current time with timezone, Asia/Shanghai is default"),
		mcp.WithString("timezone",
			mcp.Required(),
			mcp.Description("current time timezone"),
		),
	)
	// Add tool handler
	s.AddTool(tool, t.currentTimeHandler)
	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		t.logger.Error("Server error: %v\n", err)
	}

	return nil
}

func (t *mcpTimeTask) currentTimeHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	timezone, ok := request.Params.Arguments["timezone"].(string)
	if !ok {
		return mcp.NewToolResultError("timezone must be a string"), nil
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("parse timezone with error: %v", err)), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf(`current time is %s`, time.Now().In(loc))), nil
}
