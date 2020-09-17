package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Year struct {
	Ano        int  `json:"ano"`
	Bissexto   bool `json:"bissexto"`
	ProximoAno int  `json:"proximoano"`
}

var Years []Year = []Year{}

func logHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		x, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
		log.Println(fmt.Sprintf("%q", x))
		rec := httptest.NewRecorder()
		fn(rec, r)
		log.Println(fmt.Sprintf("%q", rec.Body))

		// this copies the recorded response to the response writer
		for k, v := range rec.HeaderMap {
			w.Header()[k] = v
		}
		w.WriteHeader(rec.Code)
		rec.Body.WriteTo(w)
	}
}
func MessageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "A message was received")
}

func main() {
	r := mux.NewRouter()

	headers := handlers.AllowedHeaders([]string{"X-Request", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"http://app-ano-bissexto-anobissexto.paulo-dev-apps.gncloud.nz/"})

	r.HandleFunc("/{ano}", getAno).Methods("GET")
	r.HandleFunc("/", logHandler(listarAnos)).Methods("GET")
	r.HandleFunc("/", logHandler(cadastrarAnos)).Methods("POST")
	http.Handle("/", r)
	fmt.Println("Status OK")
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		fmt.Println(t)
		return nil
	})
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(headers, methods, origins)(r)))

}

func listarAnos(w http.ResponseWriter, r *http.Request) {
	//EXIBINDO ANOS LISTADOS
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(Years)
}

func getAno(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	ano, _ := strconv.Atoi(params["ano"])
	for _, year := range Years {
		if year.Ano == ano {
			json.NewEncoder(w).Encode(year)
			return
		}
	}
	json.NewEncoder(w).Encode(&Year{})
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
