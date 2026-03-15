package initialize

import (
	"context"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/ticket/model"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func Gorm(ctx context.Context) {
	err := global.GVA_DB.WithContext(ctx).AutoMigrate(
		new(model.ScenicSpot),
		new(model.ScenicOpenTime),
		new(model.TicketProduct),
		new(model.TicketSku),
		new(model.TicketAudience),
		new(model.TicketRule),
		new(model.TicketCalendar),
		new(model.TicketOrder),
		new(model.OrderItem),
	)
	if err != nil {
		err = errors.Wrap(err, "门票插件表迁移失败")
		zap.L().Error(fmt.Sprintf("%+v", err))
	}
}
