package handlers

import (
	"log"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"

	pb "github.com/fercevik729/STLKER/eagle/protos"
)

// ControlHandler is a http.Handler
type ControlHandler struct {
	l      *log.Logger
	client pb.WatcherClient
	cache  *cache.Cache
	dbName string
}

// NewControlHandler is a constructor
func NewControlHandler(log *log.Logger, wc pb.WatcherClient, rOptions *redis.Ring, db string) *ControlHandler {
	return &ControlHandler{
		l:      log,
		client: wc,
		cache: cache.New(&cache.Options{
			Redis:      rOptions,
			LocalCache: cache.NewTinyLFU(1000, time.Minute),
		}),
		dbName: db,
	}
}
