package control

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// Registry is a registrar of consumers & subscribers
type Registry struct {
	consumers   map[string]consumer
	subscribers map[string][]subscriber
}

// ControlService is the main object that manages the control service. It is responsible for fetching
// and caching control data, and updating consumers and subscribers.
type ControlService struct {
	Registry
	logger          log.Logger
	cancel          context.CancelFunc
	requestInterval time.Duration
	data            dataProvider
	lastFetched     map[string]string
}

// consumer is an interface for something that consumes control server data updates.
type consumer interface {
	Update(io.Reader)
}

// subscriber is an interface for something that wants to be notified when a subsystem has been updated.
type subscriber interface {
	Ping()
}

// dataProvider is an interface for something that can retrieve control data. Authentication, HTTP,
// file system access, etc. lives below this abstraction layer.
type dataProvider interface {
	Get(resource, cachedETag string) (etag string, data io.Reader, err error)
}

func NewControlService(logger log.Logger, data dataProvider, opts ...Option) *ControlService {
	r := Registry{
		consumers:   make(map[string]consumer),
		subscribers: make(map[string][]subscriber),
	}
	cs := &ControlService{
		Registry:        r,
		logger:          logger,
		requestInterval: 60 * time.Second,
		data:            data,
		lastFetched:     make(map[string]string),
	}

	for _, opt := range opts {
		opt(cs)
	}

	return cs
}

func (cs *ControlService) Start(ctx context.Context) {
	ctx, cs.cancel = context.WithCancel(ctx)
	requestTicker := time.NewTicker(cs.requestInterval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-requestTicker.C:
			cs.Fetch()
		}
	}
}

func (cs *ControlService) Stop() {
	cs.cancel()
}

// controlResponse is the payload received from the control server
// type controlResponse struct {
// 	// TODO: This is a temporary and simple data format for phase 1
// 	Message string `json:"message,omitempty"`
// 	Err     string `json:"error,omitempty"`
// }

// Performs a retrieval of the latest control server data, and notifies observers of updates.
func (cs *ControlService) Fetch() error {
	_, data, err := cs.data.Get("", "")
	if err != nil {
		return fmt.Errorf("getting subsystems map: %w", err)
	}

	var subsystems map[string]string
	if err := json.NewDecoder(data).Decode(&subsystems); err != nil {
		return fmt.Errorf("decoding subsystems map: %w", err)
	}

	for subsystem, resource := range subsystems {
		cachedETag := cs.lastFetched[subsystem]
		etag, data, err := cs.data.Get(resource, cachedETag)
		if err != nil {
			return fmt.Errorf("failed to get control data: %w", err)
		}

		if cachedETag != "" && etag == cachedETag {
			// Nothing to do, skip to the next subsystem
			continue
		}

		// Consumer and subscribers notified now
		cs.update(subsystem, data)

		// Cache the last fetched version of this subsystem's data
		cs.lastFetched[subsystem] = etag
	}

	level.Debug(cs.logger).Log("msg", "control data fetch complete")

	return nil
}

func (r *Registry) RegisterConsumer(subsystem string, consumer consumer) error {
	if _, ok := r.consumers[subsystem]; ok {
		return fmt.Errorf("consumer already registered for subsystem %s", subsystem)
	}
	r.consumers[subsystem] = consumer
	return nil
}

func (r *Registry) RegisterSubscriber(subsystem string, subscriber subscriber) {
	r.subscribers[subsystem] = append(r.subscribers[subsystem], subscriber)
}

func (r *Registry) update(subsystem string, reader io.Reader) {
	// First, send to consumer, if any
	if consumer, ok := r.consumers[subsystem]; ok {
		consumer.Update(reader)
	}

	// Then send a ping to all subscribers
	for _, subscriber := range r.subscribers[subsystem] {
		subscriber.Ping()
	}
}
