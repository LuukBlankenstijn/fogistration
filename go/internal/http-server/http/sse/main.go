package sse

import (
	"context"
	"net/http"
	"sync"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/sse"
)

type SSEManager struct {
	mu       sync.RWMutex
	messages map[string]any
	subs     map[int]chan any
	nextID   int
}

func New() *SSEManager {
	return &SSEManager{
		mu:       sync.RWMutex{},
		messages: make(map[string]any),
		subs:     make(map[int]chan any),
		nextID:   0,
	}
}

func (s *SSEManager) subscribe() (int, chan any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.nextID
	s.nextID++
	ch := make(chan any, 32)
	if s.subs == nil {
		s.subs = make(map[int]chan any)
	}
	s.subs[id] = ch
	return id, ch
}
func (s *SSEManager) unsubscribe(id int) {
	s.mu.Lock()
	ch := s.subs[id]
	delete(s.subs, id)
	s.mu.Unlock()
	close(ch)
}

func (s *SSEManager) Broadcast(msg any) {
	s.mu.RLock()
	// snapshot to avoid holding lock while sending
	subs := make([]chan any, 0, len(s.subs))
	for _, ch := range s.subs {
		subs = append(subs, ch)
	}
	s.mu.RUnlock()

	for _, ch := range subs {
		select {
		case ch <- msg:
		default: /* drop if slow */
		}
	}
}

func (s *SSEManager) setMessage(path string, message any) {
	s.messages[path] = message
}

func (s *SSEManager) CreateEndpoint(api huma.API) {
	message := s.messages
	message["initMessage"] = struct{}{}

	sse.Register(api, huma.Operation{
		OperationID: "sse", Method: http.MethodGet, Path: "/api/sse", Summary: "Server-sent events",
	}, s.messages, func(ctx context.Context, _ *struct{}, send sse.Sender) {
		id, ch := s.subscribe()
		defer s.unsubscribe(id)
		// send initial data to open the channel
		send.Data(struct{}{})
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-ch:
				if err := send.Data(msg); err != nil {
					return
				}
			}
		}
	})
}

func Register[I any, O any](s *SSEManager, api huma.API, op huma.Operation, handler func(context.Context, *I) (*GetResponse[O], error)) {
	huma.Register(api, op, handler)

	if op.Method != http.MethodGet {
		logging.Warn("tried to register sse with a endpoint that is not get")
		return
	}

	update := SSEUpdate[O]{}
	s.setMessage(op.OperationID, update)
}
