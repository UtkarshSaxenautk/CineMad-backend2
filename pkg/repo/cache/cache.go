package cache

import (
	"authentication-ms/pkg/svc"
	"context"
	"encoding/json"
	"fmt"
	"github.com/allegro/bigcache/v3"
	"log"
	"time"
)

const (
	Life = 2
)

type cache struct {
	bigCache *bigcache.BigCache
	jwtCache *bigcache.BigCache
}

func NewCache(ctx context.Context) svc.Cache {
	config1 := bigcache.Config{
		// number of shards (must be a power of 2)
		Shards: 1024,
		// time after which entry can be evicted
		LifeWindow: Life * time.Minute,
		// rps * lifeWindow, used only in initial memory allocation
		MaxEntriesInWindow: 1000 * 10 * 60,
		// max entry size in bytes, used only in initial memory allocation
		MaxEntrySize: 500,
		// prints information about additional memory allocation
		Verbose: true,
		// cache will not allocate more memory than this limit, value in MB
		// if value is reached then the oldest entries can be overridden for the new ones
		// 0 value means no size limit
		HardMaxCacheSize: 8192,
		// callback fired when the oldest entry is removed because of its
		// expiration time or no space left for the new entry. Default value is nil which
		// means no callback and it prevents from unwrapping the oldest entry.
	}
	config2 := bigcache.Config{
		// number of shards (must be a power of 2)
		Shards: 1024,
		// time after which entry can be evicted
		LifeWindow: Life * time.Hour * 24,
		// rps * lifeWindow, used only in initial memory allocation
		MaxEntriesInWindow: 1000 * 10 * 60,
		// max entry size in bytes, used only in initial memory allocation
		MaxEntrySize: 500,
		// prints information about additional memory allocation
		Verbose: true,
		// cache will not allocate more memory than this limit, value in MB
		// if value is reached then the oldest entries can be overridden for the new ones
		// 0 value means no size limit
		HardMaxCacheSize: 8192,
		// callback fired when the oldest entry is removed because of its
		// expiration time or no space left for the new entry. Default value is nil which
		// means no callback and it prevents from unwrapping the oldest entry.
	}
	bigCache, err := bigcache.New(ctx, config1)

	if err != nil {
		log.Println("error in creating cache : ", err)
	}

	jwtCache, err := bigcache.New(ctx, config2)
	if err != nil {
		log.Println("error in creating cache : ", err)
	}
	return &cache{bigCache: bigCache, jwtCache: jwtCache}
}

func (c *cache) SetInCache(email string, otp string) error {
	otpBytes, err := json.Marshal(otp)
	if err != nil {
		log.Println("error in marshaling otp : ", err)
		return err
	}
	err = c.bigCache.Set(email, otpBytes)
	if err != nil {
		log.Println("error in setting value in cache : ", err)
		return err
	}
	log.Println("email and otp set ")
	return nil
}

func (c *cache) GetFromCache(email string) (string, error) {
	log.Println("email to find : ", email)
	otp, err := c.bigCache.Get(email)
	if err != nil {
		log.Println("error in getting corresponding otp of email : ", err)
		return "", err
	}
	var resOtp interface{}
	err = json.Unmarshal(otp, &resOtp)
	if err != nil {
		log.Println("error in unmarshalling otp from cache : ", err)
		return "", err
	}
	return fmt.Sprint(resOtp), nil
}

func PrintCacheStats(c *bigcache.BigCache) {
	stats := c.Stats()

	fmt.Printf("Number of Entries: %d\n", c.Len())
	fmt.Printf("Number of Hits: %d\n", stats.Hits)
	fmt.Printf("Number of Misses: %d\n", stats.Misses)
	fmt.Printf("Cache Size: %d bytes\n", c.Capacity())
}

func (c *cache) SetJwtInCache(jwt string, userID string) error {
	userIDBytes, err := json.Marshal(userID)
	if err != nil {
		log.Println("error in marshaling userID : ", err)
		return err
	}
	err = c.jwtCache.Set(jwt, userIDBytes)
	if err != nil {
		log.Println("error in setting value in cache : ", err)
		return err
	}
	PrintCacheStats(c.jwtCache)
	log.Println("jwt and userID set ")
	return nil
}

func (c *cache) GetUserIDFromJwt(jwt string) (string, error) {
	log.Println("jwt to find : ", jwt)
	otp, err := c.jwtCache.Get(jwt)
	if err != nil {
		log.Println("error in getting corresponding userID of jwt : ", err)
		return "", err
	}
	var resOtp interface{}
	err = json.Unmarshal(otp, &resOtp)
	if err != nil {
		log.Println("error in unmarshalling userID from cache : ", err)
		return "", err
	}
	return fmt.Sprint(resOtp), nil
}
