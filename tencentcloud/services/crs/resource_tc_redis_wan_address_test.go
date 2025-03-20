package crs_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svccrs "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/crs"
)

func TestAccTencentCloudRedisWanAddressResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccTencentCloudRedisInstanceWanAddressDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisWanAddress,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_redis_wan_address.redis_wan_address", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_wan_address.redis_wan_address", "wan_address"),
				),
			},
			{
				ResourceName:      "tencentcloud_redis_wan_address.redis_wan_address",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTencentCloudRedisInstanceWanAddressDestroy(s *terraform.State) error {

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svccrs.NewRedisService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_redis_wan_address" {
			continue
		}
		time.Sleep(5 * time.Second)
		has, _, instance, err := service.CheckRedisOnlineOk(ctx, rs.Primary.ID, 20*tccommon.ReadRetryTimeout)
		if err != nil {
			return err
		}
		if !has {
			return fmt.Errorf("failed to get redis")
		}
		if instance.WanAddress != nil && *instance.WanAddress != "" {
			return fmt.Errorf("redis wanAddress close failed: %v", *instance.WanAddress)
		}
	}
	return nil
}

const testAccRedisWanAddress = `

resource "tencentcloud_redis_wan_address" "redis_wan_address" {
  instance_id = "crs-dekqpd8v"
}
`
