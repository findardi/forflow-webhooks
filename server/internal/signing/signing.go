package signing

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"time"
)

const (
	secretPrefix = "whsec_"
	secretLen    = 32
	minSecretLen = 24
	maxSecretLen = 64
	sigVersion   = "v1"
)

var (
	ErrInvalidSecret       = errors.New("invalid secret")
	ErrInvalidTimestamp    = errors.New("invalid timestamp")
	ErrTimestampOutOfRange = errors.New("timestamp out of range")
	ErrNotMatchSignature   = errors.New("signature not match")
)

func GenerateSecret() []byte {
	buf := make([]byte, secretLen)
	_, _ = rand.Read(buf)

	return buf
}

func EncodeSecret(raw []byte) string {
	return secretPrefix + base64.StdEncoding.EncodeToString(raw)
}

func DecodeSecret(s string) ([]byte, error) {
	if !strings.HasPrefix(s, secretPrefix) {
		return nil, ErrInvalidSecret
	}

	raw, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(s, secretPrefix))
	if err != nil {
		return nil, ErrInvalidSecret
	}

	if len(raw) < minSecretLen || len(raw) > maxSecretLen {
		return nil, ErrInvalidSecret
	}

	return raw, nil
}

func Sign(secret []byte, msgID string, ts time.Time, body []byte) string {
	tsStr := FormatTimestamp(ts)

	mac := writeHMAC(msgID, tsStr, secret, body)

	return sigVersion + "," + base64.StdEncoding.EncodeToString(mac)
}

func SignAll(secrets [][]byte, msgID string, ts time.Time, body []byte) string {
	sigs := make([]string, len(secrets))
	for i, secret := range secrets {
		sigs[i] = Sign(secret, msgID, ts, body)
	}

	return strings.Join(sigs, " ")
}

func Verify(secret []byte, msgID, tsHeader string, sigHeader string, body []byte, now time.Time, tolerance time.Duration) error {
	ts, err := strconv.ParseInt(tsHeader, 10, 64)
	if err != nil {
		return ErrInvalidTimestamp
	}

	diff := now.Sub(time.Unix(ts, 0))

	if diff.Abs() > tolerance {
		return ErrTimestampOutOfRange
	}

	mac := writeHMAC(msgID, tsHeader, secret, body)

	for _, part := range strings.Fields(sigHeader) {
		version, value, ok := strings.Cut(part, ",")
		if !ok || version != sigVersion {
			continue
		}

		got, err := base64.StdEncoding.DecodeString(value)
		if err == nil && hmac.Equal(got, mac) {
			return nil
		}
	}

	return ErrNotMatchSignature
}

func FormatTimestamp(ts time.Time) string {
	return strconv.FormatInt(ts.Unix(), 10)
}

func writeHMAC(msgID, tsStr string, secret, body []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(msgID))
	mac.Write([]byte("."))
	mac.Write([]byte(tsStr))
	mac.Write([]byte("."))
	mac.Write(body)

	return mac.Sum(nil)
}
