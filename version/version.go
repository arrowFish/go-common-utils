package version

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/net/ghttp"
	"runtime"
)

var (
	// 这些变量将在编译时通过 -ldflags 注入

	Version   string
	GitCommit string
	BuildTime string
)

type Info struct {
	Version     string `json:"version"`
	GitCommit   string `json:"git_commit"`
	BuildTime   string `json:"build_time"`
	GoVersion   string `json:"go_version"`
	Platform    string `json:"platform"`
	ServiceName string `json:"service_name"`
}

// NewInfo 创建一个新的版本信息实例
func NewInfo(serviceName string) Info {
	return Info{
		Version:     Version,
		GitCommit:   GitCommit,
		BuildTime:   BuildTime,
		GoVersion:   runtime.Version(),
		Platform:    fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		ServiceName: serviceName,
	}
}

// String 返回格式化的版本信息字符串
func (i Info) String() string {
	return fmt.Sprintf(
		"Service: %s\nVersion: %s\nGit Commit: %s\nBuild Time: %s\nGo Version: %s\nPlatform: %s",
		i.ServiceName,
		i.Version,
		i.GitCommit,
		i.BuildTime,
		i.GoVersion,
		i.Platform,
	)
}

// JSON 返回 JSON 格式的版本信息
func (i Info) JSON() string {
	data, _ := json.MarshalIndent(i, "", "  ")
	return string(data)
}

// Print 打印版本信息
func Print(serviceName string) {
	info := NewInfo(serviceName)
	fmt.Println(info.String())
}

// PrintJSON 打印 JSON 格式的版本信息
func PrintJSON(serviceName string) {
	info := NewInfo(serviceName)
	fmt.Println(info.JSON())
}

// RegisterGoFrameHandler 注册 GoFrame 的版本信息处理器
func RegisterGoFrameHandler(server *ghttp.Server, serviceName string) {
	server.BindHandler("/version", func(r *ghttp.Request) {
		info := NewInfo(serviceName)
		r.Response.WriteJson(info)
	})
}

// RegisterGoFrameMiddleware 注册 GoFrame 的版本信息中间件
func RegisterGoFrameMiddleware(group *ghttp.RouterGroup, serviceName string) {
	group.ALL("/version", func(r *ghttp.Request) {
		info := NewInfo(serviceName)
		r.Response.WriteJson(info)
	})
}
