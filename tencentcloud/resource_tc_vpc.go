/*
Provide a resource to create a VPC.

Example Usage

```hcl
resource "tencentcloud_vpc" "foo" {
    name = "ci-temp-test-updated"
    cidr_block = "10.0.0.0/16"
	dns_servers=["119.29.29.29","8.8.8.8"]
	is_multicast=false
}
```

Import

Vpc instance can be imported, e.g.

```hcl
$ terraform import tencentcloud_vpc.test vpc-id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func resourceTencentCloudVpcInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcInstanceCreate,
		Read:   resourceTencentCloudVpcInstanceRead,
		Update: resourceTencentCloudVpcInstanceUpdate,
		Delete: resourceTencentCloudVpcInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 60),
				Description:  "The name of the VPC.",
			},
			"cidr_block": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCIDRNetworkAddress,
				Description:  "A network address block which should be a subnet of the three internal network segments (10.0.0.0/16, 172.16.0.0/12 and 192.168.0.0/16).",
			},
			"dns_servers": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set: func(v interface{}) int {
					return hashcode.String(v.(string))
				},
				Description: "The DNS server list of the VPC. And you can specify 0 to 5 servers to this list.",
			},
			"is_multicast": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicates whether VPC multicast is enabled. The default value is 'true'.",
			},

			// Computed values
			"is_default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether it is the default VPC for this region.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of VPC.",
			},
		},
	}
}

func resourceTencentCloudVpcInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc.create")()

	logId := getLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		name        string = ""
		cidrBlock   string = ""
		dnsServers         = make([]string, 0, 4)
		isMulticast bool   = true
	)
	if temp, ok := d.GetOk("name"); ok {
		name = temp.(string)
	}
	if temp, ok := d.GetOk("cidr_block"); ok {
		cidrBlock = temp.(string)
	}
	if temp, ok := d.GetOk("dns_servers"); ok {

		slice := temp.(*schema.Set).List()
		dnsServers = make([]string, 0, len(slice))
		for _, v := range slice {
			dnsServers = append(dnsServers, v.(string))
		}
		if len(dnsServers) < 1 {
			return fmt.Errorf("If dns_servers is set, then len(dns_servers) should be [1:4]")
		}
		if len(dnsServers) > 4 {
			return fmt.Errorf("If dns_servers is set, then len(dns_servers) should be [1:4]")
		}

	}
	isMulticast = d.Get("is_multicast").(bool)

	vpcId, _, err := service.CreateVpc(ctx, name, cidrBlock, isMulticast, dnsServers)
	if err != nil {
		return err
	}

	d.SetId(vpcId)

	return resourceTencentCloudVpcInstanceRead(d, meta)
}

func resourceTencentCloudVpcInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc.read")()

	logId := getLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	info, has, err := service.DescribeVpc(ctx, d.Id())
	if err != nil {
		return err
	}

	//deleted
	if has == 0 {
		log.Printf("[WARN]%s %s\n", logId, "vpc has been delete")
		d.SetId("")
		return nil
	}

	if has != 1 {
		errRet := fmt.Errorf("one vpc_id read get %d vpc info", has)
		log.Printf("[CRITAL]%s %s\n", logId, errRet.Error())
		return errRet
	}

	d.Set("name", info.name)
	d.Set("cidr_block", info.cidr)
	d.Set("dns_servers", info.dnsServers)
	d.Set("is_multicast", info.isMulticast)
	d.Set("create_time", info.createTime)
	d.Set("is_default", info.isDefault)

	return nil
}

func resourceTencentCloudVpcInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc.update")()

	logId := getLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		name        string = ""
		dnsServers         = make([]string, 0, 4)
		slice              = make([]interface{}, 0, 4)
		isMulticast bool   = true
	)

	old, now := d.GetChange("name")
	if d.HasChange("name") {
		name = now.(string)
	} else {
		name = old.(string)
	}

	old, now = d.GetChange("dns_servers")
	if d.HasChange("dns_servers") {
		slice = now.(*schema.Set).List()
		if len(slice) < 1 {
			return fmt.Errorf("If dns_servers is set, then len(dns_servers) should be [1:4]")
		}
		if len(slice) > 4 {
			return fmt.Errorf("If dns_servers is set, then len(dns_servers) should be [1:4]")
		}
	} else {
		slice = old.([]interface{})
	}

	if len(slice) > 0 {
		for _, v := range slice {
			dnsServers = append(dnsServers, v.(string))
		}
	}

	old, now = d.GetChange("is_multicast")
	if d.HasChange("is_multicast") {
		isMulticast = now.(bool)
	} else {
		isMulticast = old.(bool)
	}

	if err := service.ModifyVpcAttribute(ctx, d.Id(), name, isMulticast, dnsServers); err != nil {
		return err
	}

	return resourceTencentCloudVpcInstanceRead(d, meta)
}

func resourceTencentCloudVpcInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc.delete")()

	logId := getLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		if err := service.DeleteVpc(ctx, d.Id()); err != nil {
			if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok {
				if sdkErr.Code == VPCNotFound {
					return nil
				}
			}
			return resource.RetryableError(err)
		}
		return nil
	})

	return err
}
