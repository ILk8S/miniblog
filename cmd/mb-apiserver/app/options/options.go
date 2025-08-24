package options

import (
	"fmt"
	"time"
	"github.com/wshadm/miniblog/pkg/options"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/pflag"
	"github.com/wshadm/miniblog/internal/apiserver"
	"k8s.io/apimachinery/pkg/util/sets"
)

// 定义支持的服务器模式集合
var availableServerModes = sets.New(
	apiserver.GinServerMode,
	apiserver.GRPCServerMode,
	apiserver.GRPCGatewayServerMode,
)

// ServerOptions 包含服务器配置选项
type ServerOptions struct {
	//ServerMode 定义服务器模式：gRPC、gin、HTTP、HTTP Reverse Proxy。
	ServerMode string `json:"server-mode" mapstructure:"server-mode"`
	//JWTKey 定义JWT秘钥
	JWTKey string `json:"jwt-key" mapstructure:"jwt-key" validate:"required, min=6"`
	//Expiration定义JWT Token的过期时间
	Expiration time.Duration `json:"expiration" mapstructure:"expiration"`
	// GRPCOptions 包含 gRPC 配置选项.
	GRPCOptions *options.GRPCOptions `json:"grpc" mapstructure:"grpc"`
}

// NewServerOptions 创建带有默认值的ServerOptions 实例
func NewServerOptions() *ServerOptions {
	opts := &ServerOptions{
		ServerMode: apiserver.GRPCGatewayServerMode,
		JWTKey: "Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5",
		Expiration: 2 * time.Hour,
		GRPCOptions: options.NewGRPCOptions(),
	}
	opts.GRPCOptions.Addr = ":6666"
	return opts
}

// AddFlags 将 ServerOptions 的选项绑定到命令行标志.
// 通过使用 pflag 包，可以实现从命令行中解析这些选项的功能.
func (o *ServerOptions) AddFlags(fs *pflag.FlagSet) {
	// --server-mode是参数，默认值是o.ServerMode，最后是帮助说明描述。解析后的值会写入 &o.ServerMode
	fs.StringVar(&o.ServerMode, "server-mode", o.ServerMode, fmt.Sprintf("Server mode, available options: %v", availableServerModes.UnsortedList()))
	fs.StringVar(&o.JWTKey, "jwt-key", o.JWTKey, "JWT signing key. Must be at least 6 characters long.")
	// 绑定 JWT Token 的过期时间选项到命令行标志。
	// 参数名称为 `--expiration`，默认值为 o.Expiration
	fs.DurationVar(&o.Expiration, "expiration", o.Expiration, "The expiration duration of JWT tokens.")
	o.GRPCOptions.AddFlags(fs)
}

func (o *ServerOptions) Validate(obj any) error {
	v := validator.New()
	return v.Struct(obj)

	// 如果是 gRPC 或 gRPC-Gateway 模式，校验 gRPC 配置
	// if stringsutil.StringIn(o.ServerMode, []string{apiserver.GRPCServerMode, apiserver.GRPCGatewayServerMode}) {
	// 	errs = append(errs, o.GRPCOptions.Validate()...)
	// }
}

//Config基于ServerOption构建apiserver.Config
func (o *ServerOptions)  Config() (*apiserver.Config, error) {
	return &apiserver.Config{
		ServerMode: o.ServerMode,
		JWTKey: o.JWTKey,
		Expiration: o.Expiration,
		GRPCOptions: o.GRPCOptions,
	}, nil
}
