package initialize

import (
	"fmt"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	laService "github.com/flipped-aurora/gin-vue-admin/server/plugin/limitedActivity/service"
	ticketService "github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/service"
	"github.com/flipped-aurora/gin-vue-admin/server/task"

	"github.com/robfig/cron/v3"
)

func Timer() {
	go func() {
		var option []cron.Option
		option = append(option, cron.WithSeconds())
		// 清理DB定时任务
		_, err := global.GVA_Timer.AddTaskByFunc("ClearDB", "@daily", func() {
			err := task.ClearTable(global.GVA_DB) // 定时任务方法定在task文件包中
			if err != nil {
				fmt.Println("timer error:", err)
			}
		}, "定时清理数据库【日志，黑名单】内容", option...)
		if err != nil {
			fmt.Println("add timer error:", err)
		}

		_, err = global.GVA_Timer.AddTaskByFuncWithSecond("TicketOrderTimeoutClose", "0 */1 * * * *", func() {
			closed, closeErr := ticketService.Service.Order.CloseTimeoutUnpaidOrders(15*time.Minute, 200)
			if closeErr != nil {
				fmt.Println("ticket timeout close error:", closeErr)
				return
			}
			if closed > 0 {
				fmt.Printf("ticket timeout close success, closed=%d\n", closed)
			}
		}, "门票订单超时未支付自动关闭")
		if err != nil {
			fmt.Println("add timer error:", err)
		}

		_, err = global.GVA_Timer.AddTaskByFuncWithSecond("LimitedActivityOrderTimeoutClose", "0 */1 * * * *", func() {
			closed, closeErr := laService.Service.Order.CloseTimeoutUnpaidOrders(15*time.Minute, 200)
			if closeErr != nil {
				fmt.Println("limitedActivity timeout close error:", closeErr)
				return
			}
			if closed > 0 {
				fmt.Printf("limitedActivity timeout close success, closed=%d\n", closed)
			}
		}, "限时活动订单超时未支付自动关闭")
		if err != nil {
			fmt.Println("add timer error:", err)
		}

		// 其他定时任务定在这里 参考上方使用方法

		//_, err := global.GVA_Timer.AddTaskByFunc("定时任务标识", "corn表达式", func() {
		//	具体执行内容...
		//  ......
		//}, option...)
		//if err != nil {
		//	fmt.Println("add timer error:", err)
		//}
	}()
}
