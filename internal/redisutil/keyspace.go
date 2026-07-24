package redisutil

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func EnableKeySpaceNotifications(rdb *redis.Client) error {
	ctx := context.Background()
	err := rdb.ConfigSet(ctx, "notify-keyspace-events", "Ex").Err()
	if err != nil {
		return fmt.Errorf("failed to enable keyspace notifications: %w", err)
	}
	return nil
}
