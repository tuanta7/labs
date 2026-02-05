package redis

import (
	"context"

	goredis "github.com/redis/go-redis/v9"
)

type Location struct {
	Name      string
	Longitude float64
	Latitude  float64
	GeoHash   int64
}

type GeoCache interface {
	GeoAdd(ctx context.Context, key, member string, longitude, latitude float64) error
	GeoPos(ctx context.Context, key string, members ...string) ([]*Location, error)
	GeoSearch(ctx context.Context, longitude, latitude float64, radius float64) ([]*Location, error)
}

func (c *Client) GeoAdd(ctx context.Context, key, member string, longitude, latitude float64) error {
	return c.redisClient.GeoAdd(ctx, key, &goredis.GeoLocation{
		Name:      member,
		Longitude: longitude,
		Latitude:  latitude,
	}).Err()
}

func (c *Client) GeoPos(ctx context.Context, key string, members ...string) ([]*Location, error) {
	pos, err := c.redisClient.GeoPos(ctx, key, members...).Result()
	if err != nil {
		return nil, err
	}

	locations := make([]*Location, len(pos))
	for i, p := range pos {
		if p == nil {
			continue
		}

		locations[i] = &Location{
			Name:      members[i],
			Longitude: p.Longitude,
			Latitude:  p.Latitude,
		}
	}

	return locations, nil
}

func (c *Client) GeoSearch(ctx context.Context, key string, longitude, latitude, radius float64) ([]*Location, error) {
	loc, err := c.redisClient.GeoSearchLocation(ctx, key, &goredis.GeoSearchLocationQuery{
		GeoSearchQuery: goredis.GeoSearchQuery{
			Radius:    radius,
			Longitude: longitude,
			Latitude:  latitude,
		},
		WithHash: true, // GeoLocation's GeoHash is used here!
	}).Result()
	if err != nil {
		return nil, err
	}

	var locations []*Location
	for _, l := range loc {
		locations = append(locations, &Location{
			Name:      l.Name,
			Longitude: l.Longitude,
			Latitude:  l.Latitude,
			GeoHash:   l.GeoHash,
		})
	}

	return locations, nil
}
