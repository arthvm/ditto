/*
Copyright © 2025 Arthur Mariano
*/
package main

import (
	"github.com/joho/godotenv"

	"github.com/arthvm/ditto/cmd"
)

func main() {
	godotenv.Load()
	cmd.Execute()
}
