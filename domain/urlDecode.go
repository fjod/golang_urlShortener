package domain

import (
	C "shortUrl/cache"
	DB "shortUrl/db"
	"time"
)

// getLongUrl получить полный url по его короткой ссылке
func getLongUrl(shortUrl string, db DB.Operations) (string, error) {
	urlId := decode(shortUrl)
	url, err := db.GetUrl(urlId)
	if err != nil {
		return "", err
	}
	return url.Url, err
}

// GetLongUrl получить полный url по его короткой ссылке (с кэшированием)
func GetLongUrl(shortUrl string, db DB.Operations) (string, error) {
	cache := C.NewCacheService(time.Minute)
	cache.Key <- shortUrl
	v := <-cache.Value
	if v != "" {
		return v, nil
	}
	return getLongUrl(shortUrl, db)
}
