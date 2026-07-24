package initialize

import (
	"context"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/limitedActivity/model"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func Gorm(ctx context.Context) {
	err := global.GVA_DB.WithContext(ctx).AutoMigrate(
		new(model.Activity),
		new(model.ActivityOrder),
		new(model.ActivityOrderVerifyRecord),
	)
	if err != nil {
		err = errors.Wrap(err, "限时活动插件表迁移失败")
		zap.L().Error(fmt.Sprintf("%+v", err))
	}
}
