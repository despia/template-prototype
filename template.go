package main

import (
	"log"
	"text/template"
)

type parsedTranslation struct {
	HTML *template.Template
	Text *template.Template
}

func parseTemplate(t Translation) parsedTranslation {
	p := parsedTranslation{}
	var err error
	p.HTML, err = template.New("password changed html").Parse(t.HTML)
	if err != nil {
		log.Fatal(err)
	}
	p.Text, err = template.New("password changed text").Parse(t.Text)
	if err != nil {
		log.Fatal(err)
	}

	return p
}
