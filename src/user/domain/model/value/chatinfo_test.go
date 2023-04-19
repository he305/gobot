package value_test

import (
	"gobot/src/user/domain/model/value"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyChatIdShouldError(t *testing.T) {
	chatId := ""

	_, err := value.NewChatInfo(chatId)
	assert.NotNil(t, err)
}

func TestChatInfoOk(t *testing.T) {
	chatId := "some"

	_, err := value.NewChatInfo(chatId)
	assert.Nil(t, err)
}
