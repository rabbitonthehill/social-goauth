package oauth

type Facebook struct {
}

var _ Provider = (*Facebook)(nil)

func (p Facebook) Do() error {
	return nil
}

func (p Facebook) IDToken(token string) error {
	return nil
}

func (p Facebook) IdentityCode(token string) error {
	return nil
}
