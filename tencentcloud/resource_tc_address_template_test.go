package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudAddressTemplate_basic_and_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAddressTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAddressTemplate_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_address_template.template", "name", "test"),
					resource.TestCheckResourceAttr("tencentcloud_address_template.template", "addresses.#", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_address_template.template",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAddressTemplate_basic_update_remark,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAddressTemplateExists("tencentcloud_address_template.template"),
					resource.TestCheckResourceAttr("tencentcloud_address_template.template", "name", "test_update"),
					resource.TestCheckResourceAttr("tencentcloud_address_template.template", "addresses.#", "2"),
				),
			},
		},
	})
}

func testAccCheckAddressTemplateDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	vpcService := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
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

func testAccCheckAddressTemplateExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Address template %s is not found", n)
		}

		vpcService := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		_, has, err := vpcService.DescribeAddressTemplateById(ctx, rs.Primary.ID)
		if !has {
			return fmt.Errorf("Address template %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccAddressTemplate_basic = `
resource "tencentcloud_address_template" "template" {
  name = "test"
  addresses = ["1.1.1.1"]
}`

const testAccAddressTemplate_basic_update_remark = `
resource "tencentcloud_address_template" "template" {
  name = "test_update"
  addresses = ["1.1.1.1/24", "1.1.1.0-1.1.1.1"]
}`
