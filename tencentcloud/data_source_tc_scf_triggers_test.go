package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudScfTriggersDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccScfTriggersDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_scf_triggers.triggers")),
			},
		},
	})
}

const testAccScfTriggersDataSource = `

data "tencentcloud_scf_triggers" "triggers" {
  function_name = "testFunction"
  namespace = "testNamespace"
  order_by = "add_time"
  order = "DESC"
  filters {
		name = &lt;nil&gt;
		values = &lt;nil&gt;

  }
  }

`
