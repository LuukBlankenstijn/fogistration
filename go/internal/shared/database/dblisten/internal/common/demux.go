package common

import "sync"

type Demux struct {
	mu   sync.RWMutex
	subs map[string][]chan Mixed
}

func NewDemux(in <-chan Mixed) *Demux {
	d := &Demux{subs: make(map[string][]chan Mixed)}
	go func() {
		for m := range in {
			d.mu.RLock()
			list := append([]chan Mixed(nil), d.subs[m.Table]...)
			d.mu.RUnlock()
			for _, ch := range list {
				select { case ch <- m: default: }
			}
		}
		d.mu.RLock()
		for _, chans := range d.subs {
			for _, ch := range chans { close(ch) }
		}
		d.mu.RUnlock()
	}()
	return d
}

func (d *Demux) Subscribe(table string, buf int) <-chan Mixed {
	ch := make(chan Mixed, buf)
	d.mu.Lock()
	d.subs[table] = append(d.subs[table], ch)
	d.mu.Unlock()
	return ch
}
