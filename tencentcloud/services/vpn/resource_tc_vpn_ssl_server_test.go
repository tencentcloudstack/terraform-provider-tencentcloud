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
	resource.AddTestSweepers("tencentcloud_vpn_ssl_server", &resource.Sweeper{
		Name: "tencentcloud_vpn_ssl_server",
		F:    testSweepVpnSslServer,
	})
}

func testSweepVpnSslServer(region string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	sharedClient, err := tcacctest.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
	}
	client := sharedClient.(tccommon.ProviderMeta).GetAPIV3Conn()

	vpcService := svcvpc.NewVpcService(client)

	// Get all VPN SSL Servers using the correct filter method
	filters := make(map[string]string)
	response, err := vpcService.DescribeVpnGwSslServerByFilter(ctx, filters)
	if err != nil {
		return fmt.Errorf("get VPN SSL server list error: %s", err.Error())
	}

	if len(response) == 0 {
		return nil
	}

	// add scanning resources
	var resources, nonKeepResources []*tccommon.ResourceInstance
	for _, v := range response {
		if !tccommon.CheckResourcePersist(*v.SslVpnServerName, *v.CreateTime) {
			nonKeepResources = append(nonKeepResources, &tccommon.ResourceInstance{
				Id:   *v.SslVpnServerId,
				Name: *v.SslVpnServerName,
			})
		}
		resources = append(resources, &tccommon.ResourceInstance{
			Id:         *v.SslVpnServerId,
			Name:       *v.SslVpnServerName,
			CreateTime: *v.CreateTime,
		})
	}
	tccommon.ProcessScanCloudResources(client, resources, nonKeepResources, "CreateVpnGatewaySslServer")

	for _, v := range response {
		sslServerId := *v.SslVpnServerId
		sslServerName := *v.SslVpnServerName
		now := time.Now()
		createTime := tccommon.StringToTime(*v.CreateTime)
		interval := now.Sub(createTime).Minutes()

		if strings.HasPrefix(sslServerName, tcacctest.KeepResource) || strings.HasPrefix(sslServerName, tcacctest.DefaultResource) {
			continue
		}

		if tccommon.NeedProtect == 1 && int64(interval) < 30 {
			continue
		}

		if _, err := vpcService.DeleteVpnGatewaySslServer(ctx, sslServerId); err != nil {
			log.Printf("[ERROR] sweep VPN SSL server %s error: %s", sslServerId, err.Error())
		}
	}

	return nil
}

// TestAccTencentCloudVpnSslServerResource_basic tests basic VPN SSL Server functionality (backward compatibility)
func TestAccTencentCloudVpnSslServerResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckVpnSslServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnSslServerBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnSslServerExists("tencentcloud_vpn_ssl_server.test"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_ssl_server.test", "ssl_vpn_server_name", "tf-test-ssl-server"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_ssl_server.test", "local_address.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_ssl_server.test", "remote_address", "192.168.100.0/24"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpn_ssl_server.test", "vpn_gateway_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_vpn_ssl_server.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccTencentCloudVpnSslServerResource_withTags tests VPN SSL Server with tags
func TestAccTencentCloudVpnSslServerResource_withTags(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckVpnSslServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnSslServerWithTagsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnSslServerExists("tencentcloud_vpn_ssl_server.test"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_ssl_server.test", "ssl_vpn_server_name", "tf-test-ssl-server-tags"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_ssl_server.test", "tags.Environment", "test"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_ssl_server.test", "tags.Owner", "terraform"),
				),
			},
			{
				Config: testAccVpnSslServerWithTagsUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnSslServerExists("tencentcloud_vpn_ssl_server.test"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_ssl_server.test", "tags.Environment", "production"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_ssl_server.test", "tags.Team", "devops"),
				),
			},
		},
	})
}

// TestAccTencentCloudVpnSslServerResource_withDns tests VPN SSL Server with DNS configuration
func TestAccTencentCloudVpnSslServerResource_withDns(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckVpnSslServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnSslServerWithDnsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnSslServerExists("tencentcloud_vpn_ssl_server.test"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_ssl_server.test", "ssl_vpn_server_name", "tf-test-ssl-server-dns"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_ssl_server.test", "dns_servers.0.primary_dns", "8.8.8.8"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_ssl_server.test", "dns_servers.0.secondary_dns", "8.8.4.4"),
				),
			},
			{
				Config: testAccVpnSslServerWithDnsUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnSslServerExists("tencentcloud_vpn_ssl_server.test"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_ssl_server.test", "dns_servers.0.primary_dns", "1.1.1.1"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_ssl_server.test", "dns_servers.0.secondary_dns", "1.0.0.1"),
				),
			},
		},
	})
}

// TestAccTencentCloudVpnSslServerResource_withAccessPolicy tests VPN SSL Server with access policy control
func TestAccTencentCloudVpnSslServerResource_withAccessPolicy(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckVpnSslServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnSslServerWithAccessPolicyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnSslServerExists("tencentcloud_vpn_ssl_server.test"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_ssl_server.test", "ssl_vpn_server_name", "tf-test-ssl-server-policy"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_ssl_server.test", "access_policy_enabled", "true"),
				),
			},
		},
	})
}

// Note: SSO test is skipped because it requires whitelist approval from TencentCloud
// func TestAccTencentCloudVpnSslServerResource_withSso(t *testing.T) {
//     t.Skip("SSO feature requires whitelist approval")
// }

func testAccCheckVpnSslServerDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	conn := tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn()
	vpcService := svcvpc.NewVpcService(conn)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_vpn_ssl_server" {
			continue
		}

		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			has, _, e := vpcService.DescribeVpnSslServerById(ctx, rs.Primary.ID)
			if e != nil {
				return tccommon.RetryError(e)
			}
			if !has {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("VPN SSL server still exists: %s", rs.Primary.ID))
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func testAccCheckVpnSslServerExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("VPN SSL server %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("VPN SSL server id is not set")
		}

		conn := tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn()
		vpcService := svcvpc.NewVpcService(conn)

		var has bool
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			h, _, e := vpcService.DescribeVpnSslServerById(ctx, rs.Primary.ID)
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
			return fmt.Errorf("VPN SSL server not found: %s", rs.Primary.ID)
		}

		return nil
	}
}

// Test configurations

const testAccVpnSslServerBasicConfig = tcacctest.DefaultVpcVariable + `
resource "tencentcloud_vpn_gateway" "test" {
  name      = "tf-test-vpn-gateway"
  vpc_id    = var.vpc_id
  bandwidth = 5
  zone      = var.availability_zone
  type      = "SSL"

  tags = {
    test = "tf"
  }
}

resource "tencentcloud_vpn_ssl_server" "test" {
  vpn_gateway_id      = tencentcloud_vpn_gateway.test.id
  ssl_vpn_server_name = "tf-test-ssl-server"
  local_address       = ["10.0.200.0/24"]
  remote_address      = "192.168.100.0/24"
}
`

const testAccVpnSslServerWithTagsConfig = tcacctest.DefaultVpcVariable + `
resource "tencentcloud_vpn_gateway" "test" {
  name      = "tf-test-vpn-gateway"
  vpc_id    = var.vpc_id
  bandwidth = 5
  zone      = var.availability_zone
  type      = "SSL"

  tags = {
    test = "tf"
  }
}

resource "tencentcloud_vpn_ssl_server" "test" {
  vpn_gateway_id      = tencentcloud_vpn_gateway.test.id
  ssl_vpn_server_name = "tf-test-ssl-server-tags"
  local_address       = ["10.0.200.0/24"]
  remote_address      = "192.168.100.0/24"

  tags = {
    Environment = "test"
    Owner       = "terraform"
  }
}
`

const testAccVpnSslServerWithTagsUpdateConfig = tcacctest.DefaultVpcVariable + `
resource "tencentcloud_vpn_gateway" "test" {
  name      = "tf-test-vpn-gateway"
  vpc_id    = var.vpc_id
  bandwidth = 5
  zone      = var.availability_zone
  type      = "SSL"

  tags = {
    test = "tf"
  }
}

resource "tencentcloud_vpn_ssl_server" "test" {
  vpn_gateway_id      = tencentcloud_vpn_gateway.test.id
  ssl_vpn_server_name = "tf-test-ssl-server-tags"
  local_address       = ["10.0.200.0/24"]
  remote_address      = "192.168.100.0/24"

  tags = {
    Environment = "production"
    Team        = "devops"
  }
}
`

const testAccVpnSslServerWithDnsConfig = tcacctest.DefaultVpcVariable + `
resource "tencentcloud_vpn_gateway" "test" {
  name      = "tf-test-vpn-gateway"
  vpc_id    = var.vpc_id
  bandwidth = 5
  zone      = var.availability_zone
  type      = "SSL"

  tags = {
    test = "tf"
  }
}

resource "tencentcloud_vpn_ssl_server" "test" {
  vpn_gateway_id      = tencentcloud_vpn_gateway.test.id
  ssl_vpn_server_name = "tf-test-ssl-server-dns"
  local_address       = ["10.0.200.0/24"]
  remote_address      = "192.168.100.0/24"

  dns_servers {
    primary_dns   = "8.8.8.8"
    secondary_dns = "8.8.4.4"
  }
}
`

const testAccVpnSslServerWithDnsUpdateConfig = tcacctest.DefaultVpcVariable + `
resource "tencentcloud_vpn_gateway" "test" {
  name      = "tf-test-vpn-gateway"
  vpc_id    = var.vpc_id
  bandwidth = 5
  zone      = var.availability_zone
  type      = "SSL"

  tags = {
    test = "tf"
  }
}

resource "tencentcloud_vpn_ssl_server" "test" {
  vpn_gateway_id      = tencentcloud_vpn_gateway.test.id
  ssl_vpn_server_name = "tf-test-ssl-server-dns"
  local_address       = ["10.0.200.0/24"]
  remote_address      = "192.168.100.0/24"

  dns_servers {
    primary_dns   = "1.1.1.1"
    secondary_dns = "1.0.0.1"
  }
}
`

const testAccVpnSslServerWithAccessPolicyConfig = tcacctest.DefaultVpcVariable + `
resource "tencentcloud_vpn_gateway" "test" {
  name      = "tf-test-vpn-gateway"
  vpc_id    = var.vpc_id
  bandwidth = 5
  zone      = var.availability_zone
  type      = "SSL"

  tags = {
    test = "tf"
  }
}

resource "tencentcloud_vpn_ssl_server" "test" {
  vpn_gateway_id        = tencentcloud_vpn_gateway.test.id
  ssl_vpn_server_name   = "tf-test-ssl-server-policy"
  local_address         = ["10.0.200.0/24"]
  remote_address        = "192.168.100.0/24"
  access_policy_enabled = true
}
`
