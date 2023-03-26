package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"text/tabwriter"
	"time"

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

func delete() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Resourse to delete: ")
	resource, _ := reader.ReadString('\n')
	return strings.TrimSpace(resource)
}

func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		fmt.Print("\033[H\033[2J")
	}
}
func help() {
	// Display the menu
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	defer w.Flush()
	fmt.Fprintf(w, "%s\t%s\n", "add", "to add password")
	fmt.Fprintf(w, "%s\t%s\n", "list", "to list all records")
	fmt.Fprintf(w, "%s\t%s\n", "delete", "to delete a record")
	fmt.Fprintf(w, "%s\t%s\n", "exit", "to exit a programm")
	w.Flush()
}

func main() {

	// get current user in ordr to store a db file in $HOME

	u, err := user.Current()
	if err != nil {
		panic(err)
	}

	// Check and open the database
	dbPath := filepath.Join(u.HomeDir, "pass.db")
	_, err = os.Stat(filepath.Join(u.HomeDir, "pass.db"))
	var db *buntdb.DB
	if os.IsNotExist(err) {
		fmt.Println("DB does not exist, creating")
		db, err = buntdb.Open(dbPath)
		if err != nil {
			fmt.Println("Cannot create database")
			panic(err)
		} else {
			fmt.Printf("DB has been created successfully at path %s\n", filepath.Join(u.HomeDir, "pass.db"))
			time.Sleep(2500 * time.Millisecond)
			clearScreen()
		}
	} else {
		db, err = buntdb.Open(dbPath)
		if err != nil {
			fmt.Println("Cannot open database")
			panic(err)
		}
	}

	reader := bufio.NewReader(os.Stdin)

	// close DB after the funtion exits

	help()
	defer db.Close()

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
			fmt.Println("The record has been created successfully")
		case "delete":
			resource := delete()
			err := db.Update(func(tx *buntdb.Tx) error {
				_, err := tx.Delete(resource)
				return err
			})
			if err != nil {
				fmt.Println("The record has not been found")
			} else {
				fmt.Println("The record has been deleted")
			}
		case "list":
			// Read all values from the database
			err = db.View(func(tx *buntdb.Tx) error {
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
				defer w.Flush()
				fmt.Fprintf(w, "%s\t%s\t%s\n", "RESOURCE", "LOGIN", "PASSWORD")
				err := tx.Ascend("", func(key, value string) bool {
					username, password := strings.Split(value, "|")[0], strings.Split(value, "|")[1]
					fmt.Fprintf(w, "%s\t%s\t%s\n", key, username, password)
					return true
				})
				return err
			})
			if err != nil {
				panic(err)
			}
			continue
		case "exit":
			os.Exit(0)
		default:
			fmt.Println("no such command")
			help()
		}
	}
}
