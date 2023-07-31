/*
Provides a resource to create a vpc ipv6_cidr_block

Example Usage

```hcl
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_vpc_ipv6_cidr_block" "example" {
  vpc_id = tencentcloud_vpc.vpc.id
}
```

Import

vpc ipv6_cidr_block can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_ipv6_cidr_block.ipv6_cidr_block vpc_id
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudVpcIpv6CidrBlock() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcIpv6CidrBlockCreate,
		Read:   resourceTencentCloudVpcIpv6CidrBlockRead,
		Delete: resourceTencentCloudVpcIpv6CidrBlockDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "`VPC` instance `ID`, in the form of `vpc-f49l6u0z`.",
			},
			"ipv6_cidr_block": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ipv6 cidr block.",
			},
		},
	}
}

func resourceTencentCloudVpcIpv6CidrBlockCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_ipv6_cidr_block.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request = vpc.NewAssignIpv6CidrBlockRequest()
		vpcId   string
	)
	if v, ok := d.GetOk("vpc_id"); ok {
		vpcId = v.(string)
		request.VpcId = helper.String(vpcId)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().AssignIpv6CidrBlock(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vpc ipv6CidrBlock failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(vpcId)

	return resourceTencentCloudVpcIpv6CidrBlockRead(d, meta)
}

func resourceTencentCloudVpcIpv6CidrBlockRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_ipv6_cidr_block.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	vpcId := d.Id()

	instance, err := service.DescribeVpcById(ctx, vpcId)
	if err != nil {
		return err
	}

	if instance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcIpv6CidrBlock` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if instance.VpcId != nil {
		_ = d.Set("vpc_id", instance.VpcId)
	}

	if instance.Ipv6CidrBlock != nil {
		_ = d.Set("ipv6_cidr_block", instance.Ipv6CidrBlock)
	}

	return nil
}

func resourceTencentCloudVpcIpv6CidrBlockDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_ipv6_cidr_block.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	vpcId := d.Id()

	if err := service.DeleteVpcIpv6CidrBlockById(ctx, vpcId); err != nil {
		return err
	}

	return nil
}
