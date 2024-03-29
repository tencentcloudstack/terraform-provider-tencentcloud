package vpc

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudRouteTableAssociation() *schema.Resource {
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
	defer tccommon.LogElapsed("resource.tencentcloud_route_table_association.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var subnetId string
	if v, ok := d.GetOk("subnet_id"); ok {
		subnetId = v.(string)
	}

	d.SetId(subnetId)

	return resourceTencentCloudRouteTableAssociationUpdate(d, meta)
}

func resourceTencentCloudRouteTableAssociationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_route_table_association.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	subnetId := d.Id()

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		info VpcSubnetBasicInfo
		has  int
		e    error
	)
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		info, has, e = service.DescribeSubnet(ctx, subnetId, nil, "", "")
		if e != nil {
			return tccommon.RetryError(e)
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
	defer tccommon.LogElapsed("resource.tencentcloud_route_table_association.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := vpc.NewReplaceRouteTableAssociationRequest()

	subnetId := d.Id()

	request.SubnetId = &subnetId
	request.RouteTableId = helper.String(d.Get("route_table_id").(string))

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ReplaceRouteTableAssociation(request)
		if e != nil {
			return tccommon.RetryError(e)
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
	defer tccommon.LogElapsed("resource.tencentcloud_route_table_association.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	subnetId := d.Id()

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		info VpcSubnetBasicInfo
		has  int
		e    error
	)
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		info, has, e = service.DescribeSubnet(ctx, subnetId, nil, "", "")
		if e != nil {
			return tccommon.RetryError(e)
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
		return fmt.Errorf("Unable to find the vpc corresponding to the current subnet:  %s", subnetId)
	}

	routeTables, err := service.DescribeRouteTables(ctx, "", "", info.vpcId, nil, helper.Bool(true), "")

	if err != nil {
		log.Printf("[WARN] Describe default Route Table error: %s", err.Error())
	}

	if len(routeTables) < 1 {
		return fmt.Errorf("Unable to find the default routetable corresponding to the current vpc:  %s", info.vpcId)
	}

	defaultRoutetableId := routeTables[0].routeTableId

	request := vpc.NewReplaceRouteTableAssociationRequest()

	request.SubnetId = &subnetId
	request.RouteTableId = helper.String(defaultRoutetableId)

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ReplaceRouteTableAssociation(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update vpc routeTable failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
