package main

import (
	"butcher/butcher"
	"context"
	"fmt"
	"time"
)

type Data struct {
	Name  string
	Value int
}

type Info struct {
	Caller  string
	Message string
}

func main() {
	ctx := context.Background()
	b := butcher.NewButcher(
		butcher.Options{
			BatchType:    butcher.DivideByField,
			GroupByField: "Name",
			Goroutines:   2,
		},
		Foo(),
		Performer(),
		ErrorHandler(),
	)

	b.Run(ctx, []interface{}{
		Data{
			Name:  "A",
			Value: 1,
		},
		Data{
			Name:  "B",
			Value: 1,
		},
		Data{
			Name:  "B",
			Value: 2,
		},
		Data{
			Name:  "A",
			Value: 2,
		},
	})

}

func Foo() butcher.Executor {
	return func(ctx context.Context, provider butcher.ContentProvider, data []interface{}) {
		for _, v := range data {
			value, _ := v.(Data)

			time.Sleep(1 * time.Second)

			provider.OutputChan <- Info{
				Caller:  "Foo",
				Message: fmt.Sprintf("Name: %s, Value: %d", value.Name, value.Value),
			}
		}
	}
}

func ErrorHandler() butcher.Callback {
	return func(provider butcher.ContentProvider) {
		for err := range provider.ErrChan {
			fmt.Println(err)
		}
	}
}

func Performer() butcher.Callback {
	return func(provider butcher.ContentProvider) {
		for v := range provider.OutputChan {
			fmt.Println(v)
		}
	}
}
