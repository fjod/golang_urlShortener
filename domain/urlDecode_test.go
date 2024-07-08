package domain

import (
	"errors"
	"fmt"
	DB "shortUrl/db"
	M "shortUrl/db/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLongUrl(t *testing.T) {
	mockDB := M.NewOperations(t)

	t.Run("success", func(t *testing.T) {
		shortUrl := "ba"
		expectedLongUrl := "https://example.com/long/url"

		mockDB.On("GetUrl", 62).Return(DB.Url{Url: expectedLongUrl}, nil)

		longUrl, err := GetLongUrl(shortUrl, mockDB)
		assert.NoError(t, err)
		assert.Equal(t, expectedLongUrl, longUrl)

		mockDB.AssertExpectations(t)
	})

	mockDB = M.NewOperations(t)
	t.Run("error decoding short URL", func(t *testing.T) {
		expectedError := fmt.Errorf("failed to get URL")
		mockDB.On("GetUrl", 62).Return(DB.Url{Url: ""}, expectedError)
		longUrl, err := GetLongUrl("ba", mockDB)
		assert.Empty(t, longUrl)
		assert.Error(t, err)

		mockDB.AssertExpectations(t)
	})

	mockDB = M.NewOperations(t)
	t.Run("error getting URL from DB", func(t *testing.T) {
		shortUrl := "ba"
		expectedError := errors.New("failed to get URL")

		mockDB.On("GetUrl", 62).Return(DB.Url{Url: ""}, expectedError)

		longUrl, err := GetLongUrl(shortUrl, mockDB)
		assert.Empty(t, longUrl)
		assert.EqualError(t, err, expectedError.Error())

		mockDB.AssertExpectations(t)
	})
}
