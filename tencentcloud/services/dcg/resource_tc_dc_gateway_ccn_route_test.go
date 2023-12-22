package dcg_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	svcdcg "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dcg"
)

func TestAccTencentCloudDcgV3RouteBasic(t *testing.T) {
	t.Parallel()

	var rKey = "tencentcloud_dc_gateway_ccn_route.route"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccTencentCloudCdgRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccencentCloudDcgRouteBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudCdgRouteExists(rKey),

					resource.TestCheckResourceAttr(rKey, "cidr_block", "10.1.1.0/32"),

					resource.TestCheckResourceAttrSet(rKey, "dcg_id"),
					resource.TestCheckResourceAttrSet(rKey, "as_path.#"),
				),
			},
		},
	})
}

func testAccTencentCloudCdgRouteExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		service := svcdcg.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		items := strings.Split(rs.Primary.ID, "#")

		if len(items) != 2 {
			return fmt.Errorf("id of resource.tencentcloud_dc_gateway_ccn_route is wrong")
		}

		dcgId, routeId := items[0], items[1]

		_, has, err := service.DescribeDirectConnectGatewayCcnRoute(ctx, dcgId, routeId)

		if has == 0 {
			time.Sleep(5 * time.Second)
			_, has, err = service.DescribeDirectConnectGatewayCcnRoute(ctx, dcgId, routeId)
		}

		if has == 0 {
			time.Sleep(10 * time.Second)
			_, has, err = service.DescribeDirectConnectGatewayCcnRoute(ctx, dcgId, routeId)
		}

		if err != nil {
			return err
		}
		if has != 0 {
			return nil
		}

		return fmt.Errorf("cdg route create fail, %s(%s) not exist on the server.", routeId, dcgId)
	}
}

func testAccTencentCloudCdgRouteDestroy(s *terraform.State) error {

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcdcg.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_dc_gateway_ccn_route" {
			continue
		}

		items := strings.Split(rs.Primary.ID, "#")

		if len(items) != 2 {
			return fmt.Errorf("id of resource.tencentcloud_dc_gateway_ccn_route is wrong")
		}

		dcgId, routeId := items[0], items[1]

		_, has, err := service.DescribeDirectConnectGatewayCcnRoute(ctx, dcgId, routeId)

		if has != 0 {
			time.Sleep(5 * time.Second)
			_, has, err = service.DescribeDirectConnectGatewayCcnRoute(ctx, dcgId, routeId)
		}

		if has != 0 {
			time.Sleep(10 * time.Second)
			_, has, err = service.DescribeDirectConnectGatewayCcnRoute(ctx, dcgId, routeId)
		}

		if err != nil {
			return err
		}

		if has == 0 {
			return nil
		}

		return fmt.Errorf("cdg route delete fail, %s(%s) exists on the server.", routeId, dcgId)
	}
	return nil
}

const TestAccencentCloudDcgRouteBasic = `
resource "tencentcloud_ccn" "main" {
  name        = "ci-temp-test-ccn"
  description = "ci-temp-test-ccn-des"
  qos         = "AG"
}

resource "tencentcloud_dc_gateway" "ccn_main" {
  name                = "ci-cdg-ccn-test"
  network_instance_id = tencentcloud_ccn.main.id
  network_type        = "CCN"
  gateway_type        = "NORMAL"
}

resource "tencentcloud_dc_gateway_ccn_route" "route" {
  dcg_id     = tencentcloud_dc_gateway.ccn_main.id
  cidr_block = "10.1.1.0/32"
}

`
