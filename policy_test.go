package policy

import (
	"context"
	"fmt"
	"time"
)

func ExamplePolicy() {
	policy := NewWithRecovery(
		NewWithExclusiveExecution(
			NewWithContextValue(
				NewWithTimeout(nil, time.Minute),
				"my_key",
				"my_value",
			),
		),
		nil,
	)

	executor := NewBaseExecutor(policy, nil)

	executor.Execute(
		context.Background(),
		ActionFunc(func(ctx context.Context) error {
			time.Sleep(time.Second)
			fmt.Println("Hello, World!")
			fmt.Println("my_key=", ctx.Value("my_key"))
			return nil
		}),
	)

	// Output:
	// Hello, World!
	// my_key= my_value
}
