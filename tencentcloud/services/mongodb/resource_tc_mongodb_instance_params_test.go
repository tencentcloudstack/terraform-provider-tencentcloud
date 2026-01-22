package mongodb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudMongodbInstanceParamsResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbInstanceParams,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance_params.mongodb_instance_params", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance_params.mongodb_instance_params", "instance_params.0.key", "cmgo.crossZoneLoadBalancing"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance_params.mongodb_instance_params", "instance_params.0.value", "on"),
				),
			},
			{
				Config: testAccMongodbInstanceParamsUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance_params.mongodb_instance_params", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance_params.mongodb_instance_params", "instance_params.0.key", "cmgo.crossZoneLoadBalancing"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance_params.mongodb_instance_params", "instance_params.0.value", "off"),
				),
			},
		},
	})
}

const testAccMongodbInstanceParams = `
resource "tencentcloud_mongodb_instance_params" "mongodb_instance_params" {
  instance_id = "cmgo-c6k2v891"
  instance_params {
    key = "cmgo.crossZoneLoadBalancing"
    value = "on"
  }
}
`

const testAccMongodbInstanceParamsUpdate = `
resource "tencentcloud_mongodb_instance_params" "mongodb_instance_params" {
  instance_id = "cmgo-c6k2v891"
  instance_params {
    key = "cmgo.crossZoneLoadBalancing"
    value = "off"
  }
}
`
