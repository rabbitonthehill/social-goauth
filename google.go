package oauth

type Google struct {
}

var _ Provider = (*Google)(nil)

func (p Google) Do() error {
	return nil
}

func (p Google) IDToken(token string) error {
	return nil
}

func (p Google) IdentityCode(token string) error {
	return nil
}
