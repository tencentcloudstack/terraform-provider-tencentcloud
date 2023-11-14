package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDcdbAccountPrivilegesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbAccountPrivileges,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dcdb_account_privileges.account_privileges", "id")),
			},
			{
				ResourceName:      "tencentcloud_dcdb_account_privileges.account_privileges",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDcdbAccountPrivileges = `

resource "tencentcloud_dcdb_account_privileges" "account_privileges" {
  instance_id = "tdsql-c1nl9rpv"
  accounts {
		user = &lt;nil&gt;
		host = &lt;nil&gt;

  }
  global_privileges = &lt;nil&gt;
  database_privileges {
		privileges = &lt;nil&gt;
		database = &lt;nil&gt;

  }
  table_privileges {
		database = &lt;nil&gt;
		table = &lt;nil&gt;
		privileges = &lt;nil&gt;

  }
  column_privileges {
		database = &lt;nil&gt;
		table = &lt;nil&gt;
		column = &lt;nil&gt;
		privileges = &lt;nil&gt;

  }
  view_privileges {
		database = &lt;nil&gt;
		view = &lt;nil&gt;
		privileges = &lt;nil&gt;

  }
}

`
