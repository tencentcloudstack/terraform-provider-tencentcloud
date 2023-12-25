package tsf_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctsf "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tsf"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTsfMicroserviceResource_basic -v
func TestAccTencentCloudTsfMicroserviceResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTsfMicroserviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfMicroservice,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfMicroserviceExists("tencentcloud_tsf_microservice.microservice"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_microservice.microservice", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_microservice.microservice", "microservice_name", "test-microservice"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_microservice.microservice", "microservice_desc", "desc-microservice"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_microservice.microservice", "tags.createdBy", "terraform"),
				),
			},
			{
				ResourceName:      "tencentcloud_tsf_microservice.microservice",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTsfMicroserviceDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svctsf.NewTsfService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tsf_microservice" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		namespaceId := idSplit[0]
		microserviceId := idSplit[1]

		res, err := service.DescribeTsfMicroserviceById(ctx, namespaceId, microserviceId, "")
		if err != nil {
			return err
		}

		if res != nil {
			return fmt.Errorf("tsf microservice %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTsfMicroserviceExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		namespaceId := idSplit[0]
		microserviceId := idSplit[1]

		service := svctsf.NewTsfService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		res, err := service.DescribeTsfMicroserviceById(ctx, namespaceId, microserviceId, "")
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tsf microservice %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTsfMicroserviceVar = `
variable "namespace_id" {
	default = "` + tcacctest.DefaultNamespaceId + `"
}
`

const testAccTsfMicroservice = testAccTsfMicroserviceVar + `

resource "tencentcloud_tsf_microservice" "microservice" {
	namespace_id = var.namespace_id
	microservice_name = "test-microservice"
	microservice_desc = "desc-microservice"
	tags = {
	  "createdBy" = "terraform"
	}
}

`
