package main

import (
	"fmt"
	"groupie-tracker-filter/Handlers"
	"log"
	"net/http"
)

var err error

func main() {
	http.HandleFunc("/", Handlers.HandleHomePage)
	http.HandleFunc("/stars", Handlers.HandleStarsPage)
	http.HandleFunc("/about", Handlers.HandleAboutPage)
	http.HandleFunc("/stardetails/", Handlers.HandleStarDetailsPage)
	fs := http.FileServer(http.Dir("style"))
	http.Handle("/style/", http.StripPrefix("/style/", fs))
	fmt.Println("starting server at port 8177\n")
	err = http.ListenAndServe(":8177", nil)
	if err != nil {
		log.Fatal(err)
	}
}
