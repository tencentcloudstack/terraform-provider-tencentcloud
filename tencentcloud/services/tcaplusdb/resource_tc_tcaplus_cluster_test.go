package tcaplusdb_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctcaplusdb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tcaplusdb"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testTcaplusClusterResourceName = "tencentcloud_tcaplus_cluster"
var testTcaplusClusterResourceKey = testTcaplusClusterResourceName + ".test_cluster"

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_tcaplus_cluster
	resource.AddTestSweepers("tencentcloud_tcaplus_cluster", &resource.Sweeper{
		Name: "tencentcloud_tcaplus_cluster",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			cli, _ := tcacctest.SharedClientForRegion(r)
			client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
			service := svctcaplusdb.NewTcaplusService(client)

			clusters, err := service.DescribeClusters(ctx, "", "")
			if err != nil {
				return err
			}

			for i := range clusters {
				c := clusters[i]
				id := *c.ClusterId
				name := *c.ClusterName
				created, err := time.Parse("2006-01-02 15:04:05", *c.CreatedTime)
				if err != nil {
					created = time.Time{}
				}
				if tcacctest.IsResourcePersist(name, &created) {
					continue
				}
				_, err = service.DeleteCluster(ctx, id)
				if err != nil {
					continue
				}
			}

			return nil
		},
	})
}

func TestAccTencentCloudTcaplusClusterResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTcaplusClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTcaplusCluster,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTcaplusClusterExists(testTcaplusClusterResourceKey),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "network_type"),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "password_status"),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "api_access_id"),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "api_access_ip"),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "api_access_port"),
					resource.TestCheckResourceAttr(testTcaplusClusterResourceKey, "idl_type", "PROTO"),
					resource.TestCheckResourceAttr(testTcaplusClusterResourceKey, "cluster_name", "tf_te1_guagua"),
					resource.TestCheckResourceAttr(testTcaplusClusterResourceKey, "old_password_expire_last", "3600"),
				),
			},
			{
				ResourceName:            testTcaplusClusterResourceKey,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"old_password_expire_last", "password"},
			},

			{
				Config: testAccTcaplusClusterUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTcaplusClusterExists(testTcaplusClusterResourceKey),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "network_type"),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "password_status"),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "api_access_id"),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "api_access_ip"),
					resource.TestCheckResourceAttrSet(testTcaplusClusterResourceKey, "api_access_port"),

					resource.TestCheckResourceAttr(testTcaplusClusterResourceKey, "idl_type", "PROTO"),
					resource.TestCheckResourceAttr(testTcaplusClusterResourceKey, "cluster_name", "tf_te1_guagua_2"),
					resource.TestCheckResourceAttr(testTcaplusClusterResourceKey, "password", "aQQ2345677888"),
					resource.TestCheckResourceAttr(testTcaplusClusterResourceKey, "old_password_expire_last", "300"),
				),
			},
		},
	})
}

func testAccCheckTcaplusClusterDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testTcaplusClusterResourceName {
			continue
		}
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		service := svctcaplusdb.NewTcaplusService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		_, has, err := service.DescribeCluster(ctx, rs.Primary.ID)
		if err != nil {
			_, has, err = service.DescribeCluster(ctx, rs.Primary.ID)
		}
		if err != nil {
			return err
		}
		if !has {
			return nil
		} else {
			return fmt.Errorf("delete tcaplus cluster %s fail, still on server", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTcaplusClusterExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		service := svctcaplusdb.NewTcaplusService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		_, has, err := service.DescribeCluster(ctx, rs.Primary.ID)
		if err != nil {
			_, has, err = service.DescribeCluster(ctx, rs.Primary.ID)
		}
		if err != nil {
			return err
		}
		if has {
			return nil
		} else {
			return fmt.Errorf("tcaplus cluster %s not found on server", rs.Primary.ID)

		}
	}
}

const testAccTcaplusCluster string = tcacctest.DefaultVpcSubnets + `
resource "tencentcloud_tcaplus_cluster" "test_cluster" {
  idl_type                 = "PROTO"
  cluster_name             = "tf_te1_guagua"
  vpc_id                   = local.vpc_id
  subnet_id                = local.subnet_id
  password                 = "1qaA2k1wgvfa3ZZZ"
  old_password_expire_last = 3600
}
`
const testAccTcaplusClusterUpdate string = tcacctest.DefaultVpcSubnets + `
resource "tencentcloud_tcaplus_cluster" "test_cluster" {
  idl_type                 = "PROTO"
  cluster_name             = "tf_te1_guagua_2"
  vpc_id                   = local.vpc_id
  subnet_id                = local.subnet_id
  password                 = "aQQ2345677888"
  old_password_expire_last = 300
}
`
