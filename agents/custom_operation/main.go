package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/distum/agenty"
)

func main() {
	increment := agenty.NewOperation(incrementFunc)

	msg, err := agenty.NewProcess(
		increment,
		increment,
		increment,
	).Execute(
		context.Background(),
		agenty.NewMessage(
			agenty.UserRole,
			agenty.TextKind,
			[]byte("0"),
		),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(msg.Content()))
}

func incrementFunc(ctx context.Context, msg agenty.Message, _ *agenty.OperationConfig) (agenty.Message, error) {
	i, err := strconv.ParseInt(string(msg.Content()), 10, 10)
	if err != nil {
		return nil, err
	}
	inc := strconv.Itoa(int(i) + 1)
	return agenty.NewMessage(agenty.ToolRole, agenty.TextKind, []byte(inc)), nil
}
