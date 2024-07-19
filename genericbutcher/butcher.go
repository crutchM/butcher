package genericbutcher

import (
	"butcher/genericbutcher/customslices"
	"context"
	"sync"
)

const (
	// делим на равные батчи по порядку
	DivideByConstant = iota
	// делим на основании поля(по факту group by)
	DivideByField
)

type ContentProvider[V any] struct {
	ErrChan    chan error
	OutputChan chan V
}

func (cp *ContentProvider[V]) stop() {
	close(cp.ErrChan)
	close(cp.OutputChan)
}

type Options struct {
	BatchType    int
	GroupByField string
	Goroutines   int
}

type Callback[V any] func(provider ContentProvider[V])

type Executor[T, V any] func(ctx context.Context, provider ContentProvider[V], data []T) error

type Butcher[T, V any] struct {
	callbacks []Callback[V]
	opt       Options
	provider  ContentProvider[V]
	executor  Executor[T, V]
	wg        sync.WaitGroup
}

func NewButcher[T, V any](opts Options, executor Executor[T, V], callbacks ...Callback[V]) *Butcher[T, V] {
	return &Butcher[T, V]{
		callbacks: callbacks,
		opt:       opts,
		provider: ContentProvider[V]{
			ErrChan:    make(chan error),
			OutputChan: make(chan V),
		},
		executor: executor,
		wg:       sync.WaitGroup{},
	}

}

func (b *Butcher[T, V]) Run(ctx context.Context, data []T) {
	var batches [][]T
	switch b.opt.BatchType {
	case DivideByConstant:
		batches = customslices.DivideSlice(data, b.opt.Goroutines)
	case DivideByField:
		batches = customslices.GroupByElements(data, b.opt.GroupByField)
	}
	b.wg.Add(2)

	go b.runCallbacks(ctx)

	go b.runExecutor(ctx, batches)

	b.wg.Wait()

}

func (b *Butcher[T, V]) runCallbacks(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}
	wg := sync.WaitGroup{}
	for _, callback := range b.callbacks {
		wg.Add(1)
		go func(callback Callback[V]) {
			defer wg.Done()
			callback(b.provider)
		}(callback)
	}

	wg.Wait()
	b.wg.Done()
}

func (b *Butcher[T, V]) runExecutor(ctx context.Context, data [][]T) {
	defer b.provider.stop()
	if ctx == nil {
		ctx = context.Background()
	}
	var active chan struct{}

	wg := sync.WaitGroup{}

	// ограничиваем максимально число одновременно-работающих горутин
	if b.opt.Goroutines > 0 {
		active = make(chan struct{}, b.opt.Goroutines)
	}

	for _, batch := range data {
		if len(batch) == 0 {
			continue
		}
		wg.Add(1)

		active <- struct{}{}

		go func(batch []T) {
			defer wg.Done()
			b.provider.ErrChan <- b.executor(ctx, b.provider, batch)

			<-active
		}(batch)
	}

	wg.Wait()

	b.wg.Done()
}
