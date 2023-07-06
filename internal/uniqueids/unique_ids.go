package uniqueids

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Service struct {
	n *maelstrom.Node
}

func NewUniqueIDService() *Service {
	node := maelstrom.NewNode()

	node.Handle("generate", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return fmt.Errorf("failed to unmarshal request: %w", err)
		}

		body["type"] = "generate_ok"
		body["id"] = uuid.New().ID()

		if err := node.Reply(msg, body); err != nil {
			return fmt.Errorf("message reply failed: %w", err)
		}

		return nil
	})

	return &Service{
		n: node,
	}
}

func (s *Service) Run() error {
	if err := s.n.Run(); err != nil {
		return fmt.Errorf("uniqueIDService failed: %w", err)
	}

	return nil
}
