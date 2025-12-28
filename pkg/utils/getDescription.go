package utils

import (
	"fmt"
	"log/slog"
	"regexp"

	"github.com/shakirovformal/unu_project_api_realizer/config"
	gsr "github.com/shakirovformal/unu_project_api_realizer/pkg/google-sheet-reader"
	"github.com/shakirovformal/unu_project_api_realizer/pkg/models"
	"google.golang.org/api/sheets/v4"
)

func GetDescription(respData *sheets.ValueRange) (string, error) {

	link := fmt.Sprint(respData.Values[0][1])
	// Получаем шаблон сообщения который хотим получить
	ref, err := checkReferenceFromLinkForDescrTask(link)
	if err != nil {
		slog.Error("Ошибка при попытке мэтчинга сайта по ссылке")
		return "", models.ErrorGoogleSheet
	}
	// Описание: Шаблон + текст из таблицы
	descr := fmt.Sprintf("%s\n%s", ref, respData.Values[0][3])

	return descr, nil
}

// SiteMatcher содержит все паттерны для сопоставления
type SiteMatcherForDescr struct {
	patterns []SitePattern
}

// NewSiteMatcher создает и инициализирует SiteMatcher с предопределенными паттернами
func NewSiteMatcherForDescr() *SiteMatcher {
	return &SiteMatcher{
		patterns: []SitePattern{
			{
				Pattern: regexp.MustCompile(`maps\.app\.goo\.gl`),
				Cell:    "A3",
			},
			{
				Pattern: regexp.MustCompile(`yandex\.(ru|com)/maps`),
				Cell:    "B3",
			},
			{
				Pattern: regexp.MustCompile(`otzovik\.com`),
				Cell:    "C3",
			},
			{
				Pattern: regexp.MustCompile(`irecommend\.ru`),
				Cell:    "D3",
			},
			{
				Pattern: regexp.MustCompile(`prodoctorov\.ru`),
				Cell:    "E3",
			},
			{
				Pattern: regexp.MustCompile(`sravni\.ru`),
				Cell:    "F3",
			},
			{
				Pattern: regexp.MustCompile(`2gis\.ru`),
				Cell:    "H3",
			},
		},
	}
}

func checkReferenceFromLinkForDescrTask(link string) (string, error) {
	cfg := config.Load()
	patternSlice := NewSiteMatcherForDescr()
	siteCell, err := patternSlice.GetCellForURLDescrTask(link)
	if err != nil {
		return "", err
	}
	fmt.Println("im get sitecell for work descr", siteCell)
	resp, err := gsr.ReaderFromCell(cfg.SPREADSHEETID, "REFERENCE", siteCell)
	if err != nil {
		slog.Error("ERR", "ERROR", err)
		return "", nil
	}
	textReference := fmt.Sprint(resp.Values[0][0])
	return textReference, nil
}

// GetCellForURL возвращает ячейку для данного URL
func (sm *SiteMatcher) GetCellForURLDescrTask(url string) (string, error) {
	for _, pattern := range sm.patterns {
		if pattern.Pattern.MatchString(url) {
			return pattern.Cell, nil
		}
	}
	return "I3", models.ErrorMatchingSite // или какое-то значение по умолчанию
}
