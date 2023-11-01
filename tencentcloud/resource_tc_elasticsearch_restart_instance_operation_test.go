package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudElasticsearchRestartInstanceOperationResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticsearchRestartInstanceOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_elasticsearch_restart_instance_operation.restart_instance_operation", "id")),
			},
		},
	})
}

const testAccElasticsearchRestartInstanceOperation = `

resource "tencentcloud_elasticsearch_restart_instance_operation" "restart_instance_operation" {
  instance_id = "es-5wn36he6"
}

`
