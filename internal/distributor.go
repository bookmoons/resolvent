package internal

// Distributor is a query result distributor.
type Distributor interface {
	Distribute(response *dns.Msg, err error)
	Receive() (response *dns.Msg, err error)
}

type result struct {
	response *dns.Msg
	err      error
}

type distributor struct {
	semaphore Semaphore
}

func NewDistributor() *distributor {
	return &distributor{
		semaphore: make(chan chan result),
	}
}

func (d *distributor) Distribute(response *dns.Msg, err error) {
	for receive := range d.semaphore {
		value := result{
			response: response.Copy(),
			err:      err,
		}
		receive <- value
	}
}

func (d *distributor) Receive() (response *dns.Msg, err error) {
	receive := make(chan result)
	d.semaphore <- receive
	value := <-receive
	return value.response, value.err
}
