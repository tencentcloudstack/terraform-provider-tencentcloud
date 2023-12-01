/*
Provides a resource to create a vpc dhcp_associate_address

Example Usage

```hcl
data "tencentcloud_availability_zones" "zones" {}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones.zones.zones.0.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_vpc_dhcp_ip" "example" {
  vpc_id       = tencentcloud_vpc.vpc.id
  subnet_id    = tencentcloud_subnet.subnet.id
  dhcp_ip_name = "tf-example"
}

resource "tencentcloud_eip" "eip" {
  name = "example-eip"
}

resource "tencentcloud_vpc_dhcp_associate_address" "example" {
  dhcp_ip_id = tencentcloud_vpc_dhcp_ip.example.id
  address_ip = tencentcloud_eip.eip.public_ip
}
```

Import

vpc dhcp_associate_address can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_dhcp_associate_address.dhcp_associate_address dhcp_associate_address_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudVpcDhcpAssociateAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcDhcpAssociateAddressCreate,
		Read:   resourceTencentCloudVpcDhcpAssociateAddressRead,
		Delete: resourceTencentCloudVpcDhcpAssociateAddressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"dhcp_ip_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "`DhcpIp` unique `ID`, like: `dhcpip-9o233uri`. Must be a `DhcpIp` that is not bound to `EIP`.",
			},

			"address_ip": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Elastic public network `IP`. Must be `EIP` not bound to `DhcpIp`.",
			},
		},
	}
}

func resourceTencentCloudVpcDhcpAssociateAddressCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_dhcp_associate_address.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = vpc.NewAssociateDhcpIpWithAddressIpRequest()
		dhcpIpId  string
		addressIp string
	)
	if v, ok := d.GetOk("dhcp_ip_id"); ok {
		dhcpIpId = v.(string)
		request.DhcpIpId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("address_ip"); ok {
		addressIp = v.(string)
		request.AddressIp = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().AssociateDhcpIpWithAddressIp(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vpc DhcpAssociateAddress failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(dhcpIpId + FILED_SP + addressIp)

	return resourceTencentCloudVpcDhcpAssociateAddressRead(d, meta)
}

func resourceTencentCloudVpcDhcpAssociateAddressRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_dhcp_associate_address.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	dhcpIpId := idSplit[0]
	addressIp := idSplit[1]

	DhcpAssociateAddress, err := service.DescribeVpcDhcpAssociateAddressById(ctx, dhcpIpId, addressIp)
	if err != nil {
		return err
	}

	if DhcpAssociateAddress == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcDhcpAssociateAddress` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if DhcpAssociateAddress.DhcpIpId != nil {
		_ = d.Set("dhcp_ip_id", DhcpAssociateAddress.DhcpIpId)
	}

	if DhcpAssociateAddress.AddressIp != nil {
		_ = d.Set("address_ip", DhcpAssociateAddress.AddressIp)
	}

	return nil
}

func resourceTencentCloudVpcDhcpAssociateAddressDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_dhcp_associate_address.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	dhcpIpId := idSplit[0]

	if err := service.DeleteVpcDhcpAssociateAddressById(ctx, dhcpIpId); err != nil {
		return err
	}

	return nil
}
