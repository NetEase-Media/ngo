package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HelloController struct {
}

func NewHelloController() *HelloController {
	return &HelloController{}
}
func (controller *HelloController) GetTest(c *gin.Context) {
	c.String(http.StatusOK, "GetTest")
}
func (controller *HelloController) PostTest(c *gin.Context) {
	c.String(http.StatusOK, "PostTest")
}
func (controller *HelloController) DeleteTest(c *gin.Context) {
	c.String(http.StatusOK, "DeleteTest")
}
func (controller *HelloController) PutTest(c *gin.Context) {
	c.String(http.StatusOK, "PutTest")
}

func (controller *HelloController) OptionsTest(c *gin.Context) {
	c.String(http.StatusOK, "OptionsTest")
}
