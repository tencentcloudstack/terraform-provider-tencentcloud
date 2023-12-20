package cynosdb_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svccynosdb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cynosdb"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbWanResource_basic -v
func TestAccTencentCloudCynosdbWanResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCynosdbWanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbWan,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCynosdbWanExists("tencentcloud_cynosdb_wan.wan"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_wan.wan", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_wan.wan", "cluster_id", tcacctest.DefaultCynosdbClusterId),
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
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	cynosdbService := svccynosdb.NewCynosdbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cynosdb_wan" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cynosdb cluster %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cynosdb cluster wan id is not set")
		}
		cynosdbService := svccynosdb.NewCynosdbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
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

const testAccCynosdbWan = tcacctest.CommonCynosdb + `

resource "tencentcloud_cynosdb_wan" "wan" {
	cluster_id = var.cynosdb_cluster_id
	instance_grp_id = "cynosdbmysql-grp-lxav0p9z"
}

`
