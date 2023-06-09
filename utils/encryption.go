package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-jose/go-jose/v3"
	"os"
	"time"
)

// RSA
// RSA는 비대칭 암호로, 예를 들어 무작위로 생성된 숫자 2개 A, B가 있다고 할 때
// 무작위 정보 Information을 A로 암호화하면 B로만 해제할 수 있고, B로 암호화하면 A로만 해제할 수 있다고 하면
// A, B 둘중 하나를 임의로 정해서 공개 키로 정하고, 나머지 하나를 비밀 키로 정한다
// A를 공개한다고 하면, 내 친구가 나에게 보내고 싶은 편지를 A로 암호화 하면, 나만 B를 알고 있을 테니
// 나만 편지를 읽을 수 있다

var PRIVATE *rsa.PrivateKey

func InitKeys() {
	contents, err := os.ReadFile("./private.pem")
	if err != nil {
		panic(err.Error())
	}

	block, _ := pem.Decode(contents)

	p, err := x509.ParsePKCS8PrivateKey(block.Bytes)

	if err != nil {
		panic(err.Error())
	}

	PRIVATE = p.(*rsa.PrivateKey)
}

// ServerToClient 패킷을 서버에서 클라이언트로 보낼 준비를 함
// 1. 패킷을 서버에서 보냈음을 보장하기 위해 서버에서 자기가 만든거라고 서명을 한다
// 2. 그 서명한 거를 클라이언트 공개키로 암호화한다; 이제 클라이언트 비밀 키로만 암호 해제를 할 수 있다
// 3. 클라이언트는 해당 메세지를 받고 암호화를 해제한다. 그리고 서버가 보낸 것임 또한 확실하게 확인한다 :D
func ServerToClient(payload interface{}, claim string, clientKey rsa.PublicKey) (string, error) {
	signed, err := SignJWT(payload, claim)
	if err != nil {
		return "", err
	}

	encrypter, err := jose.NewEncrypter(jose.A128GCM, jose.Recipient{Algorithm: jose.RSA_OAEP, Key: clientKey}, nil)
	if err != nil {
		return "", err
	}

	obj, err := encrypter.Encrypt([]byte(signed))
	if err != nil {
		return "", err
	}

	return obj.FullSerialize(), nil
}

// ClientToServer 클라이언트에서 보낸 패킷을 클라이언트에서 받음
// 1. 클라이언트가 자기 비밀 키로 패킷에 서명함 -> 자기가 보냈다고 보장할 수 있음
// 2. 클라이언트가 서버 공개키로 패킷을 암호화함 -> 서버만 열람하도록 제한할 수 있음
// 3. 서버가 이를 서버 비밀키로 암호화를 해제함
// 4. 서버가 클라이언트 공개 키로 제대로 된 클라이언트가 보냈는지 확인함 :D
func ClientToServer(jwe string, claim string, clientKey *rsa.PublicKey) (interface{}, error) {
	encrypted, err := jose.ParseEncrypted(jwe)
	if err != nil {
		return nil, err
	}

	decrypted, err := encrypted.Decrypt(PRIVATE)
	if err != nil {
		return nil, err
	}
	asString := string(decrypted)

	contents, err := ParseJWT(&asString, claim, clientKey)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func SignJWTWithKey(toEncrypt interface{}, claim string, key *rsa.PrivateKey) (string, error) {
	// Define the expiration time for the token
	expirationTime := time.Now().Add(24 * time.Hour).Unix()

	// Create a new claims instance
	claims := jwt.MapClaims{}
	claims[claim] = toEncrypt
	claims["exp"] = expirationTime

	// Create a new token instance using the claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func SignJWT(toEncrypt interface{}, claim string) (string, error) {
	// Define the expiration time for the token
	expirationTime := time.Now().Add(24 * time.Hour).Unix()

	// Create a new claims instance
	claims := jwt.MapClaims{}
	claims[claim] = toEncrypt
	claims["exp"] = expirationTime

	// Create a new token instance using the claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString(PRIVATE)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ParseJWT JWT 서명을 확인하고 내용 분석
func ParseJWT(signed *string, claim string, senderKey *rsa.PublicKey) (interface{}, error) {
	// Verify the JWT signed
	VerifiedToken, err := VerifyJWT(*signed, senderKey)
	if err != nil {
		return "", err
	}

	claims := VerifiedToken.Claims.(jwt.MapClaims)

	// Extract the user ID from the signed claims
	contents, ok := claims[claim]
	if !ok {
		return "", fmt.Errorf("error: extracting %s from signed", claim)
	}

	return contents, nil
}

func VerifyJWT(tokenString string, senderKey *rsa.PublicKey) (*jwt.Token, error) {
	// Define the expected signing method and secret key
	signingMethod := jwt.SigningMethodRS256
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return senderKey, nil
	}

	// Parse the JWT token string
	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	// Verify the token signature and expiration time
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	if token.Method != signingMethod {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Method)
	}
	if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	exp, ok := token.Claims.(jwt.MapClaims)["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid expiration time")
	}
	if time.Unix(int64(exp), 0).Before(time.Now()) {
		return nil, fmt.Errorf("token has expired")
	}

	// Token is valid, return it
	return token, nil
}
