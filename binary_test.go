package main_test

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func Test_bin(t *testing.T) {
	err := ioutil.WriteFile(".test.bin", []byte{0x2, 0x0, 0xA, 0x0, 0x0, 0x0, 0x0, 0x0}, 0777)
	if err != nil {
		fmt.Println(err)
	}
}
