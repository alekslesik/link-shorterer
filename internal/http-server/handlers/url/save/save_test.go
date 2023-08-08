package save_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alekslesik/link-shorterer/internal/http-server/handlers/url/save"
	"github.com/alekslesik/link-shorterer/internal/http-server/handlers/url/save/mocks"
	"github.com/alekslesik/link-shorterer/internal/lib/logger/handlers/slogdiscard"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {
	testCases := []struct {
		name      string // Имя теста
		alias     string // Отправляемый alias
		url       string // Отправляемый URL
		respError string // Какую ошибку мы должны получить?
		mockError error  // Ошибку, которую вернёт мок
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://google.com",
			// Тут поля respError и mockError оставляем пустыми,
			// т.к. это успешный запрос
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			// Создаем объект мока стораджа
			urlSaverMock := mocks.NewURLSaver(t)

			// Если ожидается успешный ответ, значит к моку точно будет вызов
			// Либо даже если в ответе ожидаем ошибку,
			// но мок должен ответить с ошибкой, к нему тоже будет запрос:
			if tC.respError == "" || tC.mockError != nil {
				urlSaverMock.On("SaveURL", tC.url, mock.AnythingOfType("string")).
					Return(int64(1), tC.mockError).
					Once()
			}

			// Создаем наш хэндлер
			handler := save.New(slogdiscard.NewDiscardLogger(), urlSaverMock)

			// Формируем тело запроса
			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tC.url, tC.alias)

			// Создаем объект запроса
			req, err := http.NewRequest(http.MethodGet, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			// Создаем ResponseRecorder для записи ответа хэндлера
			rr := httptest.NewRecorder()

			// Обрабатываем запрос, записывая ответ в рекордер
			handler.ServeHTTP(rr, req)

			// Проверяем, что статус ответа корректный
			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp save.Response

			// Анмаршаллим тело, и проверяем что при этом не возникло ошибок
			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			// Проверяем наличие требуемой ошибки в ответе
			require.Equal(t, tC.respError, resp.Error)
		})
	}
}
