package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"../data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

//always include an http handler function This is like a factory method for each handler file
func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	//this is a factory method to determine the method of the http request

	//get ->
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}
	//post ->
	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	//put
	if r.Method == http.MethodPut {
		//expect the id of the URI
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)
		if len(g) != 1 {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		if len(g[0]) != 2 {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(rw, "Invlaid URI", http.StatusBadRequest)
			return
		}

		p.updateProducts(id, rw, r)
	}
	//catch all
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

//PUT
func (p *Products) updateProducts(id int, rw http.ResponseWriter, r *http.Request) {
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to decode the data into json", http.StatusBadRequest)
		return
	}
	//2. Handle the body by putting it into a struct
	p.l.Printf("Prod: %#v", prod)
	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}

}

//GET
func (p *Products) getProducts(rw http.ResponseWriter, h *http.Request) {
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
func (p *Products) addProduct(rw http.ResponseWriter, h *http.Request) {
	//handle the post request

	//1. Decode the json body of the request.
	prod := &data.Product{}
	err := prod.FromJSON(h.Body)
	if err != nil {
		http.Error(rw, "Unable to decode the data into json", http.StatusBadRequest)
		return
	}
	//2. Handle the body by putting it into a struct
	p.l.Printf("Prod: %#v", prod)
	data.AddProduct(prod)
}
