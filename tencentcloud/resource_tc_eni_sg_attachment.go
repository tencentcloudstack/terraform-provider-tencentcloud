/*
Provides a resource to create a eni_sg_attachment

-> **Note:** If this resource is used to bind security groups to eni, it cannot be linked to `tentcloud_eni` binding security group for simultaneous use.

Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "vpc"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "example1" {
  name        = "tf-example-sg1"
  description = "sg desc."
  project_id  = 0

  tags = {
    "example" = "test"
  }
}

resource "tencentcloud_security_group" "example2" {
  name        = "tf-example-sg2"
  description = "sg desc."
  project_id  = 0

  tags = {
    "example" = "test"
  }
}

resource "tencentcloud_eni" "example" {
  name        = "tf-example-eni"
  vpc_id      = tencentcloud_vpc.vpc.id
  subnet_id   = tencentcloud_subnet.subnet.id
  description = "eni desc."
  ipv4_count  = 1
}

resource "tencentcloud_eni_sg_attachment" "eni_sg_attachment" {
  network_interface_ids = [tencentcloud_eni.example.id]
  security_group_ids    = [
    tencentcloud_security_group.example1.id,
    tencentcloud_security_group.example2.id
  ]
}
```

Import

vpc eni_sg_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_eni_sg_attachment.eni_sg_attachment eni_sg_attachment_id
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func resourceTencentCloudEniSgAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEniSgAttachmentCreate,
		Read:   resourceTencentCloudEniSgAttachmentRead,
		Delete: resourceTencentCloudEniSgAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"network_interface_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MaxItems:    1,
				Description: "ENI instance ID. Such as:eni-pxir56ns. It Only support set one eni instance now.",
			},

			"security_group_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Security group instance ID, for example:sg-33ocnj9n, can be obtained through DescribeSecurityGroups. There is a limit of 100 instances per request.",
			},
		},
	}
}

func resourceTencentCloudEniSgAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eni_sg_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request            = vpc.NewAssociateNetworkInterfaceSecurityGroupsRequest()
		networkInterfaceId string
	)
	if v, ok := d.GetOk("network_interface_ids"); ok {
		networkInterfaceIdsSet := v.(*schema.Set).List()
		for i := range networkInterfaceIdsSet {
			networkInterfaceId = networkInterfaceIdsSet[i].(string)
			request.NetworkInterfaceIds = append(request.NetworkInterfaceIds, &networkInterfaceId)
		}
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIdsSet := v.(*schema.Set).List()
		for i := range securityGroupIdsSet {
			securityGroupIds := securityGroupIdsSet[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroupIds)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().AssociateNetworkInterfaceSecurityGroups(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vpc eniSgAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(networkInterfaceId)

	return resourceTencentCloudEniSgAttachmentRead(d, meta)
}

func resourceTencentCloudEniSgAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eni_sg_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	networkInterfaceId := d.Id()

	enis, err := service.DescribeEniById(ctx, []string{networkInterfaceId})
	if err != nil {
		return err
	}

	if len(enis) < 1 {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcEniSgAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	eni := enis[0]

	_ = d.Set("network_interface_ids", []*string{&networkInterfaceId})

	if eni.GroupSet != nil {
		_ = d.Set("security_group_ids", eni.GroupSet)
	}

	return nil
}

func resourceTencentCloudEniSgAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eni_sg_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	networkInterfaceId := d.Id()

	var securityGroupIds []string
	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIdsSet := v.(*schema.Set).List()
		for i := range securityGroupIdsSet {
			securityGroupId := securityGroupIdsSet[i].(string)
			securityGroupIds = append(securityGroupIds, securityGroupId)
		}
	}

	if err := service.DeleteVpcEniSgAttachmentById(ctx, networkInterfaceId, securityGroupIds); err != nil {
		return err
	}

	return nil
}
