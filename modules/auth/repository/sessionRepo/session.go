package sessionRepo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type sessionStore struct {
	client *redis.Client
}

func (s *sessionStore) SetUserToken(ctx context.Context, userId int64, token, subToken string, expiredTime int) error {
	signature := strings.Split(token, ".")[2]
	key := fmt.Sprintf("/users/%d/session/%s", userId, subToken)

	bytes, err := json.Marshal(signature)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = s.client.Set(ctx, key, bytes, time.Duration(expiredTime)*time.Second).Err(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *sessionStore) GetUserToken(ctx context.Context, userId int64, subToken string) (string, error) {
	key := fmt.Sprintf("/users/%d/session/%s", userId, subToken)

	result, err := s.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", errors.New("not found")
	} else if err != nil {
		return "", errors.WithStack(err)
	}

	var token string
	if err = json.Unmarshal([]byte(result), &token); err != nil {
		return "", errors.WithStack(err)
	}

	return token, nil
}

func NewSessionStore(client *redis.Client) *sessionStore {
	return &sessionStore{client: client}
}
