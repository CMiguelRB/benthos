package common

import (
	"os"
	"testing"
)

var testString string = "abcde"
var testResult string = "2e800799b9866f2f2509af45ed30a3665b0ab9d4309354ec97967f21aacfce4b"

func TestEncrypt(t *testing.T) {

	os.Setenv("ENCRYPTION_KEY", "61A7FA981C8E6FFCA386A796C1C73B17")

	encrypted, err := Encrypt(testString)

	if err != nil {
		t.Errorf("Error: %s", err)
	}

	decrypted, err := Decrypt(encrypted)

	if decrypted != testString {
		t.Errorf("%s != %s", encrypted, testResult)
	}
}

func TestDecrypt(t *testing.T) {

	os.Setenv("ENCRYPTION_KEY", "61A7FA981C8E6FFCA386A796C1C73B17")

	decrypted, err := Decrypt(testResult)

	if err != nil {
		t.Errorf("Error: %s", err)
	}

	if decrypted != testString {
		t.Errorf("%s != %s", decrypted, testString)
	}
}
