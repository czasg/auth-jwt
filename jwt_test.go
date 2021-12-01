package jwt

import (
	"reflect"
	"testing"
)

func assert(t *testing.T, v1, v2 interface{}) {
	if !reflect.DeepEqual(v1, v2) {
		t.Errorf("want [%v], but [%v]", v1, v2)
	}
}

func TestGenToken(t *testing.T) {
	t1 := Token{Claim: Claim{
		ExpiredTime: 2,
		IssuedTime:  1,
	}}
	tokenString := GenToken(&t1)
	assert(t, "eyJhbGciOiJKV1QiLCJ0eXAiOiJIUzI1NiJ9.eyJ1c3IiOiIiLCJleHAiOjIsImlhdCI6MSwibXRhIjp7fX0=.f22d997e08e632ebad756d4718052d79c1faa87aee2c182c8b56eafc3ae5bb1d", tokenString)
}

func TestParseToken(t *testing.T) {
	t1 := Token{}
	tokenString := GenToken(&t1)
	t2, err := ParseToken("Bearer " + tokenString)
	assert(t, nil, err)
	assert(t, t1, *t2)
}

func TestVerify(t *testing.T) {
	t1 := Token{}
	tokenString := GenToken(&t1)
	t2, err := Verify("Bearer " + tokenString)
	assert(t, nil, err)
	assert(t, t1, *t2)
}
