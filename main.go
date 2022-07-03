package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
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
	fileName = strings.TrimSpace(fileName)
	if !ok || fileName == "" {
		return errors.New("-fileName flag has to be specified")
	}

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, filePerm)
	if err != nil {
		return errors.New("open file error")
	}
	defer file.Close()

	operation, ok := args["operation"]
	operation = strings.TrimSpace(operation)
	if !ok || operation == "" {
		return errors.New("-operation flag has to be specified")
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
		item = strings.TrimSpace(item)
		if !ok || item == "" {
			return errors.New("-item flag has to be specified")
		}

		msg, err := add(file, item)
		if err != nil {
			return err
		}
		_, err = writer.Write([]byte(msg))
		if err != nil {
			return errors.New("writer error")
		}

		return nil
	case "findById":
		id, ok := args["id"]
		id = strings.TrimSpace(id)
		if !ok || id == "" {
			return errors.New("-id flag has to be specified")
		}

		itemBytes, err := findById(file, id)

		_, err = writer.Write(itemBytes)
		if err != nil {
			return errors.New("writer error")
		}
	default:
		return errors.New(fmt.Sprintf("Operation %s not allowed!", operation))
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

func add(file *os.File, item string) (string, error) {
	bytes, err := list(file)
	if err != nil {
		return "", err
	}

	users := make([]UserStruct, 0)
	json.Unmarshal(bytes, &users)

	var user UserStruct
	if err := json.Unmarshal([]byte(item), &user); err != nil {
		return "", err
	}

	for _, us := range users {
		if user.Id == us.Id {
			return fmt.Sprintf("Item with id %s already exists", user.Id), nil
		}
	}

	users = append(users, user)

	usersBytes, err := json.Marshal(users)
	if err != nil {
		return "", err
	}

	if _, err := file.WriteAt(usersBytes, 0); err != nil {
		return "", err
	}

	return "", nil
}

func findById(file *os.File, id string) ([]byte, error) {
	bytes, err := list(file)
	if err != nil {
		return nil, err
	}

	users := make([]UserStruct, 0)
	json.Unmarshal(bytes, &users)

	foundId := -1
	for i, us := range users {
		if id == us.Id {
			foundId = i
		}
	}

	if foundId == -1 {
		return []byte{}, nil
	}

	usersBytes, err := json.Marshal(users[foundId])
	if err != nil {
		return nil, err
	}

	return usersBytes, nil
}

func main() {

	//  go run ./main.go -operation "list" -fileName "userss.json" -item '{"id": "1", "email": "email@test.com", "age": 23}' -fileName "usersS.json"
	// userss.json	[{"id":"1","email":"test@test.com","age":34},{"id":"2","email":"tes2@test.com","age":32}]

	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
