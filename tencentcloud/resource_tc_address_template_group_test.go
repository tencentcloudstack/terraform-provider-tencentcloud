package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudAddressTemplateGroup_basic_and_update(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAddressTemplateGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAddressTemplateGroup_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_address_template_group.group", "name", "test"),
					resource.TestCheckResourceAttr("tencentcloud_address_template_group.group", "template_ids.#", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_address_template_group.group",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAddressTemplateGroup_basic_update_remark,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAddressTemplateGroupExists("tencentcloud_address_template_group.group"),
					resource.TestCheckResourceAttr("tencentcloud_address_template_group.group", "name", "test_update"),
					resource.TestCheckResourceAttr("tencentcloud_address_template_group.group", "template_ids.#", "1"),
				),
			},
		},
	})
}

func testAccCheckAddressTemplateGroupDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	vpcService := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_address_template_group" {
			continue
		}

		_, has, err := vpcService.DescribeAddressTemplateGroupById(ctx, rs.Primary.ID)
		if has {
			return fmt.Errorf("Address template group still exists")
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckAddressTemplateGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Address template group %s is not found", n)
		}

		vpcService := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		_, has, err := vpcService.DescribeAddressTemplateGroupById(ctx, rs.Primary.ID)
		if !has {
			return fmt.Errorf("Address template group %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccAddressTemplateGroup_basic = `
resource "tencentcloud_address_template" "template" {
  name = "test"
  addresses = ["1.1.1.1"]
}

resource "tencentcloud_address_template_group" "group"{
	name = "test"
	template_ids = [tencentcloud_address_template.template.id]
}
`

const testAccAddressTemplateGroup_basic_update_remark = `
resource "tencentcloud_address_template" "template" {
  name = "test"
  addresses = ["1.1.1.1"]
}

resource "tencentcloud_address_template" "templateB" {
  name = "testB"
  addresses = ["1.1.1.1/24", "1.1.1.0-1.1.1.1"]
}

resource "tencentcloud_address_template_group" "group"{
	name = "test_update"
	template_ids = [tencentcloud_address_template.templateB.id]
}
`
