package oauth

type Google struct {
	service *Service
}

func (p Google) NewGoogle(service *Service) *Google {
	return &Google{service: service}
}

func (p Google) IDToken(token string) error {
	return nil
}

func (p Google) IdentityCode(token string) error {
	return nil
}
