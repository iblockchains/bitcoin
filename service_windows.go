package main

import "fmt"

func serviceMain() (bool, error) {
	fmt.Println("待:serviceMain")
	return false, nil
}
func init() {
	winServiceMain = serviceMain
}
