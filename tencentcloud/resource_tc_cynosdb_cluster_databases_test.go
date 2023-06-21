package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbClusterPasswordComplexityResource_basic -v
func TestAccTencentCloudCynosdbClusterDatabasesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClusterDatabasesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbClusterDatabases,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterDatabasesExists("tencentcloud_cynosdb_cluster_databases.cluster_databases"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_databases.cluster_databases", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_databases.cluster_databases", "cluster_id", defaultCynosdbClusterId),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_databases.cluster_databases", "character_set", "utf8"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_databases.cluster_databases", "collate_rule", "utf8_general_ci"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_databases.cluster_databases", "db_name", "terraform-test"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_databases.cluster_databases", "description", "test"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_databases.cluster_databases", "user_host_privileges.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_databases.cluster_databases", "user_host_privileges.0.db_host", "%"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_databases.cluster_databases", "user_host_privileges.0.db_privilege", "READ_WRITE"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_databases.cluster_databases", "user_host_privileges.0.db_user_name", "root"),
				),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_cluster_databases.cluster_databases",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCynosdbClusterDatabasesUp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterDatabasesExists("tencentcloud_cynosdb_cluster_databases.cluster_databases"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_databases.cluster_databases", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_databases.cluster_databases", "cluster_id", defaultCynosdbClusterId),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_databases.cluster_databases", "character_set", "utf8"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_databases.cluster_databases", "collate_rule", "utf8_general_ci"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_databases.cluster_databases", "db_name", "terraform-test"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_databases.cluster_databases", "description", "terraform test"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_databases.cluster_databases", "user_host_privileges.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_databases.cluster_databases", "user_host_privileges.0.db_host", "%"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_databases.cluster_databases", "user_host_privileges.0.db_privilege", "READ_ONLY"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_databases.cluster_databases", "user_host_privileges.0.db_user_name", "root"),
				),
			},
		},
	})
}

func testAccCheckClusterDatabasesDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	cynosdbService := CynosdbService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cynosdb_cluster_databases" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		clusterId := idSplit[0]
		dbName := idSplit[1]

		has, err := cynosdbService.DescribeCynosdbClusterDatabasesById(ctx, clusterId, dbName)
		if err != nil {
			return err
		}
		if has == nil {
			return nil
		}
		return fmt.Errorf("cynosdb cluster databases still exists: %s", rs.Primary.ID)
	}
	return nil
}

func testAccCheckClusterDatabasesExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cynosdb cluster %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cynosdb cluster databases id is not set")
		}
		cynosdbService := CynosdbService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		clusterId := idSplit[0]
		dbName := idSplit[1]

		has, err := cynosdbService.DescribeCynosdbClusterDatabasesById(ctx, clusterId, dbName)
		if err != nil {
			return err
		}
		if has == nil {
			return fmt.Errorf("cynosdb cluster databases doesn't exist: %s", rs.Primary.ID)
		}
		return nil
	}
}

const testAccCynosdbClusterDatabases = CommonCynosdb + `

resource "tencentcloud_cynosdb_cluster_databases" "cluster_databases" {
	cluster_id = var.cynosdb_cluster_id
	db_name = "terraform-test"
	character_set = "utf8"
	collate_rule = "utf8_general_ci"
	user_host_privileges {
	  db_user_name = "root"
	  db_host = "%"
	  db_privilege = "READ_WRITE"
	}
	description = "test"
}

`

const testAccCynosdbClusterDatabasesUp = CommonCynosdb + `

resource "tencentcloud_cynosdb_cluster_databases" "cluster_databases" {
	cluster_id = var.cynosdb_cluster_id
	db_name = "terraform-test"
	character_set = "utf8"
	collate_rule = "utf8_general_ci"
	user_host_privileges {
	  db_user_name = "root"
	  db_host = "%"
	  db_privilege = "READ_ONLY"
	}
	description = "terraform test"
}

`
