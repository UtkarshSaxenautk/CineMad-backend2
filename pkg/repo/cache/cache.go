package cache

import (
	"authentication-ms/pkg/svc"
	"context"
	"encoding/json"
	"github.com/allegro/bigcache/v3"
	"log"
	"time"
)

const (
	Life = 2
)

type Otp struct {
	otp string
}

type cache struct {
	bigCache *bigcache.BigCache
}

func NewCache(ctx context.Context) svc.Cache {

	bigCache, err := bigcache.New(ctx, bigcache.DefaultConfig(time.Minute*Life))
	if err != nil {
		log.Println("error in creating cache : ", err)
	}
	return &cache{bigCache: bigCache}
}

func (c *cache) SetInCache(email string, otp string) error {
	otpBytes, err := json.Marshal(otp)
	if err != nil {
		log.Println("error in marshaling otp : ", err)
	}
	err = c.bigCache.Set(email, otpBytes)
	if err != nil {
		log.Println("error in setting value in cache : ", err)
	}
	return nil
}

func (c *cache) GetFromCache(email string) (string, error) {
	otp, err := c.bigCache.Get(email)
	if err != nil {
		log.Println("error in getting corresponding otp of email : ", err)
		return "", err
	}
	var resOtp Otp
	err = json.Unmarshal(otp, &resOtp)
	if err != nil {
		log.Println("error in unmarshalling otp from cache : ", err)
	}
	return resOtp.otp, nil
}
