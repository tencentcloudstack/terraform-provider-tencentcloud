package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_param_templates.param_templates"),
					resource.TestCheckResourceAttr("data.tencentcloud_cynosdb_param_templates.param_templates", "items.#", "1"),
				),
			},
		},
	})
}

const testAccCynosdbParamTemplatesDataSource = `

data "tencentcloud_cynosdb_param_templates" "param_templates" {
	template_names = ["keep-mysql-57-template"]
	db_modes = ["NORMAL"]
}

`
