package main

import (
	"fmt"
	"go-base/apollo"
)

func main() {

	vaule := apollo.GetConfigValue("spring.data.mongodb.uri")
	fmt.Println("测试值:" + vaule)

	vaule1 := apollo.GetConfigValue("smsMethord")
	fmt.Println("测试值:" + vaule1)

}
