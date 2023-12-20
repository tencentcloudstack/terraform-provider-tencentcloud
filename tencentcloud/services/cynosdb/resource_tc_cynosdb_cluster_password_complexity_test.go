package cynosdb_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svccynosdb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cynosdb"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbClusterPasswordComplexityResource_basic -v
func TestAccTencentCloudCynosdbClusterPasswordComplexityResource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCynosdbClusterPasswordComplexityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbClusterPasswordComplexity,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCynosdbClusterPasswordComplexityExists("tencentcloud_cynosdb_cluster_password_complexity.cluster_password_complexity"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_password_complexity.cluster_password_complexity", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_password_complexity.cluster_password_complexity", "cluster_id", tcacctest.DefaultCynosdbClusterId),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_password_complexity.cluster_password_complexity", "validate_password_length", "8"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_password_complexity.cluster_password_complexity", "validate_password_mixed_case_count", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_password_complexity.cluster_password_complexity", "validate_password_special_char_count", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_password_complexity.cluster_password_complexity", "validate_password_number_count", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_password_complexity.cluster_password_complexity", "validate_password_policy", "STRONG"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_password_complexity.cluster_password_complexity", "validate_password_dictionary.#", "3"),
				),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_cluster_password_complexity.cluster_password_complexity",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCynosdbClusterPasswordComplexityUp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCynosdbClusterPasswordComplexityExists("tencentcloud_cynosdb_cluster_password_complexity.cluster_password_complexity"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_password_complexity.cluster_password_complexity", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_password_complexity.cluster_password_complexity", "cluster_id", tcacctest.DefaultCynosdbClusterId),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_password_complexity.cluster_password_complexity", "validate_password_length", "10"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_password_complexity.cluster_password_complexity", "validate_password_mixed_case_count", "2"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_password_complexity.cluster_password_complexity", "validate_password_special_char_count", "2"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_password_complexity.cluster_password_complexity", "validate_password_number_count", "2"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_password_complexity.cluster_password_complexity", "validate_password_policy", "STRONG"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_password_complexity.cluster_password_complexity", "validate_password_dictionary.#", "2"),
				),
			},
		},
	})
}

func testAccCheckCynosdbClusterPasswordComplexityDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	cynosdbService := svccynosdb.NewCynosdbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cynosdb_cluster_password_complexity" {
			continue
		}

		has, err := cynosdbService.DescribeCynosdbClusterPasswordComplexityById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if has != nil {
			return nil
		}
		return fmt.Errorf("cynosdb cluster password complexity still exists: %s", rs.Primary.ID)
	}
	return nil
}

func testAccCheckCynosdbClusterPasswordComplexityExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cynosdb cluster %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cynosdb cluster password complexity id is not set")
		}
		cynosdbService := svccynosdb.NewCynosdbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		has, err := cynosdbService.DescribeCynosdbClusterPasswordComplexityById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if has == nil {
			return fmt.Errorf("cynosdb cluster password complexity doesn't exist: %s", rs.Primary.ID)
		}
		return nil
	}
}

const testAccCynosdbClusterPasswordComplexity = tcacctest.CommonCynosdb + `

resource "tencentcloud_cynosdb_cluster_password_complexity" "cluster_password_complexity" {
	cluster_id                           = var.cynosdb_cluster_id
	validate_password_length             = 8
	validate_password_mixed_case_count   = 1
	validate_password_special_char_count = 1
	validate_password_number_count       = 1
	validate_password_policy             = "STRONG"
	validate_password_dictionary = [
	  "cccc",
	  "xxxx",
	  "zzzz",
	]
}

`

const testAccCynosdbClusterPasswordComplexityUp = tcacctest.CommonCynosdb + `

resource "tencentcloud_cynosdb_cluster_password_complexity" "cluster_password_complexity" {
	cluster_id                           = var.cynosdb_cluster_id
	validate_password_length             = 10
	validate_password_mixed_case_count   = 2
	validate_password_special_char_count = 2
	validate_password_number_count       = 2
	validate_password_policy             = "STRONG"
	validate_password_dictionary = [
	  "cccc",
	  "xxxx",
	]
}

`
