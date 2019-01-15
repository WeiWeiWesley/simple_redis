package redis

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/garyburd/redigo/redis"
)

//P connection pool
type P struct {
	Pool *redis.Pool
}

func newPool(server string) *redis.Pool {

	return &redis.Pool{
		Wait:        true,
		MaxIdle:     20,
		MaxActive:   2000,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

// AddConnect new connection
func AddConnect(host string) *P {
	// fmt.Println("new connection", host)
	return &P{Pool: newPool(host)}
}

//CloseConn close connection
func CloseConn(pool *redis.Pool) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		pool.Close()
		os.Exit(0)
	}()
}
