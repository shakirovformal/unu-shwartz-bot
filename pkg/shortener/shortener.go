package shortener

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

type UTM struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Request struct {
	URL    string `json:"url"`
	Active bool   `json:"active"`
	UTM    []UTM  `json:"utm"`
}

func ShortLink(long_link string) (string, error) {
	api_shortener_token := "Mf8rCmtKzPbNXHxlfvNOsNgKdGOszz22xuDHor3bKFJbF58Fx4v7QrjrdDRf"
	reqBody := Request{
		URL:    long_link,
		Active: true,
		UTM: []UTM{
			{
				Key:   "utm_source",
				Value: "search",
			},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://ishort.su/api/link",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+api_shortener_token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	respJson := string(respBody)
	link := parse(respJson)
	slog.Info("Ссылка которую я получил", "LINK", link)
	return link, nil
}
