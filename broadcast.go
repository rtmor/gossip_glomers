package gglomers

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type IDs struct {
	data  []int
	cache *nodeCache

	sync.RWMutex
}

type Topology struct {
	data map[string][]string

	sync.RWMutex
}

type BroadcastService struct {
	node *maelstrom.Node

	ids      *IDs
	topology *Topology
}

func NewBroadcastService() *BroadcastService {
	srv := &BroadcastService{
		node: maelstrom.NewNode(),
		ids: &IDs{
			data:    []int{},
			cache:   &nodeCache{},
			RWMutex: sync.RWMutex{},
		},
		topology: &Topology{
			data:    map[string][]string{},
			RWMutex: sync.RWMutex{},
		},
	}

	srv.node.Handle("broadcast", srv.BroadcastHandler)
	srv.node.Handle("read", srv.ReadHandler)
	srv.node.Handle("topology", srv.TopologyHandler)

	return srv
}

func (s *BroadcastService) Run() error {
	if err := s.node.Run(); err != nil {
		return fmt.Errorf("broadcast service failed: %w", err)
	}

	return nil
}

var ErrInvalidType = errors.New("request type invalid")

type nodeCache map[int]struct{}

func (c *nodeCache) Clear() {
	*c = map[int]struct{}{}
}

func (c *nodeCache) Append(nodeID int) {
	(*c)[nodeID] = struct{}{}
}

func (s *BroadcastService) BroadcastHandler(msg maelstrom.Message) error {
	var body map[string]any

	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return fmt.Errorf("failed to unmarshal incoming request: %w", err)
	}

	var reqMessage float64

	var ok bool
	if reqMessage, ok = body["message"].(float64); !ok {
		return ErrInvalidType
	}

	s.ids.Lock()
	s.ids.data = append(s.ids.data, int(reqMessage))
	s.ids.Unlock()

	err := s.node.Reply(msg, map[string]any{
		"type": "broadcast_ok",
	})
	if err != nil {
		return fmt.Errorf("message response failed: %w", err)
	}

	return nil
}

func (s *BroadcastService) ReadHandler(msg maelstrom.Message) error {
	var body map[string]any

	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return fmt.Errorf("failed to unmarshal incoming request: %w", err)
	}

	s.ids.RLock()
	data := s.ids.data
	s.ids.RUnlock()

	response := map[string]any{
		"type":     "read_ok",
		"messages": data,
	}

	if err := s.node.Reply(msg, response); err != nil {
		return fmt.Errorf("message response failed: %w", err)
	}

	return nil
}

func (s *BroadcastService) TopologyHandler(msg maelstrom.Message) error {
	req := struct {
		Topology map[string][]string `json:"topology"`
	}{}

	if err := json.Unmarshal(msg.Body, &req); err != nil {
		return fmt.Errorf("failed to unmarshal request: %w", err)
	}

	s.ids.Lock()
	s.topology.data = req.Topology
	s.ids.Unlock()

	res := map[string]any{
		"type": "topology_ok",
	}

	if err := s.node.Reply(msg, res); err != nil {
		return fmt.Errorf("message response failed: %w", err)
	}

	return nil
}
