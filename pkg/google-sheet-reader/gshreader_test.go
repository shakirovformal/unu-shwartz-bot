package googlesheetreader

import (
	"testing"

	"github.com/shakirovformal/unu_project_api_realizer/config"
	"github.com/shakirovformal/unu_project_api_realizer/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReader(t *testing.T) {
	cfg := config.Load()
	//cfg := config.Get()
	spreadsheetId := cfg.SPREADSHEETID
	testRowPositive := []string{"2", "4", "10"}
	testRowNegative := []string{"679"}

	for _, value := range testRowPositive {
		resp, err := Reader(spreadsheetId, "BOT", value)
		if assert.NoError(t, err) {
			require.Equal(t, "убрир екб", resp.Values[0][0])
		}
	}
	for _, value := range testRowNegative {
		_, err := Reader(spreadsheetId, "BOT", value)
		if !assert.NoError(t, err) {
			require.EqualValues(t, models.ErrorGoogleSheet, err)
		}

	}

}
