package main

import "github.com/MaksymLeiber/pgen/cmd"

// Version задается при сборке
var Version = "1.0.3"

func main() {
	cmd.Version = Version
	cmd.Execute()
}
