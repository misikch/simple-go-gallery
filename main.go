package main

import (
	"fmt"
	"net/http"
)

//MVP prototype

func main() {
	h := &PhotoListHandler{
		St: NewStorage(),
		Tmpl: NewTemplates(),
	}

	http.HandleFunc("/", h.List)
	http.HandleFunc("/upload", h.Upload)

	staticHandler := http.StripPrefix(
		"/images/",
		http.FileServer(http.Dir("./images")),
	)

	http.Handle("/images/", staticHandler)

	fmt.Println("server started at 8082")
	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		fmt.Println("server starting error: ", err)
	}
}



