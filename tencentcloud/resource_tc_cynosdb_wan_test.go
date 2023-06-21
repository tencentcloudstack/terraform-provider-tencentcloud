package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbWanResource_basic -v
func TestAccTencentCloudCynosdbWanResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCynosdbWanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbWan,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCynosdbWanExists("tencentcloud_cynosdb_wan.wan"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_wan.wan", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_wan.wan", "cluster_id", defaultCynosdbClusterId),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_wan.wan", "instance_grp_id", "cynosdbmysql-grp-lxav0p9z"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_wan.wan", "wan_domain", "gz-cynosdbmysql-grp-lxav0p9z.sql.tencentcdb.com"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_wan.wan", "wan_ip"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_wan.wan", "wan_port", "23790"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_wan.wan", "wan_status", "open"),
				),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_wan.wan",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCynosdbWanDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	cynosdbService := CynosdbService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cynosdb_wan" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		clusterId := idSplit[0]
		// instanceGrpId := idSplit[1]

		has, err := cynosdbService.DescribeClusterInstanceGrps(ctx, clusterId)
		if err != nil {
			return err
		}
		if *has.Response.InstanceGrpInfoList[0].WanStatus == "closed" {
			return nil
		}
		return fmt.Errorf("cynosdb cluster wan still exists: %s", rs.Primary.ID)
	}
	return nil
}

func testAccCheckCynosdbWanExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cynosdb cluster %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cynosdb cluster wan id is not set")
		}
		cynosdbService := CynosdbService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		clusterId := idSplit[0]
		// instanceGrpId := idSplit[1]

		has, err := cynosdbService.DescribeClusterInstanceGrps(ctx, clusterId)
		if err != nil {
			return err
		}
		if *has.Response.InstanceGrpInfoList[0].WanStatus != "open" {
			return fmt.Errorf("cynosdb cluster wan doesn't exist: %s", rs.Primary.ID)
		}
		return nil
	}
}

const testAccCynosdbWan = CommonCynosdb + `

resource "tencentcloud_cynosdb_wan" "wan" {
	cluster_id = var.cynosdb_cluster_id
	instance_grp_id = "cynosdbmysql-grp-lxav0p9z"
}

`
