package rum_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	svcrum "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/rum"
)

// go test -i; go test -test.run TestAccTencentCloudRumTawInstanceResource_basic -v
func TestAccTencentCloudRumTawInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckRumTawInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRumTawInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRumTawInstanceExists("tencentcloud_rum_taw_instance.tawInstance"),
					resource.TestCheckResourceAttr("tencentcloud_rum_taw_instance.tawInstance", "area_id", "1"),
					resource.TestCheckResourceAttr("tencentcloud_rum_taw_instance.tawInstance", "charge_status", "1"),
					resource.TestCheckResourceAttr("tencentcloud_rum_taw_instance.tawInstance", "charge_type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_rum_taw_instance.tawInstance", "cluster_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_rum_taw_instance.tawInstance", "data_retention_days", "30"),
					resource.TestCheckResourceAttr("tencentcloud_rum_taw_instance.tawInstance", "instance_desc", "instanceDesc"),
					resource.TestCheckResourceAttr("tencentcloud_rum_taw_instance.tawInstance", "instance_name", "instanceName"),
					resource.TestCheckResourceAttr("tencentcloud_rum_taw_instance.tawInstance", "instance_status", "2"),
				),
			},
			{
				ResourceName:      "tencentcloud_rum_taw_instance.tawInstance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckRumTawInstanceDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcrum.NewRumService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_rum_taw_instance" {
			continue
		}

		instance, err := service.DescribeRumTawInstance(ctx, rs.Primary.ID)
		if instance != nil {
			return fmt.Errorf("rum instance %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckRumTawInstanceExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svcrum.NewRumService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		instance, err := service.DescribeRumTawInstance(ctx, rs.Primary.ID)
		if instance == nil {
			return fmt.Errorf("rum instance %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccRumTawInstance = `

resource "tencentcloud_rum_taw_instance" "tawInstance" {
	area_id = "1"
	charge_type = "1"
	data_retention_days = "30"
	instance_name = "instanceName"
	tags = {
	  createdBy = "terraform"
	}
	instance_desc = "instanceDesc"
}

`
