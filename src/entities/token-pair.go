package entities

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type ParsedToken struct {
	UserID  uint    `json:"userId"`
	EmailID string  `json:"emailId"`
	User    JWTUser `json:"user"`
}
