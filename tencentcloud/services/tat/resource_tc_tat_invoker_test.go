package tat_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctat "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tat"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTatInvokerResource_basic -v
func TestAccTencentCloudTatInvokerResource_basic(t *testing.T) {
	// t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTatInvokerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTatInvoker,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTatInvokerExists("tencentcloud_tat_invoker.invoker"),
					resource.TestCheckResourceAttr("tencentcloud_tat_invoker.invoker", "name", "invoker-test"),
					resource.TestCheckResourceAttr("tencentcloud_tat_invoker.invoker", "type", "SCHEDULE"),
					resource.TestCheckResourceAttr("tencentcloud_tat_invoker.invoker", "instance_ids.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tat_invoker.invoker", "username", "root"),
					resource.TestCheckResourceAttr("tencentcloud_tat_invoker.invoker", "schedule_settings.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tat_invoker.invoker", "schedule_settings.0.policy", "ONCE"),
				),
			},
			{
				ResourceName:      "tencentcloud_tat_invoker.invoker",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTatInvokerDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svctat.NewTatService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tat_invoker" {
			continue
		}

		invoker, err := service.DescribeTatInvoker(ctx, rs.Primary.ID)
		if invoker != nil {
			return fmt.Errorf("tat invoker %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckTatInvokerExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svctat.NewTatService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		invoker, err := service.DescribeTatInvoker(ctx, rs.Primary.ID)
		if invoker == nil {
			return fmt.Errorf("tat invoker %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTatInvokerVar = `
variable "instance_id" {
  default = "` + tcacctest.DefaultInstanceId + `"
}
`
const testAccTatInvoker = testAccTatInvokerVar + testAccTatCommand + `

resource "tencentcloud_tat_invoker" "invoker" {
	name          = "invoker-test"
	type          = "SCHEDULE"
	command_id    = tencentcloud_tat_command.command.id
	instance_ids  = [var.instance_id,]
	username      = "root"
	# parameters = ""
	schedule_settings {
		policy = "ONCE"
		# recurrence = ""
		invoke_time = "2099-11-17T16:00:00Z"
  
	}
  }

`
