package domain

import DB "shortUrl/db"

// GetShortUrl сохранить новый url в базу и вернуть его короткую ссылку
func GetShortUrl(longUrl string, db DB.Operations) (string, error) {
	newUrlId, err := db.GetUrlId()
	if err != nil {
		return "", err
	}
	shortUrl := encode(newUrlId)
	err = db.SetUrl(longUrl, newUrlId)
	return shortUrl, err
}
