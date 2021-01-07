package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"encoding/csv"
	"os"
	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
)

type ProductTy struct {
	Name string `json:"name"`
	ImageURL string `json:"imageURL"`
	Description string `json:"description"`
	Price string `json:"price"`
	TotalReviews string `json:"totalReviews"`
}

type Crwl struct {
	Url string `json:"url"`
	Product ProductTy `json:"product"`
}

type UrlType struct {
	Url string `json:"url"`
}

var baseURL = "http://localhost:5000/"


func crawData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	crwl := Crwl{
		Url:"",
		Product:ProductTy{
			Name:"",
			ImageURL:"",
			Description:"",
			Price:"",
			TotalReviews:"",
		},
	}
	c := colly.NewCollector()

	//for product title
	c.OnHTML("#title", func(e *colly.HTMLElement) {
		crwl.Product.Name = e.ChildText("span")
	})

	//for product description
	c.OnHTML("#feature-bullets", func(e *colly.HTMLElement) {
		crwl.Product.Description = e.ChildText("li")
	})

	//for product price
	c.OnHTML("#edition_0_price", func(e *colly.HTMLElement) {
		crwl.Product.Price = e.ChildText("span")
	})

	//for product reviews
	c.OnHTML("#acrCustomerReviewLink", func(e *colly.HTMLElement) {
		crwl.Product.TotalReviews = e.ChildText("span")
	})

	//for product ImgUrl
	// c.OnHTML("#imgTagWrapperId",func(e *colly.HTMLElement){
	// 	crwl.Product.ImageURL = e.ChildAttr("img","src")
	// })
	
	var urlData UrlType
	_ = json.NewDecoder(r.Body).Decode(&urlData)
	crwl.Product.ImageURL = urlData.Url
	crwl.Url = urlData.Url
	c.Visit(urlData.Url)
	fmt.Println("Scrapping data")
	json.NewEncoder(w).Encode(crwl)
	jsonVal, _ := json.Marshal(crwl)
	_, _ = http.Post(baseURL+"add", "application/json", bytes.NewBuffer(jsonVal))
}

func saveData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rData Crwl
	_ = json.NewDecoder(r.Body).Decode(&rData)
	file, err := os.OpenFile("data.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write([]string{
		"{"+rData.Url,
		rData.Product.Name,
		rData.Product.ImageURL,
		rData.Product.Description,
		rData.Product.Price,
		rData.Product.TotalReviews+"}",
	})
	json.NewEncoder(w).Encode(rData)
	fmt.Println("Adding Data")
}

func handleRequests() {

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/",crawData)
	r.HandleFunc("/add",saveData)
	log.Fatal(http.ListenAndServe(":5000",r))
}

func main() {
	handleRequests()
}