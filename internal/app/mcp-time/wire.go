// Copyright 2025 fsyyft-ai
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

//go:build wireinject
// +build wireinject

package mcp_time

import (
	"github.com/google/wire"

	"github.com/fsyyft-ai/mcp-time/internal/config"
	"github.com/fsyyft-ai/mcp-time/internal/task"
)

func wireServer(cfg *config.Config) (task.MCPTimeTask, func(), error) {
	// wire.Build 函数用于声明依赖关系图，将所有组件连接在一起。
	// panic 调用会在编译时被 wire 工具替换为实际的依赖注入代码。
	panic(wire.Build(
		ProviderSet,
		task.ProviderSet,
	))
}
