package main

import (
	"PromocodesUpload/internal/config"
	"PromocodesUpload/internal/logger"
	"PromocodesUpload/internal/services"
	"fmt"
	"github.com/orandin/slog-gorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	_, b, _, _ = runtime.Caller(0)
	rootPath   = filepath.Dir(filepath.Dir(filepath.Dir(b)))
)

func main() {
	// init Config
	cgf := config.MustLoad(rootPath)

	// init Logger
	log := logger.NewFileSLogger(
		fmt.Sprintf("%s/logs", rootPath),
		cgf.LogConfig.DefaultLogFile,
		cgf.Env,
	)
	log.SetUpLogger()

	// init connection
	connection, err := gorm.Open(mysql.Open(cgf.GormDns()), &gorm.Config{
		Logger: slogGorm.New(
			slogGorm.WithLogger(log.Logger),
			//slogGorm.SetLogLevel(slogGorm.SlowQueryLogType, log.GetLevel()),
		),
	})
	if err != nil {
		log.Logger.Error(err.Error())
		os.Exit(0)
	}
	// main
	for {
		importPromoCodesService := services.NewImportPromoCodes(
			fmt.Sprintf("%s/%s", rootPath, cgf.SourceDir), cgf.PackPromo, connection, log.Logger)
		importPromoCodesService.Execute()
		time.Sleep(5 * time.Second)
	}
}
