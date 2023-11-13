package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDbbrainSqlFilterResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbbrainSqlFilter,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dbbrain_sql_filter.sql_filter", "id")),
			},
			{
				ResourceName:      "tencentcloud_dbbrain_sql_filter.sql_filter",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDbbrainSqlFilter = `

resource "tencentcloud_dbbrain_sql_filter" "sql_filter" {
  instance_id = &lt;nil&gt;
  session_token {
		user = &lt;nil&gt;
		password = &lt;nil&gt;

  }
  sql_type = &lt;nil&gt;
  filter_key = &lt;nil&gt;
  max_concurrency = &lt;nil&gt;
  duration = &lt;nil&gt;
  product = &lt;nil&gt;
  status = &lt;nil&gt;
}

`
