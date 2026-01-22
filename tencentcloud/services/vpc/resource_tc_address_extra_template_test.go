package vpc_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudAddressExtraTemplate_basic_and_update(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckAddressExtraTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAddressExtraTemplate_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_address_extra_template.foo", "name", "demo"),
				),
			},
			{
				ResourceName:      "tencentcloud_address_extra_template.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAddressExtraTemplate_basic_update_name,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAddressTemplateExists("tencentcloud_address_extra_template.foo"),
					resource.TestCheckResourceAttr("tencentcloud_address_extra_template.foo", "name", "hello"),
				),
			},
		},
	})
}

func testAccCheckAddressExtraTemplateDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	vpcService := svcvpc.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_address_template" {
			continue
		}

		_, has, err := vpcService.DescribeAddressTemplateById(ctx, rs.Primary.ID)
		if has {
			return fmt.Errorf("address template still exists")
		}
		if err != nil {
			return err
		}
	}
	return nil
}

const testAccAddressExtraTemplate_basic = `
resource "tencentcloud_address_extra_template" "foo" {
  name = "demo"

  addresses_extra {
    address     = "10.0.0.1"
    description = "create by terraform"
  }

  addresses_extra {
    address     = "10.0.1.0/24"
    description = "delete by terraform"
  }

  addresses_extra {
    address     = "10.0.0.1-10.0.0.100"
    description = "modify by terraform"
  }

  tags = {
    createBy = "terraform"
    deleteBy = "terraform"
  }

}`

const testAccAddressExtraTemplate_basic_update_name = `
resource "tencentcloud_address_extra_template" "foo" {
  name = "hello"

  addresses_extra {
    address     = "10.0.0.1"
    description = "create by terraform"
  }

  addresses_extra {
    address     = "10.0.1.0/24"
    description = "delete by terraform"
  }

  addresses_extra {
    address     = "10.0.0.1-10.0.0.100"
    description = "modify by terraform"
  }

  tags = {
    createBy = "terraform"
    deleteBy = "terraform"
  }

}`
