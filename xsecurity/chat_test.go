package xsecurity

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestDecryptChatMessage(t *testing.T) {
	message := uuid.NewString()
	key := Hash32(time.Now().String())
	t.Log("Raw:", message)
	encryptedMessage, err := EncryptChatMessage(message, key)
	if err != nil {
		t.Error(err)
	}
	t.Log("Encrypted:", encryptedMessage)
	decryptedMessage, err := DecryptChatMessage(encryptedMessage, key)
	if err != nil {
		t.Error(err)
	}
	t.Log("Decrypted:", decryptedMessage)
	if decryptedMessage != message {
		t.Log(message)
		t.Log(decryptedMessage)
		t.Failed()
	}

	_, err = DecryptChatMessage(encryptedMessage+"a", key)
	if err == nil {
		t.Fatal("Should fail")
	}

	decryptedMessage, _ = DecryptChatMessage(encryptedMessage, Hash32(time.Now().String()+"a"))
	if decryptedMessage == message {
		t.Fatal("Should not be decrypted correctly")
	}
}
