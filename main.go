package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()

	//	print ctx
	fmt.Printf("hello, %#v\n", ctx)
}
