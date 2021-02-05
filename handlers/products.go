package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"../data"
	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

//PUT
func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	err = data.UpdateProduct(id, &prod)

	//specific error to defined issue
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
	//any other type of error
	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}

}

//GET
func (p *Products) GetProducts(rw http.ResponseWriter, h *http.Request) {
	//route to return the inforamtion from teh db
	//;in this case we are using the products.go file from dta module
	lp := data.GetProducts()
	err := lp.ToJSON(rw) //this is the same as:
	/*
		d, err := json.Marshall(lp)
	*/
	if err != nil {
		http.Error(rw, "Unable to marshal jsonlp", http.StatusInternalServerError)
		return
	}
	//rw.Write(d)
}

//POST
func (p *Products) PostProduct(rw http.ResponseWriter, h *http.Request) {
	//handle the post request
	prod := h.Context().Value(KeyProduct{}).(data.Product)
	//2. Handle the body by putting it into a struct
	data.AddProduct(&prod)
}

//this is a key
type KeyProduct struct{}

//this function is meant to extract and marshal JSON data from
//a request
func (p Products) MiddleWareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		//this is basic middleware function using gorillamux
		prod := data.Product{}
		err := prod.FromJSON(r.Body)

		if err != nil {
			http.Error(rw, "Unable to unmarshal JSON", http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
