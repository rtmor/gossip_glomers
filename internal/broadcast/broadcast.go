package broadcast

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"golang.org/x/sync/errgroup"
)

type IDs struct {
	data []int

	mu sync.RWMutex
}

type Topology struct {
	data map[string][]string

	mu sync.RWMutex
}

type Service struct {
	node *maelstrom.Node

	ids      *IDs
	topology *Topology
}

var ErrInvalidType = errors.New("request type invalid")

func NewBroadcastService() *Service {
	srv := &Service{
		node: maelstrom.NewNode(),
		ids: &IDs{
			data: []int{},
		},
		topology: &Topology{
			data: map[string][]string{},
		},
	}

	srv.node.Handle("broadcast", srv.BroadcastHandler)
	srv.node.Handle("read", srv.ReadHandler)
	srv.node.Handle("topology", srv.TopologyHandler)

	return srv
}

func (s *Service) Run() error {
	if err := s.node.Run(); err != nil {
		return fmt.Errorf("broadcast service failed: %w", err)
	}

	return nil
}

func (s *Service) BroadcastHandler(msg maelstrom.Message) error {
	var body map[string]any

	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return fmt.Errorf("failed to unmarshal incoming request: %w", err)
	}

	var reqMessage float64

	var ok bool
	if reqMessage, ok = body["message"].(float64); !ok {
		return ErrInvalidType
	}

	s.ids.mu.Lock()
	defer s.ids.mu.Unlock()

	s.ids.data = append(s.ids.data, int(reqMessage))

	b := broadcastAll(s.node)
	if err := b(msg.Src, body); err != nil {
		return err
	}

	err := s.node.Reply(msg, map[string]any{
		"type": "broadcast_ok",
	})
	if err != nil {
		return fmt.Errorf("message response failed: %w", err)
	}

	return nil
}

func (s *Service) ReadHandler(msg maelstrom.Message) error {
	var body map[string]any

	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return fmt.Errorf("failed to unmarshal incoming request: %w", err)
	}

	s.ids.mu.RLock()
	defer s.ids.mu.RUnlock()

	data := s.ids.data

	response := map[string]any{
		"type":     "read_ok",
		"messages": data,
	}

	if err := s.node.Reply(msg, response); err != nil {
		return fmt.Errorf("message response failed: %w", err)
	}

	return nil
}

func (s *Service) TopologyHandler(msg maelstrom.Message) error {
	req := struct {
		Topology map[string][]string `json:"topology"`
	}{}

	if err := json.Unmarshal(msg.Body, &req); err != nil {
		return fmt.Errorf("failed to unmarshal request: %w", err)
	}

	s.ids.mu.Lock()
	defer s.ids.mu.Unlock()

	s.topology.data = req.Topology

	res := map[string]any{
		"type": "topology_ok",
	}

	if err := s.node.Reply(msg, res); err != nil {
		return fmt.Errorf("message response failed: %w", err)
	}

	return nil
}

func broadcastAll(n *maelstrom.Node) func(src string, body map[string]any) error {
	return func(src string, body map[string]any) error {
		var g errgroup.Group

		for _, dst := range n.NodeIDs() {
			dst := dst
			if dst == src || n.ID() == dst {
				continue
			}

			g.Go(func() error {
				if err := n.Send(dst, body); err != nil {
					return fmt.Errorf("sending broadcast reply failed: %w", err)
				}
				return nil
			})

			if err := g.Wait(); err != nil {
				return err
			}
		}

		return nil
	}
}
