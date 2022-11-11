package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	tplBuf, err := ioutil.ReadFile("template.html")
	check(err)
	tpl := string(tplBuf)

	t, err := template.New("webpage").Parse(tpl)
	check(err)

	// add your new files here
	var contentFilePrefixes []string

	now := time.Now()
	sourceFolder := "contents/"
	targetFolder := "www/"

	err = filepath.Walk(sourceFolder, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Printf("ERROR: error walking \"%s\" folder: %s", sourceFolder, err)
		}
		if info.IsDir() {
			// ensure that this directory exists in the target directory, too
			dirPath := strings.TrimPrefix(path, sourceFolder)
			dirPath = targetFolder + dirPath
			dirPathStat, err := os.Stat(dirPath)
			if err != nil || !dirPathStat.IsDir() {
				err := os.Mkdir(dirPath, os.ModePerm)
				if err != nil {
					log.Printf("ERROR: could not create \"%s\": %s. This might cause issues later.", dirPath, err)
				} else {
					log.Printf("created folder \"%s\" for later", dirPath)
				}
			}
			return nil
		}
		path = strings.TrimSuffix(path, filepath.Ext(path))
		path = strings.TrimPrefix(path, sourceFolder)
		contentFilePrefixes = append(contentFilePrefixes, path)
		return nil
	})
	check(err)

	if len(contentFilePrefixes) == 0 {
		log.Fatalf("ERROR: \"%s\" folder is empty!", sourceFolder)
	} else {
		log.Println("found prefixes:", strings.Join(contentFilePrefixes, ", "))
	}

	for _, prefix := range contentFilePrefixes {
		filename := sourceFolder + prefix + ".html"
		log.Println("processing " + prefix + " (" + filename + ")")
		content, err := os.ReadFile(filename)
		if err != nil {
			log.Println("ERROR: file \"" + filename + "\" could not be opened. this is not a fatal error. " +
				"output file was not created for the target \"" + prefix + "\".")
			continue
		}
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
		// second replace step
		secondTpl := output.String()

		t2, err := template.New("webpage").Parse(secondTpl)
		check(err)
		var finalOutput bytes.Buffer
		err = t2.Execute(&finalOutput, data)
		check(err)
		log.Printf("resulting size: %0.2f KB", float32(len(finalOutput.String()))/1024.0)
		err = ioutil.WriteFile(targetFolder+prefix+".html", finalOutput.Bytes(), 0644)
		check(err)
	}
}
