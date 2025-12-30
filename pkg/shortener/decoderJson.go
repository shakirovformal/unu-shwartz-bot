package shortener

import (
	"encoding/json"
	"fmt"
	"log"
)

type Response struct {
	Data Data `json:"data"`
}

type Data struct {
	URL             string      `json:"url"`
	Hash            string      `json:"hash"`
	Name            *string     `json:"name"`
	ShortURL        string      `json:"short_url"`
	Visits          *int        `json:"visits"`
	StartAt         *int64      `json:"start_at"`
	EndAt           *int64      `json:"end_at"`
	CreatedAt       int64       `json:"created_at"`
	Active          bool        `json:"active"`
	AllowCountries  interface{} `json:"allow_countries"`
	TrackingPixelID interface{} `json:"tracking_pixel_id"`
	DomainID        interface{} `json:"domain_id"`
	UTM             []UTMResp   `json:"utm"`
	MaxVisits       *int        `json:"max_visits"`
	UserID          int         `json:"user_id"`
}

type UTMResp struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func parse(jsonData string) string {
	var resp Response

	if err := json.Unmarshal([]byte(jsonData), &resp); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Short URL:", resp.Data.ShortURL)
	return resp.Data.ShortURL
}
