package tpulsar_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixTdmqProInstancesDataSource_basic -v
func TestAccTencentCloudNeedFixTdmqProInstancesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqProInstancesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_pro_instances.pro_instances"),
				),
			},
			{
				Config: testAccTdmqProInstancesDataSourcelFilter,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_pro_instances.pro_instances_filter"),
				),
			},
		},
	})
}

const testAccTdmqProInstancesDataSource = `
data "tencentcloud_tdmq_pro_instances" "pro_instances" {
}
`

const testAccTdmqProInstancesDataSourcelFilter = `
data "tencentcloud_tdmq_pro_instances" "pro_instances_filter" {
  filters {
    name   = "InstanceName"
    values = ["keep"]
  }
}
`
