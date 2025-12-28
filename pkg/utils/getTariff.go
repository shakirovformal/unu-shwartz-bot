package utils

import (
	"regexp"

	"github.com/shakirovformal/unu_project_api_realizer/pkg/models"
)

// SitePattern хранит шаблон URL и соответствующую ячейку
type SitePatternForTariff struct {
	PatternTariff *regexp.Regexp
	TariffID      int
}

// SiteMatcher содержит все паттерны для сопоставления
type SiteMatcherForTariff struct {
	patterns []SitePatternForTariff
}

// NewSiteMatcher создает и инициализирует SiteMatcher с предопределенными паттернами
func NewSiteMatcherForTariff() *SiteMatcherForTariff {
	return &SiteMatcherForTariff{
		patterns: []SitePatternForTariff{
			{
				PatternTariff: regexp.MustCompile(`maps\.app\.goo\.gl`),
				TariffID:      37,
			},
			{
				PatternTariff: regexp.MustCompile(`yandex\.(ru|com)/maps`),
				TariffID:      36,
			},
			{
				PatternTariff: regexp.MustCompile(`otzovik\.com`),
				TariffID:      40,
			},
			{
				PatternTariff: regexp.MustCompile(`irecommend\.ru`),
				TariffID:      40,
			},
			{
				PatternTariff: regexp.MustCompile(`2gis\.ru`),
				TariffID:      38,
			},
		},
	}
}

// GetCellForURL возвращает ячейку для данного URL
func (sm *SiteMatcherForTariff) GetCellForURLGetTariff(url string) (int, error) {
	for _, pattern := range sm.patterns {
		if pattern.PatternTariff.MatchString(url) {
			return pattern.TariffID, nil
		}
	}
	return 45, models.ErrorMatchingSite // или какое-то значение по умолчанию
}

func GetTariff(link string) (int, error) {
	patternSlice := NewSiteMatcherForTariff()
	siteCell, err := patternSlice.GetCellForURLGetTariff(link)
	if err != nil {
		return 45, nil
	}

	return siteCell, nil
}
