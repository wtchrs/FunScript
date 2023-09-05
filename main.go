package main

import (
	"fmt"
	"funscript/repl"
	"os"
	"os/user"
)

const LOGO = `
███████╗██╗   ██╗███╗   ██╗███████╗ ██████╗██████╗ ██╗██████╗ ████████╗
██╔════╝██║   ██║████╗  ██║██╔════╝██╔════╝██╔══██╗██║██╔══██╗╚══██╔══╝
█████╗  ██║   ██║██╔██╗ ██║███████╗██║     ██████╔╝██║██████╔╝   ██║   
██╔══╝  ██║   ██║██║╚██╗██║╚════██║██║     ██╔══██╗██║██╔═══╝    ██║   
██║     ╚██████╔╝██║ ╚████║███████║╚██████╗██║  ██║██║██║        ██║   
╚═╝      ╚═════╝ ╚═╝  ╚═══╝╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝╚═╝        ╚═╝   
`

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf(LOGO)
	fmt.Printf("Hello %s! This is the FunScript programming language!!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
