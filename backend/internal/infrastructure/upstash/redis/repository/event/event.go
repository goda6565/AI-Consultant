package event

import (
	"context"
	"fmt"
	"time"

	actionValue "github.com/goda6565/ai-consultant/backend/internal/domain/action/value"
	"github.com/goda6565/ai-consultant/backend/internal/domain/event/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/event/repository"
	eventValue "github.com/goda6565/ai-consultant/backend/internal/domain/event/value"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/redis/go-redis/v9"
)

type RedisEventRepository struct {
	client *redis.Client
}

func NewRedisEventRepository(client *redis.Client) repository.EventRepository {
	return &RedisEventRepository{
		client: client,
	}
}

func (c *RedisEventRepository) Create(ctx context.Context, event *entity.Event) error {
	return c.client.XAdd(ctx, &redis.XAddArgs{
		Stream: fmt.Sprintf("stream:%s", event.ProblemID.Value()),
		Values: map[string]interface{}{
			"eventId":    event.ID,
			"eventType":  event.EventType,
			"actionType": event.ActionType,
			"message":    event.Message,
		},
	}).Err()
}

func (c *RedisEventRepository) FindAllByProblemID(ctx context.Context, problemID sharedValue.ID) ([]entity.Event, error) {
	res, err := c.client.XRange(ctx, fmt.Sprintf("stream:%s", problemID.Value()), "0", "+").Result()
	if err != nil {
		return nil, err
	}

	events := make([]entity.Event, 0, len(res))
	for _, msg := range res {
		newEvent, err := msgToEntity(msg, problemID)
		if err != nil {
			return nil, err
		}
		events = append(events, *newEvent)
	}
	return events, nil
}

func (c *RedisEventRepository) FindAllByProblemIDAsStream(ctx context.Context, problemID sharedValue.ID) (<-chan entity.Event, error) {
	ch := make(chan entity.Event)

	go func() {
		defer close(ch)
		lastID := "$"

		for {
			select {
			case <-ctx.Done():
				return
			default:
				res, err := c.client.XRead(ctx, &redis.XReadArgs{
					Streams: []string{fmt.Sprintf("stream:%s", problemID.Value()), lastID},
					Block:   5 * time.Second,
					Count:   100,
				}).Result()

				if err != nil {
					if err == redis.Nil {
						continue
					}
					time.Sleep(500 * time.Millisecond)
					continue
				}

				for _, s := range res {
					for _, msg := range s.Messages {
						newEvent, err := msgToEntity(msg, problemID)
						if err != nil {
							continue
						}
						ch <- *newEvent
						lastID = msg.ID
					}
				}
			}
		}
	}()

	return ch, nil
}

func (c *RedisEventRepository) DeleteAllByProblemID(ctx context.Context, problemID sharedValue.ID) (numDeleted int64, err error) {
	return c.client.Del(ctx, fmt.Sprintf("stream:%s", problemID.Value())).Result()
}

func msgToEntity(msg redis.XMessage, problemID sharedValue.ID) (*entity.Event, error) {
	id, err := sharedValue.NewID(getStringValue(msg.Values["eventId"]))
	if err != nil {
		return nil, err
	}

	eventType, err := eventValue.NewEventType(getStringValue(msg.Values["eventType"]))
	if err != nil {
		return nil, err
	}
	actionType, err := actionValue.NewActionType(getStringValue(msg.Values["actionType"]))
	if err != nil {
		return nil, err
	}
	message, err := eventValue.NewMessage(getStringValue(msg.Values["message"]))
	if err != nil {
		return nil, err
	}
	return entity.NewEvent(id, problemID, eventType, actionType, *message), nil
}

func getStringValue(v interface{}) string {
	return fmt.Sprintf("%v", v)
}
