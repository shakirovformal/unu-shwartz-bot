package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/stretchr/testify/assert/yaml"
)

type Config map[string][]string

// ProjectMatcher для сопоставления проектов с папками
type ProjectMatcher struct {
	config Config
	cache  map[string]*regexp.Regexp // Кэш скомпилированных регулярных выражений
}

// NewProjectMatcher создает новый матчер
func NewProjectMatcher() *ProjectMatcher {
	return &ProjectMatcher{
		cache: make(map[string]*regexp.Regexp),
	}
}

// LoadConfig загружает конфигурацию из файла
func (pm *ProjectMatcher) LoadConfig() error {
	data, err := os.ReadFile("/home/rinat/Desktop/unu_project/go/config/config.conf")
	if err != nil {
		return fmt.Errorf("ошибка чтения конфига: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("ошибка парсинга JSON: %v", err)
	}

	pm.config = config
	return nil
}

// FindFolderForProject находит папку для проекта
func (pm *ProjectMatcher) FindFolderForProject(projectName string) (string, error) {
	if pm.config == nil {
		return "", fmt.Errorf("конфигурация не загружена")
	}

	// Приводим к нижнему регистру для case-insensitive поиска
	projectName = strings.ToLower(strings.TrimSpace(projectName))

	for folder, patterns := range pm.config {
		for _, pattern := range patterns {
			// Проверяем точное совпадение
			if strings.ToLower(pattern) == projectName {
				return folder, nil
			}

			// Проверяем как регулярное выражение
			if matched, err := pm.matchRegex(pattern, projectName); err != nil {
				return "", err
			} else if matched {
				return folder, nil
			}
		}
	}

	return "", fmt.Errorf("папка для проекта '%s' не найдена", projectName)
}

// matchRegex проверяет совпадение по регулярному выражению
func (pm *ProjectMatcher) matchRegex(pattern, projectName string) (bool, error) {
	// Используем кэш для избежания повторной компиляции
	if re, exists := pm.cache[pattern]; exists {
		return re.MatchString(projectName), nil
	}

	// Компилируем регулярное выражение
	re, err := regexp.Compile("(?i)" + pattern) // (?i) - case insensitive
	if err != nil {
		return false, fmt.Errorf("неверное регулярное выражение '%s': %v", pattern, err)
	}

	// Сохраняем в кэш
	pm.cache[pattern] = re

	return re.MatchString(projectName), nil
}

// MoveProjectToFolder перемещает проект в соответствующую папку
func (pm *ProjectMatcher) MoveProjectToFolder(projectName, sourcePath string) error {
	folder, err := pm.FindFolderForProject(projectName)
	if err != nil {
		return err
	}

	// Создаем папку назначения, если её нет
	destDir := filepath.Join(".", folder)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("ошибка создания папки %s: %v", destDir, err)
	}

	// Определяем исходный и целевой пути
	sourceFile := sourcePath
	destFile := filepath.Join(destDir, filepath.Base(sourcePath))

	// Перемещаем файл
	if err := os.Rename(sourceFile, destFile); err != nil {
		return fmt.Errorf("ошибка перемещения файла: %v", err)
	}

	fmt.Printf("Проект '%s' перемещен в папку '%s'\n", projectName, folder)
	return nil
}

func (pm *ProjectMatcher) LoadConfigYAML(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("ошибка чтения конфига: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("ошибка парсинга YAML: %v", err)
	}

	pm.config = config
	return nil
}
