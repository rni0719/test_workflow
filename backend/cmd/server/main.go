package main

import (
	"log"
	"os"

	"github.com/rni0719/test_workflow/internal/app"
	"github.com/rni0719/test_workflow/pkg/config"
)

func main() {
	// 設定の読み込み
	cfg := config.LoadConfig()
	
	// アプリケーションの初期化
	app := app.NewApp(cfg.DatabaseURL)
	
	// サーバーの起動
	app.Run(":" + cfg.Port)
}
