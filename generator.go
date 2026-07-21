package main

import (
	"bytes"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	tpl, err := os.ReadFile("template.html")
	check(err)

	t, err := template.New("webpage").Parse(string(tpl))
	check(err)

	sourceFolder := "contents/"
	targetFolder := "www/"
	var prefixes []string

	err = filepath.Walk(sourceFolder, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Printf("ERROR: walking %q: %s", sourceFolder, err)
			return nil
		}
		if info.IsDir() {
			dir := targetFolder + strings.TrimPrefix(path, sourceFolder)
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				os.MkdirAll(dir, os.ModePerm)
			}
			return nil
		}
		p := strings.TrimPrefix(path, sourceFolder)
		p = strings.TrimSuffix(p, filepath.Ext(p))
		prefixes = append(prefixes, p)
		return nil
	})
	check(err)

	if len(prefixes) == 0 {
		log.Fatalf("contents/ folder is empty")
	}
	log.Println("found:", strings.Join(prefixes, ", "))

	for _, p := range prefixes {
		content, err := os.ReadFile(sourceFolder + p + ".html")
		if err != nil {
			log.Printf("ERROR: skipping %q: %s", p, err)
			continue
		}
		data := struct {
			Title   string
			Content template.HTML
		}{
			Title:   title(p),
			Content: template.HTML(content),
		}
		var buf bytes.Buffer
		err = t.Execute(&buf, data)
		check(err)
		out := targetFolder + p + ".html"
		os.WriteFile(out, buf.Bytes(), 0644)
		log.Printf("wrote %s (%.0f KB)", out, float32(buf.Len())/1024)
	}
}

func title(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
