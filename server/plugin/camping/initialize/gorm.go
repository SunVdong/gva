package initialize

import (
	"context"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/camping/model"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func Gorm(ctx context.Context) {
	err := global.GVA_DB.WithContext(ctx).AutoMigrate(
		new(model.Venue),
		new(model.VenueOpenTime),
		new(model.VenueTimeslot),
		new(model.VenueCalendar),
		new(model.VenueReservation),
	)
	if err != nil {
		err = errors.Wrap(err, "露营插件表迁移失败")
		zap.L().Error(fmt.Sprintf("%+v", err))
	}
}
