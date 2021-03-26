package rest

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"strconv"
)

func ExtractID(request *http.Request) (int64, error) {
	idParam := chi.URLParam(request, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	return id, err
}

func WriteAsJson(writer http.ResponseWriter, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return err
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(bytes)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}
