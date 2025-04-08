package utils

import (
	"encoding/json"
	"github.com/tuanpnt17/eleven-golang-projects/go-bookstore/pkg/models"
	"io"
	"net/http"
)

func ParseBody(request *http.Request, book *models.Book) {
	if body, err := io.ReadAll(request.Body); err == nil {
		if unmarshalErr := json.Unmarshal(body, book); unmarshalErr == nil {
			return
		}
	}
}
