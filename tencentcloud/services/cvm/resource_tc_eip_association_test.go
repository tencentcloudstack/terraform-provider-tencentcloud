package cvm_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svccvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"

	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudEipAssociationWithInstance -v
func TestAccTencentCloudEipAssociationWithInstance(t *testing.T) {
	t.Parallel()
	id := "tencentcloud_eip_association.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckEipAssociationDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { tcacctest.AccStepPreConfigSetTempAKSK(t, tcacctest.ACCOUNT_TYPE_COMMON) },
				Config:    testAccTencentCloudEipAssociationWithInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipAssociationExists(id),
					resource.TestCheckResourceAttrSet(id, "eip_id"),
					resource.TestCheckResourceAttrSet(id, "instance_id"),
					resource.TestCheckNoResourceAttr(id, "network_interface_id"),
					resource.TestCheckNoResourceAttr(id, "private_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_eip.foo", "public_ip"),
					resource.TestCheckResourceAttr("tencentcloud_eip.foo", "name", tcacctest.DefaultInsName),
					resource.TestCheckResourceAttr("tencentcloud_eip.foo", "status", "UNBIND"),
				),
			},
			{
				ResourceName:      id,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// go test -i; go test -test.run TestAccTencentCloudEipAssociationWithNetworkInterface -v
func TestAccTencentCloudEipAssociationWithNetworkInterface(t *testing.T) {
	t.Parallel()
	id := "tencentcloud_eip_association.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckEipAssociationDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { tcacctest.AccStepPreConfigSetTempAKSK(t, tcacctest.ACCOUNT_TYPE_COMMON) },
				Config:    testAccTencentCloudEipAssociationWithNetworkInterface,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipAssociationExists(id),
					resource.TestCheckResourceAttrSet(id, "eip_id"),
					resource.TestCheckResourceAttrSet(id, "network_interface_id"),
					resource.TestCheckResourceAttrSet(id, "private_ip"),
					resource.TestCheckNoResourceAttr(id, "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_eip.foo", "public_ip"),
					resource.TestCheckResourceAttr("tencentcloud_eip.foo", "name", tcacctest.DefaultInsName),
					resource.TestCheckResourceAttr("tencentcloud_eip.foo", "status", "UNBIND"),
				),
			},
		},
	})
}

func testAccCheckEipAssociationDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	vpcService := svccvm.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_eip_association" {
			continue
		}

		associationId, err := svccvm.ParseEipAssociationId(rs.Primary.ID)

		if err != nil {
			return err
		}

		eip, err := vpcService.DescribeEipById(ctx, associationId.EipId)
		if err != nil {
			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				eip, err = vpcService.DescribeEipById(ctx, associationId.EipId)
				if err != nil {
					return tccommon.RetryError(err)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		if eip == nil {
			return nil
		}
		if eip.InstanceId != nil && *eip.InstanceId == associationId.InstanceId {
			return fmt.Errorf("eip association still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckEipAssociationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("eip association %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("eip association id is not set")
		}
		vpcService := svccvm.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		associationId, err := svccvm.ParseEipAssociationId(rs.Primary.ID)

		if err != nil {
			return err
		}
		eip, err := vpcService.DescribeEipById(ctx, associationId.EipId)
		if err != nil {
			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				eip, err = vpcService.DescribeEipById(ctx, associationId.EipId)
				if err != nil {
					return tccommon.RetryError(err)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		if eip == nil || (eip.InstanceId != nil && *eip.InstanceId != associationId.InstanceId) {
			return fmt.Errorf("eip %s is not found", associationId.EipId)
		}
		return nil
	}
}

const testAccTencentCloudEipAssociationWithInstance = tcacctest.DefaultInstanceVariable + `
resource "tencentcloud_eip" "foo" {
  name = var.instance_name
}

resource "tencentcloud_instance" "foo" {
  instance_name      = var.instance_name
  availability_zone  = var.availability_cvm_zone
  image_id           = data.tencentcloud_images.default.images.0.image_id
  instance_type      = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type   = "CLOUD_PREMIUM"
}

resource "tencentcloud_eip_association" "foo" {
  eip_id      = tencentcloud_eip.foo.id
  instance_id = tencentcloud_instance.foo.id
}
`

const testAccTencentCloudEipAssociationWithNetworkInterface = tcacctest.DefaultVpcVariable + `
resource "tencentcloud_eip" "foo" {
  name = var.instance_name
}

resource "tencentcloud_eni" "foo" {
  name        = var.instance_name
  vpc_id      = var.vpc_id
  subnet_id   = var.subnet_id
  description = var.instance_name
  ipv4_count  = 1
}

resource "tencentcloud_eip_association" "foo" {
  eip_id               = tencentcloud_eip.foo.id
  network_interface_id = tencentcloud_eni.foo.id
  private_ip           = tencentcloud_eni.foo.ipv4_info.0.ip
}
`
