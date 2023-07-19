package repoimplement

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"main/model"
	"main/repository"

	"github.com/redis/go-redis/v9"
)

type commentRepoImplement struct {
	Db *sql.DB
}

type commentRepoWithCacheImplement struct {
	Db  *sql.DB
	Rdb *redis.Client
}

func NewCommentRepo(db *sql.DB) repository.CommentRepo {
	return &commentRepoImplement{
		Db: db,
	}
}

func NewCommentRepoWithCache(db *sql.DB, rdb *redis.Client) repository.CommentRepoWithCache {
	return &commentRepoWithCacheImplement{
		Db:  db,
		Rdb: rdb,
	}
}

func (commentRepo commentRepoWithCacheImplement) GetAllComment(id int) ([]*model.Comment, error) {
	key := fmt.Sprintf("comment-product-id--%d", id)
	values, err := commentRepo.Rdb.LRange(context.Background(), key, 0, -1).Result()
	if err != nil || len(values) == 0 {

		rows, err := commentRepo.Db.Query("select account_id, id, text from comments where product_id = $1", id)
		if err != nil {
			return nil, err
		}

		var comments []*model.Comment
		for rows.Next() {
			var accountId, commentId int
			var accountEmail, text string

			err := rows.Scan(&accountId, &commentId, &text)
			if err != nil {
				return nil, err
			}
			commentRepo.Db.QueryRow("select email from accounts where id = $1", accountId).Scan(&accountEmail)

			var comment model.Comment
			comment.Id = commentId
			comment.Comment = text
			comment.Email = accountEmail
			comment.ProductId = id

			comments = append(comments, &comment)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}

		go func() {
			//lưu vào redis
			var commentsString []interface{}
			for _, comment := range comments {
				jsonValue, err := json.Marshal(comment)
				if err != nil {
					panic(err)
				}
				commentsString = append(commentsString, string(jsonValue))
			}
			commentRepo.Rdb.RPush(context.Background(), key, commentsString...)
		}()

		defer rows.Close()

		return comments, nil
	}
	var commentsFromRedis []*model.Comment
	for _, value := range values {
		var comment model.Comment
		err := json.Unmarshal([]byte(value), &comment)
		if err != nil {
			panic(err)
		}
		commentsFromRedis = append(commentsFromRedis, &comment)
	}
	return commentsFromRedis, nil

}

func (commentRepo commentRepoImplement) Create(comment *model.Comment) error {

	accountId, errEmail := commentRepo.getAccountId(comment.Email)
	if errEmail != nil || accountId == 0 {
		return errEmail
	}

	_, errComment := commentRepo.Db.Exec(
		"insert into comments(account_id, text, product_id) values($1, $2, $3)",
		accountId, comment.Comment, comment.ProductId,
	)

	if errComment != nil {
		return errComment
	}

	return nil
}

func (commentRepo commentRepoImplement) Update(comment *model.Comment) error {
	accountId, err := commentRepo.getAccountId(comment.Email)
	if err != nil || accountId == 0 {
		return err
	}

	result, errWhenUpdate := commentRepo.Db.Exec(
		"update comments set text = $1 where id = $2 and account_id = $3",
		comment.Comment, comment.Id, accountId,
	)
	//lỗi khi update
	if errWhenUpdate != nil {
		return errWhenUpdate
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return errors.New("cập nhật không thành công")
	}

	return nil
}

func (commentRepo commentRepoImplement) Delete(comment *model.Comment) error {
	accountId, errEmail := commentRepo.getAccountId(comment.Email)
	if errEmail != nil || accountId == 0 {
		return errEmail
	}

	_, err := commentRepo.Db.Exec("delete from comments where account_id = $1 and id = $2", accountId, comment.Id)
	if err != nil {
		return err
	}

	return nil
}

func (commentRepo commentRepoImplement) getAccountId(email string) (int, error) {
	var accountId int

	err := commentRepo.Db.QueryRow("select id from accounts where email = $1", email).Scan(&accountId)
	if err != nil {
		return 0, err
	}
	return accountId, nil
}
