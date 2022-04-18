package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudCbsSnapshotPolicyAttachment(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCbsSnapshotPolicyAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCbsSnapshotPolicyAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCbsSnapshotPolicyAttachmentExists("tencentcloud_cbs_snapshot_policy_attachment.foo"),
					resource.TestCheckResourceAttrSet("tencentcloud_cbs_snapshot_policy_attachment.foo", "storage_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cbs_snapshot_policy_attachment.foo", "snapshot_policy_id"),
				),
			},
		},
	})
}

func testAccCheckCbsSnapshotPolicyAttachmentDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	cbsService := CbsService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cbs_snapshot_policy_attachment" {
			continue
		}
		id := rs.Primary.ID
		idSplit := strings.Split(id, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("tencentcloud_cbs_snapshot_policy_attachment id is illegal: %s", id)
		}
		storageId := idSplit[0]
		policyId := idSplit[1]
		policy, err := cbsService.DescribeAttachedSnapshotPolicy(ctx, storageId, policyId)
		if err != nil {
			return err
		}
		if policy != nil {
			return errors.New("cbs snapshot policy attachment still exists")
		}
	}

	return nil
}

func testAccCheckCbsSnapshotPolicyAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cbs snapshot policy attachment %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return errors.New("cbs snapshot policy attachment id is not set")
		}
		id := rs.Primary.ID
		idSplit := strings.Split(id, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("tencentcloud_cbs_snapshot_policy_attachment id is illegal: %s", id)
		}
		storageId := idSplit[0]
		policyId := idSplit[1]
		cbsService := CbsService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		policy, err := cbsService.DescribeAttachedSnapshotPolicy(ctx, storageId, policyId)
		if err != nil {
			return err
		}
		if policy == nil {
			return errors.New("cbs snapshot policy attachment not exists")
		}
		return nil
	}
}

const testAccCbsSnapshotPolicyAttachmentConfig = defaultVpcVariable + `
resource "tencentcloud_cbs_storage" "foo" {
  availability_zone = var.availability_zone
  storage_size      = 100
  storage_type      = "CLOUD_PREMIUM"
  storage_name      = var.instance_name
}

resource "tencentcloud_cbs_snapshot_policy" "policy" {
  snapshot_policy_name = "tf-test-snapshot-policy"
  repeat_weekdays      = [0, 3]
  repeat_hours         = [0]
  retention_days       = 30
}

resource "tencentcloud_cbs_snapshot_policy_attachment" "foo" {
  storage_id = tencentcloud_cbs_storage.foo.id 
  snapshot_policy_id = tencentcloud_cbs_snapshot_policy.policy.id
}
`
