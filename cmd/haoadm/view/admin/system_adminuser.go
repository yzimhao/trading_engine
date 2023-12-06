package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/cmd/haoadm/models"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

func AdminuserList(ctx *gin.Context) {
	db := app.Database().NewSession()
	defer db.Close()

	if ctx.Request.Method == "GET" {
		page := utils.S2Int(ctx.Query("page"))
		limit := utils.S2Int(ctx.Query("limit"))
		searchParams := ctx.Query("searchParams")

		if page <= 0 {
			page = 1
		}
		if limit <= 0 {
			limit = 10
		}
		offset := (page - 1) * limit

		data := []models.Adminuser{}

		q := db.Table(new(models.Adminuser))

		cond := q.Conds()
		err := q.OrderBy("id desc").Limit(limit, offset).Find(&data)
		if err != nil {
			render(ctx, 1, err.Error(), 0, "")
			return
		}

		total, _ := db.Table(new(models.Adminuser)).And(cond).Count()
		if ctx.Query("api") == "1" {
			render(ctx, 0, "", int(total), data)
		} else {
			ctx.HTML(200, "system_adminuser_list", gin.H{
				"searchParams": searchParams,
			})
		}
		return
	}
}

func AdminuserAdd(ctx *gin.Context) {
	id := utils.S2Int64(ctx.Query("id"))

	db := app.Database().NewSession()
	defer db.Close()

	if ctx.Request.Method == "GET" {
		data := models.Adminuser{Id: id}

		if id > 0 {
			db.Table(new(models.Adminuser)).Get(&data)
		}

		ctx.HTML(200, "system_adminuser_add", gin.H{
			"data":             data,
			"adminuser_status": models.AdminuserStatus.List(models.AdminuserStatusDisable),
		})
	} else {
		data := models.Adminuser{
			Id:       id,
			Username: ctx.PostForm("username"),
			Password: ctx.PostForm("password"),
			Email:    ctx.PostForm("email"),
			Mobile:   ctx.PostForm("mobile"),
			// Role:     ctx.PostForm("role"),
			Status: models.AdminuserStatus(ctx.PostForm("status")),
		}

		var err error
		if id > 0 {
			_, err = db.Table(new(models.Adminuser)).Where("id=?", id).Update(&data)
		} else {
			_, err = db.Table(new(models.Adminuser)).Insert(&data)
		}

		if err != nil {
			utils.ResponseFailJson(ctx, err.Error())
			return
		}
		utils.ResponseOkJson(ctx, "")
	}
}
