/*
Provides a resource to create a VPC routing table.

Example Usage

```hcl
resource "tencentcloud_vpc" "foo" {
    name = "ci-temp-test"
    cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_route_table" "foo" {
   vpc_id = "${tencentcloud_vpc.foo.id}"
   name = "ci-temp-test-rt"
}
```

Import

Vpc routetable instance can be imported, e.g.

```hcl
$ terraform import tencentcloud_route_table.test route_table_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func resourceTencentCloudVpcRouteTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcRouteTableCreate,
		Read:   resourceTencentCloudVpcRouteTableRead,
		Update: resourceTencentCloudVpcRouteTableUpdate,
		Delete: resourceTencentCloudVpcRouteTableDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of VPC to which the route table should be associated.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 60),
				Description:  "The name of routing table.",
			},
			// Computed values
			"subnet_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "ID list of the subnets associated with this route table.",
			},
			"route_entry_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "ID list of the routing entries.",
			},
			"is_default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether it is the default routing table.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the routing table.",
			},
		},
	}
}

func resourceTencentCloudVpcRouteTableCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_route_table.create")()

	logId := getLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		vpcId string = ""
		name  string = ""
	)
	if temp, ok := d.GetOk("vpc_id"); ok {
		vpcId = temp.(string)
		if len(vpcId) < 1 {
			return fmt.Errorf("vpc_id should be not empty string")
		}
	}
	if temp, ok := d.GetOk("name"); ok {
		name = temp.(string)
	}

	routeTableId, err := service.CreateRouteTable(ctx, name, vpcId)
	if err != nil {
		return err
	}
	d.SetId(routeTableId)

	return resourceTencentCloudVpcRouteTableRead(d, meta)
}

func resourceTencentCloudVpcRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_route_table.read")()

	logId := getLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	info, has, err := service.DescribeRouteTable(ctx, d.Id())
	if err != nil {
		return err
	}
	//deleted
	if has == 0 {
		d.SetId("")
		return nil
	}
	if has != 1 {
		errRet := fmt.Errorf("one route_table_id read get %d route_table info", has)
		log.Printf("[CRITAL]%s %s", logId, errRet.Error())
		return errRet
	}

	routeEntryIds := make([]string, 0, len(info.entryInfos))
	for _, v := range info.entryInfos {
		tfRouteEntryId := fmt.Sprintf("%d.%s", v.routeEntryId, d.Id())
		routeEntryIds = append(routeEntryIds, tfRouteEntryId)
	}

	d.Set("vpc_id", info.vpcId)
	d.Set("name", info.name)
	d.Set("subnet_ids", info.subnetIds)
	d.Set("route_entry_ids", routeEntryIds)
	d.Set("is_default", info.isDefault)
	d.Set("create_time", info.createTime)

	return nil
}

func resourceTencentCloudVpcRouteTableUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_route_table.update")()

	logId := getLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		name string = ""
	)

	if temp, ok := d.GetOk("name"); ok {
		name = temp.(string)
	}

	err := service.ModifyRouteTableAttribute(ctx, d.Id(), name)
	if err != nil {
		return err
	}

	return resourceTencentCloudVpcRouteTableRead(d, meta)
}

func resourceTencentCloudVpcRouteTableDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_route_table.delete")()

	logId := getLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		if err := service.DeleteRouteTable(ctx, d.Id()); err != nil {
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
