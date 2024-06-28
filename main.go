package main

import "C"
import (
	"go_demo/mit"
	"go_demo/tool"
)

//export Run
func Run() {

	mit.Run()
}
func main() {
	tool.Base64()
}
