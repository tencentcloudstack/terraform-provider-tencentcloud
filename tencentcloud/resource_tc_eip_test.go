package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudEip_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEipBasicWithName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipExists("tencentcloud_eip.foo"),
					resource.TestCheckResourceAttr("tencentcloud_eip.foo", "name", "gateway_eip"),
					resource.TestCheckResourceAttr("tencentcloud_eip.foo", "status", "UNBIND"),
					resource.TestCheckResourceAttrSet("tencentcloud_eip.foo", "public_ip"),
				),
			},
			{
				Config: testAccEipBasicWithNewName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipExists("tencentcloud_eip.foo"),
					resource.TestCheckResourceAttr("tencentcloud_eip.foo", "name", "new_name"),
					resource.TestCheckResourceAttr("tencentcloud_eip.foo", "status", "UNBIND"),
					resource.TestCheckResourceAttrSet("tencentcloud_eip.foo", "public_ip"),
				),
			},
			{
				Config: testAccEipBasicWithoutName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipExists("tencentcloud_eip.bar"),
					resource.TestCheckResourceAttr("tencentcloud_eip.bar", "status", "UNBIND"),
					resource.TestCheckResourceAttrSet("tencentcloud_eip.bar", "public_ip"),
				),
			},
			{
				ResourceName:      "tencentcloud_eip.bar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckEipExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("eip %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("eip id is not set")
		}

		vpcService := VpcService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		eip, err := vpcService.DescribeEipById(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				eip, err = vpcService.DescribeEipById(ctx, rs.Primary.ID)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if eip == nil {
			return fmt.Errorf("eip id is not found")
		}
		return nil
	}
}

func testAccCheckEipDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	vpcService := VpcService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_eip" {
			continue
		}

		eip, err := vpcService.DescribeEipById(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				eip, err = vpcService.DescribeEipById(ctx, rs.Primary.ID)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if eip != nil {
			return fmt.Errorf("eip still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

const testAccEipBasicWithName = `
resource "tencentcloud_eip" "foo" {
	name = "gateway_eip"
}
`
const testAccEipBasicWithNewName = `
resource "tencentcloud_eip" "foo" {
	name = "new_name"
}
`

const testAccEipBasicWithoutName = `
resource "tencentcloud_eip" "bar" {
}
`
