package dcg_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcdcg "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dcg"

	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudDcgV3InstancesBasic(t *testing.T) {
	t.Parallel()

	var rKey = "tencentcloud_dc_gateway.ccn_main"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccTencentCloudCdgInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccencentCloudDcgInstancesBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudCdgInstanceExists(rKey),

					resource.TestCheckResourceAttr(rKey, "name", "ci-cdg-ccn-test"),
					resource.TestCheckResourceAttr(rKey, "network_type", "CCN"),
					resource.TestCheckResourceAttr(rKey, "gateway_type", "NORMAL"),

					resource.TestCheckResourceAttrSet(rKey, "create_time"),
					resource.TestCheckResourceAttrSet(rKey, "cnn_route_type"),
					resource.TestCheckResourceAttrSet(rKey, "enable_bgp"),
					resource.TestCheckResourceAttrSet(rKey, "network_instance_id"),
				),
			},

			{
				ResourceName:      rKey,
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config: TestAccencentCloudDcgInstancesUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudCdgInstanceExists(rKey),

					resource.TestCheckResourceAttr(rKey, "name", "ci-cdg-ccn-test-update"),
					resource.TestCheckResourceAttr(rKey, "network_type", "CCN"),
					resource.TestCheckResourceAttr(rKey, "gateway_type", "NORMAL"),

					resource.TestCheckResourceAttrSet(rKey, "create_time"),
					resource.TestCheckResourceAttrSet(rKey, "cnn_route_type"),
					resource.TestCheckResourceAttrSet(rKey, "enable_bgp"),
					resource.TestCheckResourceAttrSet(rKey, "network_instance_id"),
				),
			},
		},
	})
}

func testAccTencentCloudCdgInstanceExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svcdcg.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		_, has, err := service.DescribeDirectConnectGateway(ctx, rs.Primary.ID)

		if err != nil {
			return err
		}
		if has != 0 {
			return nil
		}

		return fmt.Errorf("cdg create fail, %s not exist on the server.", rs.Primary.ID)
	}
}

func testAccTencentCloudCdgInstanceDestroy(s *terraform.State) error {

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcdcg.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_dc_gateway" {
			continue
		}
		_, has, err := service.DescribeDirectConnectGateway(ctx, rs.Primary.ID)

		if has > 0 {
			time.Sleep(5 * time.Second)
			_, has, err = service.DescribeDirectConnectGateway(ctx, rs.Primary.ID)
		}

		if has > 0 {
			time.Sleep(10 * time.Second)
			_, has, err = service.DescribeDirectConnectGateway(ctx, rs.Primary.ID)
		}

		if err != nil {
			return err
		}

		if has == 0 {
			return nil
		}

		return fmt.Errorf("cdg delete fail, %s  exists on the server.", rs.Primary.ID)
	}
	return nil
}

const TestAccencentCloudDcgInstancesBasic = `
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
`
const TestAccencentCloudDcgInstancesUpdate = `
resource "tencentcloud_ccn" "main" {
  name        = "ci-temp-test-ccn"
  description = "ci-temp-test-ccn-des"
  qos         = "AG"
}

resource "tencentcloud_dc_gateway" "ccn_main" {
  name                = "ci-cdg-ccn-test-update"
  network_instance_id = tencentcloud_ccn.main.id
  network_type        = "CCN"
  gateway_type        = "NORMAL"
}
`
