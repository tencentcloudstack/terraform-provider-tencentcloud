package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

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
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_environment_attributes.environment_attributes")),
			},
		},
	})
}

const testAccTdmqEnvironmentAttributesDataSource = `

data "tencentcloud_tdmq_environment_attributes" "environment_attributes" {
    cluster_id = ""
              }

`
