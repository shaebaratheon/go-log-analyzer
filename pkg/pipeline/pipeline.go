package pipeline

import (
	"context"
	"sync"
	"time"
)

type LogEntry struct {
	Timestamp time.Time
	Level     string
	Service   string
	Message   string
	Fields    map[string]interface{}
}

type Transformer interface {
	Transform(entry *LogEntry) (*LogEntry, error)
}

type Sink interface {
	Write(batch []*LogEntry) error
	Close() error
}

type StreamPipeline struct {
	workers      int
	batchSize    int
	flushTimeout time.Duration
	inChan       chan *LogEntry
	transformers []Transformer
	sink         Sink
	wg           sync.WaitGroup
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewStreamPipeline(workers, batchSize int, flushTimeout time.Duration, sink Sink) *StreamPipeline {
	ctx, cancel := context.WithCancel(context.Background())
	return &StreamPipeline{
		workers:      workers,
		batchSize:    batchSize,
		flushTimeout: flushTimeout,
		inChan:       make(chan *LogEntry, batchSize*4),
		sink:         sink,
		ctx:          ctx,
		cancel:       cancel,
	}
}

func (p *StreamPipeline) AddTransformer(t Transformer) {
	p.transformers = append(p.transformers, t)
}

func (p *StreamPipeline) Start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.workerLoop()
	}
}

func (p *StreamPipeline) Submit(entry *LogEntry) {
	p.inChan <- entry
}

func (p *StreamPipeline) workerLoop() {
	defer p.wg.Done()
	batch := make([]*LogEntry, 0, p.batchSize)
	ticker := time.NewTicker(p.flushTimeout)
	defer ticker.Stop()

	flush := func() {
		if len(batch) > 0 {
			_ = p.sink.Write(batch)
			batch = make([]*LogEntry, 0, p.batchSize)
		}
	}

	for {
		select {
		case <-p.ctx.Done():
			flush()
			return
		case entry, ok := <-p.inChan:
			if !ok {
				flush()
				return
			}
			curr := entry
			var err error
			for _, t := range p.transformers {
				curr, err = t.Transform(curr)
				if err != nil || curr == nil {
					break
				}
			}
			if curr != nil {
				batch = append(batch, curr)
				if len(batch) >= p.batchSize {
					flush()
				}
			}
		case <-ticker.C:
			flush()
		}
	}
}

func (p *StreamPipeline) Stop() {
	close(p.inChan)
	p.cancel()
	p.wg.Wait()
	p.sink.Close()
}
