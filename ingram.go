package ingram

import "github.com/go-playground/validator/v10"

const apiEndpoint = "https://api.ingrammicro.com"

type Ingram struct {
	clientID     string
	clientSecret string
	isSandbox    bool
	endpoint     string
	token        *Token
	validate     *validator.Validate
}

type OptionFunc func(i *Ingram) error

func WithOAuthCredentials(clientID, clientSecret string) OptionFunc {
	return func(i *Ingram) error {
		i.clientID = clientID
		i.clientSecret = clientSecret
		return nil
	}
}

func EnableSandbox() OptionFunc {
	return func(i *Ingram) error {
		i.isSandbox = true
		return nil
	}
}

func New(options ...OptionFunc) (*Ingram, error) {
	i := &Ingram{
		validate: validator.New(),
	}

	for _, v := range options {
		err := v(i)
		if err != nil {
			return nil, err
		}
	}

	i.endpoint = apiEndpoint
	if i.isSandbox {
		i.endpoint += "/sandbox"
	}

	return i, nil
}
