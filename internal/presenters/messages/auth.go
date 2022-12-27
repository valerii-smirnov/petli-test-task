package messages

type SignUpRequestBody struct {
	Email    string `json:"email" binding:"required,email" example:"your@email.com"`
	Password string `json:"password" binding:"required" example:"yousupersecretpassword"`
}

type SignInRequestBody struct {
	Email    string `json:"email" binding:"required,email" example:"your@email.com"`
	Password string `json:"password" binding:"required" example:"yousupersecretpassword"`
}

type SignInResponseBody struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzIxODY4MDUsImlhdCI6MTY3MjEwMDQwNSwianRpIjoiMjkyNWRlYjItNDljMi00NTdjLWJmN2UtZWM4N2UyYzhhOTRhIiwidXNlci1pZCI6IjhkZDZjOGJiLTg2YmMtNDVhOC1iNzhmLThkZTQxMWEzZWJlMyJ9.qF5Gye0jAkmXMvJKLCCUjFWJVjuM3C3-L4eQkEZTf3Q"`
}
