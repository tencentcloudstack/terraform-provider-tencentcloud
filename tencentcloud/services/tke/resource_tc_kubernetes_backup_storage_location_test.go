package tke_test

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctke "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tke"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	testTkeBackupStorageLocationResourceKey = "tencentcloud_kubernetes_backup_storage_location.test_case_backup_storage_location"

	backupStorageLocationNameTemplate = "tf-test-case-backup-storage-location-%d"
	backupLocationBucketTemplate      = "tke-backup-tf-test-case-%d"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_backup_storage_location
	resource.AddTestSweepers("tencentcloud_backup_storage_location", &resource.Sweeper{
		Name: "tencentcloud_backup_storage_location",
		F:    testBackupStorageLocationSweep,
	})
}

func TestAccTencentCloudKubernetesBackupStorageLocationResource_Basic(t *testing.T) {
	t.Parallel()

	randIns := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNum := randIns.Intn(100)
	backupStorageLocationName := fmt.Sprintf(backupStorageLocationNameTemplate, randomNum)
	backupLocationBucket := fmt.Sprintf(backupLocationBucketTemplate, randomNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckBackupStorageLocationDestroy,
		Steps: []resource.TestStep{
			{
				Config: getTestAccTkeBackupStorageLocationConfig(backupStorageLocationName, backupLocationBucket),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(testTkeBackupStorageLocationResourceKey),
					resource.TestCheckResourceAttr(testTkeBackupStorageLocationResourceKey, "name", backupStorageLocationName),
					resource.TestCheckResourceAttr(testTkeBackupStorageLocationResourceKey, "storage_region", "ap-guangzhou"),
					resource.TestCheckResourceAttrSet(testTkeBackupStorageLocationResourceKey, "bucket"),
				),
			},
		},
	})
}

func testBackupStorageLocationSweep(region string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	cli, err := tcacctest.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	log.Printf("testServerlessNodePoolSweep region %s", region)

	client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := svctke.NewTkeService(client)

	// delete all backup storage location
	locations, err := service.DescribeBackupStorageLocations(ctx, []string{})
	if err != nil {
		return err
	}

	for _, l := range locations {
		deleteLocation := l.Name
		if err = service.DeleteBackupStorageLocation(ctx, *deleteLocation); err != nil {
			return fmt.Errorf("error deleting backup storage location: %s", err)
		}
	}

	return nil
}

func testAccCheckBackupStorageLocationDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	client := tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn()
	service := svctke.NewTkeService(client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_kubernetes_backup_storage_location" {
			continue
		}

		locations, err := service.DescribeBackupStorageLocations(ctx, []string{rs.Primary.ID})
		if err != nil {
			return err
		}
		if len(locations) > 0 {
			return fmt.Errorf("backup storage location still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func getTestAccTkeBackupStorageLocationConfig(name, bucket string) string {
	return fmt.Sprintf(testBackupStorageLocationConfigTemplate, bucket, name)
}

const (
	testBackupStorageLocationConfigTemplate = `
data "tencentcloud_user_info" "info" {}
locals {
  app_id = data.tencentcloud_user_info.info.app_id
  uin = data.tencentcloud_user_info.info.uin
  owner_uin = data.tencentcloud_user_info.info.owner_uin
}
resource "tencentcloud_cos_bucket" "back_up_bucket" {
  bucket = "%s-${local.app_id}"
}
resource "tencentcloud_kubernetes_backup_storage_location" "test_case_backup_storage_location" {
  name            = "%s"
  storage_region  = "ap-guangzhou"
  bucket          = tencentcloud_cos_bucket.back_up_bucket.bucket
}
`
)
