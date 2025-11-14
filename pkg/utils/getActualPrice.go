package utils

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GetActualPriceFromTariffID(tariffID int) int {
	tariffIDString := strconv.Itoa(tariffID)
	priceNotRound := GetTariffPrice(tariffIDString)
	price := normalizePrice(priceNotRound)
	return price
}

func normalizePrice(price int) int {

	if price%5 == 0 {
		return price + 5 // Всегда прибавляем 5, если уже кратно 5
	}
	return price + (5 - price%5)
}

func fetchAvgTariffPrice(tariffID string) (string, error) {
	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Подготавливаем данные формы
	formData := url.Values{}
	formData.Set("tarif_id", tariffID)

	// Создаем запрос
	req, err := http.NewRequestWithContext(ctx, "POST",
		"https://unu.im/include/ajax_avg_tarif_price.php",
		strings.NewReader(formData.Encode()))
	if err != nil {
		return "", fmt.Errorf("ошибка создания запроса: %v", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:144.0) Gecko/20100101 Firefox/144.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Priority", "u=0")

	// Устанавливаем referrer
	req.Header.Set("Referer", fmt.Sprintf("https://unu.im/tasks/add?tarif_id=%s", tariffID))

	// Создаем HTTP клиент
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Выполняем запрос
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	// Проверяем статус код
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("неверный статус код: %d, тело ответа: %s", resp.StatusCode, string(body))
	}

	pattern := `Средняя стоимость по системе\s*&mdash;\s*<strong>(\d+\.\d+)</strong>`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(string(body))
	if len(matches) < 2 {
		return "", fmt.Errorf("не удалось найти среднюю стоимость в ответе")
	}

	// Преобразуем строку в float
	priceStr := matches[1]
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return "", fmt.Errorf("ошибка преобразования цены: %v", err)
	}

	return fmt.Sprint(int(price)), nil
}

// Ваша основная функция
func GetTariffPrice(tariffID string) int {
	result, err := fetchAvgTariffPrice(tariffID)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return 0
	}

	fmt.Printf("Результат: %s\n", result)
	tarifID, err := strconv.Atoi(result)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return 0
	}
	return tarifID
}
