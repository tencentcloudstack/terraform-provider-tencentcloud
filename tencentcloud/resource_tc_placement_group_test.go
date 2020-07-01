package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudPlacementGroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPlacementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPlacementGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPlacementGroupExists("tencentcloud_placement_group.placement"),
					resource.TestCheckResourceAttr("tencentcloud_placement_group.placement", "name", "tf-test-placement"),
					resource.TestCheckResourceAttr("tencentcloud_placement_group.placement", "type", "HOST"),
					resource.TestCheckResourceAttrSet("tencentcloud_placement_group.placement", "cvm_quota_total"),
					resource.TestCheckResourceAttrSet("tencentcloud_placement_group.placement", "current_num"),
					resource.TestCheckResourceAttrSet("tencentcloud_placement_group.placement", "create_time"),
				),
			},
		},
	})
}

func testAccCheckPlacementGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("placement group %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("placement group id is not set")
		}

		cvmService := CvmService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		placement, err := cvmService.DescribePlacementGroupById(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				placement, err = cvmService.DescribePlacementGroupById(ctx, rs.Primary.ID)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if placement == nil {
			return fmt.Errorf("placement group id is not found")
		}
		return nil
	}
}

func testAccCheckPlacementGroupDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cvmService := CvmService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_placement_group" {
			continue
		}

		placement, err := cvmService.DescribePlacementGroupById(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				placement, err = cvmService.DescribePlacementGroupById(ctx, rs.Primary.ID)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if placement != nil {
			return fmt.Errorf("placement group still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

const testAccPlacementGroup = `
resource "tencentcloud_placement_group" "placement" {
	name = "tf-test-placement"
	type = "HOST"
}
`
