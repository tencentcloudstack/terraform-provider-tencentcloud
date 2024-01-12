package cbs_test

import (
	"log"
	"strings"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	localcbs "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cbs"

	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("tencentcloud_cbs_snapshot", &resource.Sweeper{
		Name: "tencentcloud_cbs_snapshot",
		F:    testSweepCbsSnapshot,
	})
}

func testSweepCbsSnapshot(region string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	sharedClient, err := tcacctest.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
	}
	client := sharedClient.(tccommon.ProviderMeta)

	cbsService := localcbs.NewCbsService(client.GetAPIV3Conn())

	instances, err := cbsService.DescribeSnapshotsByFilter(ctx, nil)
	if err != nil {
		return fmt.Errorf("get instance list error: %s", err.Error())
	}

	for _, v := range instances {
		instanceId := v.SnapshotId
		instanceName := v.SnapshotName

		now := time.Now()

		createTime := tccommon.StringToTime(*v.CreateTime)
		interval := now.Sub(createTime).Minutes()
		if strings.HasPrefix(*instanceName, tcacctest.KeepResource) || strings.HasPrefix(*instanceName, tcacctest.DefaultResource) {
			continue
		}
		// less than 30 minute, not delete
		if tccommon.NeedProtect == 1 && int64(interval) < 30 {
			continue
		}

		if err = cbsService.DeleteSnapshot(ctx, *instanceId); err != nil {
			log.Printf("[ERROR] sweep instance %s error: %s", *instanceId, err.Error())
		}
	}

	return nil
}

func TestAccTencentCloudCbsSnapshot(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCbsSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCbsSnapshot,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSnapshotExists("tencentcloud_cbs_snapshot.snapshot"),
					resource.TestCheckResourceAttrSet("tencentcloud_cbs_snapshot.snapshot", "storage_id"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_snapshot.snapshot", "snapshot_name", "tf-test-snapshot"),
				),
			},
			{
				Config: testAccCbsSnapshot_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_cbs_snapshot.snapshot", "snapshot_name", "tf-test-snapshot-update"),
				),
			},
			{
				ResourceName:      "tencentcloud_cbs_snapshot.snapshot",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCbsSnapshotDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	cbsService := localcbs.NewCbsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cbs_snapshot" {
			continue
		}

		snapshot, err := cbsService.DescribeSnapshotById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if snapshot != nil {
			return fmt.Errorf("cbs snapshot still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckSnapshotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cbs snapshot %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cbs snapshot id is set")
		}
		cbsService := localcbs.NewCbsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		snapshot, err := cbsService.DescribeSnapshotById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if snapshot == nil {
			return fmt.Errorf("cbs snapshot is not exist")
		}
		return nil
	}
}

const testAccCbsSnapshot = `
resource "tencentcloud_cbs_storage" "storage" {
	availability_zone = "ap-guangzhou-3"
	storage_size      = 50
	storage_type      = "CLOUD_PREMIUM"
	storage_name      = "tf-test-storage"
}

resource "tencentcloud_cbs_snapshot" "snapshot" {
	storage_id    = tencentcloud_cbs_storage.storage.id
	snapshot_name = "tf-test-snapshot"
}
`

const testAccCbsSnapshot_update = `
resource "tencentcloud_cbs_storage" "storage" {
	availability_zone = "ap-guangzhou-3"
	storage_size      = 50
	storage_type      = "CLOUD_PREMIUM"
	storage_name      = "tf-test-storage"
}

resource "tencentcloud_cbs_snapshot" "snapshot" {
	storage_id    = tencentcloud_cbs_storage.storage.id
	snapshot_name = "tf-test-snapshot-update"
}
`
