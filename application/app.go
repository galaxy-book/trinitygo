/**
 * @ Author: Daniel Tan
 * @ Date: 2020-04-07 09:41:11
 * @ LastEditTime: 2020-07-29 19:03:26
 * @ LastEditors: Daniel Tan
 * @ Description:
 * @ FilePath: /trinitygo/application/app.go
 * @
 */
package application

import (
	"context"

	"github.com/PolarPanda611/trinitygo/keyword"
	truntime "github.com/PolarPanda611/trinitygo/runtime"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"

	"github.com/PolarPanda611/trinitygo/conf"

	"github.com/jinzhu/gorm"
	"github.com/kataras/golog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Application global app interface
type Application interface {
	IsLogSelfCheck() bool
	Logger() *golog.Logger
	RuntimeKeys() []truntime.RuntimeKey
	Conf() conf.Conf
	Keyword() keyword.Keyword
	ContextPool() *ContextPool
	DB() *gorm.DB
	Enforcer() *casbin.Enforcer
	InstallDB(f func() *gorm.DB)
	ControllerPool() *ControllerPool
	InstancePool() *InstancePool
	UseInterceptor(interceptor ...grpc.UnaryServerInterceptor) Application
	UseMiddleware(middleware ...gin.HandlerFunc) Application
	RegRuntimeKey(runtime ...truntime.RuntimeKey) Application
	InitGRPC()
	InitHTTP()
	InitRouter()
	GetGRPCServer() *grpc.Server
	ServeGRPC()
	ServeHTTP()
	ResponseFactory() func(status int, res interface{}, runtime map[string]string) interface{}
}

// DecodeGRPCRuntimeKey  decode runtime key from ctx
func DecodeGRPCRuntimeKey(ctx context.Context, runtimeKey []truntime.RuntimeKey) map[string]string {
	runtimeKeyMap := make(map[string]string)
	if ctx != nil {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			for _, v := range runtimeKey {
				runtimeKeyMap[v.GetKeyName()] = md[v.GetKeyName()][0]
			}
		}
	}
	return runtimeKeyMap
}

// DecodeHTTPRuntimeKey decode http runtime
func DecodeHTTPRuntimeKey(c *gin.Context, runtimeKey []truntime.RuntimeKey) map[string]string {
	runtimeKeyMap := make(map[string]string)
	if c != nil {
		for _, v := range runtimeKey {
			runtimeKeyMap[v.GetKeyName()] = c.GetString(v.GetKeyName())
		}
	}
	return runtimeKeyMap
}
