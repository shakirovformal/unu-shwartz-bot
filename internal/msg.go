package internal

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/shakirovformal/unu_project_api_realizer/config"
	"github.com/shakirovformal/unu_project_api_realizer/internal/pkg"
)

var cfg = config.Load()
var req = pkg.NewSender()

type UserState struct {
	State     string
	Data      map[string]interface{}
	CreatedAt time.Time
	Command   string
}

const (
	STATE_WAIT_FOLDER_NAME  = "wait_folder_name"
	STATE_WAIT_INPUT_ROWS   = "wait_input_rows"
	STATE_WAIT_FOLDER_ID    = "wait_folder_id"
	STATE_IDLE              = "idle"
	STATE_WAIT_TASK_NUMBERS = "wait_task_numbers"
)

var userStates = make(map[int64]*UserState)
var stateMutex sync.RWMutex

func setState(chatID int64, state *UserState) {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	state.CreatedAt = time.Now()
	userStates[chatID] = state
}

func getState(chatID int64) (*UserState, bool) {
	stateMutex.RLock()
	defer stateMutex.RUnlock()
	state, exists := userStates[chatID]
	return state, exists
}

func clearState(chatID int64) {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	delete(userStates, chatID)
}

func welcomeMessage(ctx context.Context, b *bot.Bot, update *models.Update) {

	slog.Info(fmt.Sprintf("User '%s' wrote '%s' for will start work", update.Message.Chat.Username, update.Message.Text))
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Привет!\nЧтобы посмотреть список доступных команд, введи /help",
	})
}

func helpMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	slog.Info(fmt.Sprintf("User '%s' wrote '%s' for get help information", update.Message.Chat.Username, update.Message.Text))
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: `Список доступных команд:
--------------------
/help - помощь по командам
/docs - полная развертка по командам
-----------------------------
Команды связанные с балансом:
/balance - посмотреть баланс
-----------------------------
Команды связанные с папками:
/create_folder - создать папку с названием
/delete_folder - удалить папку
-----------------------------
Команды связанные с задачами:
/create_task - создать задачу
/delete_task - удалить задачу или задачи
-----------------------------
Сервисные команды. Обычно использует только разработчик, но если тебе мой уважаемый читатель интересно, то потыкай, здесь ты точно ничего не сломаешь =)
/get_folders - посмотреть существующие папки(Используется для разработки и модификации нашего бота)
-----------------------------
Остальные команды в разработке 🙂
Связаться с разработчиком: @tatarkazawarka`,
	})
}
func docsMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	slog.Info(fmt.Sprintf("User '%s' wrote '%s' for get help information", update.Message.Chat.Username, update.Message.Text))
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: `Полная развертка по командам:
-----------------------------
Команды связанные с балансом:
/balance - Позволяет посмотреть баланс в кошельке. ВНИМАНИЕ: отображается баланс без учета замороженных средств. Поэтому при желании включить много задач. Просьба проверить наличие доступных средств к оплате.
-----------------------------
Команды связанные с папками:
/create_folder - Команда даёт возможность создать папку с нужным вам названием, это гораздо удобнее, если вы работаете в команде и для вашего коллеги по работе нужно создать папку удалённо(например с телефона)
/delete_folder - Команда удаляет папку. Но для того чтобы удалить папку, нужно передать значение номера папки. Скоро появится команда которая выдаёт список папок с их номерами, поэтому сначала нужно будет узнать номер папки
-----------------------------
Команды связанные с задачами:
/create_task - Создание задачи из гугл таблицы. ВАЖНО: список строк для обработки передаётся через знак "-". Пример: 2-15. Если ввести некорректное значение, то ничего не получится и скорее всего лучше написать разработчику
/delete_task - Удаление задач или списка задач. Список передаётся в одну строку разделенным пробелом. Пример: 12345 54321 15243 51423
-----------------------------
Сервисные команды.
/get_folders - Отдаёт список папок которые созданы (Находится в этом блоке, потому что команда находится в тестовом режиме и не всегда работает корректно)
-----------------------------
Остальные команды в разработке 🙂
Связаться с разработчиком: @tatarkazawarka`,
	})
}

func checkBalance(ctx context.Context, b *bot.Bot, update *models.Update) {
	slog.Info(fmt.Sprintf("User '%s' wrote '%s' for check balance wallet", update.Message.Chat.Username, update.Message.Text))

	ctxWT, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	balance, err := req.CheckBalance(ctxWT)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("Ошибка получения баланса. error:%s", err),
		})
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("Баланс вашего кошелька: %.2f", balance),
	})
}

func getFoldersId(ctx context.Context, b *bot.Bot, update *models.Update) {
	slog.Info(fmt.Sprintf("User '%s' wrote '%v' for get folder list id", update.Message.Chat.Username, update.Message.Text))
	ctxWT, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	folder_list, err := req.GetFolders(ctxWT)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("Ошибка получения папок. error:%s", err),
		})
	}
	result_text := "Список папок:"
	for _, value := range folder_list {
		result_text += fmt.Sprintf("\nID: %v. Название: %s", value.ID, value.Name)
	}
	result_text += "\nP.S Помните, что эта информация для вас бесполезна и используется только для разработчика 😉"
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("Ваши папки: %s", result_text),
	})

}

func handleFolderNameInput(ctx context.Context, b *bot.Bot, update *models.Update, state *UserState) {
	chatID := update.Message.Chat.ID
	folderName := strings.TrimSpace(update.Message.Text)

	if len(folderName) == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Имя папки не может быть пустым. Введите имя еще раз:",
		})
		return
	}

	// Показываем что начали обработку
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   fmt.Sprintf("Создаю папку '%s'...", folderName),
	})

	// Создаем папку
	ctxWT, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	folder_id, err := req.CreateFolder(ctxWT, folderName)

	if err != nil {
		slog.Error("Ошибка создания папки:", "ERROR:", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   fmt.Sprintf("❌ Ошибка при создании папки: %v", err),
		})
	} else {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   fmt.Sprintf("✅ Папка '%s' успешно создана!\nID: %d", folderName, folder_id),
		})
	}

	clearState(chatID)
}

func createFolder(ctx context.Context, b *bot.Bot, update *models.Update) {
	slog.Info(fmt.Sprintf("User '%s' wrote '%s' for create folder", update.Message.Chat.Username, update.Message.Text))
	chatID := update.Message.Chat.ID

	setState(chatID, &UserState{
		State:   STATE_WAIT_FOLDER_NAME,
		Data:    make(map[string]interface{}),
		Command: "create_folder",
	})
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Пожалуйста, введите имя для папки которую хотим создать:",
	})

}

func deleteFolder(ctx context.Context, b *bot.Bot, update *models.Update) {
	slog.Info(fmt.Sprintf("User '%s' wrote '%s' for delete folder", update.Message.Chat.Username, update.Message.Text))
	chatID := update.Message.Chat.ID

	setState(chatID, &UserState{
		State:   STATE_WAIT_FOLDER_ID,
		Data:    make(map[string]interface{}),
		Command: "delete_folder",
	})
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Пожалуйста, введите ID для папки которую хотим удалить:",
	})

}

func handleDeleteFolderIdInput(ctx context.Context, b *bot.Bot, update *models.Update, state *UserState) {
	chatID := update.Message.Chat.ID
	folderId := strings.TrimSpace(update.Message.Text)

	if len(folderId) == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "ID папки не может быть пустым. Введите ID еще раз:",
		})
		return
	}

	// Показываем что начали обработку
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   fmt.Sprintf("Удаляю папку '%s'...", folderId),
	})

	folderIdInt, err := strconv.Atoi(folderId)
	ctxWT, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err = req.DeleteFolder(ctxWT, folderIdInt)
	if err != nil {
		slog.Error("Ошибка создания папки:", "ERROR:", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   fmt.Sprintf("❌ Ошибка при удалении папки: %v", err),
		})
	} else {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "✅ Папка успешно удалена!\n",
		})
	}

	clearState(chatID)
}

func createTask(ctx context.Context, b *bot.Bot, update *models.Update) {
	// ctxWT, cancel := context.WithTimeout(ctx, time.Second*30)
	// defer cancel()
	slog.Info(fmt.Sprintf("User '%s' wrote '%s' for create folder", update.Message.Chat.Username, update.Message.Text))
	chatID := update.Message.Chat.ID
	// FIXME: Сделать здесь логику, чтобы при входе в данную функцию, сначала проверялась очередь.
	// Есть ли незавершенные задачи? Если есть, нужно ли обработать их в первую очередь или оставить на потом?
	// db := database.NewDB(cfg.DB_HOST, cfg.DB_PASSWORD, cfg.DB_DB)
	// rdb := db.Connect(db)
	// stringUnfullfilled, err := db.CheckUnfullfilledRows(ctxWT, rdb)
	// if err != nil {
	// 	slog.Error("Простите, произошла какая-то неизвестная ошибка с базой данных, пожалуйста поправьте")
	// }
	// fmt.Println("Незавершенные задачи в базе:", stringUnfullfilled)
	// if len(stringUnfullfilled) > 0 {
	// 	b.SendMessage(ctx, &bot.SendMessageParams{
	// 		ChatID: update.Message.Chat.ID,
	// 		Text:   fmt.Sprintf("Дело в том, что перед тем как создать новые задачи, давайте разберёмся со старыми. Я сходил в базу данных и нашёл строки, которые по каким-то либо причинам не были обработаны. Вот список %v", stringUnfullfilled),
	// 	})
	// }

	// Проверили что задач нет, запрашиваем у клиента номера строк для выполнения
	setState(chatID, &UserState{
		State:   STATE_WAIT_INPUT_ROWS,
		Data:    make(map[string]interface{}),
		Command: "create_task",
	})
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Пожалуйста, введи номера строк для начала работы: Пример: 2-15(Не забывайте, что строка с номером 1, сервисная, на ней находятся названия колонок)",
	})

}
func handleTaskRowInput(ctx context.Context, b *bot.Bot, update *models.Update, state *UserState) {
	listTasks := []int64{}
	chatID := update.Message.Chat.ID
	task_list_message := strings.TrimSpace(update.Message.Text)

	if len(task_list_message) == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Строки не могут быть пустым сообщением...",
		})
		return
	}
	rows := strings.Split(update.Message.Text, "-")
	beginRowString, endRowString := rows[0], rows[1]

	beginRowInt, err := strconv.Atoi(beginRowString)
	if err != nil {
		slog.Error(fmt.Sprintf("Ошибка конвертации значения для начальной строки... Просьба проверить корректность. Пользователь %s попытался ввёл: %s что привело к данной ошибке", update.Message.Chat.Username, update.Message.Text))
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Простите, вы ввели некорректное значение для начала работы. Пожалуйста, ориентируйтесь на пример который я вам показал",
		})
		clearState(chatID)
		return

	}
	endRowInt, err := strconv.Atoi(endRowString)
	if err != nil {
		slog.Error(fmt.Sprintf("Ошибка конвертации значения для начальной строки... Просьба проверить корректность. Пользователь %s попытался ввёл: %s что привело к данной ошибке", update.Message.Chat.Username, update.Message.Text))
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Простите, вы ввели некорректное значение для начала работы. Пожалуйста, ориентируйтесь на пример который я вам показал",
		})
		clearState(chatID)
		return
	}

	// Показываем что начали обработку
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Создаю задачи...",
	})

	ctxWT, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()
	for i := beginRowInt; i < endRowInt; i++ {
		// Создаем task
		taskObject, err := req.AddTask(ctxWT, update.Message.ID, beginRowString)
		if err != nil {
			slog.Error("Ошибка создания задачи:", "ERROR:", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   fmt.Sprintf("❌ Ошибка при создании: %v", err),
			})
		} else {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text: fmt.Sprintf("✅  Задача #%d успешно создана!\n Название: %s\nЦена: %f\n№Созданной задачи: %d\nСсылка на задачу: https://unu.im/tasks/edit/%d",
					i, taskObject.Name, taskObject.Price, taskObject.TaskId, taskObject.TaskId),
			})
		}
		listTasks = append(listTasks, taskObject.TaskId)
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   fmt.Sprintf("Список всех выполненных задач: %d", listTasks),
	})
	clearState(chatID)
}

// Удаление задачи
func deleteTask(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID

	setState(chatID, &UserState{
		State:   STATE_WAIT_TASK_NUMBERS,
		Data:    make(map[string]interface{}),
		Command: "delete_task",
	})
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Пожалуйста, введи номер задачи для удаления (Более полное описание команды можно посмотреть в документации /docs): Пример: 12345 или 12345 54321",
	})
}

func handleDeleteTaskIdInput(ctx context.Context, b *bot.Bot, update *models.Update, state *UserState) {
	listTasks := []int{}
	chatID := update.Message.Chat.ID
	task_list_message := strings.TrimSpace(update.Message.Text)

	if len(task_list_message) == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Строки не могут быть пустым сообщением...",
		})
		return
	}
	tasksStringSlice := strings.Split(update.Message.Text, " ")

	for _, value := range tasksStringSlice {
		task, _ := strconv.Atoi(value)
		listTasks = append(listTasks, task)
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Удаляю задачи...",
	})
	ctxWT, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	for _, value := range listTasks {
		if err := req.DelTask(ctxWT, value); err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   fmt.Sprintf("Задача %d не удалена по причине %v", value, err),
			})
		}
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Задачи успешно удалены",
	})
	clearState(chatID)
}
func handler(ctx context.Context, b *bot.Bot, update *models.Update) {

	chatID := update.Message.Chat.ID
	state, exists := getState(chatID)

	if !exists {
		// Обычное сообщение, не связанное с состоянием
		return
	}

	// Проверяем время жизни состояния (максимум 5 минут)
	if time.Since(state.CreatedAt) > 5*time.Minute {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Время сессии истекло. Начните заново.",
		})
		clearState(chatID)
		return
	}

	switch state.State {
	case STATE_WAIT_FOLDER_NAME:
		handleFolderNameInput(ctx, b, update, state)
	case STATE_WAIT_FOLDER_ID:
		handleDeleteFolderIdInput(ctx, b, update, state)
	case STATE_WAIT_INPUT_ROWS:
		handleTaskRowInput(ctx, b, update, state)
	case STATE_WAIT_TASK_NUMBERS:
		handleDeleteTaskIdInput(ctx, b, update, state)
	default:
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Неизвестное состояние.",
		})
		clearState(chatID)
	}

}
