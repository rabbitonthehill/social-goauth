package oauth

type Facebook struct {
	service *Service
}

// https://developers.facebook.com/docs/apps/for-business#field
// https://developers.facebook.com/docs/apps/for-business#api
// https://developers.facebook.com/docs/apps/business-manager#create-business

const (
	FacebookApiEndpoint        = "https://api.facebook.com"
	FacebookApiVideoEndpoint   = "https://api-video.facebook.com"
	FacebookApiReadEndpoint    = "https://api-read.facebook.com"
	FacebookGraphEndpoint      = "https://graph.facebook.com"
	FacebookGraphVideoEndpoint = "https://graph-video.facebook.com"
	FacebookWWWEndpoint        = "https://www.facebook.com"
)

func NewFacebook(service *Service) *Facebook {
	return &Facebook{service: service}
}
func (p Facebook) IDToken(token string) error {
	return nil
}

func (p Facebook) IdentityCode(token string) error {
	return nil
}
