package clickhouse

import (
	"github.com/ClickHouse/clickhouse-go/v2"
)

type Client struct {
	clickhouse.Conn
}

func NewClient(addr []string) (*Client, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: addr,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		Conn: conn,
	}, nil
}

func (c *Client) Close() error {
	return c.Conn.Close()
}
