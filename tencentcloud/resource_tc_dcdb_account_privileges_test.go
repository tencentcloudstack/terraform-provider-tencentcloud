package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudDCDBAccountPrivilegesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDCDBAccountPrivileges_basic, defaultDcdbInstanceId, "%"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCDBAccountPrivilegesExists("tencentcloud_dcdb_account_privileges.account_privileges"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_account_privileges.account_privileges", "account.#"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_account_privileges.account_privileges", "account.0.user", "tf_test"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_account_privileges.account_privileges", "account.0.host", "%"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_account_privileges.account_privileges", "global_privileges.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_account_privileges.account_privileges", "table_privileges.#"),
				),
			},
			{
				ResourceName:      "tencentcloud_dcdb_account_privileges.account_privileges",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckDCDBAccountPrivilegesExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("dcdb account privileges  %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("dcdb account privileges id is not set")
		}

		dcdbService := DcdbService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		ret, err := dcdbService.DescribeDcdbAccountPrivilegesById(ctx, rs.Primary.ID, nil, nil, nil, nil)
		if err != nil {
			return err
		}

		if ret.InstanceId == nil {
			return fmt.Errorf("dcdb account privileges not found, instanceId: %v", rs.Primary.ID)
		}

		return nil
	}
}

const testAccDCDBAccountPrivileges_basic = `

resource "tencentcloud_dcdb_account_privileges" "account_privileges" {
  instance_id = "%s"
  account {
		user = "tf_test"
		host = "%s"
  }
  global_privileges = ["SELECT","INSERT","CREATE"]
  database_privileges {
		privileges = ["SELECT","INSERT","UPDATE","INDEX","CREATE"]
		database = "tf_test_db"
  }

  table_privileges {
		database = "tf_test_db"
		table = "tf_test_table"
		privileges = ["SELECT","INSERT","UPDATE","DROP","CREATE"]

  }

}

`
