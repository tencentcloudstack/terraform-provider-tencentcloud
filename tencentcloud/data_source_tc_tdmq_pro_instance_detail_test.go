package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixTdmqProInstanceDetailDataSource_basic -v
func TestAccTencentCloudNeedFixTdmqProInstanceDetailDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqProInstanceDetailDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_pro_instance_detail.pro_instance_detail"),
				),
			},
		},
	})
}

const testAccTdmqProInstanceDetailDataSource = `
data "tencentcloud_tdmq_pro_instance_detail" "pro_instance_detail" {
  cluster_id = "pulsar-5z3g4227qnwr"
}
`
