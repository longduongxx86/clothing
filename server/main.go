package main

import (
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"main/adapters/controller"
	server_middleware "main/adapters/middleware"
	"main/driver"

	"github.com/labstack/echo/v4"
)

const (
	hostPsql     = "localhost"
	portPsql     = "5432"
	userPsql     = "postgres"
	passwordPsql = "Nguyenthuyngan00"
	dbnamePsql   = "clothing"
	maxIdleConns = 10
	maxOpenConns = 100

	addressRedis  = "localhost:6379"
	passwordRedis = ""
	dbRedis       = 0
)

func main() {
	db := driver.ConnectPsql(hostPsql, portPsql, userPsql, passwordPsql, dbnamePsql, maxIdleConns, maxOpenConns)
	rdb := driver.ConnectRedis(addressRedis, passwordRedis, dbRedis)

	e := echo.New()
	e.Use(echo_middleware.CORS())

	auth := e.Group("/auth")
	auth.POST("/login", controller.NewControllerAuth().Login(db))
	auth.PATCH("/re_password", controller.NewControllerAuth().ResetPassword(db))
	auth.PUT("/signin", controller.NewControllerAuth().SignIn(db))

	productAdmin := e.Group("/admin/products")
	productAdmin.Use(server_middleware.JWTMiddleware)
	productAdmin.POST("/create", controller.NewControllerProduct().Create(db))
	productAdmin.PATCH("/update/:id", controller.NewControllerProduct().Update(db))
	productAdmin.DELETE("/delete/:id", controller.NewControllerProduct().Delete(db))

	product := e.Group("/products")
	product.GET("/:id", controller.NewControllerProduct().GetProductDetail(db, rdb))
	product.GET("", controller.NewControllerProduct().GetListProductWithCondition(db, rdb))

	commentMDW := e.Group("comments")
	commentMDW.Use(server_middleware.JWTMiddleware)
	commentMDW.POST("/create", controller.NewControllerComment().Create(db))
	commentMDW.PUT("/update/:id", controller.NewControllerComment().Update(db))
	commentMDW.DELETE("/delete/:id", controller.NewControllerComment().Delete(db))

	comment := e.Group("/comments")
	comment.GET("/:product_id", controller.NewControllerComment().GetAllComment(db, rdb))

	bill := e.Group("/bill")
	bill.Use(server_middleware.JWTMiddleware)
	bill.POST("", controller.NewControllerBill().Create(db))
	bill.GET("", controller.NewControllerBill().GetAllBills(db))
	bill.GET("/:bill-id", controller.NewControllerBill().GetBillDetail(db))

	e.Logger.Fatal(e.Start(":1323"))
}
