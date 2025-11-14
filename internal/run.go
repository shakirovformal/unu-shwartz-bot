package internal

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
)

func RunBot() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(cfg.TG_TOKEN, opts...)
	if err != nil {
		panic(err)
	}
	slog.Info("BOT STARTED")
	// delete_folder
	// TODO: move_task
	// TODO: get_tasks
	// TODO: get_reports
	// TODO: approve_report
	// TODO: reject_report
	// TODO: get_expenses
	// add_task
	// TODO: task_limit_add
	// TODO: task_limit_sub
	// TODO: edit_task
	// del_task
	// TODO: get_tariffs
	// TODO: get_countries
	// TODO: task_pause
	// TODO: task_play
	// TODO: task_to_top
	// TODO: add_blacklist
	// TODO: add_whitelist
	// TODO: get_blacklist
	// TODO: delete_user_blacklist

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, welcomeMessage)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, helpMessage)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/docs", bot.MatchTypeExact, docsMessage)

	b.RegisterHandler(bot.HandlerTypeMessageText, "/balance", bot.MatchTypeExact, checkBalance)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/get_folders", bot.MatchTypeExact, getFoldersId)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/create_folder", bot.MatchTypeExact, createFolder)

	b.RegisterHandler(bot.HandlerTypeMessageText, "/delete_folder", bot.MatchTypeExact, deleteFolder)

	b.RegisterHandler(bot.HandlerTypeMessageText, "/create_task", bot.MatchTypeExact, createTask)

	b.RegisterHandler(bot.HandlerTypeMessageText, "/delete_task", bot.MatchTypeExact, deleteTask)

	b.Start(ctx)
}
