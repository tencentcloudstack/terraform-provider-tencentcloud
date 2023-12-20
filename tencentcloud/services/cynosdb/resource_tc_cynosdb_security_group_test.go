package cynosdb_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svccynosdb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cynosdb"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
)

func TestAccTencentCloudCynosdbSecurityGroupResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCynosdbSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbSecurityGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCynosdbSecurityGroupExists("tencentcloud_cynosdb_security_group.foo"),
				),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_security_group.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCynosdbSecurityGroupDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svccynosdb.NewCynosdbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cynosdb_security_group" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		clusterId := idSplit[0]
		instanceGroupType := idSplit[1]
		grpsResponse, err := service.DescribeClusterInstanceGrps(ctx, clusterId)
		if err != nil {
			return err
		}
		instanceGrpInfoList := grpsResponse.Response.InstanceGrpInfoList
		if len(instanceGrpInfoList) == 0 {
			return fmt.Errorf("Not fount instanceGrpInfoList")
		}

		var securityGroups []*cynosdb.SecurityGroup

		for _, instanceGrpInfo := range instanceGrpInfoList {
			if *instanceGrpInfo.Type != strings.ToLower(instanceGroupType) {
				continue
			}
			securityGroups, err = service.DescribeCynosdbSecurityGroups(ctx, *instanceGrpInfo.InstanceGrpId)
			if err != nil {
				return err
			}
		}

		if len(securityGroups) != 0 {
			return fmt.Errorf("cynosdb sg instance still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCynosdbSecurityGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cynosdb readonly instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cynosdb readonly instance id is not set")
		}
		service := svccynosdb.NewCynosdbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		clusterId := idSplit[0]
		instanceGroupType := idSplit[1]
		grpsResponse, err := service.DescribeClusterInstanceGrps(ctx, clusterId)
		if err != nil {
			return err
		}
		instanceGrpInfoList := grpsResponse.Response.InstanceGrpInfoList
		if len(instanceGrpInfoList) == 0 {
			return fmt.Errorf("Not fount instanceGrpInfoList")
		}

		var securityGroups []*cynosdb.SecurityGroup

		for _, instanceGrpInfo := range instanceGrpInfoList {
			if *instanceGrpInfo.Type != strings.ToLower(instanceGroupType) {
				continue
			}
			securityGroups, err = service.DescribeCynosdbSecurityGroups(ctx, *instanceGrpInfo.InstanceGrpId)
			if err != nil {
				return err
			}
		}

		if len(securityGroups) == 0 {
			return fmt.Errorf("cynosdb readonly instance doesn't exist: %s", rs.Primary.ID)
		}
		return nil
	}
}

const testAccCynosdbSecurityGroup = tcacctest.CommonCynosdb + `
resource "tencentcloud_cynosdb_security_group" "foo" {
  cluster_id = var.cynosdb_cluster_id
  security_group_ids = [var.cynosdb_cluster_security_group_id]
  instance_group_type = "RO"
}
`
