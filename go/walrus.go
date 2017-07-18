/******************************************************************************
 *
 * Walrus Test Framework
 * Copyright Ryan Goodfellow 2017, all rights reserved
 * Apache 2.0
 *
 *****************************************************************************/

package walrustf

import (
	"fmt"
	"github.com/go-redis/redis"
)

type Client struct {
	collector, test, participant string
	conn                         *redis.Client
	counter                      int
}

func NewClient(collector, test, participant string) (*Client, error) {
	c := &Client{
		collector:   collector,
		test:        test,
		participant: participant,
		conn: redis.NewClient(
			&redis.Options{Addr: fmt.Sprintf("%s:6379", collector)}),
	}
	_, err := c.conn.Ping().Result()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) Error(format string, args ...interface{}) error {
	return c.report("error", format, args...)
}

func (c *Client) Warning(format string, args ...interface{}) error {
	return c.report("warning", format, args...)
}

func (c *Client) Ok(format string, args ...interface{}) error {
	return c.report("ok", format, args...)
}

func (c *Client) report(level, format string, args ...interface{}) error {

	_, err := c.conn.Ping().Result()
	if err != nil {
		return err
	}

	t, err := c.conn.Time().Result()
	if err != nil {
		return err
	}

	msg := fmt.Sprintf(format, args...)

	key := fmt.Sprintf(
		"%s:%s:%d:%d",
		c.test,
		c.participant,
		t.Unix(),
		t.UnixNano(),
	)

	value := fmt.Sprintf(
		"%s:::%s",
		level,
		msg)

	err = c.conn.Set(key, value, 0).Err()
	if err != nil {
		return nil
	}

	return nil

}
