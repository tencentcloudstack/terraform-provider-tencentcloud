package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudCbsSnapshotsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCbsSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCbsSnapshotsDataSource,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSnapshotExists("tencentcloud_cbs_snapshot.snapshot"),
					resource.TestCheckResourceAttr("data.tencentcloud_cbs_snapshots.snapshots", "snapshot_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cbs_snapshots.snapshots", "snapshot_list.0.snapshot_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_cbs_snapshots.snapshots", "snapshot_list.0.snapshot_name", "tf-test-snapshot"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cbs_snapshots.snapshots", "snapshot_list.0.storage_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_cbs_snapshots.snapshots", "snapshot_list.0.storage_size", "50"),
					resource.TestCheckResourceAttr("data.tencentcloud_cbs_snapshots.snapshots", "snapshot_list.0.availability_zone", "ap-guangzhou-3"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cbs_snapshots.snapshots", "snapshot_list.0.percent"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cbs_snapshots.snapshots", "snapshot_list.0.create_time"),
					resource.TestCheckResourceAttr("data.tencentcloud_cbs_snapshots.snapshots", "snapshot_list.0.encrypt", "false"),
				),
			},
		},
	})
}

const testAccCbsSnapshotsDataSource = `
resource "tencentcloud_cbs_storage" "storage" {
	availability_zone = "ap-guangzhou-3"
	storage_size      = 50
	storage_type      = "CLOUD_PREMIUM"
	storage_name      = "tf-test-storage"
}

resource "tencentcloud_cbs_snapshot" "snapshot" {
	storage_id    = "${tencentcloud_cbs_storage.storage.id}"
	snapshot_name = "tf-test-snapshot"
}

data "tencentcloud_cbs_snapshots" "snapshots" {
	snapshot_id = "${tencentcloud_cbs_snapshot.snapshot.id}"
}
`
