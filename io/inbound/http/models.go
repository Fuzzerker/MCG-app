package inboundhttp

type UserRequest struct {
	Username string `json:"username" required:"true" minLength:"6" description:"Username of the user"`
	Password string `json:"password" required:"true" minLength:"6" description:"Password the user will use to log in"`
}

type LoginResponse struct {
	Token string `json:"token" descripiton:"access token generated for the given credentials.  Should be sent as a bearer token on all future requests"`
}

type Empty struct {
}
