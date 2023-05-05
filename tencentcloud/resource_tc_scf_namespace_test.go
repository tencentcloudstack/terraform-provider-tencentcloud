package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pkg/errors"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_scf_namespace
	resource.AddTestSweepers("tencentcloud_scf_namespace", &resource.Sweeper{
		Name: "tencentcloud_scf_namespace",
		F: func(r string) error {
			logId := getLogId(contextNil)
			ctx := context.WithValue(context.TODO(), logIdKey, logId)
			cli, _ := sharedClientForRegion(r)
			client := cli.(*TencentCloudClient).apiV3Conn
			service := ScfService{client: client}
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
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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

		service := ScfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

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
		client := testAccProvider.Meta().(*TencentCloudClient).apiV3Conn
		service := ScfService{client: client}

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
