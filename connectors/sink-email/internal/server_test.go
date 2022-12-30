package internal

import (
	"fmt"
	stdMail "net/mail"
	"testing"
)

func TestA(t *testing.T) {
	addrs, err := stdMail.ParseAddressList("a@b.com,")
	fmt.Println(addrs)
	fmt.Println(err)
}
