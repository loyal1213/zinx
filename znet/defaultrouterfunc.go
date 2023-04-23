package znet

import (
	"bytes"
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"path"
	"runtime"
	"strings"
	"time"
)

const (
	//开始追踪堆栈信息的层数
	StackBegin = 3
	//追踪到最后的层数
	StackEnd = 5
)

//用来存放一些RouterSlicesMode下的路由可用的默认中间件

// RouterRecovery 如果使用NewDefaultRouterSlicesServer方法初始化的获得的server将自带这个函数
// 作用是接收业务执行上产生的panic并且尝试记录现场信息
func RouterRecovery(request ziface.IRequest) {
	defer func() {
		if err := recover(); err != nil {
			panicInfo := getInfo(StackBegin)
			//记录错误
			zlog.Ins().ErrorF("MsgId:%d Handler panic: info:%s err:%v", request.GetMsgID(), panicInfo, err)

			//fmt.Printf("MsgId:%d Handler panic: info:%s err:%v", request.GetMsgID(), panicInfo, err)

			//应该回传一个错误的
			//request.GetConnection().SendMsg()
		}

	}()
	request.RouterSlicesNext()
}

// RouterTime 简单累计所有路由组的耗时，不启用
func RouterTime(request ziface.IRequest) {
	now := time.Now()
	request.RouterSlicesNext()
	duration := time.Since(now)
	fmt.Println(duration.String())
}

func getInfo(ship int) (infoStr string) {

	panicInfo := new(bytes.Buffer)
	//也可以不指定终点层数即i := ship;; i++ 通过if！ok 结束循环，但是会一直追到最底层报错信息
	for i := ship; i <= StackEnd; i++ {
		pc, file, lineNo, ok := runtime.Caller(i)
		if !ok {
			break
		}
		funcname := runtime.FuncForPC(pc).Name()
		filename := path.Base(file)
		funcname = strings.Split(funcname, ".")[1]
		fmt.Fprintf(panicInfo, "funcname:%s filename:%s LineNo:%d\n", funcname, filename, lineNo)
	}
	return panicInfo.String()

}
