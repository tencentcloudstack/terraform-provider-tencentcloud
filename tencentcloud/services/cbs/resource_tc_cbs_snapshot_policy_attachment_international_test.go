package cbs_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	localcbs "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cbs"
)

func TestAccTencentCloudInternationalCbsResource_snapshotPolicyAttachment(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccInternationalCheckCbsSnapshotPolicyAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInternationalCbsSnapshotPolicyAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCbsSnapshotPolicyAttachmentExists("tencentcloud_cbs_snapshot_policy_attachment.foo"),
					resource.TestCheckResourceAttrSet("tencentcloud_cbs_snapshot_policy_attachment.foo", "storage_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cbs_snapshot_policy_attachment.foo", "snapshot_policy_id"),
				),
			},
		},
	})
}

func testAccInternationalCheckCbsSnapshotPolicyAttachmentDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	cbsService := localcbs.NewCbsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cbs_snapshot_policy_attachment" {
			continue
		}
		id := rs.Primary.ID
		idSplit := strings.Split(id, tccommon.FILED_SP)
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

func testAccInternationalCheckCbsSnapshotPolicyAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cbs snapshot policy attachment %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return errors.New("cbs snapshot policy attachment id is not set")
		}
		id := rs.Primary.ID
		idSplit := strings.Split(id, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("tencentcloud_cbs_snapshot_policy_attachment id is illegal: %s", id)
		}
		storageId := idSplit[0]
		policyId := idSplit[1]
		cbsService := localcbs.NewCbsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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

const testAccInternationalCbsSnapshotPolicyAttachmentConfig = tcacctest.DefaultVpcVariable + `
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
