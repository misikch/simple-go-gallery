package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/misikch/simple-go-gallery/config"
)

//MVP prototype

func main() {
	cfg, err := config.InitConfig("simple_go_gallery")
	if err != nil {
		fmt.Println("init config error", err)
		return
	}

	db, err := sql.Open("mysql", cfg.MysqlMaster.DSN)
	if err != nil {
		fmt.Println("failed to setup mysql connection", err)
		return
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("failed to ping database", err)
		return
	}

	h := &PhotoListHandler{
		St: NewStorage(db),
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
	err = http.ListenAndServe(":8082", nil)
	if err != nil {
		fmt.Println("server starting error: ", err)
	}
}



