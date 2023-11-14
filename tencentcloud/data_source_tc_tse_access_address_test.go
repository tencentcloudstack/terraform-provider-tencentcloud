package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTseAccessAddressDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseAccessAddressDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tse_access_address.access_address")),
			},
		},
	})
}

const testAccTseAccessAddressDataSource = `

data "tencentcloud_tse_access_address" "access_address" {
  instance_id = "ins-xxxxxx"
  vpc_id = "vpc-xxxxxx"
  subnet_id = "subnet-xxxxxx"
  workload = "pushgateway"
  engine_region = "ap-guangzhou"
                  tags = {
    "createdBy" = "terraform"
  }
}

`
