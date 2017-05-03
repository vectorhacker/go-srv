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
	existing map[string]int
	previous map[string]int
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
		previous: map[string]int{},
		existing: map[string]int{},
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

				hosts, err := net.LookupHost(w.target)

				func() {
					w.m.Lock()
					defer w.m.Unlock()
					w.previous = w.existing

					newExisting := map[string]int{}

					for index, addr := range addrs {
						host := hosts[index]
						newExisting[host] = int(addr.Port)
					}

					w.existing = newExisting
					
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

		for addr, port := range w.previous {
			if _, ok := w.existing[addr]; !ok {
				updates = append(updates, &naming.Update{
					Op:   naming.Delete,
					Addr: formatAddress(addr, port),
				})
			}
		}

		for addr, port := range w.existing {
			if _, ok := w.previous[addr]; !ok {
				updates = append(updates, &naming.Update{
					Op:   naming.Add,
					Addr: formatAddress(addr, port),
				})
			}
		}

		return updates, nil
	}
}

// Close stops the Watcher's internal loop
func (w *watcher) Close() {
	w.stopChan <- true
}

func formatAddress(addr string, port int) string {
	return fmt.Sprintf("%s:%d", addr, port)
}
