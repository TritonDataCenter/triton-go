package triton

import (
	"sync"

	"github.com/hashicorp/go-multierror"
)

type TritonStateBag interface {
	Get(string) interface{}
	GetOk(string) (interface{}, bool)
	Put(string, interface{})
	Remove(string)

	Client() *Client

	AppendError(error)
	ErrorsOrNil() error
}

// TritonStateBag implements StateBag by using a normal map underneath
// protected by a RWMutex.
type basicTritontateBag struct {
	TritonClient *Client

	errors *multierror.Error
	data   map[string]interface{}

	l    sync.RWMutex
	once sync.Once
}

func (b *basicTritontateBag) Client() *Client {
	b.l.RLock()
	defer b.l.RUnlock()

	return b.TritonClient
}

func (b *basicTritontateBag) AppendError(err error) {
	b.l.Lock()
	defer b.l.Unlock()

	b.errors = multierror.Append(b.errors, err)
}

func (b *basicTritontateBag) ErrorsOrNil() error {
	b.l.RLock()
	defer b.l.RUnlock()

	if b.errors == nil {
		return nil
	}

	return b.errors.ErrorOrNil()
}

func (b *basicTritontateBag) Get(k string) interface{} {
	result, _ := b.GetOk(k)
	return result
}

func (b *basicTritontateBag) GetOk(k string) (interface{}, bool) {
	b.l.RLock()
	defer b.l.RUnlock()

	result, ok := b.data[k]
	return result, ok
}

func (b *basicTritontateBag) Put(k string, v interface{}) {
	b.l.Lock()
	defer b.l.Unlock()

	// Make sure the map is initialized one time, on write
	b.once.Do(func() {
		b.data = make(map[string]interface{})
	})

	// Write the data
	b.data[k] = v
}

func (b *basicTritontateBag) Remove(k string) {
	b.l.Lock()
	defer b.l.Unlock()

	delete(b.data, k)
}
