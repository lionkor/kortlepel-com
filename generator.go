package main

import (
  "fmt"
	"html/template"
	"log"
  "strings"
  "bytes"
  "io/ioutil"
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

  content_file_prefixes := []string { "contacts", "index", "projects" };

  for _, prefix := range content_file_prefixes {
    filename := "contents/" + prefix + ".html"
    fmt.Println("processing " + prefix + "(" + filename + ")")
    content, err := ioutil.ReadFile(filename)
    check(err)
    fmt.Println("read", len(content), "characters")

    data := struct {
      Title   string
      Content template.HTML
    }{
      Title: strings.Title(prefix),
      Content: template.HTML(content),
    }

    var output bytes.Buffer
    err = t.Execute(&output, data)
    check(err)
    fmt.Println("resulting size:", len(output.String()))
    err = ioutil.WriteFile(prefix + ".html", output.Bytes(), 0644)
	  check(err)
  }
}
