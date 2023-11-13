package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbParamTemplatesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbParamTemplatesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cdb_param_templates.param_templates")),
			},
		},
	})
}

const testAccCdbParamTemplatesDataSource = `

data "tencentcloud_cdb_param_templates" "param_templates" {
  engine_versions = &lt;nil&gt;
  engine_types = &lt;nil&gt;
  template_names = &lt;nil&gt;
  template_ids = &lt;nil&gt;
  total_count = &lt;nil&gt;
  items {
		template_id = &lt;nil&gt;
		name = &lt;nil&gt;
		description = &lt;nil&gt;
		engine_version = &lt;nil&gt;
		template_type = &lt;nil&gt;

  }
}

`
