package controllers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"github.com/yzimhao/trading_engine/v2/internal/persistence"
)

type UserController struct {
	assetRepo persistence.AssetRepository
}

type inUserContext struct {
	fx.In
	AssetRepo persistence.AssetRepository
}

func NewUserController(in inUserContext) *UserController {
	return &UserController{
		assetRepo: in.AssetRepo,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Captcha  string `json:"captcha,omitempty"`
}

// @Summary user login
// @Description user login
// @ID v1.user.login
// @Tags user
// @Accept json
// @Produce json
// @Param args body LoginRequest true "args"
// @Success 200 {string} any
// @Router /api/v1/user/login [post]
func (u *UserController) Login(ctx *gin.Context) {
	// implement in jwt middleware
}

type RegisterRequest struct {
	Username       string `json:"username" binding:"required"`
	Password       string `json:"password" binding:"required"`
	RepeatPassword string `json:"repeat_password" binding:"required"`
	Email          string `json:"email" binding:"required"`
	Captcha        string `json:"captcha,omitempty"`
}

// @Summary user register
// @Description user register
// @ID v1.user.register
// @Tags user
// @Accept json
// @Produce json
// @Param args body RegisterRequest true "args"
// @Success 200 {string} any
// @Router /api/v1/user/register [post]
func (u *UserController) Register(ctx *gin.Context) {
	//TODO
}
