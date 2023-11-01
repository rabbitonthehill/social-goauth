package oauth

type Line struct {
}

var _ Provider = (*Line)(nil)

func (p Line) Do() error {
	return nil
}

func (p Line) IDToken(token string) error {
	return nil
}

func (p Line) IdentityCode(token string) error {
	return nil
}
