package main

import (
	"terra9.it/vadovia/internal"
)

func main() {
	internal.Version.Print()
	NewApp().Run()
}
