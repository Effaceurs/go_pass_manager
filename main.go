package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"os/user"
	"path/filepath"

	"github.com/tidwall/buntdb"
)

type Password struct {
	resource string
	username string
	password string
}

func add() Password {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Resourse: ")
	resource, _ := reader.ReadString('\n')
	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	cm := Password{resource: strings.TrimSpace(resource), username: strings.TrimSpace(username), password: strings.TrimSpace(password)}
	return cm
}

func main() {

	// get current user in ordr to store a db file in $HOME

	u, err := user.Current()
	if err != nil {
		panic(err)
	}

	// Open the database

	db, err := buntdb.Open(filepath.Join(u.HomeDir, "pass.db"))
	reader := bufio.NewReader(os.Stdin)

	// close DB after the funtion exits

	defer db.Close()
	// Display the menu
	fmt.Println("Options:")
	fmt.Println("add - to add password")
	fmt.Println("list - to list all records")
	fmt.Println("exit - to exit a programm")
	for {

		// MENU

		// user prompt command
		fmt.Print("> ")
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)
		switch command {
		case "add":
			combination := add()
			if err := db.Update(func(tx *buntdb.Tx) error {
				_, _, err := tx.Set(combination.resource, fmt.Sprintf("%s|%s", combination.username, combination.password), nil)
				return err
			}); err != nil {
				panic(err)
			}
			fmt.Println("Password has been stored")
			return
		case "delete":
			fmt.Println("delete")
		case "list":
			// Read all values from the database
			err = db.View(func(tx *buntdb.Tx) error {
				err := tx.Ascend("", func(key, value string) bool {
					username, password := strings.Split(value, "|")[0], strings.Split(value, "|")[1]
					fmt.Printf("%s = %s: %s\n", key, username, password)
					return true
				})
				return err
			})
			if err != nil {
				panic(err)
			}
			continue
		case "help":
			fmt.Println("help")
		case "exit":
			os.Exit(0)
		default:
			fmt.Println("no such command")
		}
	}
}
