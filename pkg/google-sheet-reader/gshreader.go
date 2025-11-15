package googlesheetreader

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/shakirovformal/unu_project_api_realizer/pkg/models"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func googleServiceConstructor(ctx context.Context) *sheets.Service {

	svc, err := sheets.NewService(ctx, option.WithCredentialsFile("config/creds.json"))
	if err != nil {
		slog.Error("Err is:", "ERROR", err)
	}
	return svc
}

// Используем значения у результата resp.Values[0][index]:
// 0  - название проекта
// 1  - ссылка
// 2 - гендерный пол
// 3 - текст отзыва
// 5 - дата публикации
func Reader(spreadsheetId, spreadsheetName string, rowNumber string) (*sheets.ValueRange, error) {
	readRange := fmt.Sprintf("%s!A%s:F%s", spreadsheetName, rowNumber, rowNumber)

	ctx := context.Background()
	svc := googleServiceConstructor(ctx)

	resp, err := svc.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
		return nil, models.ErrorGoogleSheet
	}



	return resp, nil
}
func ReaderFromCell(spreadsheetId, spreadsheetName string, cell string) (*sheets.ValueRange, error) {

	readRange := fmt.Sprintf("%s!%s:%s", spreadsheetName, cell, cell)

	ctx := context.Background()
	svc := googleServiceConstructor(ctx)
	resp, err := svc.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
		return nil, models.ErrorGoogleSheet
	}

	return resp, nil
}
