package session

import (
	"errors"
	"fmt"

	"github.com/varunmbet/goblack"
	"github.com/varunmbet/redis"
)

// RedisProvider redis session provider
type RedisProvider struct {
	goblack.IProvider
	//Client       *redis.Cacher
	maxLifeTime  int64
	idLength     int
	cookieOption goblack.CookieOptions
	redisOption  RedisOptions
}

type RedisOptions struct {
	Client *redis.Cacher
	Prefix string
}

// Init initialize provider
// Init initialize provider
func (p *RedisProvider) Init(options goblack.SessionOptions) error {
	if options != (goblack.SessionOptions{}) {
		p.cookieOption = goblack.NewCookieOption(options.Cookieoption)
		if options.Session_lifetime > 0 {
			p.maxLifeTime = options.Session_lifetime
		} else {
			p.maxLifeTime = 60 * 30
		}
		if options.IDLength > 0 {
			p.idLength = options.IDLength
		} else {
			p.idLength = 16
		}
		if options.Providerconfig != nil {
			p.redisOption = options.Providerconfig.(RedisOptions)
			if p.redisOption.Prefix == "" {
				p.redisOption.Prefix = "token:"
			}
		} else {
			fmt.Printf("请设置redis session的redis")
		}

	}
	return nil
}

// Exist check session id is exist
func (p *RedisProvider) Exist(sid string) bool {
	has, err := p.redisOption.Client.Exists(sid)
	if err != nil {
		return false
	}
	return has
}

// Start begin session manager
func (p *RedisProvider) Start(c *goblack.Context, withsid ...string) (*goblack.Session, error) {
	//var session *goblack.Session
	if len(withsid) > 0 {
		sessid := withsid[0]
		session, err := p.Load(sessid)
		if err != nil {
			return nil, err
		}
		return session, nil
	}
	sid := c.GetHeader("x-token") //c.GetCookie(p.cookieOption.Cookie_key)
	if len(sid) > 0 && p.Exist(p.redisOption.Prefix+sid) {
		return p.Load(sid)
	}

	return nil, errors.New("没有生成token")
	/*	sid = goblack.CreateSessionID(p.idLength)
		session, err := p.Load(p.redisOption.Prefix + sid)
		if err != nil {
			return nil, err
		}
	*/
	//c.SetHeader("newtoken", sid)
	/*
			c.SetCookie(
				p.cookieOption.Cookie_key,
				sid,
				p.cookieOption.LifeTime,
				p.cookieOption.Path,
				p.cookieOption.Domain,
				p.cookieOption.Secure,
				p.cookieOption.HttpOnly,
			)
		p.ResetTtl(p.redisOption.Prefix + sid)

		return session, nil*/
}

func (p *RedisProvider) ResetTtl(sid string) {
	p.redisOption.Client.Expire(sid, p.maxLifeTime)
}

// Read read session data from provider
func (p *RedisProvider) Load(sid string) (*goblack.Session, error) {
	sessid := p.redisOption.Prefix + sid
	if !p.Exist(sessid) {
		if err := p.redisOption.Client.Set(sessid, "", p.maxLifeTime); err != nil {
			return nil, err
		}
	}
	p.ResetTtl(sessid)

	var data map[interface{}]interface{}
	raw, err := p.redisOption.Client.GetString(sessid)
	if err != nil {
		return nil, err
	}

	if len(raw) == 0 {
		data = make(map[interface{}]interface{})
	} else {
		data, err = goblack.DecodeGob([]byte(raw))
		if err != nil {
			return nil, err
		}
	}

	return goblack.NewSession(p, sid, data)
}

// Write write session data to provider
func (p *RedisProvider) Save(sid string, data map[interface{}]interface{}) error {
	encoded, err := goblack.EncodeGob(data)
	if err != nil {
		return err
	}
	return p.redisOption.Client.Set(p.redisOption.Prefix+sid, string(encoded), p.maxLifeTime)
}

// Destroy destroy session id from provider
func (p *RedisProvider) Destroy(sid string) error {
	return p.redisOption.Client.Del(p.redisOption.Prefix + sid)
}

func NewRedisProvider() *RedisProvider {
	return &RedisProvider{}
}
