package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudElasticsearchRestartKibanaOperationResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticsearchRestartKibanaOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_elasticsearch_restart_kibana_operation.restart_kibana_operation", "id")),
			},
		},
	})
}

const testAccElasticsearchRestartKibanaOperation = `

resource "tencentcloud_elasticsearch_restart_kibana_operation" "restart_kibana_operation" {
  instance_id = "es-5wn36he6"
}

`
