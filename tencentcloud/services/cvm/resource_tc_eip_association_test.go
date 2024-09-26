package cvm_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svccvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudEipAssociationWithInstance -v
func TestAccTencentCloudEipAssociationWithInstance(t *testing.T) {
	t.Parallel()
	id := "tencentcloud_eip_association.example"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckEipAssociationDestroy,
		Steps: []resource.TestStep{
			{
				//PreConfig: func() { tcacctest.AccStepPreConfigSetTempAKSK(t, tcacctest.ACCOUNT_TYPE_COMMON) },
				Config: testAccTencentCloudEipAssociationWithInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipAssociationExists(id),
					resource.TestCheckResourceAttrSet(id, "eip_id"),
					resource.TestCheckResourceAttrSet(id, "instance_id"),
					resource.TestCheckNoResourceAttr(id, "network_interface_id"),
					resource.TestCheckNoResourceAttr(id, "private_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_eip.example", "public_ip"),
					resource.TestCheckResourceAttr("tencentcloud_eip.example", "name", "tf-example"),
					resource.TestCheckResourceAttr("tencentcloud_eip.example", "status", "UNBIND"),
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
	id := "tencentcloud_eip_association.example"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckEipAssociationDestroy,
		Steps: []resource.TestStep{
			{
				//PreConfig: func() { tcacctest.AccStepPreConfigSetTempAKSK(t, tcacctest.ACCOUNT_TYPE_COMMON) },
				Config: testAccTencentCloudEipAssociationWithNetworkInterface,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipAssociationExists(id),
					resource.TestCheckResourceAttrSet(id, "eip_id"),
					resource.TestCheckResourceAttrSet(id, "network_interface_id"),
					resource.TestCheckResourceAttrSet(id, "private_ip"),
					resource.TestCheckNoResourceAttr(id, "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_eip.example", "public_ip"),
					resource.TestCheckResourceAttr("tencentcloud_eip.example", "name", "tf-example"),
					resource.TestCheckResourceAttr("tencentcloud_eip.example", "status", "UNBIND"),
				),
			},
		},
	})
}

func testAccCheckEipAssociationDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	vpcService := svcvpc.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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
		vpcService := svcvpc.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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

const testAccTencentCloudEipAssociationWithInstance = `
# create eip
resource "tencentcloud_eip" "example" {
  name = "tf-example"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = "ap-guangzhou-6"
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create cvm
resource "tencentcloud_instance" "example" {
  instance_name     = "tf_example"
  availability_zone = "ap-guangzhou-6"
  image_id          = "img-9qrfy1xt"
  instance_type     = "SA3.MEDIUM4"
  system_disk_type  = "CLOUD_HSSD"
  system_disk_size  = 100
  hostname          = "example"
  project_id        = 0
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id

  data_disks {
    data_disk_type = "CLOUD_HSSD"
    data_disk_size = 50
    encrypt        = false
  }

  tags = {
    tagKey = "tagValue"
  }
}

resource "tencentcloud_eip_association" "example" {
  eip_id      = tencentcloud_eip.example.id
  instance_id = tencentcloud_instance.example.id
}
`

const testAccTencentCloudEipAssociationWithNetworkInterface = `
# create eip
resource "tencentcloud_eip" "example" {
  name = "tf-example"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = "ap-guangzhou-6"
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create eni
resource "tencentcloud_eni" "example" {
  name        = "tf-example"
  vpc_id      = tencentcloud_vpc.vpc.id
  subnet_id   = tencentcloud_subnet.subnet.id
  description = "remark."
  ipv4_count  = 1
}

# create cvm
resource "tencentcloud_instance" "example" {
  instance_name     = "tf_example"
  availability_zone = "ap-guangzhou-6"
  image_id          = "img-9qrfy1xt"
  instance_type     = "SA3.MEDIUM4"
  system_disk_type  = "CLOUD_HSSD"
  system_disk_size  = 100
  hostname          = "example"
  project_id        = 0
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id

  data_disks {
    data_disk_type = "CLOUD_HSSD"
    data_disk_size = 50
    encrypt        = false
  }

  tags = {
    tagKey = "tagValue"
  }
}

resource "tencentcloud_eip_association" "example" {
  eip_id               = tencentcloud_eip.example.id
  network_interface_id = tencentcloud_eni.example.id
  private_ip           = tencentcloud_eni.example.ipv4_info[0].ip
}
`
