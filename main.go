package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type Users []User

type Arguments map[string]string

func Perform(args Arguments, writer io.Writer) error {
	if args["fileName"] == "" {
		return fmt.Errorf("-fileName flag has to be specified")
	}

	switch args["operation"] {
	case "list":
		return listItems(args, writer)
	case "add":
		return AddItem(args, writer)
	case "remove":
		return RemoveItem(args, writer)
	case "findById":
		return FindItem(args, writer)
	case "":
		return fmt.Errorf("-operation flag has to be specified")
	default:
		return fmt.Errorf("Operation %s not allowed!", args["operation"])
	}
}

func listItems(args Arguments, writer io.Writer) error {
	file, err := os.OpenFile(args["fileName"], os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	writer.Write(content)
	return nil
}

func AddItem(args Arguments, writer io.Writer) error {
	if args["item"] == "" {
		return fmt.Errorf("-item flag has to be specified")
	}
	users := Users{}
	newUser := &User{}
	json.Unmarshal([]byte(args["item"]), newUser)

	file, err := os.OpenFile(args["fileName"], os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	content, errRead := ioutil.ReadAll(file)
	if errRead != nil {
		return errRead
	}

	json.Unmarshal(content, &users)

	for _, u := range users {
		if u.ID == newUser.ID {
			writer.Write([]byte(fmt.Sprintf("Item with id %s already exists", u.ID)))
			return nil
		}
	}
	users = append(users, *newUser)
	data, errMarsh := json.Marshal(users)
	if errMarsh != nil {
		return errMarsh
	}

	file.Write(data)
	return nil
}

func RemoveItem(args Arguments, writer io.Writer) error {
	if args["id"] == "" {
		return fmt.Errorf("-id flag has to be specified")
	}

	file, err := os.OpenFile(args["fileName"], os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	users := Users{}
	content, errRead := ioutil.ReadAll(file)
	if errRead != nil {
		return errRead
	}

	json.Unmarshal(content, &users)
	var index int
	var found bool
	for i, u := range users {
		if u.ID == args["id"] {
			found = true
			index = i
			break
		}
	}
	if found {
		users = append(users[:index], users[index+1:]...)
		data, _ := json.Marshal(users)
		err := ioutil.WriteFile(args["fileName"], data, 0o644)
		if err != nil {
			panic(err)
		}
	} else {
		writer.Write([]byte(fmt.Sprintf("Item with id %s not found", args["id"])))
	}
	return nil
}

func FindItem(args Arguments, writer io.Writer) error {
	if args["id"] == "" {
		return fmt.Errorf("-id flag has to be specified")
	}

	file, err := os.OpenFile(args["fileName"], os.O_RDONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	users := Users{}
	content, errRead := ioutil.ReadAll(file)
	if errRead != nil {
		return errRead
	}

	json.Unmarshal(content, &users)
	for _, u := range users {
		if u.ID == args["id"] {
			data, err := json.Marshal(u)
			if err != nil {
				return err
			}
			writer.Write(data)
			return nil
		}
	}
	writer.Write([]byte(""))
	return nil
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}

func parseArgs() Arguments {
	idFlag := flag.String("id", "", "id of item")
	operationFlag := flag.String("operation", "", "the operation type")
	itemFlag := flag.String("item", "", "the item info")
	fileNameFlag := flag.String("fileName", "", "the file name")

	flag.Parse()

	newArgumet := Arguments{
		"id":        *idFlag,
		"operation": *operationFlag,
		"item":      *itemFlag,
		"filename":  *fileNameFlag,
	}

	return newArgumet
}
