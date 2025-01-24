package bitmap

import (
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {
	fmt.Println(hash("0x0000001000000001"))
	fmt.Println(hash("0x0000001000000001") % 2000)
	fmt.Println(hash("0x0000002000000001"))
	fmt.Println(hash("0x0000002000000001") % 2000)
	fmt.Println(hash("0x0000003000000001"))
	fmt.Println(hash("0x0000003000000001") % 2000)
}
