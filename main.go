package main

import (
	"encoding/json"
	"errors"
	"flag"
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
	id        string
	operation string
	item      string
	fName     string
)

func init() {
	flag.StringVar(&id, "id", "", "argument -id")
	flag.StringVar(&operation, "operation", "", "argument -operation")
	flag.StringVar(&item, "item", "", "argument -item")
	flag.StringVar(&fName, "fileName", "users.json", "argument -fileName")
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

	operation, ok := args["operation"]
	if !ok {
		return errors.New("invalid args")
	}

	switch operation {
	case "list":
		itemBytes, err := list(file)

		_, err = writer.Write(itemBytes)
		if err != nil {
			return errors.New("writer error")
		}
	case "add":
		item, ok := args["item"]
		if !ok {
			return errors.New("invalid args")
		}

		err := add(file, item)
		if err != nil {
			return err
		}

		return nil
	default:
		return errors.New("unknown operation")
	}

	return nil
}

func parseArgs() Arguments {
	flag.Parse()

	args := Arguments{
		"id":        id,
		"operation": operation,
		"item":      item,
		"fileName":  fName,
	}

	return args
}

func list(file *os.File) ([]byte, error) {
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.New("read file error")
	}

	return bytes, nil
}

func add(file *os.File, item string) error {
	bytes, err := list(file)
	if err != nil {
		return err
	}

	var users []UserStruct
	if err := json.Unmarshal(bytes, &users); err != nil {
		return err
	}

	var user UserStruct
	if err := json.Unmarshal([]byte(item), &user); err != nil {
		return err
	}
	users = append(users, user)
	usersBytes, err := json.Marshal(users)
	if err != nil {
		return err
	}

	if _, err := file.WriteAt(usersBytes, 0); err != nil {
		return err
	}

	return nil
}

func main() {

	//  go run ./main.go -operation "list" -fileName "userss.json" -item '{"id": "1", "email": "email@test.com", "age": 23}' -fileName "usersS.json"
	// userss.json	[{"id":"1","email":"test@test.com","age":34},{"id":"2","email":"tes2@test.com","age":32}]

	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
