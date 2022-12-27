package domain

type SignUp struct {
	Email    string
	Password string
}

type SingIn struct {
	Email    string
	Password string
}

type Token string
