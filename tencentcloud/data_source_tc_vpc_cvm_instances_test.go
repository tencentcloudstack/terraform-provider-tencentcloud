package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpcCvmInstancesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcCvmInstancesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_cvm_instances.cvm_instances")),
			},
		},
	})
}

const testAccVpcCvmInstancesDataSource = `

data "tencentcloud_vpc_cvm_instances" "cvm_instances" {
  filters {
    name   = "vpc-id"
    values = ["vpc-lh4nqig9"]
  }
}
`
