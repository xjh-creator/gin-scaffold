package component_name

import (
	"gin-scaffold/internal/pkg/logger"
	"gin-scaffold/internal/pkg/setting"
	"github.com/jinzhu/gorm"
)

var (
	ServerSetting   *setting.ServerSettingS
	AppSetting      *setting.AppSettingS
	EmailSetting    *setting.EmailSettingS
	JWTSetting      *setting.JWTSettingS
	DatabaseSetting *setting.DatabaseSettingS
	Logger          *logger.Logger
)

var (
	DBEngine *gorm.DB
)
