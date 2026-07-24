package redisutil

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func StartExpirySubscriber(rdb *redis.Client, dbIndex int) {
	ctx := context.Background()
	channel := fmt.Sprintf("__keyevent@%d__:expired", dbIndex)
	pubsub := rdb.Subscribe(ctx, channel)
	go func() {
		defer pubsub.Close()
		ch := pubsub.Channel()
		log.Printf("expiry subscriber listening on %s", channel)
		for msg := range ch {
			HandleExpiredKey(ctx, rdb, msg.Payload)
		}
		log.Println("expiry subscriber channel closed")
	}()
}

func HandleExpiredKey(ctx context.Context, rdb *redis.Client, expiredKey string) {
	parts := strings.Split(expiredKey, ":")
	if len(parts) != 3 || parts[0] != "seat_lock" {
		return
	}
	showID, err := uuid.Parse(parts[1])
	if err != nil {
		return
	}
	seatID, err := uuid.Parse(parts[2])
	if err != nil {
		return
	}
	hashKey := fmt.Sprintf("show_seats:%s", showID.String())
	if err := rdb.HSet(ctx, hashKey, seatID.String(), "available").Err(); err != nil {
		log.Printf("warning: failed to update seat hash after expiry (show %s seat %s): %v", showID, seatID, err)
	}
}
