package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudRouteTableAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRouteTableAssociationCreate,
		Read:   resourceTencentCloudRouteTableAssociationRead,
		Update: resourceTencentCloudRouteTableAssociationUpdate,
		Delete: resourceTencentCloudRouteTableAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"subnet_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Subnet instance ID, such as `subnet-3x5lf5q0`. This can be queried using the DescribeSubnets API.",
			},
			"route_table_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The route table instance ID, such as `rtb-azd4dt1c`.",
			},
		},
	}
}

func resourceTencentCloudRouteTableAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_route_table_association.create")()
	defer inconsistentCheck(d, meta)()

	var subnetId string
	if v, ok := d.GetOk("subnet_id"); ok {
		subnetId = v.(string)
	}

	d.SetId(subnetId)

	return resourceTencentCloudRouteTableAssociationUpdate(d, meta)
}

func resourceTencentCloudRouteTableAssociationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_route_table_association.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	subnetId := d.Id()

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		info VpcSubnetBasicInfo
		has  int
		e    error
	)
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		info, has, e = service.DescribeSubnet(ctx, subnetId, nil, "", "")
		if e != nil {
			return retryError(e)
		}

		// deleted
		if has == 0 {
			d.SetId("")
			return nil
		}

		if has != 1 {
			errRet := fmt.Errorf("one subnet_id read get %d subnet info", has)
			log.Printf("[CRITAL]%s %s", logId, errRet.Error())
			return resource.NonRetryableError(errRet)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if has == 0 {
		return nil
	}

	_ = d.Set("subnet_id", subnetId)
	if info.routeTableId != "" {
		_ = d.Set("route_table_id", info.routeTableId)
	}

	return nil
}

func resourceTencentCloudRouteTableAssociationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_route_table_association.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := vpc.NewReplaceRouteTableAssociationRequest()

	subnetId := d.Id()

	request.SubnetId = &subnetId
	request.RouteTableId = helper.String(d.Get("route_table_id").(string))

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().ReplaceRouteTableAssociation(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update vpc routeTable failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudRouteTableAssociationRead(d, meta)
}

func resourceTencentCloudRouteTableAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_route_table_association.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
