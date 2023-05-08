package cookies

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
)

// Will verify vadility and size of cookie
var (
	ErrValueTooLong = errors.New("cookie value too long")
	ErrInvalidValue = errors.New("invalid cookie value")
)

// Function to encode cookie
func Write(w http.ResponseWriter, cookie http.Cookie) error {
	// Encode cookie value using base64.
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(cookie.Value))

	// Check length of cookie contents, error if it is more than 4096 bytes.
	if len(cookie.String()) > 4096 {
		return ErrValueTooLong
	}

	// Send cookie to client
	http.SetCookie(w, &cookie)

	return nil
}

// Function to decode cookie
func Read(r *http.Request, name string) (string, error) {
	// Reads the cookie
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	/* Decodes the encoded cookie value,
	   If not a valid encoded value,
	   Return the ErrInvalidValue error.*/
	value, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return "", ErrInvalidValue
	}

	// Returns decoded cookie value.
	return string(value), nil
}

// Tampering proof cookie
func WriteSigned(w http.ResponseWriter, cookie http.Cookie, secretKey []byte) error {
	// Calculate a HMAC signature of the cookie name and value.
	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(cookie.Name))
	mac.Write([]byte(cookie.Value))
	signature := mac.Sum(nil)

	// Prepend the cookie value with the HMAC signature.
	cookie.Value = string(signature) + cookie.Value

	// Base64-encode the new cookie value and write the cookie.
	return Write(w, cookie)
}

func ReadSigned(r *http.Request, name string, secretKey []byte) (string, error) {
	// Read in the signed value from the cookie.
	// "{signature}{original value}".
	signedValue, err := Read(r, name)
	if err != nil {
		return "", err
	}

	// A SHA256 HMAC signature has a fixed length of 32 bytes.
	// Ensure that the length is at least this long.
	if len(signedValue) < sha256.Size {
		return "", ErrInvalidValue
	}

	// Split the signature and original cookie value.
	signature := signedValue[:sha256.Size]
	value := signedValue[sha256.Size:]

	// Recalculate the HMAC signature.
	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(name))
	mac.Write([]byte(value))
	expectedSignature := mac.Sum(nil)

	// Check recalculated signature matches the signature received.
	if !hmac.Equal([]byte(signature), expectedSignature) {
		return "", ErrInvalidValue
	}

	// Return the original cookie value.
	return value, nil
}
