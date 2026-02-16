package location

import (
	"context"

	"github.com/tuanta7/k6noz/services/internal/domain"
	"github.com/tuanta7/k6noz/services/pkg/redis"
)

type Cache struct {
	redis redis.GeoCache
}

func NewCache(cache redis.GeoCache) *Cache {
	return &Cache{redis: cache}
}

func (c *Cache) GetLocation(ctx context.Context, driverID string) (*domain.Location, error) {
	locations, err := c.redis.GeoPos(ctx, driverID)
	if err != nil {
		return nil, err
	}

	return &domain.Location{
		Latitude:  locations[0].Latitude,
		Longitude: locations[0].Longitude,
	}, nil
}

func (c *Cache) SetLocation(ctx context.Context, driverID string, location domain.Location) error {
	return c.redis.GeoAdd(ctx, "drivers:geo", driverID, location.Longitude, location.Latitude)
}
