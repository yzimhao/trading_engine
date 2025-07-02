package product

import (
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
)

type productModule struct {
	router      *provider.Router
	productRepo persistence.ProductRepository
}

func newProductModule(router *provider.Router, repo persistence.ProductRepository) {
	p := productModule{
		productRepo: repo,
		router:      router,
	}
	p.registerRouter()
}

func (p *productModule) registerRouter() {
	productGroup := p.router.APIv1.Group("/product")
	productGroup.GET("/", p.query)
	productGroup.GET("/:symbol", p.detail)

}

// @Summary product list
// @Description get product list
// @ID v1.product.list
// @Tags product
// @Accept json
// @Produce json
// @Success 200 {string} any
// @Router /api/v1/product [get]
func (p *productModule) query(c *gin.Context) {
	//TODO implement
	p.router.ResponseOk(c, nil)
}

// @Summary product detail
// @Description get product detail
// @ID v1.product
// @Tags product
// @Accept json
// @Produce json
// @Query param string true "symbol"
// @Success 200 {string} any
// @Router /api/v1/product/:symbol [get]
func (p *productModule) detail(c *gin.Context) {
	//TODO implement
	p.router.ResponseOk(c, nil)
}
