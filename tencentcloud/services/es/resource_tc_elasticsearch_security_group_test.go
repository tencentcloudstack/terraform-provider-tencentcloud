package es_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svces "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/es"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudElasticsearchSecurityGroupResource_basic -v
func TestAccTencentCloudElasticsearchSecurityGroupResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckElasticsearchSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticsearchSecurityGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticsearchSecurityGroupExists("tencentcloud_elasticsearch_security_group.security_group"),
					resource.TestCheckResourceAttrSet("tencentcloud_elasticsearch_security_group.security_group", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_elasticsearch_security_group.security_group", "security_group_ids.#"),
				),
			},
			{
				ResourceName:      "tencentcloud_elasticsearch_security_group.security_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccElasticsearchSecurityGroupUp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticsearchSecurityGroupExists("tencentcloud_elasticsearch_security_group.security_group"),
					resource.TestCheckResourceAttrSet("tencentcloud_elasticsearch_security_group.security_group", "instance_id"),
				),
			},
		},
	})
}

func testAccCheckElasticsearchSecurityGroupDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	elasticsearchService := svces.NewElasticsearchService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_elasticsearch_security_group" {
			continue
		}

		instance, err := elasticsearchService.DescribeInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				instance, err = elasticsearchService.DescribeInstanceById(ctx, rs.Primary.ID)
				if err != nil {
					return tccommon.RetryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if instance != nil && len(instance.SecurityGroups) > 0 {
			return fmt.Errorf("elasticsearch securityGroup still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckElasticsearchSecurityGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("elasticsearch instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("elasticsearch instance id is not set")
		}
		elasticsearchService := svces.NewElasticsearchService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		instance, err := elasticsearchService.DescribeInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				instance, err = elasticsearchService.DescribeInstanceById(ctx, rs.Primary.ID)
				if err != nil {
					return tccommon.RetryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if instance == nil {
			return fmt.Errorf("elasticsearch securityGroup is not found")
		}
		return nil
	}
}

const testAccElasticsearchSecurityGroup = tcacctest.DefaultEsVariables + `

resource "tencentcloud_elasticsearch_security_group" "security_group" {
    instance_id        = var.instance_id
    security_group_ids = [
        var.security_group_id
    ]
}

`

const testAccElasticsearchSecurityGroupUp = tcacctest.DefaultEsVariables + `

resource "tencentcloud_elasticsearch_security_group" "security_group" {
    instance_id        = var.instance_id
}

`
