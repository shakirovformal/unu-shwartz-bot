package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	api "github.com/shakirovformal/unu_api"
	"github.com/shakirovformal/unu_project_api_realizer/config"
	gsr "github.com/shakirovformal/unu_project_api_realizer/pkg/google-sheet-reader"
	"github.com/shakirovformal/unu_project_api_realizer/pkg/models"
	"github.com/shakirovformal/unu_project_api_realizer/pkg/utils"
)

// Methods:
// get_folders
// create_folder
// del_folder
// move_task
// get_tasks
// get_reports
// approve_report
// reject_report
// get_expenses
// add_task
// task_limit_add
// task_limit_sub
// edit_task
// del_task
// get_tariffs
// get_countries
// task_pause
// task_play
// task_to_top
// add_blacklist
// add_whitelist
// get_blacklist
// delete_user_blacklist

var cfg = config.Load()
var client = api.NewClient(cfg.UNU_URL, cfg.UNU_TOKEN)

type Sender struct {
	c *api.Client
}

func NewSender() *Sender {
	return &Sender{
		c: client,
	}
}

func (s *Sender) CheckBalance(ctx context.Context) (float64, error) {

	resp, err := s.c.Get_balance(ctx)
	if err != nil {
		return 0, err
	}
	return resp.Balance, nil
}
func (s *Sender) GetFolders(ctx context.Context) ([]struct {
	ID   json.Number "json:\"id\""
	Name string      "json:\"name\""
}, error) {
	resp, err := s.c.Get_folders(ctx)
	if err != nil {
		return nil, err
	}
	return resp.Folders, nil
}

func (s *Sender) CreateFolder(ctx context.Context, folder_name string) (int64, error) {
	resp, err := s.c.Create_folder(ctx, folder_name)
	if err != nil {
		return 0, err
	}
	folder_id, err := resp.FolderID.Int64()
	if err != nil {
		return 0, err
	}
	return folder_id, nil
}
func (s *Sender) DeleteFolder(ctx context.Context, folder_id int) error {
	_, err := s.c.Del_folder(ctx, folder_id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Sender) AddTask(ctx context.Context, userID int, rowWork string) (*models.AddedTask, error) {
	//"""Обработка в функции идёт только 1 строки"""

	//Пойти в таблицу и получить строку
	cfg := config.Load()
	resp, err := gsr.Reader(cfg.SPREADSHEETID, cfg.SHEETLIST, rowWork)
	if err != nil {
		slog.Error("Ошибка получения данных из таблицы")
	}
	if len(fmt.Sprint(resp.Values[0][3])) > 2300 {
		return nil, models.LongMessage
	}

	// projectName := fmt.Sprint(resp.Values[0][0])
	// datepublic := fmt.Sprint(resp.Values[0][5])
	slog.Info("Ошибки до этого момента не произошло")
	//Получить имя для задачи
	name, err := utils.GetName(resp)
	if err != nil {
		slog.Error("Ошибка получения имени для задачи", "ERROR", err)
		return nil, err
	}

	//Получить описание задания
	descr, err := utils.GetDescription(resp)
	if err != nil {
		slog.Error("Ошибка получения описания для задачи", "ERROR", err)
		return nil, err
	}
	//Получить ссылку для задания
	link := fmt.Sprint(resp.Values[0][1])
	//Получить данные: что нужно для выполнения задания
	needReport, err := gsr.ReaderFromCell(cfg.SPREADSHEETID, "REFERENCE", "H1")
	if err != nil {
		slog.Error("Ошибка получения данных из таблицы")
	}
	need_for_report := fmt.Sprint(needReport.Values[0][0])
	//понять, какой тариф выбрать
	tarif_id, err := utils.GetTariff(link)
	if err != nil {
		slog.Error("Ошибка получения id тарифа для задания")
	}
	//получить стоимость задания
	priceInt := utils.GetActualPriceFromTariffID(tarif_id)
	price := float64(priceInt)
	//понять в какую папку сохранить задание

	matcher := utils.NewProjectMatcher()
	stringProject := fmt.Sprint(resp.Values[0][0])
	if err := matcher.LoadConfig(); err != nil {
		fmt.Println(err)
	}

	folder_idString, err := matcher.FindFolderForProject(stringProject)
	if err != nil {
		slog.Error("Ошибка в поиске подходящей папке для сохранения", "ERROR", err)
	}
	folder_id, err := strconv.Atoi(folder_idString)
	if err != nil {
		slog.Error("Ошибка в конвертации id папки из строки в число", "ERROR", err)
	}
	//Необходимость скриншота (по умолчанию всегда True)
	var need_screen bool = true
	//Время на выполнение 72 часа
	time_for_work := 72
	//Время на проверку 120 часов
	time_for_check := 120
	//Получить значение гендерного пола для задания
	stringGender := fmt.Sprint(resp.Values[0][2])
	targeting_gender := utils.CheckGenderInt(stringGender)
	//Выбрать страну: Россия для задания
	targeting_geo_country_id := 1

	slog.Info("Host for connect", "INFO", cfg.DB_HOST)

	// db := database.NewDB(cfg.DB_HOST, cfg.DB_PASSWORD, cfg.DB_DB)
	// rdb := db.Connect(db)
	// rowObj := models.NewRowObject(userID, projectName, link, targeting_gender, descr, datepublic)
	// db.AddRow(ctx, rdb, rowWork, rowObj)

	respApi, err := s.c.Add_task(ctx, name, descr, link, need_for_report, price, tarif_id, folder_id,
		need_screen, false, time_for_work, time_for_check, 0, 0, 0, 0, 0, 0, "", "", 0, 0, targeting_gender, 0, 0, targeting_geo_country_id, 0, 0, 0, "")
	if err != nil {
		return nil, err
	}
	fmt.Println(respApi)
	task_id, err := respApi.TaskID.Int64()
	if err != nil {
		return nil, models.ErrorUnmarshallJSON
	}
	// db.DelRow(ctx, rdb, rowWork)
	return &models.AddedTask{
		Name:   name,
		Price:  price,
		TaskId: task_id,
	}, nil
}

func (s *Sender) DelTask(ctx context.Context, task_id int) error {
	resp, err := s.c.Del_task(ctx, task_id)
	if err != nil {
		return err
	}
	if !resp.Success {
		return errors.New("error with deleting task")
	}

	return nil
}
