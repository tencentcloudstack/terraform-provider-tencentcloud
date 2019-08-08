package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudCbsSnapshotPolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCbsSnapshotPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCbsSnapshotPolicy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSnapshotPolicyExists("tencentcloud_cbs_snapshot_policy.snapshot_policy"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_snapshot_policy.snapshot_policy", "snapshot_policy_name", "tf-test-snapshot-policy"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_snapshot_policy.snapshot_policy", "repeat_weekdays.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_snapshot_policy.snapshot_policy", "repeat_weekdays.0", "0"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_snapshot_policy.snapshot_policy", "repeat_weekdays.1", "3"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_snapshot_policy.snapshot_policy", "repeat_hours.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_snapshot_policy.snapshot_policy", "repeat_hours.0", "0"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_snapshot_policy.snapshot_policy", "retention_days", "30"),
				),
			},
			{
				Config: testAccCbsSnapshotPolicy_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_cbs_snapshot_policy.snapshot_policy", "snapshot_policy_name", "tf-snapshot-policy-update"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_snapshot_policy.snapshot_policy", "repeat_weekdays.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_snapshot_policy.snapshot_policy", "repeat_weekdays.0", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_snapshot_policy.snapshot_policy", "repeat_weekdays.1", "4"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_snapshot_policy.snapshot_policy", "repeat_hours.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_snapshot_policy.snapshot_policy", "repeat_hours.0", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_snapshot_policy.snapshot_policy", "retention_days", "7"),
				),
			},
			{
				ResourceName:      "tencentcloud_cbs_snapshot_policy.snapshot_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCbsSnapshotPolicyDestroy(s *terraform.State) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	cbsService := CbsService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cbs_snapshot_policy" {
			continue
		}

		_, err := cbsService.DescribeSnapshotPolicyById(ctx, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("cbs snapshot policy still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckSnapshotPolicyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := GetLogId(nil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cbs snapshot policy %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cbs snapshot policy id is not set")
		}
		cbsService := CbsService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		_, err := cbsService.DescribeSnapshotPolicyById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		return nil
	}
}

const testAccCbsSnapshotPolicy = `
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

resource "tencentcloud_cbs_snapshot_policy" "snapshot_policy" {
  snapshot_policy_name = "tf-test-snapshot-policy"
  repeat_weekdays      = [0, 3]
  repeat_hours         = [0]
  retention_days       = 30
}
`

const testAccCbsSnapshotPolicy_update = `
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

resource "tencentcloud_cbs_snapshot_policy" "snapshot_policy" {
  snapshot_policy_name = "tf-snapshot-policy-update"
  repeat_weekdays      = [1, 4]
  repeat_hours         = [1]
  retention_days       = 7
}
`
