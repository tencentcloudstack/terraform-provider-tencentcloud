package vpc_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"
)

func TestAccTencentCloudInternationalVpcResource_instance(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccInternationalCheckVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInternationalVpcConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccInternationalCheckVpcExists("tencentcloud_vpc.foo"),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "cidr_block", "172.16.0.0/16"),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "name", "tf-vpc"),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "is_multicast", "true"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc.foo", "default_route_table_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_vpc.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccInternationalCheckVpcExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svcvpc.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		_, has, err := service.DescribeVpc(ctx, rs.Primary.ID, "", "")
		if err != nil {
			return err
		}
		if has > 0 {
			return nil
		}

		return fmt.Errorf("vpc %s not exists", rs.Primary.ID)
	}
}

func testAccInternationalCheckVpcDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcvpc.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_vpc" {
			continue
		}
		time.Sleep(5 * time.Second)
		_, has, err := service.DescribeVpc(ctx, rs.Primary.ID, "", "")
		if err != nil {
			return err
		}
		if has == 0 {
			return nil
		}
		return fmt.Errorf("vpc %s still exists", rs.Primary.ID)
	}

	return nil
}

const testAccInternationalVpcConfig = `
resource "tencentcloud_vpc" "foo" {
  name       = "tf-vpc"
  cidr_block = "172.16.0.0/16"
}
`
