package orders

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGenerateOrderId(t *testing.T) {
	Convey("生成订单号", t, func() {
		order_id := generate_order_id("A")
		t.Logf("order_id: %s", order_id)
		So(order_id, ShouldNotBeEmpty)
		So(len(order_id), ShouldEqual, 17)
	})
}
