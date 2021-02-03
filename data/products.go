package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"` //the "-" will emit the values form the output
	UpdatedOn   string  `json:"-"` //the "-" will emit the values form the output
	DeletedOn   string  `json:"-"` //the "-" will emit the values form the output
}

func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r) //new decoder takes a reader
	return e.Decode(p)
}

type Products []*Product

// ^ this is just a custom tupe

//capture logic for encoding json in this method
//this can be something like a middleware
func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

//data function to actually perform the query
func GetProducts() Products {
	return productList
}

var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffe",
		Price:       2.45,
		SKU:         "adfczlk",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Capuccino",
		Description: "Short and strong coffe",
		Price:       4.50,
		SKU:         "ertlkj",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}

//add a product to the array
func AddProduct(p *Product) {
	//generate the id for the data
	p.ID = getNextID()
	productList = append(productList, p)
}

//update product
func UpdateProduct(id int, p *Product) error {
	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}

	p.ID = id
	productList[pos] = p
	return nil
}

var ErrProductNotFound = fmt.Errorf("Product not found")

func findProduct(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}
	return nil, -1, ErrProductNotFound
}

func getNextID() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}
