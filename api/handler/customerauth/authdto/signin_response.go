package authdto

type CustomerAuthSigninOutput struct {
	Body *CustomerAuthSigninOutputBody
}

// TGen service.SessionOutput [reverse]
type CustomerAuthSigninOutputBody struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}
