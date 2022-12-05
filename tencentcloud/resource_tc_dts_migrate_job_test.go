package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	resource.AddTestSweepers("tencentcloud_dts_migrate_job", &resource.Sweeper{
		Name: "tencentcloud_dts_migrate_job",
		F:    testSweepDtsMigrateJob,
	})
}

// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_dts_migrate_job
func testSweepDtsMigrateJob(r string) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cli, _ := sharedClientForRegion(r)
	dtsService := DtsService{client: cli.(*TencentCloudClient).apiV3Conn}
	param := map[string]interface{}{}

	ret, err := dtsService.DescribeDtsMigrateJobsByFilter(ctx, param)
	if err != nil {
		return err
	}

	for _, v := range ret {
		delId := *v.JobId

		if strings.HasPrefix(*v.JobName, keepResource) || strings.HasPrefix(*v.JobName, defaultResource) {
			continue
		}

		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			err := dtsService.DeleteDtsMigrateJobById(ctx, delId)
			if err != nil {
				return retryError(err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("[ERROR] sweeper tencentcloud_dts_migrate_job:[%v] failed! reason:[%s]", delId, err.Error())
		}
	}
	return nil
}

func TestAccTencentCloudDtsMigrateJobResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDtsMigrateJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDtsMigrateJob,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDtsMigrateJobExists("tencentcloud_dts_migrate_job.migrate_job"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_migrate_job.migrate_job", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dts_migrate_job.migrate_job", "src_database_type", "mysql"),
					resource.TestCheckResourceAttr("tencentcloud_dts_migrate_job.migrate_job", "src_region", "ap-guangzhou"),
					resource.TestCheckResourceAttr("tencentcloud_dts_migrate_job.migrate_job", "dst_database_type", "cynosdbmysql"),
					resource.TestCheckResourceAttr("tencentcloud_dts_migrate_job.migrate_job", "dst_region", "ap-guangzhou"),
					resource.TestCheckResourceAttr("tencentcloud_dts_migrate_job.migrate_job", "instance_class", "small"),
					resource.TestCheckResourceAttr("tencentcloud_dts_migrate_job.migrate_job", "job_name", "tf_test_migration_job"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_migrate_job.migrate_job", "tags.#"),
				),
			},
			{
				ResourceName:      "tencentcloud_dts_migrate_job.migrate_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckDtsMigrateJobDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	dtsService := DtsService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_dts_migrate_job" {
			continue
		}

		job, err := dtsService.DescribeDtsMigrateJob(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if job != nil {
			status := *job.TradeInfo.TradeStatus
			if status != "isolated" && status != "offlined" {
				return fmt.Errorf("DTS migrate job still exist, Id: %v, status:%s", rs.Primary.ID, status)
			}
		}
	}
	return nil
}

func testAccCheckDtsMigrateJobExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		dtsService := DtsService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("DTS migrate job %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("DTS migrate job id is not set")
		}

		job, err := dtsService.DescribeDtsMigrateJob(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if job == nil {
			return fmt.Errorf("DTS migrate job not found, Id: %v", rs.Primary.ID)
		}
		return nil
	}
}

const testAccDtsMigrateJob = `

resource "tencentcloud_dts_migrate_job" "migrate_job" {
  src_database_type = "mysql"
  dst_database_type = "cynosdbmysql"
  src_region = "ap-guangzhou"
  dst_region = "ap-guangzhou"
  instance_class = "small"
  job_name = "tf_test_migration_job"
  tags {
	tag_key = "aaa"
	tag_value = "bbb"
  }
}

`
