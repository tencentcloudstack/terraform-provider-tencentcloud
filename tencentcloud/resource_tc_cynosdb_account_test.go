package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbAccountResource_basic -v
func TestAccTencentCloudCynosdbAccountResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCynosdbAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbAccount,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCynosdbAccountExists("tencentcloud_cynosdb_account.account"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_account.account", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account.account", "cluster_id", defaultCynosdbClusterId),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account.account", "account_name", "terraform_test"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account.account", "description", "test"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account.account", "host", "%"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account.account", "max_user_connections", "1"),
				),
			},
			{
				ResourceName:            "tencentcloud_cynosdb_account.account",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password"},
			},
			{
				Config: testAccCynosdbAccountUp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCynosdbAccountExists("tencentcloud_cynosdb_account.account"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_account.account", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account.account", "cluster_id", defaultCynosdbClusterId),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account.account", "account_name", "terraform_test"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account.account", "description", "terraform test"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account.account", "host", "%"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_account.account", "max_user_connections", "2"),
				),
			},
		},
	})
}

func testAccCheckCynosdbAccountDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	cynosdbService := CynosdbService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cynosdb_account" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		clusterId := idSplit[0]
		accountName := idSplit[1]
		host := idSplit[2]

		has, err := cynosdbService.DescribeCynosdbAccountById(ctx, clusterId, accountName, host)
		if err != nil {
			return err
		}
		if has == nil {
			return nil
		}
		return fmt.Errorf("cynosdb cluster account still exists: %s", rs.Primary.ID)
	}
	return nil
}

func testAccCheckCynosdbAccountExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cynosdb cluster %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cynosdb cluster account id is not set")
		}
		cynosdbService := CynosdbService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		clusterId := idSplit[0]
		accountName := idSplit[1]
		host := idSplit[2]

		has, err := cynosdbService.DescribeCynosdbAccountById(ctx, clusterId, accountName, host)
		if err != nil {
			return err
		}
		if has == nil {
			return fmt.Errorf("cynosdb cluster account doesn't exist: %s", rs.Primary.ID)
		}
		return nil
	}
}

const testAccCynosdbAccount = CommonCynosdb + `

resource "tencentcloud_cynosdb_account" "account" {
	cluster_id = var.cynosdb_cluster_id
	account_name = "terraform_test"
	account_password = "Password@1234"
	host = "%"
	description = "test"
	max_user_connections = 1
}

`

const testAccCynosdbAccountUp = CommonCynosdb + `

resource "tencentcloud_cynosdb_account" "account" {
	cluster_id = var.cynosdb_cluster_id
	account_name = "terraform_test"
	account_password = "Password@1234"
	host = "%"
	description = "terraform test"
	max_user_connections = 2
}

`
