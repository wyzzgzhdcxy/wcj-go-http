package main

import (
	"fmt"
	"github.com/wyzzgzhdcxy/wcj-go-http/lib"
)

func main() {
	addr, err := lib.Start(":38472")
	if err != nil {
		fmt.Println("启动失败:", err)
		return
	}
	fmt.Println("启动成功，地址:", addr)
}
