package domain

import DB "shortUrl/db"

// GetLongUrl получить полный url по его короткой ссылке
func GetLongUrl(shortUrl string, db DB.Operations) (string, error) {
	urlId := decode(shortUrl)
	url, err := db.GetUrl(urlId)
	if err != nil {
		return "", err
	}
	return url.Url, err
}
