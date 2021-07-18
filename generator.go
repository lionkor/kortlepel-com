package main

import (
	"html/template"
	"io"
	"log"
	"os"
)

func main() {
	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	templ_file, err := os.Open("template.html")
	check(err)
	tpl_buf := make([]byte, 0)
	_, err = templ_file.Read(tpl)
	check(err)
	tpl := string(tpl_buf)

	t, err := template.New("webpage").Parse(tpl)
	check(err)

	data := struct {
		Title   string
		Content template.HTML
	}{
		Title:   "Kortlepel.com",
		Content: `Hello, <b>Dashi</b>!!`,
	}

	dashafile, err := os.OpenFile("dasha.html", os.O_WRONLY|os.O_CREATE, 0600)
	err = t.Execute(io.Writer(dashafile), data)
	check(err)

	noItems := struct {
		Title   string
		Content template.HTML
	}{
		Title:   "Kortlepel.com",
		Content: `Hello, <b>Lion</b>!!`,
	}

	lionfile, err := os.OpenFile("lion.html", os.O_WRONLY|os.O_CREATE, 0600)
	err = t.Execute(io.Writer(lionfile), noItems)
	check(err)
}
