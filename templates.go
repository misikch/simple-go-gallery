package main

import "html/template"

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

func NewTemplates() *template.Template {
	tmpl := template.Must(template.New("list").Parse(imagesTmpl))

	return tmpl
}
