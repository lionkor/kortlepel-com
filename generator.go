package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

func main() {
	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	tpl_buf, err := ioutil.ReadFile("template.html")
	check(err)
	tpl := string(tpl_buf)

	t, err := template.New("webpage").Parse(tpl)
	check(err)

	// add your new files here
	content_file_prefixes := []string{"contact", "index", "projects", "404"}

	now := time.Now()

	for _, prefix := range content_file_prefixes {
		filename := "contents/" + prefix + ".html"
		fmt.Println("processing " + prefix + " (" + filename + ")")
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println("ERROR: file \"" + filename + "\" could not be opened. this is not a fatal error. " +
				"output file was not created for the target \"" + prefix + "\".")
			continue
		}
		fmt.Println("read", len(content), "characters")

		data := struct {
			Title   string
			Content template.HTML
			Date    string
		}{
			Title:   strings.Title(prefix),
			Content: template.HTML(content),
			Date:    fmt.Sprintf("%02d/%02d/%04d", now.Day(), now.Month(), now.Year()),
		}

		var output bytes.Buffer
		err = t.Execute(&output, data)
		check(err)
		fmt.Println("resulting size:", len(output.String()))
		err = ioutil.WriteFile("www/"+prefix+".html", output.Bytes(), 0644)
		check(err)
	}
}
