// Package main
// Wrote by yijian on 2024/01/04
package main

import (
	"fmt"
)
import (
	"github.com/eyjian/gomooon/utils"
)

func main() {
	testIsResidentIdentityCardNumber()
}

func testIsResidentIdentityCardNumber() {
	id := "371522199402189127"
	if utils.IsResidentIdentityCardNumber(id) {
		fmt.Printf("%s is ResidentIdentityCardNumber\n", id)
	} else {
		fmt.Printf("%s is not ResidentIdentityCardNumber\n", id)
	}

	id = "152822199008120030"
	if utils.IsResidentIdentityCardNumber(id) {
		fmt.Printf("%s is ResidentIdentityCardNumber\n", id)
	} else {
		fmt.Printf("%s is not ResidentIdentityCardNumber\n", id)
	}
}
