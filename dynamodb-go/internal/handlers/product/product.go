package product

import (
	product "dynamodb-go/internal/controllers"
	EntityProduct "dynamodb-go/internal/entities/product"
	"dynamodb-go/internal/handlers"
	"dynamodb-go/internal/repository/adapter"
	Rules "dynamodb-go/internal/rules"
	RulesProduct "dynamodb-go/internal/rules/product"
	httpStatus "dynamodb-go/utils/http"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type Handler struct{
	handlers.Interface
	Controller product.Interface
	Rules Rules.Interface
}

func NewHandler(repository adapter.Interface) handlers.Interface{
	return &Handler{
		Controller: product.NewController(repository),
		Rules: RulesProduct.NewRules(),
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	if chi.URLParam(r, "ID") != ""{
		h.getOne(w, r)
	} else {
		h.getAll(w, r)
	}
}

func (h *Handler) getOne(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.Parse(chi.URLParam(r, "ID"))
	if err != nil{
		httpStatus.StatusBadRequest(w, r, errors.New("id is not uuid valid"))
		return
	}
	response, err := h.Controller.ListOne(ID)
	if err != nil{
		httpStatus.StatusInternalServerError(w, r, err)
		return
	}
	httpStatus.StatusOK(w, r, response)
}

func (h *Handler) getAll(w http.ResponseWriter, r *http.Request){
	response, err := h.Controller.ListAll()
	if err != nil{
		httpStatus.StatusInternalServerError(w, r, err)
		return
	}

	httpStatus.StatusOK(w, r, response)
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request){
	productBody, err := h.getBodyAndValidate(r, uuid.Nil)
	if err != nil{
		httpStatus.StatusBadRequest(w, r, err)
		return
	}

	ID, err := h.Controller.Create(productBody)

	if err != nil{
		httpStatus.StatusInternalServerError(w, r, err)
		return
	}
	httpStatus.StatusOK(w, r, map[string]interface{}{"id":ID.String()})
}

func (h *Handler) Put(w http.ResponseWriter, r *http.Request){
	ID, err := uuid.Parse(chi.URLParam(r, "ID"))
	if err != nil{
		httpStatus.StatusBadRequest(w, r, errors.New("id is not uuid valid"))
		return
	}
	productBody, err := h.getBodyAndValidate(r, ID)
	if err != nil{
		httpStatus.StatusBadRequest(w, r, err)
		return
	}

	if err := h.Controller.Update(ID, productBody); err != nil{
		httpStatus.StatusInternalServerError(w, r, err)
		return
	}

	httpStatus.StatusNoContent(w, r)

}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.Parse(chi.URLParam(r, "ID"))
	if err != nil{
		httpStatus.StatusBadRequest(w, r, errors.New("id is not uuid valid"))
		return
	}

	if err := h.Controller.Remove(ID); err != nil{
		httpStatus.StatusInternalServerError(w, r, err)
		return
	}

	httpStatus.StatusNoContent(w, r)
}

func (h *Handler) Options(w http.ResponseWriter, r *http.Request){
	httpStatus.StatusNoContent(w, r)
}

func (h *Handler) getBodyAndValidate(r *http.Request, ID uuid.UUID) (*EntityProduct.Product, error) {
	productBody := &EntityProduct.Product{}
	body, err := h.Rules.ConvertIoReaderToStruct(r.Body, productBody)
	if err != nil {
		return &EntityProduct.Product{}, errors.New("body is required")
	}

	productParsed, err := EntityProduct.InterfaceToModel(body)
	if err != nil {
		return &EntityProduct.Product{}, errors.New("error on convert body to model")
	}

	setDefaultValues(productParsed, ID)

	return productParsed, h.Rules.Validate(productParsed)
}

func setDefaultValues(product *EntityProduct.Product, ID uuid.UUID){
	product.UpdatedAt = time.Now()
	if ID == uuid.Nil{
		product.ID = uuid.New()
		product.CreatedAt = time.Now()
	} else {
		product.ID = ID
	}
}