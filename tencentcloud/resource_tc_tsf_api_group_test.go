package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTsfApiGroupResource_basic -v
func TestAccTencentCloudNeedFixTsfApiGroupResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTsfApiGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApiGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfApiGroupExists("tencentcloud_tsf_api_group.api_group"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_api_group.api_group", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_api_group.api_group", "group_name", "terraform_test_group"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_api_group.api_group", "group_context", "/terraform-test"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_api_group.api_group", "auth_type", "none"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_api_group.api_group", "description", "terraform-test"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_api_group.api_group", "group_type", "ms"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_api_group.api_group", "namespace_name_key_position", "path"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_api_group.api_group", "service_name_key_position", "path"),
				),
			},
			{
				ResourceName:      "tencentcloud_tsf_api_group.api_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTsfApiGroupDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tsf_api_group" {
			continue
		}

		res, err := service.DescribeTsfApiGroupById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res != nil {
			return fmt.Errorf("tsf api group %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTsfApiGroupExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		res, err := service.DescribeTsfApiGroupById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tsf api group %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTsfApiGroup = `

resource "tencentcloud_tsf_api_group" "api_group" {
	group_name = "terraform_test_group"
	group_context = "/terraform-test"
	auth_type = "none"
	description = "terraform-test"
	group_type = "ms"
	gateway_instance_id = "gw-ins-i6mjpgm8"
	# namespace_name_key = "path"
	# service_name_key = "path"
	namespace_name_key_position = "path"
	service_name_key_position = "path"
  }

`
