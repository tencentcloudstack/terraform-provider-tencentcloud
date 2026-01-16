package apm_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudApmInstancesDataSource_basic -v
func TestAccTencentCloudApmInstancesDataSource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccApmInstancesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_apm_instances.instances"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_apm_instances.instances", "instance_list.#"),
				),
			},
		},
	})
}

const testAccApmInstancesDataSource = `
data "tencentcloud_apm_instances" "instances" {
}
`
