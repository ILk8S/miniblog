package main

import (
	"fmt"
	"os"

	// 导入 automaxprocs 包，可以在程序启动时自动设置 GOMAXPROCS 配置，
	// 使其与 Linux 容器的 CPU 配额相匹配。
	// 这避免了在容器中运行时，因默认 GOMAXPROCS 值不合适导致的性能问题，
	// 确保 Go 程序能够充分利用可用的 CPU 资源，避免 CPU 浪费。
	"github.com/wshadm/miniblog/cmd/mb-apiserver/app"
	_ "go.uber.org/automaxprocs"
)

func main() {
	command := app.NewMiniBlogCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
	fmt.Println("Hello world")
}
