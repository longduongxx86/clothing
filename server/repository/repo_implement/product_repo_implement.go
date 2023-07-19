package repoimplement

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"main/model"
	"main/repository"

	"github.com/redis/go-redis/v9"
)

type ProductRepoImplementWithCache struct {
	Db  *sql.DB
	Rdb *redis.Client
}
type ProductRepoImplementWithoutCache struct {
	Db *sql.DB
}

func NewProductRepoWithCache(db *sql.DB, rdb *redis.Client) repository.ProductRepoWithCache {
	return ProductRepoImplementWithCache{
		Db:  db,
		Rdb: rdb,
	}
}

func NewProductRepoWithoutCache(db *sql.DB) repository.ProductRepoWithoutCache {
	return ProductRepoImplementWithoutCache{
		Db: db,
	}
}

func (productRepo ProductRepoImplementWithoutCache) Create(product *model.Product) error {
	var productCheck model.Product
	result := productRepo.Db.QueryRow(`select title from products where title = $1`, product.Title).Scan(&productCheck.Title)

	if result != nil {
		if result.Error() == "sql: no rows in result set" {
			query := fmt.Sprintf(`insert into products("category_id", "title", "price", "thumbnail", "description") values(%d, '%s', %d, '%s', '%s')`, product.CategoryId, product.Title, product.Price, product.Thumbnail, product.Description)
			productRepo.Db.QueryRow(query)
			return nil
		} else {
			return result
		}
	} else {
		return errors.New("product exists")
	}
}

func (productRepo ProductRepoImplementWithoutCache) Update(product *model.Product) error {
	fmt.Println("product: ", product)

	var productCheck model.Product
	checkExistsQuery := fmt.Sprintf(`select id from product where id = %d`, product.Id)
	result := productRepo.Db.QueryRow(checkExistsQuery).Scan(&productCheck.Title)

	fmt.Println("result: ", result)

	if result != nil {
		return result
	}

	var err error

	if product.Title != "" {
		if _, err = productRepo.Db.Exec(`update products set title = $1 where id = $2`, product.Title, product.Id); err != nil {
			fmt.Println("title")
			return err
		}
	}
	if product.Price != 0 {
		if _, err = productRepo.Db.Exec(`update products set price = $1 where id = $2`, product.Price, product.Id); err != nil {
			fmt.Println("price: ", err)
			return err
		}
	}
	if product.Thumbnail != "" {
		if _, err = productRepo.Db.Exec(`update products set thumbnail = $1 where id = $2`, product.Thumbnail, product.Id); err != nil {
			fmt.Println("thumbnail")
			return err
		}
	}
	if product.Description != "" {
		if _, err = productRepo.Db.Exec(`update products set description = $1 where id = $2`, product.Description, product.Id); err != nil {
			fmt.Println("description")
			return err
		}
	}

	return nil
}

func (productRepo ProductRepoImplementWithoutCache) DeleteProduct(idProduct string) error {
	_, err := productRepo.Db.Exec(`delete from products where id = $1`, idProduct)

	if err != nil {
		return err
	}
	return nil
}

func (productRepo ProductRepoImplementWithCache) GetProductDetail(id string) (*model.Product, error) {
	var product model.Product
	result := productRepo.Db.QueryRow(
		"select products.id, products.title, products.price, products.thumbnail, products.description, category.name as category_name from products left join category on products.category_id = category.id where products.id = $1",
		id,
	)

	err := result.Scan(&product.Id, &product.Title, &product.Price, &product.Thumbnail, &product.Description, &product.Category)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (productRepo ProductRepoImplementWithCache) GetListProductWithCondition(condition *model.ConditionSelected) ([]*model.Product, error) {
	var result *sql.Rows
	var err error

	if condition.Category != 0 {
		result, err = productRepo.Db.Query(
			"select products.id, products.title, products.price, products.thumbnail, products.description, products.category_id, category.name from products left join category on products.category_id = category.id where price >= $1 and price <= $2 and title like $3 and category_id = $4",
			condition.Min, condition.Max, "%"+condition.Title+"%", condition.Category,
		)
	} else {
		result, err = productRepo.Db.Query(
			"select products.id, products.title, products.price, products.thumbnail, products.description, products.category_id, category.name from products left join category on products.category_id = category.id where price >= $1 and price <= $2 and title like $3 ",
			condition.Min, condition.Max, "%"+condition.Title+"%",
		)
	}

	defer result.Close()

	if err != nil {
		return nil, err
	}

	productList := make([]*model.Product, 0)

	for result.Next() {
		var product model.Product

		err := result.Scan(
			&product.Id,
			&product.Title,
			&product.Price,
			&product.Thumbnail,
			&product.Description,
			&product.CategoryId,
			&product.Category,
		)
		if err != nil {
			log.Fatal(err)
			break
		}
		productList = append(productList, &product)
	}

	return productList, nil
}
