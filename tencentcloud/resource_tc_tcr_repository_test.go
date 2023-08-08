package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func init() {
	resource.AddTestSweepers("tencentcloud_tcr_repository", &resource.Sweeper{
		Name: "tencentcloud_tcr_repository",
		F:    testSweepTCRRepository,
	})
}

// go test -v ./tencentcloud -sweep=ap-shanghai -sweep-run=tencentcloud_tcr_repository
func testSweepTCRRepository(r string) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cli, _ := sharedClientForRegion(r)
	tcrService := TCRService{client: cli.(*TencentCloudClient).apiV3Conn}

	var filters []*tcr.Filter
	filters = append(filters, &tcr.Filter{
		Name:   helper.String("RegistryName"),
		Values: []*string{helper.String(defaultTCRInstanceName)},
	})

	instances, err := tcrService.DescribeTCRInstances(ctx, "", filters)
	if err != nil {
		return err
	}

	if len(instances) == 0 {
		return nil
	}

	instanceId := *instances[0].RegistryId
	// the non-keep namespace will be removed directly when run sweeper tencentcloud_tcr_namespace
	// so... only need to care about the repos under the keep namespace
	repos, err := tcrService.DescribeTCRRepositories(ctx, instanceId, "", "")
	if err != nil {
		return err
	}

	for i := range repos {
		n := repos[i]
		names := strings.Split(*n.Name, "/")
		if len(names) != 2 {
			continue
		}
		repoName := names[1]
		if isResourcePersist(repoName, nil) {
			continue
		}
		err = tcrService.DeleteTCRRepository(ctx, instanceId, *n.Namespace, repoName)
		if err != nil {
			continue
		}
	}
	return nil
}

func TestAccTencentCloudTcrRepository_basic_and_update(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTCRRepositoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTCRRepository_basic,
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-shanghai")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_tcr_repository.mytcr_repository", "name", "test"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_repository.mytcr_repository", "brief_desc", "111"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_repository.mytcr_repository", "description", "111111111111111111111111111111111111"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_repository.mytcr_repository", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_repository.mytcr_repository", "update_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_repository.mytcr_repository", "is_public"),
				),
				Destroy: false,
			},
			{
				ResourceName:      "tencentcloud_tcr_repository.mytcr_repository",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTCRRepository_basic_update_remark,
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-shanghai")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTCRRepositoryExists("tencentcloud_tcr_repository.mytcr_repository"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_repository.mytcr_repository", "brief_desc", "2222"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_repository.mytcr_repository", "description", "211111111111111111111111111111111111"),
				),
			},
		},
	})
}

func testAccCheckTCRRepositoryDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	tcrService := TCRService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tcr_repository" {
			continue
		}
		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) != 3 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}

		instanceId := items[0]
		namespaceName := items[1]
		repositoryName := items[2]
		_, has, err := tcrService.DescribeTCRRepositoryById(ctx, instanceId, namespaceName, repositoryName)
		if has {
			return fmt.Errorf("TCR repository still exists")
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckTCRRepositoryExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("TCR repository %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("TCR repository id is not set")
		}
		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) != 3 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}

		instanceId := items[0]
		namespaceName := items[1]
		repositoryName := items[2]
		tcrService := TCRService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		_, has, err := tcrService.DescribeTCRRepositoryById(ctx, instanceId, namespaceName, repositoryName)
		if !has {
			return fmt.Errorf("TCR repository %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTCRRepository_basic = defaultTCRInstanceData + `

resource "tencentcloud_tcr_repository" "mytcr_repository" {
  instance_id	 = local.tcr_id
  namespace_name = var.tcr_namespace
  name 	         = "test"
  brief_desc 	 = "111"
  description	 = "111111111111111111111111111111111111"
}`

const testAccTCRRepository_basic_update_remark = defaultTCRInstanceData + `
resource "tencentcloud_tcr_repository" "mytcr_repository" {
  instance_id 	 = local.tcr_id
  namespace_name = var.tcr_namespace
  name			 = "test"
  brief_desc 	 = "2222"
  description	 = "211111111111111111111111111111111111"
}`
