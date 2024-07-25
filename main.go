package butcher

//package butcher
//
//import (
//	"butcher/butcher"
//	"butcher/genericbutcher"
//	"context"
//	"errors"
//	"fmt"
//	"time"
//)
//
//type Data struct {
//	Name  string
//	Value int
//}
//
//type Info struct {
//	Caller  string
//	Message string
//}
//
//func main() {
//	ctx := context.Background()
//	b := genericbutcher.NewButcher[Data, Info](
//		genericbutcher.Options{
//			BatchType:    butcher.DivideByField,
//			GroupByField: "Name",
//			Goroutines:   2,
//		},
//		Foo(),
//		Performer(),
//		ErrorHandler(),
//	)
//
//	b.Run(ctx, []Data{
//		Data{
//			Name:  "A",
//			Value: 1,
//		},
//		Data{
//			Name:  "B",
//			Value: 1,
//		},
//		Data{
//			Name:  "B",
//			Value: 2,
//		},
//		Data{
//			Name:  "A",
//			Value: 2,
//		},
//	})
//
//}
//
//func Foo() genericbutcher.Executor[Data, Info] {
//	return func(ctx context.Context, provider genericbutcher.ContentProvider[Info], data []Data) error {
//		for _, v := range data {
//			if v.Name == "A" {
//				return errors.New("some error")
//
//			}
//			time.Sleep(1 * time.Second)
//
//			provider.OutputChan <- Info{
//				Caller:  "Foo",
//				Message: fmt.Sprintf("Name: %s, Value: %d", v.Name, v.Value),
//			}
//		}
//
//		return nil
//	}
//}
//
//func ErrorHandler() genericbutcher.Callback[Info] {
//	return func(provider genericbutcher.ContentProvider[Info]) {
//		for err := range provider.ErrChan {
//			if err != nil {
//				fmt.Println(err)
//			}
//		}
//	}
//}
//
//func Performer() genericbutcher.Callback[Info] {
//	return func(provider genericbutcher.ContentProvider[Info]) {
//		for v := range provider.OutputChan {
//			fmt.Println(v)
//		}
//	}
//}
