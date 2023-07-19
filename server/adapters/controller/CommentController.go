package controller

import (
	"fmt"
	"main/adapters/ports/incoming"
	"main/driver"
	"main/model"
	repoimplement "main/repository/repo_implement"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type controllerComment struct{}

func NewControllerComment() controllerComment {
	return controllerComment{}
}

func (controllerComment) GetAllComment(db *driver.PostgresDB, rdb *driver.RedisDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		productId, errorWhenConvertId := strconv.Atoi(c.Param("product_id"))
		if errorWhenConvertId != nil {
			return c.JSON(http.StatusBadRequest, model.Message{
				Code:   http.StatusBadRequest,
				Text:   "id is not int",
				Output: errorWhenConvertId.Error(),
			})
		}

		comments, err := repoimplement.NewCommentRepoWithCache(db.SQL, rdb.Redis).GetAllComment(productId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, model.Message{
				Code:   http.StatusBadRequest,
				Text:   "something went wrong",
				Output: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, model.Message{
			Code:   http.StatusOK,
			Text:   "(y)",
			Output: comments,
		})
	}
}

func (controllerComment) Create(db *driver.PostgresDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		email := c.Get("email")

		incomingComment := new(incoming.CommentCreated)
		err := c.Bind(incomingComment)
		if err != nil {
			return c.JSON(http.StatusBadRequest, model.Message{
				Code:   http.StatusBadRequest,
				Text:   "không đúng định dạng comment",
				Output: err.Error(),
			})
		}

		comment := new(model.Comment)
		comment.Email = fmt.Sprintf("%v", email)
		comment.Comment = incomingComment.Comment
		comment.ProductId = incomingComment.ProductId

		if err := repoimplement.NewCommentRepo(db.SQL).Create(comment); err != nil {
			return c.JSON(http.StatusBadRequest, model.Message{
				Code:   http.StatusBadRequest,
				Text:   "error when create comment",
				Output: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, model.Message{
			Code:   http.StatusOK,
			Text:   "Bình luận đã được gửi thành công!",
			Output: comment,
		})
	}
}

func (controllerComment) Update(db *driver.PostgresDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		commentId, errorWhenConvertId := strconv.Atoi(c.Param("id"))

		if errorWhenConvertId != nil {
			return c.JSON(http.StatusBadRequest, model.Message{
				Code:   http.StatusBadRequest,
				Text:   "comment id is not int",
				Output: errorWhenConvertId.Error(),
			})
		}

		var commentUpdate incoming.CommentUpdated

		if errWhenBin := c.Bind(&commentUpdate); errWhenBin != nil {
			return c.JSON(http.StatusBadRequest, model.Message{
				Code:   http.StatusBadRequest,
				Text:   "error when binding comment updated content",
				Output: errWhenBin.Error(),
			})
		}

		var comment = model.Comment{
			Comment: commentUpdate.Comment,
			Id:      commentId,
			Email:   fmt.Sprintf("%v", c.Get("email")),
		}

		err := repoimplement.NewCommentRepo(db.SQL).Update(&comment)
		if err != nil {
			return c.JSON(http.StatusBadRequest, model.Message{
				Code:   http.StatusBadRequest,
				Text:   "error when update comment",
				Output: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, model.Message{
			Code:   http.StatusOK,
			Text:   "update successfully",
			Output: comment,
		})
	}
}

func (controllerComment) Delete(db *driver.PostgresDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		commentId, errorWhenConvertId := strconv.Atoi(c.Param("id"))

		if errorWhenConvertId != nil {
			return c.JSON(http.StatusBadRequest, model.Message{
				Code:   http.StatusBadRequest,
				Text:   "comment id is not int",
				Output: errorWhenConvertId.Error(),
			})
		}

		var comment = model.Comment{
			Id:    commentId,
			Email: fmt.Sprintf("%v", c.Get("email")),
		}

		if err := repoimplement.NewCommentRepo(db.SQL).Delete(&comment); err != nil {
			return c.JSON(http.StatusBadRequest, model.Message{
				Code:   http.StatusBadRequest,
				Text:   "error when delete",
				Output: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, model.Message{
			Code:   http.StatusOK,
			Text:   "deletd",
			Output: comment.Id,
		})
	}
}
