package domain

import (
	C "shortUrl/cache"
	DB "shortUrl/db"
	"time"
)

// GetShortUrl сохранить новый url в базу и вернуть его короткую ссылку
func GetShortUrl(longUrl string, db DB.Operations) (string, error) {
	cache := C.NewCacheService(time.Minute)
	newUrlId, err := db.GetUrlId()
	if err != nil {
		return "", err
	}
	shortUrl := encode(newUrlId)
	err = db.SetUrl(longUrl, newUrlId)
	cache.Add <- C.Kvp{Key: shortUrl, Value: longUrl}
	return shortUrl, err
}
