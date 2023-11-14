package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbParamTemplatesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbParamTemplatesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_param_templates.param_templates")),
			},
		},
	})
}

const testAccCynosdbParamTemplatesDataSource = `

data "tencentcloud_cynosdb_param_templates" "param_templates" {
  engine_versions = &lt;nil&gt;
  template_names = &lt;nil&gt;
  template_ids = &lt;nil&gt;
  db_modes = &lt;nil&gt;
  offset = &lt;nil&gt;
  limit = &lt;nil&gt;
  products = &lt;nil&gt;
  template_types = &lt;nil&gt;
  engine_types = &lt;nil&gt;
  order_by = &lt;nil&gt;
  order_direction = &lt;nil&gt;
  total_count = &lt;nil&gt;
  items {
		id = &lt;nil&gt;
		template_name = &lt;nil&gt;
		template_description = &lt;nil&gt;
		engine_version = &lt;nil&gt;
		db_mode = &lt;nil&gt;
		param_info_set {
			current_value = &lt;nil&gt;
			default = &lt;nil&gt;
			enum_value = &lt;nil&gt;
			max = &lt;nil&gt;
			min = &lt;nil&gt;
			param_name = &lt;nil&gt;
			need_reboot = &lt;nil&gt;
			description = &lt;nil&gt;
			param_type = &lt;nil&gt;
		}

  }
}

`
