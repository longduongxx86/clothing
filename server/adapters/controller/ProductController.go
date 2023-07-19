package controller

import (
	"context"
	"encoding/json"
	"main/adapters/ports/incoming"
	"main/driver"
	"main/model"
	repoimplement "main/repository/repo_implement"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type controllerProduct struct{}

func NewControllerProduct() controllerProduct {
	return controllerProduct{}
}

func (controllerProduct) Create(db *driver.PostgresDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		roleIdFromContext := c.Get("role_id")

		roleId := int(roleIdFromContext.(float64))

		if roleId == 1 {
			var incomingNewProduct incoming.ProductCreate

			if err := c.Bind(&incomingNewProduct); err != nil {
				return err
			}
			newProduct := model.Product{
				CategoryId:  incomingNewProduct.CategoryId,
				Title:       incomingNewProduct.Title,
				Price:       incomingNewProduct.Price,
				Thumbnail:   incomingNewProduct.Thumbnail,
				Description: incomingNewProduct.Description,
			}

			if err := repoimplement.NewProductRepoWithoutCache(db.SQL).Create(&newProduct); err != nil {
				return c.JSON(http.StatusBadRequest, model.Message{
					Code:   http.StatusBadRequest,
					Text:   "something went wrong",
					Output: err.Error(),
				})
			}

			return c.JSON(http.StatusBadRequest, model.Message{
				Code:   http.StatusOK,
				Text:   "success",
				Output: newProduct,
			})
		} else {
			return c.JSON(http.StatusBadRequest, model.Message{
				Code: http.StatusBadRequest,
				Text: "bạn đéo có quyền",
			})
		}
	}

}

func (controllerProduct) Update(db *driver.PostgresDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		roleIdFromContext := c.Get("role_id")

		roleId := int(roleIdFromContext.(float64))
		if roleId == 1 {
			productId, errorWhenConvertId := strconv.Atoi(c.Param("id"))
			if errorWhenConvertId != nil {
				return c.JSON(http.StatusBadRequest, model.Message{
					Code:   http.StatusBadRequest,
					Text:   "id is not int",
					Output: errorWhenConvertId.Error(),
				})
			}
			var incomingProductUpdate incoming.ProductUpdate

			if err := c.Bind(&incomingProductUpdate); err != nil {
				return err
			}
			productUpdated := model.Product{
				Id:          productId,
				Title:       incomingProductUpdate.Title,
				Price:       incomingProductUpdate.Price,
				Thumbnail:   incomingProductUpdate.Thumbnail,
				Description: incomingProductUpdate.Description,
			}

			if err := repoimplement.NewProductRepoWithoutCache(db.SQL).Update(&productUpdated); err != nil {
				return c.JSON(http.StatusBadRequest, model.Message{
					Code:   http.StatusBadRequest,
					Text:   "something went wrong",
					Output: err,
				})
			}

			return c.JSON(http.StatusOK, model.Message{
				Code:   http.StatusOK,
				Text:   "successfull",
				Output: productUpdated,
			})
		} else {
			return c.JSON(http.StatusBadRequest, model.Message{
				Code: http.StatusBadRequest,
				Text: "bạn đéo có quyền",
			})
		}
	}

}

func (controllerProduct) Delete(db *driver.PostgresDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		roleIdFromContext := c.Get("role_id")

		roleId := int(roleIdFromContext.(float64))

		if roleId == 1 {
			if err := repoimplement.NewProductRepoWithoutCache(db.SQL).DeleteProduct(c.Param("id")); err != nil {
				return c.JSON(http.StatusBadRequest, model.Message{
					Code:   http.StatusBadRequest,
					Text:   "something went wrong",
					Output: err,
				})
			}

			return c.JSON(http.StatusOK, model.Message{
				Code: http.StatusOK,
				Text: "successfull",
			})
		} else {
			return c.JSON(http.StatusBadRequest, model.Message{
				Code: http.StatusBadRequest,
				Text: "bạn đéo có quyền",
			})
		}
	}
}

// -_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_

func (controllerProduct) GetProductDetail(db *driver.PostgresDB, rdb *driver.RedisDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		productId := c.Param("id")

		productStringRedis := rdb.Redis.Get(context.Background(), productId).Val()

		if productStringRedis == "" { //trong redis không có
			product, errorPostgres := repoimplement.NewProductRepoWithCache(db.SQL, rdb.Redis).GetProductDetail(productId)

			if errorPostgres != nil {
				return c.JSON(http.StatusBadRequest, model.Message{
					Code:   http.StatusBadRequest,
					Text:   errorPostgres.Error(),
					Output: errorPostgres.Error(),
				})
			}

			listProduct := [1]*model.Product{
				product,
			}

			go setToRedis(rdb, productId, listProduct[:])

			return c.JSON(http.StatusOK, model.Message{
				Code:   http.StatusOK,
				Text:   "lấy thành công",
				Output: product,
			})
		} else {
			product := new(model.Product)
			err := json.Unmarshal([]byte(productStringRedis), &product)

			if err != nil {
				return c.JSON(http.StatusBadRequest, model.Message{
					Code:   http.StatusBadRequest,
					Text:   "có trong redis, nhưng unmarshall thất bại",
					Output: err.Error(),
				})
			}
			return c.JSON(http.StatusOK, model.Message{
				Code:   http.StatusOK,
				Text:   "lấy thành công",
				Output: product,
			})
		}
	}
}

func (controllerProduct) GetListProductWithCondition(db *driver.PostgresDB, rdb *driver.RedisDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var min, max, category int
		var title, minParam, maxParam, categoryParam string
		var errMin, errMax, errCategory error

		minParam = c.QueryParams().Get("min")
		maxParam = c.QueryParams().Get("max")
		categoryParam = c.QueryParams().Get("category")

		if minParam == "" {
			min, errMin = strconv.Atoi(minParam)
			if errMin != nil {
				min = 0
			}
		} else {
			min, _ = strconv.Atoi(minParam)
		}

		if maxParam == "" {
			max, errMax = strconv.Atoi(maxParam)
			if errMax != nil {
				max = 999999999
			}
		} else {
			max, _ = strconv.Atoi(maxParam)
		}

		if categoryParam != "" {
			category, errCategory = strconv.Atoi(categoryParam)
			if errCategory != nil {
				category = 0
			}
		} else {
			category, _ = strconv.Atoi(categoryParam)
		}

		title = c.QueryParams().Get("title")

		condition := model.ConditionSelected{
			Min:      min,
			Max:      max,
			Title:    title,
			Category: category,
		}

		productList, err := repoimplement.NewProductRepoWithCache(db.SQL, rdb.Redis).GetListProductWithCondition(&condition)

		if err != nil {
			return c.JSON(http.StatusBadRequest, model.Message{
				Code:   http.StatusBadRequest,
				Text:   "something went wrong: " + err.Error(),
				Output: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, model.Message{
			Code:   http.StatusOK,
			Text:   "Ok",
			Output: productList,
		})
	}
}

func setToRedis(rdb *driver.RedisDB, key string, product []*model.Product) {
	if len(product) == 1 {
		jsonProduct, err := json.Marshal(product[0])
		if err != nil {
			return
		}

		rdb.Redis.Set(context.Background(), key, string(jsonProduct), 7*24*60*60*time.Second).Err() //lưu trong 1 tuần
	} else {
		return
	}
}
