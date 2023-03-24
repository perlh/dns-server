package util

import "fmt"

func PrintError(err error) {
	if err != nil {
		fmt.Println("error is ", err)
	}
}
