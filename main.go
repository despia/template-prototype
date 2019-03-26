package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbClient *mongo.Client

// var operations = map[string]func([]string){
// 	"add": func(command []string) {
// 		fmt.Println("test")
// 	},
// }

var fakeUserDB = map[string]string{
	"mail": "test@test.com",
	"name": "testname",
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	dbClient, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("db connected..")
	// loop()

	fmt.Println("building example data")
	english := Translation{
		HTML:    "en html template for {{ .name }} to {{ .mail }}",
		Text:    "en text template for {{ .name }} to {{ .mail }}",
		Expects: []string{"mail", "name"},
	}
	german := Translation{
		HTML:    "de html template for {{ .name }} to {{ .mail }}",
		Text:    "de text template for {{ .name }} to {{ .mail }}",
		Expects: []string{"mail", "name"},
	}
	something := Translation{
		HTML:    "rnd html template for {{ .name }} to {{ .mail }}",
		Text:    "rnd text template for {{ .name }} to {{ .mail }}",
		Expects: []string{"mail", "name"},
	}

	err = createTemplate("password changed", english)
	if err != nil {
		fmt.Println(err)
	}

	err = createTranslation("password changed", "DE", german)
	if err != nil {
		fmt.Println(err)
	}

	err = createTranslation("password changed", "RND", something)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("get and create tempaltes")
	tde, err := getTranslationOrDefault("password changed", "DE")
	if err != nil {
		log.Fatal(err)
	}
	parsedde := parseTemplate(tde)

	fmt.Println("execute templates")
	err = parsedde.HTML.Execute(os.Stdout, fakeUserDB)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel = getContext()
	defer cancel()
	fmt.Println("dropping example collection")
	getCollection().Drop(ctx)
}

// func loop() {
// 	reader := bufio.NewReader(os.Stdin)
// 	fmt.Println("Please input your command.")

// 	for true {
// 		fmt.Print("<< ")
// 		input, _ := reader.ReadString('\n')
// 		input = strings.TrimSuffix(input, "\n")

// 		if input == "exit" {
// 			os.Exit(0)
// 		}

// 		cmd := strings.Split(input, " ")

// 		if len(cmd) < 2 {
// 			fmt.Println("Error: missing command")
// 		} else {
// 			if f, exists := operations[cmd[0]]; exists {
// 				f(cmd[:1])
// 			}
// 		}
// 	}
// }
