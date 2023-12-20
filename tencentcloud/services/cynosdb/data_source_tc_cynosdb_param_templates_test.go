package cynosdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCynosdbParamTemplatesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbParamTemplatesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_param_templates.param_templates"),
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
