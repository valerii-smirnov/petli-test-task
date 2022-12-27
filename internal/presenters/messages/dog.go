package messages

type DogResponseBody struct {
	ID    string `json:"id" example:"c23bca5a-640a-4f61-bb7b-5f69b1ede69d"`
	Name  string `json:"name" example:"Spike"`
	Sex   string `json:"sex" example:"male|female"`
	Age   uint   `json:"age" example:"5"`
	Breed string `json:"breed" example:"Bulldog"`
	Image string `json:"image" example:"https://cdn.w600.comps.canstockphoto.com/shepherd-cartoon-dog-vector-clipart_csp67503106.jpg"`
}

type DogListResponseBody []DogResponseBody

type CreateOrUpdateDogRequestBody struct {
	Name  string `json:"name" binding:"required,min=3,max=30" example:"Spike"`
	Sex   string `json:"sex" binding:"required,oneof=male female" example:"male|female"`
	Age   uint   `json:"age" binding:"required,min=0,max=30" example:"5"`
	Breed string `json:"breed" binding:"required,max=30" example:"Bulldog"`
	Image string `json:"image" binding:"required,url,max=1024" example:"https://cdn.w600.comps.canstockphoto.com/shepherd-cartoon-dog-vector-clipart_csp67503106.jpg"`
}

type ReactionRequestBody struct {
	Liker  string `json:"liker" binding:"required,uuid" example:"c23bca5a-640a-4f61-bb7b-5f69b1ede69d"`
	Liked  string `json:"liked" binding:"required,uuid" example:"c23bca5a-640a-4f61-bb7b-5f69b1ede69d"`
	Action string `json:"action" binding:"required,oneof=like dislike" example:"like|dislike"`
}
