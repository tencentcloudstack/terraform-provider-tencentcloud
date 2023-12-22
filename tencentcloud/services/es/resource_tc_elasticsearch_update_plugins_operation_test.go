package es_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudElasticsearchUpdatePluginsOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticsearchUpdatePluginsOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_elasticsearch_update_plugins_operation.update_plugins_operation", "id")),
			},
		},
	})
}

const testAccElasticsearchUpdatePluginsOperation = `

resource "tencentcloud_elasticsearch_update_plugins_operation" "update_plugins_operation" {
	instance_id = "es-5wn36he6"
	install_plugin_list = ["analysis-pinyin"]
	force_restart = false
	force_update = true
}

`
