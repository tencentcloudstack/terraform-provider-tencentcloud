package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -i; go test -test.run TestAccTencentCloudTsfGroupResource_basic -v
func TestAccTencentCloudTsfGroupResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTsfGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfGroupExists("tencentcloud_tsf_group.group"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_group.group", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_group.group", "application_id", "application-ym9mxmza"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_group.group", "namespace_id", "namespace-vjlkzkgy"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_group.group", "group_name", "terraform-test"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_group.group", "cluster_id", "cluster-ym9mxm3a"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_group.group", "group_desc", "terraform group desc"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_group.group", "group_resource_type", "DEF"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_group.group", "alias", "terraform desc"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_group.group", "tags.createdBy", "terraform"),
				),
			},
			{
				ResourceName:      "tencentcloud_tsf_group.group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTsfGroupDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tsf_group" {
			continue
		}

		res, err := service.DescribeTsfGroupById(ctx, rs.Primary.ID)
		if err != nil {
			ee, ok := err.(*sdkErrors.TencentCloudSDKError)
			if !ok {
				return err
			}
			if ee.Code == "ResourceNotFound.GroupNotExist" {
				return nil
			}
			return err
		}

		if res != nil {
			return fmt.Errorf("tsf group %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTsfGroupExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		res, err := service.DescribeTsfGroupById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tsf group %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTsfGroup = `

resource "tencentcloud_tsf_group" "group" {
	application_id = "application-ym9mxmza"
	namespace_id = "namespace-vjlkzkgy"
	group_name = "terraform-test"
	cluster_id = "cluster-ym9mxm3a"
	group_desc = "terraform group desc"
	group_resource_type = "DEF"
	alias = "terraform desc"
	tags = {
	  "createdBy" = "terraform"
	}
}

`
