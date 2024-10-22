package bot

import (
	"fmt"
	"strings"
)

func getHelpMessage() string {
	fmt.Println("got the help message")
	return `help is not implemented yet`
}

func ping() string {
	return "pong"
}

func addPlayer(discordStmt string) string {
	arr := strings.Split(discordStmt, " ")
	fmt.Printf("arr: %v", arr)
	return "haven't implemented yet"
}