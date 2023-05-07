package cookies

import (
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
