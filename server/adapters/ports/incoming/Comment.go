package incoming

type CommentCreated struct {
	Comment   string `json:"comment" form:"comment"`
	ProductId int    `json:"productId" form:"productId"`
}

type CommentUpdated struct {
	Comment   string `json:"comment" form:"comment"`
}
