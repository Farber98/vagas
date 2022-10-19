package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pagarme/internal/models"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHello(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	controller := &HelloController{}

	if assert.NoError(t, controller.Hello(c)) {

		errorResponseModel := &models.MsgResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), errorResponseModel)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, OK_HELLO, errorResponseModel.Message)
		}
	}

}
