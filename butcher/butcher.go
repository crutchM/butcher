package butcher

import (
	"butcher/butcher/internal/slices"
	"context"
	"sync"
)

const (
	// делим на равные батчи по порядку
	DivideByConstant = iota
	// делим на основании поля(по факту group by)
	DivideByField
)

type ContentProvider struct {
	ErrChan    chan error
	OutputChan chan interface{}
}

func (cp *ContentProvider) stop() {
	close(cp.ErrChan)
	close(cp.OutputChan)
}

type Options struct {
	BatchType    int
	GroupByField string
	Goroutines   int
}

type Callback func(provider ContentProvider)

type Executor func(ctx context.Context, provider ContentProvider, data []interface{})

type Butcher struct {
	callbacks []Callback
	opt       Options
	provider  ContentProvider
	executor  Executor
	wg        sync.WaitGroup
}

func NewButcher(opts Options, executor Executor, callbacks ...Callback) *Butcher {
	return &Butcher{
		callbacks: callbacks,
		opt:       opts,
		provider: ContentProvider{
			ErrChan:    make(chan error),
			OutputChan: make(chan interface{}),
		},
		executor: executor,
		wg:       sync.WaitGroup{},
	}

}

func (b *Butcher) Run(ctx context.Context, data interface{}) {
	var batches [][]interface{}
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

func (b *Butcher) runCallbacks(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}
	wg := sync.WaitGroup{}
	for _, callback := range b.callbacks {
		wg.Add(1)
		go func(callback Callback) {
			defer wg.Done()
			callback(b.provider)
		}(callback)
	}

	wg.Wait()
	b.wg.Done()
}

func (b *Butcher) runExecutor(ctx context.Context, data [][]interface{}) {
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

		go func(batch []interface{}) {
			defer wg.Done()
			b.executor(ctx, b.provider, batch)

			<-active
		}(batch)
	}

	wg.Wait()

	b.wg.Done()
}
