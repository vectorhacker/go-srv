package srv

import (
	"net"

	"sync"
	"time"

	"fmt"

	"google.golang.org/grpc/naming"
)

type watcher struct {
	target   string
	existing []*net.SRV
	m        sync.RWMutex
	stopChan chan bool
	errChan  chan error
}

// NewWatcher returns a naming.Watcher that watches for changes on the SRV
// records of the target address string
func NewWatcher(target string) naming.Watcher {
	w := &watcher{
		stopChan: make(chan bool),
		errChan:  make(chan error, 10),
		target:   target,
		existing: []*net.SRV{},
	}

	w.start()

	return w
}

func (w *watcher) start() {
	go func() {
		for {
			select {
			case <-time.After(5 * time.Second):
				_, addrs, err := net.LookupSRV("", "", w.target)
				if err != nil {
					w.errChan <- err
					continue
				}

				func() {
					w.m.Lock()
					defer w.m.Unlock()

					w.existing = addrs
				}()
			case <-w.stopChan:
				return
			}
		}
	}()
}

// Next returns the next set up *naming.Updates or an error
func (w *watcher) Next() ([]*naming.Update, error) {
	select {
	case err := <-w.errChan:
		return nil, err
	default:
		w.m.RLock()
		defer w.m.RUnlock()

		updates := []*naming.Update{}

		for _, addr := range w.existing {

			updates = append(updates, &naming.Update{
				Addr: fmt.Sprintf("%s:%d", w.target, addr.Port),
				Op:   naming.Add,
			})
		}

		return updates, nil
	}
}

// Close stops the Watcher's internal loop
func (w *watcher) Close() {
	w.stopChan <- true
}
