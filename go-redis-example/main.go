package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go-redis-example/example"
)

func main() {
	defer func() {
		_ = example.RedisCli.Close()
	}()

	argsProg := os.Args
	var argsWithoutProg []string
	if len(argsProg) > 0 {
		argsWithoutProg = os.Args[1:]
		fmt.Printf("输入参数:\n%s\n----------\n", strings.Join(argsWithoutProg, "\n"))
	}
	ctx := context.Background()
	runExample := argsWithoutProg[0]
	exampleParams := argsWithoutProg[1:]

	switch runExample {
	case "Ex01":
		example.Ex01(ctx, exampleParams)
	case "Ex02":
		example.Ex02(ctx)
	case "Ex03":
		example.Ex03(ctx)
	case "Ex04":
		example.Ex04(ctx)
	case "Ex05":
		example.Ex05(ctx, exampleParams)
	case "Ex06":
		example.Ex06(ctx, exampleParams)
	case "Ex07":
		example.Ex07(ctx)
	default:
		panic(fmt.Sprintf("not support type: %s", runExample))
	}
}
