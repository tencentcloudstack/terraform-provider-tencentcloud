package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func TestAccTencentCloudDbbrainSqlFilterResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDbbrainSqlFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDbbrainSqlFilter(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDbbrainSqlFilterExists("tencentcloud_dbbrain_sql_filter.sql_filter"),
					resource.TestCheckResourceAttrSet("tencentcloud_dbbrain_sql_filter.sql_filter", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dbbrain_sql_filter.sql_filter", "session_token.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_dbbrain_sql_filter.sql_filter", "session_token.0.user"),
					resource.TestCheckResourceAttrSet("tencentcloud_dbbrain_sql_filter.sql_filter", "session_token.0.password"),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_sql_filter.sql_filter", "sql_type", "SELECT"),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_sql_filter.sql_filter", "filter_key", "test"),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_sql_filter.sql_filter", "max_concurrency", "10"),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_sql_filter.sql_filter", "duration", "3600"),
				),
			},
			{
				Config: testAccDbbrainSqlFilter_update(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDbbrainSqlFilterExists("tencentcloud_dbbrain_sql_filter.sql_filter"),
					resource.TestCheckResourceAttrSet("tencentcloud_dbbrain_sql_filter.sql_filter", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dbbrain_sql_filter.sql_filter", "session_token.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_dbbrain_sql_filter.sql_filter", "session_token.0.user"),
					resource.TestCheckResourceAttrSet("tencentcloud_dbbrain_sql_filter.sql_filter", "session_token.0.password"),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_sql_filter.sql_filter", "sql_type", "SELECT"),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_sql_filter.sql_filter", "filter_key", "test"),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_sql_filter.sql_filter", "max_concurrency", "10"),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_sql_filter.sql_filter", "duration", "3600"),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_sql_filter.sql_filter", "product", "mysql"),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_sql_filter.sql_filter", "status", "TERMINATED"),
				),
			},
			// {
			// 	ResourceName:            "tencentcloud_dbbrain_sql_filter.sql_filter",
			// 	ImportState:             true,
			// 	ImportStateVerify:       true,
			// 	ImportStateVerifyIgnore: []string{"duration", "session_token"},
			// },
		},
	})
}

func testAccCheckDbbrainSqlFilterDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	dbbrainService := DbbrainService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_dbbrain_sql_filter" {
			continue
		}

		ids := strings.Split(rs.Primary.ID, FILED_SP)
		if len(ids) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		instanceId := helper.String(ids[0])
		filterId := helper.String(ids[1])

		filter, err := dbbrainService.DescribeDbbrainSqlFilter(ctx, instanceId, filterId)
		if err != nil {
			return err
		}

		if filter != nil {
			return fmt.Errorf("Dbbrain sql filter still exist, Id: %v", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckDbbrainSqlFilterExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		dbbrainService := DbbrainService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("Dbbrain sql filter  %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Dbbrain sql filter id is not set")
		}

		ids := strings.Split(rs.Primary.ID, FILED_SP)
		if len(ids) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		instanceId := helper.String(ids[0])
		filterId := helper.String(ids[1])

		filter, err := dbbrainService.DescribeDbbrainSqlFilter(ctx, instanceId, filterId)
		if err != nil {
			return err
		}

		if filter == nil {
			return fmt.Errorf("Dbbrain sql filter not found, Id: %v", rs.Primary.ID)
		}
		return nil
	}
}

func testAccDbbrainSqlFilter() string {
	return fmt.Sprintf(`%s

resource "tencentcloud_dbbrain_sql_filter" "sql_filter" {
  instance_id = local.mysql_id
  session_token {
    user = "keep_dbbrain"
	password = "Test@123456#"
  }
  sql_type = "SELECT"
  filter_key = "test"
  max_concurrency = 10
  duration = 3600
}
`, CommonPresetMysql)
}

func testAccDbbrainSqlFilter_update() string {
	return fmt.Sprintf(`%s
resource "tencentcloud_dbbrain_sql_filter" "sql_filter" {
  instance_id = local.mysql_id
  session_token {
    user = "keep_dbbrain"
	password = "Test@123456#"
  }
  sql_type = "SELECT"
  filter_key = "test"
  max_concurrency = 10
  duration = 3600
  product = "mysql"
  status = "TERMINATED"
}
`, CommonPresetMysql)
}
