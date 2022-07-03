package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

const filePerm = 0644

type Arguments map[string]string

type UserStruct struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

var (
	operation string
	item      string
	fileName  string
)

func init() {
	flag.StringVar(&operation, "operation", "", "argument -operation")
	flag.StringVar(&item, "item", "", "argument -item")
	flag.StringVar(&fileName, "fileName", "users.json", "argument -fileName")
}

func Perform(args Arguments, writer io.Writer) error {
	fileName, ok := args["fileName"]
	if !ok {
		return errors.New("invalid args")
	}

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, filePerm)
	if err != nil {
		return errors.New("open file error")
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return errors.New("read file error")
	}

	var us []UserStruct
	if err := json.Unmarshal(bytes, &us); err != nil {
		return err
	}
	var s UserStruct
	if err := json.Unmarshal([]byte(args["item"]), &s); err != nil {
		return err
	}
	us = append(us, s)
	fmt.Println(us)

	_, err = writer.Write(bytes)
	if err != nil {
		return errors.New("writer error")
	}

	return nil
}

func parseArgs() Arguments {
	flag.Parse()

	args := Arguments{
		"id":        "",
		"operation": operation,
		"item":      item,
		"fileName":  fileName,
	}

	fmt.Println(flag.Args())
	fmt.Println(operation)
	fmt.Println(item)
	fmt.Println(fileName)

	return args
}

func main() {

	//  go run ./main.go -operation "list" -fileName "userss.json" -item '{"id": "1", "email": "email@test.com", "age": 23}' -fileName "usersS.json"
	// userss.json	[{"id":"1","email":"test@test.com","age":34},{"id":"2","email":"tes2@test.com","age":32}]

	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
