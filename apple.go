package oauth

type Apple struct {
}

var _ Provider = (*Apple)(nil)

func (p Apple) Do() error {
	return nil
}

func (p Apple) IDToken(token string) error {
	return nil
}

func (p Apple) IdentityCode(token string) error {
	return nil
}
