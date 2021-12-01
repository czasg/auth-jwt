package jwt

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/czasg/go-fill"
	"strings"
	"time"
)

type Header struct {
	Alg string `json:"alg" default:"JWT"`
	Typ string `json:"typ" default:"HS256"`
}

type Claim struct {
	User        string                 `json:"usr"`
	ExpiredTime int64                  `json:"exp"`
	IssuedTime  int64                  `json:"iat"`
	Meta        map[string]interface{} `json:"mta"`
}

type Token struct {
	Header
	Claim
	Digest string `json:"digest"`
}

func toBase64String(v interface{}) string {
	body, _ := json.Marshal(v)
	return base64.StdEncoding.EncodeToString(body)
}

func toHeader(header string) (Header, error) {
	var h Header
	body, err := base64.StdEncoding.DecodeString(header)
	if err != nil {
		return h, err
	}
	err = json.Unmarshal(body, &h)
	return h, err
}

func toClaim(claim string) (Claim, error) {
	var c Claim
	body, err := base64.StdEncoding.DecodeString(claim)
	if err != nil {
		return c, err
	}
	err = json.Unmarshal(body, &c)
	return c, err
}

func GenToken(token *Token) string {
	_ = fill.Fill(token, fill.OptDefault, fill.OptSilent)
	if token.IssuedTime <= 0 {
		token.IssuedTime = time.Now().Unix()
	}
	if token.ExpiredTime < token.IssuedTime {
		token.ExpiredTime = time.Now().Add(time.Minute * 15).Unix()
	}
	header := toBase64String(token.Header)
	claim := toBase64String(token.Claim)
	alg := sha256.New()
	alg.Write([]byte(header + claim + header + claim))
	token.Digest = hex.EncodeToString(alg.Sum(nil))
	return fmt.Sprintf(
		"%s.%s.%s",
		header,
		claim,
		token.Digest,
	)
}

func ParseToken(tokenString string) (*Token, error) {
	if !strings.HasPrefix(tokenString, "Bearer ") {
		return nil, errors.New("invalid token[bearer token err]")
	}
	tokens := strings.Split(strings.TrimPrefix(tokenString, "Bearer "), ".")
	if len(tokens) != 3 {
		return nil, errors.New("invalid token[token struct err]")
	}
	header, err := toHeader(tokens[0])
	if err != nil {
		return nil, err
	}
	claim, err := toClaim(tokens[1])
	if err != nil {
		return nil, err
	}
	return &Token{Header: header, Claim: claim, Digest: tokens[2]}, nil
}

func Verify(tokenString string) (*Token, error) {
	if !strings.HasPrefix(tokenString, "Bearer ") {
		return nil, errors.New("invalid token[bearer token err]")
	}
	tokens := strings.Split(strings.TrimPrefix(tokenString, "Bearer "), ".")
	if len(tokens) != 3 {
		return nil, errors.New("invalid token[token err]")
	}
	alg := sha256.New()
	alg.Write([]byte(tokens[0] + tokens[1] + tokens[0] + tokens[1]))
	v := hex.EncodeToString(alg.Sum(nil))
	if v != tokens[2] {
		return nil, errors.New("invalid token[secret err]")
	}
	token, err := ParseToken(tokenString)
	if err != nil {
		return nil, err
	}
	if token.ExpiredTime < time.Now().Unix() {
		return nil, errors.New("invalid token[expire err]")
	}
	return token, nil
}
