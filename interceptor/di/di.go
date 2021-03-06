/**
 * @ Author: Daniel Tan
 * @ Date: 2020-04-05 01:51:55
 * @ LastEditTime: 2020-08-17 22:10:37
 * @ LastEditors: Daniel Tan
 * @ Description:
 * @ FilePath: /trinitygo/interceptor/di/di.go
 * @
 */

package di

import (
	"context"
	"reflect"
	"strings"

	"github.com/PolarPanda611/trinitygo/application"

	"google.golang.org/grpc"
)

// New new DI interceptor
func New(app application.Application) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		runtimeKeyMap := application.DecodeGRPCRuntimeKey(ctx, app.RuntimeKeys())
		tContext := app.ContextPool().Acquire(app, runtimeKeyMap, nil)
		method := strings.Split(info.FullMethod, "/") // /user.UserService/GetUserByID
		defer func() {
			if tContext.AutoFree() {
				//release trinity go context obj
				app.ContextPool().Release(tContext)
			}
		}()

		controller, _, toFreeInstance := app.ControllerPool().GetController(method[1], tContext, app, nil)
		defer func() {
			for _, v := range toFreeInstance {
				app.InstancePool().Release(v)
			}
		}()
		currentMethod, ok := reflect.TypeOf(controller).MethodByName(method[2])
		if !ok {
			panic("controller has no method ")
		}
		var inParam []reflect.Value
		inParam = append(inParam, reflect.ValueOf(controller))
		inParam = append(inParam, reflect.ValueOf(ctx))
		inParam = append(inParam, reflect.ValueOf(req))
		res := currentMethod.Func.Call(inParam)
		if len(res) != 2 {
			panic("wrong res type")
		}
		if res[1].Interface() != nil {
			tContext.SafeRollback()
			return nil, res[1].Interface().(error)
		}
		if app.Conf().GetAtomicRequest() {
			tContext.SafeCommit()
		}
		return res[0].Interface(), nil
	}
}
