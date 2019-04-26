package redis

import (
	// "reflect"
	"math"
	"fmt"
	"strconv"
	// "encoding/json"
	"github.com/garyburd/redigo/redis"
	_ "sync"
	"time"
)

//types
type (
	Stream struct {
		HitTime		int        `bson:"hittime" json:"hittime"`
		ResponseTime	int        `bson:"responsetime" json:"responsetime"`
	}

	PoolConn struct {
		pool *redis.Pool
	}
)

func New() *PoolConn {
	fmt.Println("New")//
	return &PoolConn{
		pool: &redis.Pool{
			MaxIdle:     80,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", "localhost:6379")
				if err != nil {
					return nil, err
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		},
	}
}

func (c *PoolConn) Add(key int, receiver int) error {
	// fmt.Println("Adding")//
	// var err error
	// key, err := json.Marshal(stream)
	// if err != nil {
	// 	return err
	// }
	conn := c.pool.Get()
	defer conn.Close()
	
	// Search check
	values, err := redis.Strings(conn.Do("LRANGE", key, 0, -1))
	if err != nil {
		return err
	} else {
		if len(values) == 0 {
			if _, err := conn.Do("LPUSH", key, receiver); err != nil {
				return err
			}
			return nil
		}
		// fmt.Println(values[0], reflect.TypeOf(values[0]))
		if s, err := strconv.ParseFloat(values[0], 64); err == nil {
			if _, err := conn.Do("LREM", key, -1, s); err != nil {
				return err
			}
			if _, err := conn.Do("LPUSH", key, math.Min(float64(receiver), s)); err != nil {
				return err
			}
		}
		return nil
	}
	// if _, err := conn.Do("LPUSH", key, math.Min(float64(receiver), float64(values))); err != nil {
	// // if _, err := conn.Do("SET", key, receiver); err != nil {
	// 	return err
	// }
	return nil
}

func (c *PoolConn) Get(hitpoint int) ([]string, error) {
	// out, err := json.Marshal(hitpoint)
	// if err != nil {
	// 	return nil, err
	// }
	conn := c.pool.Get()
	defer conn.Close()
	values, err := redis.Strings(conn.Do("LRANGE", hitpoint, 0, -1))
	// values, err := redis.Strings(conn.Do("GET", string(hitpoint)))
	if err != nil {
		return nil, err
	}
	return values, nil
}