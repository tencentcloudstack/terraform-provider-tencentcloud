package trabbit_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctdmq "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdmq"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -i; go test -test.run TestAccTencentCloudTdmqRabbitmqVipInstanceResource_basic -v
func TestAccTencentCloudTdmqRabbitmqVipInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY)
		},
		CheckDestroy: testAccCheckTdmqRabbitmqVipInstanceDestroy,
		Providers:    tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRabbitmqVipInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTdmqRabbitmqVipInstanceExists("tencentcloud_tdmq_rabbitmq_vip_instance.example"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_vip_instance.example", "id"),
				),
			},
			{
				Config: testAccTdmqRabbitmqVipInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTdmqRabbitmqVipInstanceExists("tencentcloud_tdmq_rabbitmq_vip_instance.example"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_vip_instance.example", "id"),
				),
			},
		},
	})
}

func testAccCheckTdmqRabbitmqVipInstanceExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("rabbitmq vip instance  %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("rabbitmq vip instance id is not set")
		}

		service := svctdmq.NewTdmqService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		id := rs.Primary.ID

		ret, err := service.DescribeTdmqRabbitmqVipInstanceById(ctx, id)
		if err != nil {
			return err
		}

		if ret == nil {
			return fmt.Errorf("tdmq rabbitmq vip instance not found, id: %v", id)
		}

		return nil
	}
}

func testAccCheckTdmqRabbitmqVipInstanceDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svctdmq.NewTdmqService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tdmq_rabbitmq_vip_instance" {
			continue
		}

		id := rs.Primary.ID
		ret, err := service.DescribeTdmqRabbitmqVipInstanceById(ctx, id)
		if err != nil {
			code := err.(*sdkErrors.TencentCloudSDKError).Code
			if code == "InternalError" || code == "FailedOperation" {
				return nil
			}

			return err
		}

		if ret != nil {
			return fmt.Errorf("tdmq rabbitmq vip instance exist, id: %v", id)
		}
	}

	return nil
}

const testAccTdmqRabbitmqVipInstance = tcacctest.DefaultVpcSubnets + `
data "tencentcloud_availability_zones" "zones" {}

resource "tencentcloud_tdmq_rabbitmq_vip_instance" "example" {
  zone_ids                              = [data.tencentcloud_availability_zones.zones.zones.0.id]
  vpc_id                                = local.vpc_id
  subnet_id                             = local.subnet_id
  cluster_name                          = "tf-example-rabbitmq-vip-instance"
  node_spec                             = "rabbit-vip-basic-1"
  node_num                              = 1
  storage_size                          = 200
  enable_create_default_ha_mirror_queue = false
  auto_renew_flag                       = true
  time_span                             = 1
}
`

const testAccTdmqRabbitmqVipInstanceUpdate = tcacctest.DefaultVpcSubnets + `
data "tencentcloud_availability_zones" "zones" {}

resource "tencentcloud_tdmq_rabbitmq_vip_instance" "example" {
  zone_ids                              = [data.tencentcloud_availability_zones.zones.zones.0.id]
  vpc_id                                = local.vpc_id
  subnet_id                             = local.subnet_id
  cluster_name                          = "tf-example-rabbitmq-vip-instance-update"
  node_spec                             = "rabbit-vip-basic-1"
  node_num                              = 1
  storage_size                          = 200
  enable_create_default_ha_mirror_queue = false
  auto_renew_flag                       = true
  time_span                             = 1
}
`
