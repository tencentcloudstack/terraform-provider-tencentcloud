package dc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTencentCloudDcV3InstancesBasic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudDcInstances,
				Check: resource.ComposeTestCheckFunc(
					//name filter
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dc_instances.name_select"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dc_instances.name_select", "instance_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dc_instances.name_select", "name"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudDcInstances = `
data tencentcloud_dc_instances  name_select {
    name ="x"
}
`
