package bot

import "fmt"

func getHelpMessage() string {
	fmt.Println("got the help message")
	return `help is not implemented yet`
}

func ping() string {
	return "pong"
}