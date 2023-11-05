package http

import (
	"encoding/json"
	"log"
	"net/http"
)

type response struct {
	Status int         `json:"status"`
	Result interface{} `json:"result"`
}

func newResponse(data interface{}, status int) *response {
	return &response{
		Status: status,
		Result: data,
	}
}

func (r *response) bytes() []byte {
	data, _ := json.Marshal(r)
	return data
}

func (r *response) string() string {
	return string(r.bytes())
}

func (r *response) sendResponse(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(r.Status)
	_, _ = w.Write(r.bytes())
	log.Println(r.string())
}


//200
func StatusOK(w http.ResponseWriter, r *http.Request, data interface{}){
	newResponse(data, http.StatusOK).sendResponse(w, r)
}
//204
func StatusNoContent(w http.ResponseWriter, r *http.Request){
	newResponse(nil, http.StatusNoContent).sendResponse(w, r)
}
//400
func StatusBadRequest(w http.ResponseWriter, r *http.Request, err error){
	data := map[string]interface{}{"error":err.Error()}
	newResponse(data, http.StatusBadRequest).sendResponse(w, r)
}
//404
func StatusNotFound(w http.ResponseWriter, r *http.Request, err error){
	data := map[string]interface{}{"error":err.Error()}
	newResponse(data, http.StatusNotFound).sendResponse(w, r)
}
//405
func StatusMethodNotAllowed(w http.ResponseWriter, r *http.Request){
	newResponse(nil, http.StatusMethodNotAllowed).sendResponse(w, r)
}
//409
func StatusConflict(w http.ResponseWriter, r *http.Request, err error){
	data := map[string]interface{}{"error":err.Error()}
	newResponse(data, http.StatusOK).sendResponse(w, r)
}
//500
func StatusInternalServerError(w http.ResponseWriter, r *http.Request, err error){
	data := map[string]interface{}{"error":err.Error()}
	newResponse(data, http.StatusOK).sendResponse(w, r)
}