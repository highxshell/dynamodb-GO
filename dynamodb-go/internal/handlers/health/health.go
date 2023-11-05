package health

import (
	"dynamodb-go/internal/handlers"
	"dynamodb-go/internal/repository/adapter"
	httpStatus "dynamodb-go/utils/http"
	"errors"
	"net/http"
)

type Handler struct{
	handlers.Interface
	Repository adapter.Interface
}

func NewHandler(repository adapter.Interface) handlers.Interface{
	return &Handler{
		Repository: repository,
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	if !h.Repository.Health(){
		httpStatus.StatusInternalServerError(w, r, errors.New("relational database not alive"))
		return
	}
	httpStatus.StatusOK(w, r, "Service OK")
}

func (h *Handler) Put(w http.ResponseWriter, r *http.Request) {
	httpStatus.StatusMethodNotAllowed(w, r)
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	httpStatus.StatusMethodNotAllowed(w, r)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	httpStatus.StatusMethodNotAllowed(w, r)
}

func (h *Handler) Options(w http.ResponseWriter, r *http.Request) {
	httpStatus.StatusMethodNotAllowed(w, r)
}