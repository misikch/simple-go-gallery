package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

//MVP prototype

//language=html
var imagesTmpl = `
<html>
	<body>
		<div>
			<form action="/upload" method="post" enctype="multipart/form-data">
				Image <input type="file" name="my_file"><br />
				<input type="submit" value="Upload">
			</form>
		</div>
		<br />
		<div>
			{{range.Items}}
				<div>
					<img src="/images/{{.Path}}_160.jpg" /><br />
				</div>
			{{end}}	
		</div>
	</body>
</html>
`

type Photo struct {
	ID int
	UserID int
	Path string
}

var (
	items = []*Photo{}
	userID = 0
)

func List(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("list").Parse(imagesTmpl))
	tmplArgs := struct {
		Items []*Photo
	}{
		Items: items,
	}

	err := tmpl.Execute(w, tmplArgs)
	if err != nil {
		log.Printf("failed to execute template: %v", err)
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func Upload(w http.ResponseWriter, r *http.Request) {
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

	items = append(
		items,
		&Photo{
			Path: md5Sum,
			UserID: userID,
		})

	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	http.HandleFunc("/", List)
	http.HandleFunc("/upload", Upload)

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



