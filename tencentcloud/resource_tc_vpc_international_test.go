package tencentcloud

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudInternationalVpcResource_instance(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
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
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
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
