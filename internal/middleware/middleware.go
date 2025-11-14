package middleware

import (
	"context"
	"reflect"
	"runtime"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/shakirovformal/unu_project_api_realizer/pkg/logger"
)

func MW(next bot.HandlerFunc) bot.HandlerFunc {

	l := logger.NewLogger()
	return func(ctx context.Context, bot *bot.Bot, update *models.Update) {

		l.Printf("timestamp=%s, func=%v", time.Now().Format("2006-01-02 15:04:05"), runtime.FuncForPC(reflect.ValueOf(next).Pointer()).Name())
		
	}
}
