package crs_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svccrs "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/crs"

	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudRedisAccountResource_basic -v
func TestAccTencentCloudRedisAccountResource_basic(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccTencentCloudRedisAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisAccount,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudRedisAccountExists("tencentcloud_redis_account.account"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_account.account", "id"),
					resource.TestCheckResourceAttr("tencentcloud_redis_account.account", "instance_id", tcacctest.DefaultCrsInstanceId),
					resource.TestCheckResourceAttr("tencentcloud_redis_account.account", "account_name", "account_test"),
					resource.TestCheckResourceAttr("tencentcloud_redis_account.account", "remark", "redis desc"),
					resource.TestCheckResourceAttr("tencentcloud_redis_account.account", "readonly_policy.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_redis_account.account", "privilege", "r"),
				),
			},
			{
				ResourceName:            "tencentcloud_redis_account.account",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password"},
			},
		},
	})
}

func testAccTencentCloudRedisAccountExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		items := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(items) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}
		instanceId := items[0]
		accountName := items[1]

		service := svccrs.NewRedisService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		account, err := service.DescribeRedisAccountById(ctx, instanceId, accountName)
		if err != nil {
			return err
		}
		if account == nil {
			return fmt.Errorf("redis account %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

func testAccTencentCloudRedisAccountDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svccrs.NewRedisService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_redis_account" {
			continue
		}
		time.Sleep(5 * time.Second)

		items := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(items) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}
		instanceId := items[0]
		accountName := items[1]

		account, err := service.DescribeRedisAccountById(ctx, instanceId, accountName)
		if err != nil {
			return err
		}
		if account != nil {
			return fmt.Errorf("redis account %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

const testAccRedisAccountVar = `
variable "instance_id" {
	default = "` + tcacctest.DefaultCrsInstanceId + `"
}
`

const testAccRedisAccount = testAccRedisAccountVar + `

resource "tencentcloud_redis_account" "account" {
	instance_id 	 = var.instance_id
	account_name 	 = "account_test"
	account_password = "test1234"
	remark 			 = "redis desc"
	readonly_policy  = ["master"]
	privilege 		 = "r"
  }

`
