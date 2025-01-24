package websocket

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"net/http"
)

type Discover interface {
	Register(addr string) error
	BoundUser(uid string) error
	RelieveUser(uid string) error
	Transpond(msg interface{}, uid ...string) error
}

type nopDiscover struct {
	serverAddr string
}

func (d *nopDiscover) Register(addr string) error {
	return nil
}

func (d *nopDiscover) BoundUser(uid string) error {
	return nil
}

func (d *nopDiscover) RelieveUser(uid string) error {
	return nil
}

func (d *nopDiscover) Transpond(msg interface{}, uid ...string) error {
	return nil
}

type redisDiscover struct {
	addr         string
	auth         http.Header
	srvKey       string
	boundUserKey string
	redis        *redis.Redis
	clients      map[string]Client
}

func NewRedisDiscover(auth http.Header, srvKey string, redisCfg redis.RedisConf) *redisDiscover {
	return &redisDiscover{
		auth:         auth,
		srvKey:       srvKey,
		boundUserKey: fmt.Sprintf("%s.%s", srvKey, "boundUserKey"),
		redis:        redis.MustNewRedis(redisCfg),
		clients:      make(map[string]Client),
	}
}

func (d *redisDiscover) Register(addr string) error {
	d.addr = addr
	go d.redis.Set(d.srvKey, addr)
	return nil
}

func (d *redisDiscover) BoundUser(uid string) error {
	hexists, err := d.redis.Hexists(d.boundUserKey, uid)
	if err != nil {
		return err
	}
	if hexists {
		return nil
	}

	// 绑定
	return d.redis.Hset(d.boundUserKey, uid, d.addr)
}

func (d *redisDiscover) RelieveUser(uid string) error {
	_, err := d.redis.Hdel(d.boundUserKey, uid)
	return err
}

func (d *redisDiscover) Transpond(msg interface{}, uids ...string) error {
	for _, uid := range uids {
		addr, err := d.redis.Hget(d.boundUserKey, uid)
		if err != nil {
			return err
		}
		srvClient, ok := d.clients[addr]
		if !ok {
			srvClient = d.createClient(addr)
			d.clients[addr] = srvClient
		}

		fmt.Println("redis transpond -》 ", addr, " uid ", uid)

		if err := d.send(srvClient, msg, uid); err != nil {
			return err
		}
	}
	return nil
}

func (d *redisDiscover) send(client Client, msg interface{}, uid string) error {
	return client.Send(Message{
		FrameType:    FrameTranspond,
		Data:         msg,
		TranspondUid: uid,
	})
}

func (d *redisDiscover) createClient(addr string) Client {
	return NewClient(addr, WithHeader(d.auth))
}
