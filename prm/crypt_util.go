package prm

import (
	//"crypto/tls"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
)

// createUffer creates a hex string of an AES Block encoded string
// Relies on the global uffer_key. The uffer is basically our protection
// mechanism between forms so a user can't just send form data
func createUffer(plainstring string, ufferKey string) string {
	key := []byte(ufferKey)

	// Make sure we pad out to the block size
	texttemp := []byte(plainstring)
	var plaintext []byte
	pad := 0
	var smallpad uint8

	// We always pad so we can check padding when we decrypt
	if len(texttemp)%aes.BlockSize != 0 {
		pad = aes.BlockSize - (len(texttemp) % aes.BlockSize)
	} else {
		pad = aes.BlockSize
	}

	for i := 0; i < len(texttemp); i++ {
		plaintext = append(plaintext, texttemp[i])
	}

	smallpad = uint8(pad)

	for i := 0; i < pad; i++ {
		plaintext = append(plaintext, smallpad)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Print(err)
		return ("")
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Print(err)
		return ("")
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return hex.EncodeToString(ciphertext)

}

// decrypt_uffer decrypts a hex string generated from the encrypt_uffer function
func decryptUffer(hextext string, ufferKey string) string {
	key := []byte(ufferKey)

	ciphertext, _ := hex.DecodeString(hextext)
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Print(err)
		return ("")
	}

	if len(ciphertext) < aes.BlockSize {
		log.Print(err)
		return ("")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		log.Print(err)
		return ("")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	padding := int(ciphertext[len(ciphertext)-1])
	plaintext := fmt.Sprintf("%s", ciphertext[0:len(ciphertext)-padding])

	return plaintext
}

// CreatePasswordHash creates a Linux salty hash for the changing of passwords
func CreatePasswordHash(password string) (string, []byte) {

	h := sha1.New()
	salt := make([]byte, 4)
	_, err := rand.Read(salt)

	if err != nil {
		log.Print(err)
		return "", salt
	}

	io.WriteString(h, password)
	final := "{SSHA}" + base64.StdEncoding.EncodeToString(h.Sum(salt))

	return final, salt
}
