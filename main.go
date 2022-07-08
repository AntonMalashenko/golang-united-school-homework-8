package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

var operationArg = flag.String("operation", "", "")
var fileNameArg = flag.String("fileName", "", "")
var itemArg = flag.String("item", "", "")
var userIdArg = flag.String("id", "", "")

const ADD = "add"
const LIST = "list"
const REMOVE = "remove"
const FIND_BY_ID = "findById"

const FilePermission = 0644

var allowedOperations = []string{
	ADD,
	LIST,
	REMOVE,
	FIND_BY_ID,
}

type Arguments map[string]string

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func ValidateArgs(args Arguments) error {
	operation := args["operation"]
	if operation == "" {
		return fmt.Errorf("-operation flag has to be specified")
	}
	if !isExists(allowedOperations, operation) {
		return fmt.Errorf("Operation %s not allowed!", operation)
	}
	if args["fileName"] == "" {
		return fmt.Errorf("-fileName flag has to be specified")
	}
	if operation == ADD {
		if args["item"] == "" {
			return fmt.Errorf("-item flag has to be specified")
		}
	} else if operation == FIND_BY_ID || operation == REMOVE {
		if args["id"] == "" {
			return fmt.Errorf("-id flag has to be specified")
		}
	}
	return nil
}

func List(arguments Arguments, writer io.Writer) error {
	file, err := os.OpenFile(arguments["fileName"], os.O_RDWR, FilePermission)
	defer file.Close()
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	//var users []User
	//err = json.Unmarshal(bytes, &users)
	//if err != nil {
	//	return err
	//}
	_, err = writer.Write(bytes)
	return nil
}

func Add(arguments Arguments, writer io.Writer) error {
	file, err := os.OpenFile(arguments["fileName"], os.O_RDWR|os.O_CREATE, FilePermission)
	defer file.Close()
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	var users []User
	if len(bytes) > 0 {
		err = json.Unmarshal(bytes, &users)
		if err != nil {
			return err
		}
	}

	var itemToAdd User
	json.Unmarshal([]byte(arguments["item"]), &itemToAdd)

	for _, user := range users {
		if user.Id == itemToAdd.Id {
			result := fmt.Sprintf("Item with id %s already exists", itemToAdd.Id)
			_, err = writer.Write([]byte(result))
			return nil
		}
	}

	users = append(users, itemToAdd)
	bytes, err = json.Marshal(users)
	ioutil.WriteFile(arguments["fileName"], bytes, FilePermission)
	_, err = writer.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func Remove(arguments Arguments, writer io.Writer) error {
	file, err := os.OpenFile(arguments["fileName"], os.O_RDWR, FilePermission)
	defer file.Close()
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	var users []User
	err = json.Unmarshal(bytes, &users)
	if err != nil {
		return err
	}

	userId := arguments["id"]
	var usersToSave []User
	var found bool

	for _, user := range users {
		if user.Id == userId {
			found = true
		} else {
			usersToSave = append(usersToSave, user)
		}
	}
	if !found {
		result := fmt.Sprintf("Item with id %s not found", userId)
		_, err = writer.Write([]byte(result))
	} else {
		bytes, err = json.Marshal(usersToSave)
		ioutil.WriteFile(arguments["fileName"], bytes, FilePermission)
		_, err = writer.Write(bytes)
	}
	if err != nil {
		return err
	}
	return nil
}

func FindById(arguments Arguments, writer io.Writer) error {
	file, err := os.OpenFile(arguments["fileName"], os.O_RDWR, FilePermission)
	defer file.Close()
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	var users []User
	err = json.Unmarshal(bytes, &users)
	if err != nil {
		return err
	}
	userId := arguments["id"]
	var result interface{}
	for _, user := range users {
		if user.Id == userId {
			result = user
		}
	}
	if result == nil {
		_, err = writer.Write([]byte{})
	} else {
		bytes, err = json.Marshal(result)
		_, err = writer.Write(bytes)
	}
	if err != nil {
		return err
	}
	return nil
}

func Perform(args Arguments, writer io.Writer) error {
	err := ValidateArgs(args)
	if err != nil {
		return err
	}

	switch args["operation"] {
	case LIST:
		return List(args, writer)
	case REMOVE:
		return Remove(args, writer)
	case ADD:
		return Add(args, writer)
	case FIND_BY_ID:
		return FindById(args, writer)
	}

	return nil
}

func isExists(list []string, value string) bool {
	for _, v := range list {
		if value == v {
			return true
		}
	}
	return false
}

func parseArgs() (Arguments, error) {
	flag.Parse()
	args := Arguments{
		"operation": *operationArg,
		"fileName":  *fileNameArg,
		"item":      *itemArg,
		"userId":    *userIdArg,
	}
	return args, nil
}

func main() {
	args, err := parseArgs()
	if err != nil {
		panic(err)
	}
	err = Perform(args, os.Stdout)
	if err != nil {
		panic(err)
	}
}
