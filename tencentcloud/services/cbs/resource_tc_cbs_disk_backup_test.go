package cbs_test

import (
	"context"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"

	localcbs "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cbs"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_cbs_disk_backup
	resource.AddTestSweepers("tencentcloud_cbs_disk_backup", &resource.Sweeper{
		Name: "tencentcloud_cbs_disk_backup",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			cli, _ := tcacctest.SharedClientForRegion(r)
			request := cbs.NewDescribeDiskBackupsRequest()
			resp, err := cli.(tccommon.ProviderMeta).GetAPIV3Conn().UseCbsClient().DescribeDiskBackups(request)
			if err != nil {
				return err
			}
			diskBuckups := resp.Response.DiskBackupSet
			for _, diskBuckup := range diskBuckups {
				created, err := time.Parse("2006-01-02 15:04:05", *diskBuckup.CreateTime)
				if err != nil {
					created = time.Now()
				}
				name := *diskBuckup.DiskBackupName
				if tcacctest.IsResourcePersist(name, &created) {
					continue
				}
				buckupId := *diskBuckup.DiskBackupId
				client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
				ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
				service := localcbs.NewCbsService(client)
				err = service.DeleteCbsDiskBackupById(ctx, buckupId)
				if err != nil {
					continue
				}
			}

			return nil
		},
	})
}

func TestAccTencentCloudCbsDiskBackupResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCbsDiskBackup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cbs_disk_backup.disk_backup", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_cbs_disk_backup.disk_backup",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCbsDiskBackup = `
resource "tencentcloud_cbs_disk_backup" "disk_backup" {
  disk_id = "disk-r69pg9vw"
  disk_backup_name = "test-disk-backup"
}
`
