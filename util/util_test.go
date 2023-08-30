package util

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/jonada182/cover-letter-ai-api/mocks"
	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/assert"
)

func TestLoadEnvFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUtil := mocks.NewMockUtil(ctrl)

	t.Run("ValidFile", func(t *testing.T) {
		// Prepare a temporary file with mock content
		tmpFile, err := os.CreateTemp("", "envfile")
		assert.NoError(t, err)
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		_, err = tmpFile.WriteString("OPENAI_KEY=MyKey\nMONGO_URI=some_uri\n")
		assert.NoError(t, err)

		mockUtil.EXPECT().LoadEnvFile(tmpFile.Name()).Return(nil)

		err = mockUtil.LoadEnvFile(tmpFile.Name())
		assert.NoError(t, err)
	})

	t.Run("InvalidFile", func(t *testing.T) {
		mockUtil.EXPECT().LoadEnvFile("nonexistent_file").Return(errors.New("file not found"))

		err := mockUtil.LoadEnvFile("nonexistent_file")
		assert.Error(t, err)
		assert.Equal(t, "file not found", err.Error())
	})

	t.Run("InvalidLine", func(t *testing.T) {
		mockUtil.EXPECT().LoadEnvFile(gomock.Any()).DoAndReturn(
			func(filename string) error {
				return LoadEnvFile(filename) // Use actual function for this test
			})

		tmpFile, err := os.CreateTemp("", "envfile")
		assert.NoError(t, err)
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		_, err = tmpFile.WriteString("OPENAI_KEY=MyKey\nInvalidLine\n")
		assert.NoError(t, err)

		err = mockUtil.LoadEnvFile(tmpFile.Name())
		fmt.Println(err.Error())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid line in .env file")
	})
}
