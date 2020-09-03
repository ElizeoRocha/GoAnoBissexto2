package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type Year struct {
	Ano        int  `json:"ano"`
	Bissexto   bool `json:"bissexto"`
	ProximoAno int  `json:"proximoano"`
}

var Years []Year

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", listarAnos).Methods("GET")
	r.HandleFunc("/", cadastrarAnos).Methods("POST")
	http.Handle("/", r)
	fmt.Println("Status OK")
	http.ListenAndServe(":8080", nil)
}

func listarAnos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(Years)
}
func cadastrarAnos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//OBTENDO POST
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
	}
	//EXECUTANDO FUNÇÕES DO ANO
	var novoAno Year
	json.Unmarshal(body, &novoAno)
	novoAno.Bissexto = CalculoBissexto(novoAno.Ano)
	novoAno.ProximoAno = ProximoAnoBissexto(novoAno.Ano)
	Years = append(Years, novoAno)

	//EXIBINDO ANOS LISTADOS
	encoder := json.NewEncoder(w)
	encoder.Encode(Years)
}

//VERIFICAR SE O ANO É BISSEXTO
func CalculoBissexto(year int) bool {
	var numZero = 0
	var result bool
	if numZero == (year % 4) {
		if numZero == (year % 100) {
			if numZero == (year % 400) {
				result = true
			} else {
				result = false
			}
		} else {
			result = true
		}
	} else {
		result = false
	}
	return result
}

//PROXIMO ANO BISSEXTO
func ProximoAnoBissexto(ano int) int {
	ano = ano + 1
	for CalculoBissexto(ano) != true {
		ano = ano + 1
	}
	return ano
}
