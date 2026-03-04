/*
Copyright © 2025 Arthur Mariano
*/
package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/arthvm/ditto/cmd"
)

func main() {
	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "warning: failed to load .env: %v\n", err)
		}
	}
	cmd.Execute()
}
