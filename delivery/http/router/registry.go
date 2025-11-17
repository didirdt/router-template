/*
 * Copyright (c) 2022 Randy Ardiansyah https://github.com/randyardiansyah25/<repo>
 *
 * Created Date: Wednesday, 16/03/2022, 10:32:08
 * Author: Randy Ardiansyah
 *
 * Filename: /home/Documents/workspace/go/src/router-template/delivery/router/registry.go
 * Project : /home/Documents/workspace/go/src/router-template/delivery/router
 *
 * HISTORY:
 * Date                  	By                 	Comments
 * ----------------------	-------------------	--------------------------------------------------------------------------------------------------------------------
 */

package router

import (
	"net/http"
	"router-template/delivery/http/handler"
	"router-template/delivery/http/handler/bni"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(router *gin.Engine) {
	router.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, "Router Template V0.0.0")
	})
	router.GET("/list", handler.GetEmployeListHandler)
	router.POST("/employee", handler.CreateEmployeeHandler)
	router.GET("/employee/:id", handler.GetEmployeeHandler)
	router.PUT("/update_employee", handler.UpdateEmployeeHandler)
	router.DELETE("/delete_employee", handler.DeleteEmployeeHandler)

	router.POST("/topup", handler.TopupBalance)
	router.POST("/send_balance", handler.SendBalance)
	router.POST("/send_notif", handler.SendNotif)
	router.GET("/get_notif/:id", handler.GetNotif)

	router.POST("/firebase_test", handler.TestFirebase)
	router.POST("/pay_with_va", handler.PayWithVa)
	router.POST("/pay_with_qris", handler.PayWithQris)

	router.GET("/grcp/notes", handler.GetNotes)
	router.GET("/grcp/products", handler.GetProducts)

	router_bni := router.Group("/bni")
	router_bni.POST("/token", bni.GetTokenBni)
}
