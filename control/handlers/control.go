package handlers

import (
	r "github.com/fercevik729/STLKER/control/repository"
	"log/slog"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"

	pb "github.com/fercevik729/STLKER/grpc/protos"
)

// ControlHandler is a http.Handler
type ControlHandler struct {
	l        *slog.Logger
	client   pb.WatcherClient
	cache    *cache.Cache
	portRepo *r.PortfolioRepository
	secRepo  *r.SecurityRepository
	userRepo *r.UserRepository
}

// NewControlHandler constructs a new ControlHandler struct
func NewControlHandler(log *slog.Logger, wc pb.WatcherClient, rOptions *redis.Ring,
	portRepo *r.PortfolioRepository, secRepo *r.SecurityRepository, userRepo *r.UserRepository) *ControlHandler {
	// Check if redis options were presented
	var c *cache.Cache
	if rOptions != nil {
		c = cache.New(&cache.Options{
			Redis:      rOptions,
			LocalCache: cache.NewTinyLFU(1000, time.Minute)})
	} else {
		c = nil
	}
	return &ControlHandler{
		l:        log,
		client:   wc,
		cache:    c,
		portRepo: portRepo,
		secRepo:  secRepo,
		userRepo: userRepo,
	}
}
