package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTsfClusterResource_basic -v
func TestAccTencentCloudTsfClusterResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTsfClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfCluster,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfClusterExists("tencentcloud_tsf_cluster.cluster"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_cluster.cluster", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_cluster.cluster", "cluster_name", "terraform-test"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_cluster.cluster", "cluster_type", "V"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_cluster.cluster", "vpc_id", "vpc-2hfyray3"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_cluster.cluster", "cluster_desc", "cluster desc"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_cluster.cluster", "tsf_region_id", "ap-guangzhou"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_cluster.cluster", "tags.createdBy", "terraform"),
				),
			},
			{
				ResourceName:      "tencentcloud_tsf_cluster.cluster",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTsfClusterDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tsf_cluster" {
			continue
		}

		res, err := service.DescribeTsfClusterById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res != nil {
			return fmt.Errorf("tsf cluster %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTsfClusterExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		res, err := service.DescribeTsfClusterById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tsf cluster %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTsfCluster = `

resource "tencentcloud_tsf_cluster" "cluster" {
	cluster_name = "terraform-test"
	cluster_type = "V"
	vpc_id = "vpc-2hfyray3"
	# cluster_cidr = ""
	cluster_desc = "cluster desc"
	tsf_region_id = "ap-guangzhou"
  
	tags = {
	  "createdBy" = "terraform"
	}
}

`
