package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

func main() {

	// file, err := os.Open("../bin/classServer.exe")
	file, err := os.Open("C:/Users/Public/classServer/classServer.exe")
	if err != nil {
		fmt.Println("err:", err)
	}

	data := md5.New()
	if _, err := io.Copy(data, file); err != nil {
		fmt.Println("Error copying to hash")
	}

	fmt.Printf("Sum is: %x\n", data.Sum(nil))

}
