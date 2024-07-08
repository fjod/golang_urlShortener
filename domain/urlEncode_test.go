package domain

import (
	"fmt"
	M "shortUrl/db/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetShortUrl(t *testing.T) {
	mockDB := M.NewOperations(t)

	t.Run("success", func(t *testing.T) {
		longUrl := "https://example.com/long/url"
		expectedShortUrl := "ba"

		mockDB.On("GetUrlId").Return(62, nil)
		mockDB.On("SetUrl", longUrl, 62).Return(nil)

		shortUrl, err := GetShortUrl(longUrl, mockDB)
		assert.NoError(t, err)
		assert.Equal(t, expectedShortUrl, shortUrl)

		mockDB.AssertExpectations(t)
	})

	mockDB = M.NewOperations(t)
	t.Run("error getting URL ID", func(t *testing.T) {
		longUrl := "https://example.com/long/url"
		expectedError := fmt.Errorf("failed to get URL ID")

		mockDB.On("GetUrlId").Return(0, expectedError)

		shortUrl, err := GetShortUrl(longUrl, mockDB)
		assert.Empty(t, shortUrl)
		assert.EqualError(t, err, expectedError.Error())

		mockDB.AssertExpectations(t)
	})

	mockDB = M.NewOperations(t)
	t.Run("error setting URL", func(t *testing.T) {
		longUrl := "https://example.com/long/url"
		expectedError := fmt.Errorf("failed to set URL")

		mockDB.On("GetUrlId").Return(62, nil)
		mockDB.On("SetUrl", longUrl, 62).Return(expectedError)

		shortUrl, err := GetShortUrl(longUrl, mockDB)
		assert.Equal(t, "ba", shortUrl)
		assert.EqualError(t, err, expectedError.Error())

		mockDB.AssertExpectations(t)
	})
}
