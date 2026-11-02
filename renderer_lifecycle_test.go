package tea

import (
	"io"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

type controlledRendererTicker struct {
	ch              chan time.Time
	entered         chan struct{}
	release         chan struct{}
	stops           atomic.Int32
	resets          atomic.Int32
	stopped         atomic.Bool
	resetBeforeStop atomic.Bool
}

func (t *controlledRendererTicker) ticks() <-chan time.Time { return t.ch }

func (t *controlledRendererTicker) Stop() {
	if t.stops.Add(1) == 1 {
		close(t.entered)
		<-t.release
		t.stopped.Store(true)
	}
}

func (t *controlledRendererTicker) Reset(time.Duration) {
	// The first Reset starts the injected ticker. Only restart must follow Stop.
	if t.resets.Add(1) > 1 && !t.stopped.Load() {
		t.resetBeforeStop.Store(true)
	}
}

type lifecycleRenderer struct {
	nilRenderer
	flushed chan struct{}
}

func (r lifecycleRenderer) flush(bool) error {
	r.flushed <- struct{}{}
	return nil
}

func TestRendererStopJoinsBeforeRestart(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ticker := &controlledRendererTicker{
			ch:      make(chan time.Time, 1),
			entered: make(chan struct{}),
			release: make(chan struct{}),
		}
		r := lifecycleRenderer{flushed: make(chan struct{}, 8)}
		p := NewProgram(&testModel{}, WithOutput(io.Discard))
		p.renderer, p.ticker = r, ticker
		p.startRenderer()
		restarted := make(chan struct{})
		go func() {
			p.stopRenderer(true)
			p.startRenderer()
			close(restarted)
		}()
		<-ticker.entered
		// Stop is held at an exact barrier. Wait establishes that every other
		// goroutine has blocked, rather than relying on scheduler timing.
		synctest.Wait()
		select {
		case <-restarted:
			t.Error("renderer restarted before the old ticker Stop completed")
		default:
		}
		if ticker.resets.Load() != 1 {
			t.Error("ticker Reset ran while the old Stop was blocked")
		}
		close(ticker.release)
		<-restarted
		synctest.Wait()
		if ticker.resetBeforeStop.Load() {
			t.Error("ticker Reset preceded completion of Stop")
		}
		if ticker.resets.Load() != 2 {
			t.Errorf("Reset calls = %d, want 2", ticker.resets.Load())
		}
		// Restart must also consume new ticks and flush, not merely return.
		ticker.ch <- time.Time{}
		synctest.Wait()
		select {
		case <-r.flushed:
		default:
			t.Error("restarted renderer did not flush an injected tick")
		}
		p.stopRenderer(true)
	})
}
