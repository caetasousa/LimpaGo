package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"phresh-go/api/dto"
)

// escreverJSON serializa dados como JSON e escreve na resposta.
func escreverJSON(w http.ResponseWriter, status int, dados interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(dados)
}

// escreverErro serializa e escreve um erro de domínio na resposta.
func escreverErro(w http.ResponseWriter, err error) {
	status, resposta := dto.MapearErroDominio(err)
	escreverJSON(w, status, resposta)
}

// lerJSON decodifica o corpo da requisição em destino.
func lerJSON(r *http.Request, destino interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(destino)
}

// lerParamInteiro lê um parâmetro de path inteiro pelo nome.
func lerParamInteiro(r *http.Request, nome string) (int, error) {
	return strconv.Atoi(chi.URLParam(r, nome))
}

// lerQueryInteiro lê um parâmetro de query inteiro com valor padrão.
func lerQueryInteiro(r *http.Request, nome string, padrao int) int {
	val := r.URL.Query().Get(nome)
	if val == "" {
		return padrao
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return padrao
	}
	return n
}
