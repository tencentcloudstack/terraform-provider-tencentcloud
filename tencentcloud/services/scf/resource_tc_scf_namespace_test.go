package scf_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcscf "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/scf"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pkg/errors"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_scf_namespace
	resource.AddTestSweepers("tencentcloud_scf_namespace", &resource.Sweeper{
		Name: "tencentcloud_scf_namespace",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			cli, _ := tcacctest.SharedClientForRegion(r)
			client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
			service := svcscf.NewScfService(client)
			info, err := service.DescribeNamespaces(ctx)
			if err != nil {
				return err
			}
			for _, v := range info {
				name := *v.Name
				if !strings.Contains(name, "ci-test-scf") {
					continue
				}
				if err := service.DeleteNamespace(ctx, name); err != nil {
					continue
				}
			}
			return nil
		},
	})
}

func TestAccTencentCloudScfNamespace_basic(t *testing.T) {
	t.Parallel()
	var nsId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckScfNamespaceDestroy(&nsId),
		Steps: []resource.TestStep{
			{
				Config: testAccScfNamespaceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScfNamespaceExists("tencentcloud_scf_namespace.foo", &nsId),
					resource.TestCheckResourceAttr("tencentcloud_scf_namespace.foo", "namespace", "ci-test-scf"),
					resource.TestCheckResourceAttr("tencentcloud_scf_namespace.foo", "description", "test1"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_namespace.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_namespace.foo", "modify_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_namespace.foo", "type"),
				),
			},
			{
				Config: testAccScfNamespaceBasicUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScfNamespaceExists("tencentcloud_scf_namespace.foo", &nsId),
					resource.TestCheckResourceAttr("tencentcloud_scf_namespace.foo", "description", "test2"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_namespace.foo", "modify_time"),
				),
			},
		},
	})
}

func testAccCheckScfNamespaceExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return errors.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no scf namespace id is set")
		}

		service := svcscf.NewScfService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		fn, err := service.DescribeNamespace(context.TODO(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if fn == nil {
			return errors.Errorf("scf namespace not found: %s", rs.Primary.ID)
		}

		*id = rs.Primary.ID

		return nil
	}
}

func testAccCheckScfNamespaceDestroy(id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn()
		service := svcscf.NewScfService(client)

		namespace, err := service.DescribeNamespace(context.TODO(), *id)
		if err != nil {
			return err
		}

		if namespace != nil {
			return fmt.Errorf("scf namespace still exists")
		}

		return nil
	}
}

const testAccScfNamespaceBasic = `
resource "tencentcloud_scf_namespace" "foo" {
  namespace   = "ci-test-scf"
  description = "test1"
}
`

const testAccScfNamespaceBasicUpdate = `
resource "tencentcloud_scf_namespace" "foo" {
  namespace   = "ci-test-scf"
  description = "test2"
}
`
