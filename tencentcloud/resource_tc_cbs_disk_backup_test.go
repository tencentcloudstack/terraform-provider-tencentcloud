package tencentcloud

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_cbs_disk_backup
	resource.AddTestSweepers("tencentcloud_cbs_disk_backup", &resource.Sweeper{
		Name: "tencentcloud_cbs_disk_backup",
		F: func(r string) error {
			logId := getLogId(contextNil)
			cli, _ := sharedClientForRegion(r)
			request := cbs.NewDescribeDiskBackupsRequest()
			resp, err := cli.(*TencentCloudClient).apiV3Conn.UseCbsClient().DescribeDiskBackups(request)
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
				if isResourcePersist(name, &created) {
					continue
				}
				buckupId := *diskBuckup.DiskBackupId
				client := cli.(*TencentCloudClient).apiV3Conn
				ctx := context.WithValue(context.TODO(), logIdKey, logId)
				service := CbsService{client}
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
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
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
