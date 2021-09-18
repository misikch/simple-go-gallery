package main

import (
	"html/template"
	"log"
	"net/http"
)

var (
	userID = 0
)

type PhotoListHandler struct {
	St *StDb
	Tmpl *template.Template
}

func (h *PhotoListHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.St.GetPhotos(userID)
	if err != nil {
		log.Println("failed to get Items", err)
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}

	err = h.Tmpl.Execute(w, struct {
		Items []*Photo
	}{
		Items: items,
	})
	if err != nil {
		log.Println("failed to execute template", err)
		http.Error(w, "failed to execute template", http.StatusInternalServerError)
		return
	}
}

func (h PhotoListHandler) Upload(w http.ResponseWriter, r *http.Request) {
	uploadData, _, err := r.FormFile("my_file")
	if err != nil {
		log.Printf("failed to parse uploaded form file: %v", err)
		http.Error(w, "parse form file 'my_file' error", http.StatusInternalServerError)
		return
	}
	defer uploadData.Close()

	md5Sum, err := SaveFile(uploadData)
	if err != nil {
		log.Printf("failed to save uploaded form file: %v", err)
		http.Error(w, "save file error", http.StatusInternalServerError)
		return
	}

	realFile := "./images/" + md5Sum + ".jpg"
	err = MakeThumbnail(realFile, md5Sum)
	if err != nil {
		log.Printf("failed to create thumbnail for uploaded form file: %v", err)
		http.Error(w, "thumbnail file creating error", http.StatusInternalServerError)
		return
	}

	err = h.St.Add(&Photo{
		Path: md5Sum,
		UserID: userID,
	})
	if err != nil {
		log.Printf("failed to store file info into storage: %v", err)
		http.Error(w, "failed to store file info into storage", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}


