package repository

import "main/model"

type CommentRepo interface {
	Create(*model.Comment) error
	Update(*model.Comment) error
	Delete(*model.Comment) error
}
type CommentRepoWithCache interface {
	GetAllComment(id int) ([]*model.Comment,error)
}
