package app

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wshadm/miniblog/cmd/mb-apiserver/app/options"
)

var configFile string  //配置文件路径

func NewMiniBlogCommand() *cobra.Command {
	ops := options.NewServerOptions()
	cmd := &cobra.Command {
		//指定命令的名族，这个名字会出现在帮助信息中
		Use: "mb-apiserver",
		Short: "A mini blog show best practices for develop a full-featured Go project",
		// 命令的详细描述
		Long: "A mini blog show best practices for develop a full-featured Go project.",
		// 命令出错时，不打印帮助信息。设置为 true 可以确保命令出错时一眼就能看到错误信息
		SilenceUsage: true,
		// 指定调用 cmd.Execute() 时，执行的 Run 函数
		RunE: func(cmd *cobra.Command, args []string) error {
			
			
			return nil
		},
		// 设置命令运行时的参数检查，不需要指定命令行参数。例如：./miniblog param1 param2
		Args: cobra.NoArgs,
	}
	// 初始化配置函数，在每个命令运行时调用
	cobra.OnInitialize(onInitialize)
	// cobra 支持持久性标志(PersistentFlag)，该标志可用于它所分配的命令以及该命令下的每个子命令
	// 推荐使用配置文件来配置应用，便于管理配置项
	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", filePath(), "配置文件路径")
	// 将 ServerOptions 中的选项绑定到命令标志
	ops.AddFlags(cmd.PersistentFlags())
	return cmd
}

func run(opts *options.ServerOptions) error {
	if err := viper.Unmarshal(&opts); err != nil {
		return  err
	}
	//对命令行选项值进行校验
	if err := opts.Validate(&opts); err != nil {
		return  err
	}
	cfg, err := opts.Config()
	if err != nil {
		return err
	}
	// 创建服务器实例.
	// 注意这里是联合服务器，因为可能同时启动多个不同类型的服务器.
	server, err := cfg.NewUnionServer()
	if err != nil {
		return err
	}
	return server.Run()
}