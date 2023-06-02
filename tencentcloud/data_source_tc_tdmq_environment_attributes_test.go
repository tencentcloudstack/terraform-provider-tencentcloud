package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTdmqEnvironmentAttributesDataSource_basic -v
func TestAccTencentCloudTdmqEnvironmentAttributesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqEnvironmentAttributesDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_environment_attributes.environment_attributes"),
				),
			},
		},
	})
}

const testAccTdmqEnvironmentAttributesDataSource = `
data "tencentcloud_tdmq_environment_attributes" "environment_attributes" {
    environment_id = "keep-ns"
    cluster_id     = "pulsar-9n95ax58b9vn"
}
`
