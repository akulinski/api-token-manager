package api

import "testing"

func TestGenerateToken(t *testing.T)  {
	tokenString := GenerateToken("test")

	_, err := ValidateJwt(TokenModel{Token: tokenString})

	if err!=nil{
		t.Error("Failed to parse generated token")
	}
}

func TestValidateJwt(t *testing.T) {

	_, err := ValidateJwt(TokenModel{Token: "someRandomString"})

	if err == nil{
		t.Error("Should not validate random string as token")

	}
}