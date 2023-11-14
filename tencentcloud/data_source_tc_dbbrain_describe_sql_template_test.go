package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDbbrainDescribeSqlTemplateDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbbrainDescribeSqlTemplateDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_describe_sql_template.describe_sql_template")),
			},
		},
	})
}

const testAccDbbrainDescribeSqlTemplateDataSource = `

data "tencentcloud_dbbrain_describe_sql_template" "describe_sql_template" {
  instance_id = ""
  schema = ""
  sql_text = ""
  product = ""
      }

`
