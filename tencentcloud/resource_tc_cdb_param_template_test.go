package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbParamTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbParamTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdb_param_template.param_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_cdb_param_template.param_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCdbParamTemplate = `

resource "tencentcloud_cdb_param_template" "param_template" {
  name = &lt;nil&gt;
  description = &lt;nil&gt;
  engine_version = &lt;nil&gt;
  template_id = &lt;nil&gt;
  param_list {
		name = &lt;nil&gt;
		current_value = &lt;nil&gt;

  }
  template_type = "HIGH_STABILITY"
  engine_type = "InnoDB"
}

`
