package main

import "fmt"

func serviceMain() (bool, error) {
	fmt.Println("Unfinished:serviceMain")
	return false, nil
}
func init() {
	winServiceMain = serviceMain
}
