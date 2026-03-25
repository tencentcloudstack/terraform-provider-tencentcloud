package vpn_test

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("tencentcloud_vpn_ssl_client", &resource.Sweeper{
		Name: "tencentcloud_vpn_ssl_client",
		F:    testSweepVpnSslClient,
	})
}

func testSweepVpnSslClient(region string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	sharedClient, err := tcacctest.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
	}
	client := sharedClient.(tccommon.ProviderMeta).GetAPIV3Conn()

	vpcService := svcvpc.NewVpcService(client)

	// Get all VPN SSL Clients
	filters := make(map[string]string)
	response, err := vpcService.DescribeVpnGwSslClientByFilter(ctx, filters)
	if err != nil {
		return fmt.Errorf("get VPN SSL client list error: %s", err.Error())
	}

	if len(response) == 0 {
		return nil
	}

	// Sweep test clients
	for _, v := range response {
		sslClientId := *v.SslVpnClientId
		sslClientName := *v.Name

		// Skip protected resources
		if strings.HasPrefix(sslClientName, tcacctest.KeepResource) || strings.HasPrefix(sslClientName, tcacctest.DefaultResource) {
			continue
		}

		// Only sweep terraform test clients
		if !strings.Contains(sslClientName, "terraform-test-client") {
			continue
		}

		// Check certificate begin time to avoid deleting very new resources in parallel tests
		if v.CertBeginTime != nil {
			certBeginTime := tccommon.StringToTime(*v.CertBeginTime)
			if time.Since(certBeginTime).Minutes() < 5 {
				continue
			}
		}

		if _, err := vpcService.DeleteVpnGatewaySslClient(ctx, sslClientId); err != nil {
			log.Printf("[ERROR] sweep VPN SSL client %s error: %s", sslClientId, err.Error())
		}
	}

	return nil
}

func TestAccTencentCloudVpnSslClientResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckVpnSslClientDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnSslClient_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnSslClientExists("tencentcloud_vpn_ssl_client.client"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpn_ssl_client.client", "ssl_vpn_server_id"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_ssl_client.client", "ssl_vpn_client_name", "terraform-test-client"),
				),
			},
			{
				ResourceName:      "tencentcloud_vpn_ssl_client.client",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudVpnSslClientResource_withTags(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckVpnSslClientDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnSslClient_withTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnSslClientExists("tencentcloud_vpn_ssl_client.client"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpn_ssl_client.client", "ssl_vpn_server_id"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_ssl_client.client", "ssl_vpn_client_name", "terraform-test-client-tags"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_ssl_client.client", "tags.Environment", "test"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_ssl_client.client", "tags.Owner", "terraform"),
				),
			},
		},
	})
}

func TestAccTencentCloudVpnSslClientResource_updateTags(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckVpnSslClientDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnSslClient_withTagsUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnSslClientExists("tencentcloud_vpn_ssl_client.client_update"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_ssl_client.client_update", "tags.Environment", "test"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_ssl_client.client_update", "tags.Owner", "terraform"),
				),
			},
			{
				Config: testAccVpnSslClient_updateTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnSslClientExists("tencentcloud_vpn_ssl_client.client_update"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_ssl_client.client_update", "tags.Environment", "production"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_ssl_client.client_update", "tags.Team", "devops"),
				),
			},
			{
				Config: testAccVpnSslClient_updateTagsEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnSslClientExists("tencentcloud_vpn_ssl_client.client_update"),
					resource.TestCheckNoResourceAttr("tencentcloud_vpn_ssl_client.client_update", "tags.Environment"),
					resource.TestCheckNoResourceAttr("tencentcloud_vpn_ssl_client.client_update", "tags.Team"),
				),
			},
		},
	})
}

func testAccCheckVpnSslClientDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	conn := tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn()
	vpcService := svcvpc.NewVpcService(conn)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_vpn_ssl_client" {
			continue
		}

		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			has, _, e := vpcService.DescribeVpnSslClientById(ctx, rs.Primary.ID)
			if e != nil {
				return tccommon.RetryError(e)
			}
			if !has {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("VPN SSL client still exists: %s", rs.Primary.ID))
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func testAccCheckVpnSslClientExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("VPN SSL client %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("VPN SSL client id is not set")
		}

		conn := tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn()
		vpcService := svcvpc.NewVpcService(conn)

		var has bool
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			h, _, e := vpcService.DescribeVpnSslClientById(ctx, rs.Primary.ID)
			if e != nil {
				return tccommon.RetryError(e)
			}
			has = h
			return nil
		})

		if err != nil {
			return err
		}

		if !has {
			return fmt.Errorf("VPN SSL client not found: %s", rs.Primary.ID)
		}

		return nil
	}
}

const testAccVpnSslClient_basic = tcacctest.DefaultVpnDataSource + `
# SSL Server for testing - you need to have an SSL Server available
# Either use an existing one or create one for testing
variable "vpn_ssl_server_id" {
  default = "vpns-test-placeholder"  # Replace with actual SSL server ID in testing environment
}

resource "tencentcloud_vpn_ssl_client" "client" {
  ssl_vpn_server_id   = var.vpn_ssl_server_id
  ssl_vpn_client_name = "terraform-test-client"
}
`

const testAccVpnSslClient_withTags = tcacctest.DefaultVpnDataSource + `
variable "vpn_ssl_server_id" {
  default = "vpns-test-placeholder"  # Replace with actual SSL server ID in testing environment
}

resource "tencentcloud_vpn_ssl_client" "client" {
  ssl_vpn_server_id   = var.vpn_ssl_server_id
  ssl_vpn_client_name = "terraform-test-client-tags"

  tags = {
    Environment = "test"
    Owner       = "terraform"
  }
}
`

const testAccVpnSslClient_withTagsUpdate = tcacctest.DefaultVpnDataSource + `
variable "vpn_ssl_server_id" {
  default = "vpns-test-placeholder"  # Replace with actual SSL server ID in testing environment
}

resource "tencentcloud_vpn_ssl_client" "client_update" {
  ssl_vpn_server_id   = var.vpn_ssl_server_id
  ssl_vpn_client_name = "terraform-test-client-update-tags"

  tags = {
    Environment = "test"
    Owner       = "terraform"
  }
}
`

const testAccVpnSslClient_updateTags = tcacctest.DefaultVpnDataSource + `
variable "vpn_ssl_server_id" {
  default = "vpns-test-placeholder"  # Replace with actual SSL server ID in testing environment
}

resource "tencentcloud_vpn_ssl_client" "client_update" {
  ssl_vpn_server_id   = var.vpn_ssl_server_id
  ssl_vpn_client_name = "terraform-test-client-update-tags"

  tags = {
    Environment = "production"
    Team        = "devops"
  }
}
`

const testAccVpnSslClient_updateTagsEmpty = tcacctest.DefaultVpnDataSource + `
variable "vpn_ssl_server_id" {
  default = "vpns-test-placeholder"  # Replace with actual SSL server ID in testing environment
}

resource "tencentcloud_vpn_ssl_client" "client_update" {
  ssl_vpn_server_id   = var.vpn_ssl_server_id
  ssl_vpn_client_name = "terraform-test-client-update-tags"
}
`
