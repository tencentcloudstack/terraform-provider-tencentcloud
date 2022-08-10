package tencentcloud

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_cbs_storage
	resource.AddTestSweepers("tencentcloud_cbs_storage", &resource.Sweeper{
		Name: "tencentcloud_cbs_storage",
		F: func(r string) error {
			logId := getLogId(contextNil)
			ctx := context.WithValue(context.TODO(), logIdKey, logId)
			cli, _ := sharedClientForRegion(r)
			client := cli.(*TencentCloudClient).apiV3Conn

			service := CbsService{client}

			disks, err := service.DescribeDisksByFilter(ctx, nil)

			if err != nil {
				return err
			}

			for i := range disks {
				disk := disks[i]
				id := *disk.DiskId
				if disk.DiskName == nil {
					continue
				}
				name := *disk.DiskName
				created, err := time.Parse("2006-01-02 15:04:05", *disk.CreateTime)
				if err != nil {
					created = time.Now()
				}
				if isResourcePersist(name, &created) {
					continue
				}
				if *disk.DiskState == CBS_STORAGE_STATUS_ATTACHED {
					continue
				}
				err = service.DeleteDiskById(ctx, id)
				if err != nil {
					continue
				}

			}

			return nil
		},
	})
}

func TestAccTencentCloudCbsStorage_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCbsStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCbsStorage_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageExists("tencentcloud_cbs_storage.storage_basic"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_basic", "storage_name", "tf-storage-basic"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_basic", "storage_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_basic", "storage_size", "50"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_basic", "availability_zone", "ap-guangzhou-3"),
				),
			},
			{
				ResourceName:            "tencentcloud_cbs_storage.storage_basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete"},
			},
		},
	})
}

func TestAccTencentCloudCbsStorage_full(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCbsStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCbsStorage_full,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageExists("tencentcloud_cbs_storage.storage_full"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_full", "storage_name", "tf-storage-full"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_full", "storage_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_full", "storage_size", "50"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_full", "availability_zone", "ap-guangzhou-3"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_full", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_full", "encrypt", "false"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_full", "tags.test", "tf"),
				),
			},
			{
				Config: testAccCbsStorage_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageExists("tencentcloud_cbs_storage.storage_full"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_full", "storage_name", "tf-storage-update"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_full", "storage_size", "60"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_full", "availability_zone", "ap-guangzhou-3"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_full", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_full", "encrypt", "false"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_full", "tags.test", "tf-test"),
				),
			},
		},
	})
}

// Prepaid Disks has quota every period
func TestAccTencentCloudNeedFixCbsStorage_prepaid(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCbsStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCbsStorage_prepaid,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageExists("tencentcloud_cbs_storage.storage_prepaid"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_prepaid", "storage_name", "tf-storage-prepaid"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_prepaid", "charge_type", "PREPAID"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_prepaid", "prepaid_renew_flag", "NOTIFY_AND_AUTO_RENEW"),
				),
			},
			{
				Config: testAccCbsStorage_prepaidupdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageExists("tencentcloud_cbs_storage.storage_prepaid"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_prepaid", "storage_name", "tf-storage-prepaid"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_prepaid", "charge_type", "PREPAID"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_prepaid", "prepaid_renew_flag", "NOTIFY_AND_MANUAL_RENEW"),
				),
			},
		},
	})
}

func TestAccTencentCloudCbsStorage_upgrade(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCbsStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCbsStorage_upgrade,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageExists("tencentcloud_cbs_storage.storage_upgrade"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_upgrade", "storage_name", "tf-storage-upgrade"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_upgrade", "charge_type", "POSTPAID_BY_HOUR"),
				),
			},
			{
				SkipFunc: func() (bool, error) {
					fmt.Printf("Step1 should skip because Prepaid Disks refund has quota every period\n")
					return true, nil
				},
				Config: testAccCbsStorage_upgradeupdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageExists("tencentcloud_cbs_storage.storage_upgrade"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_upgrade", "storage_name", "tf-storage-upgrade"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_upgrade", "charge_type", "PREPAID"),
				),
			},
		},
	})
}

func testAccCheckCbsStorageDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	cbsService := CbsService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cbs_storage" {
			continue
		}

		storage, err := cbsService.DescribeDiskById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if storage != nil {
			return fmt.Errorf("cbs storage still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckStorageExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cbs storage %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cbs storage id is not set")
		}
		cbsService := CbsService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		storage, err := cbsService.DescribeDiskById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if storage == nil {
			return fmt.Errorf("cbs storage is not exist")
		}
		return nil
	}
}

const testAccCbsStorage_basic = `
resource "tencentcloud_cbs_storage" "storage_basic" {
	storage_type      = "CLOUD_PREMIUM"
	storage_name      = "tf-storage-basic"
	storage_size      = 50
	availability_zone = "ap-guangzhou-3"
}
`

const testAccCbsStorage_full = `
resource "tencentcloud_cbs_storage" "storage_full" {
	storage_type      = "CLOUD_PREMIUM"
	storage_name      = "tf-storage-full"
	storage_size      = 50
	availability_zone = "ap-guangzhou-3"
	project_id = 0
	encrypt = false
	tags = {
		test = "tf"
	}
}
`
const testAccCbsStorage_update = `
resource "tencentcloud_cbs_storage" "storage_full" {
	storage_type      = "CLOUD_PREMIUM"
	storage_name      = "tf-storage-update"
	storage_size      = 60
	availability_zone = "ap-guangzhou-3"
	project_id = 0
	encrypt = false
	tags = {
		test = "tf-test"
	}
}
`

const testAccCbsStorage_prepaid = `
resource "tencentcloud_cbs_storage" "storage_prepaid" {
	storage_type      = "CLOUD_PREMIUM"
	storage_name      = "tf-storage-prepaid"
	storage_size      = 50
	availability_zone = "ap-guangzhou-3"
	charge_type			= "PREPAID"
	prepaid_renew_flag = "NOTIFY_AND_AUTO_RENEW"
	prepaid_period = 1
	project_id = 0
	encrypt = false
	tags = {
		test = "tf"
	}
	force_delete = true
}
`
const testAccCbsStorage_prepaidupdate = `
resource "tencentcloud_cbs_storage" "storage_prepaid" {
	storage_type      = "CLOUD_PREMIUM"
	storage_name      = "tf-storage-prepaid"
	storage_size      = 50
	charge_type			= "PREPAID"
	prepaid_renew_flag = "NOTIFY_AND_MANUAL_RENEW"
	prepaid_period = 1
	availability_zone = "ap-guangzhou-3"
	project_id = 0
	encrypt = false
	tags = {
		test = "tf"
	}
	force_delete = true
}
`

const testAccCbsStorage_upgrade = `
resource "tencentcloud_cbs_storage" "storage_upgrade" {
	storage_type      = "CLOUD_PREMIUM"
	storage_name      = "tf-storage-upgrade"
	storage_size      = 50
	availability_zone = "ap-guangzhou-3"
	charge_type       = "POSTPAID_BY_HOUR"
}
`

const testAccCbsStorage_upgradeupdate = `
resource "tencentcloud_cbs_storage" "storage_upgrade" {
	storage_type      = "CLOUD_PREMIUM"
	storage_name      = "tf-storage-upgrade"
	storage_size      = 50
	availability_zone = "ap-guangzhou-3"
	charge_type			= "PREPAID"
	prepaid_renew_flag = "NOTIFY_AND_MANUAL_RENEW"
	prepaid_period = 1
	force_delete = true
}
`
