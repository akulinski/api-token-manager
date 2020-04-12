package domain

import "time"

type Token struct {
	Id       string    `json:"id" bson:"_id,omitempty"`
	IssuedAt time.Time `json:"issuedAt" bson:"issuedAt,omitempty"`
	Issuer   string    `json:"issuer" bson:"issuer,omitempty"`
	UserID   string    `json:"userID" bson:"userID"`
	Token    string    `json:"token" bson:"token"`
	Expired  bool      `json:"expired" bson:"expired,omitempty"`
	Revoked  bool      `json:"revoked" bson:"revoked,omitempty"`
}

type TokenModel struct {
	Token       string    `json:"token"`
	GeneratedAt time.Time `json:"generatedAt"`
}

type UserModel struct{
	UserIdentificator string `json:"UserIdentificator"`
}