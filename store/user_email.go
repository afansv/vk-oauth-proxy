package store

import (
	"time"

	"github.com/jellydator/ttlcache/v3"
)

type UserEmail struct {
	cache *ttlcache.Cache[int, string]
}

func (store UserEmail) Start() {
	store.cache.Start()
}

func NewUserEmail(ttl time.Duration) *UserEmail {
	cache := ttlcache.New[int, string](
		ttlcache.WithTTL[int, string](ttl),
	)
	return &UserEmail{cache: cache}
}

func (store UserEmail) Set(userID int, email string) {
	if email == "" || userID == 0 {
		return
	}
	store.cache.Set(userID, email, ttlcache.DefaultTTL)
}

func (store UserEmail) Get(userID int) (email string, ok bool) {
	if userID == 0 {
		return
	}
	emailPtr := store.cache.Get(userID)
	if emailPtr == nil {
		return "", false
	}

	return emailPtr.Value(), true
}
