package ccn

import (
	"context"
	"fmt"
	"log"
	"strings"

	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCcnRouteTableAssociateInstanceConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCcnRouteTableAssociateInstanceConfigCreate,
		Read:   resourceTencentCloudCcnRouteTableAssociateInstanceConfigRead,
		Update: resourceTencentCloudCcnRouteTableAssociateInstanceConfigUpdate,
		Delete: resourceTencentCloudCcnRouteTableAssociateInstanceConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"ccn_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the CCN.",
			},
			"route_table_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Ccn instance route table ID.",
			},
			"instances": {
				Required:    true,
				Type:        schema.TypeSet,
				Description: "Instances list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instances ID.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Cloud networking supports instance types: VPC, DIRECTCONNECT, BMVPC, EDGE, EDGE_TUNNEL, EDGE_VPNGW, VPNGW.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCcnRouteTableAssociateInstanceConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn_route_table_associate_instance_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	ccnId := d.Get("ccn_id").(string)
	routeTableId := d.Get("route_table_id").(string)

	d.SetId(strings.Join([]string{ccnId, routeTableId}, tccommon.FILED_SP))

	return resourceTencentCloudCcnRouteTableAssociateInstanceConfigUpdate(d, meta)
}

func resourceTencentCloudCcnRouteTableAssociateInstanceConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn_route_table_associate_instance_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	ccnId := items[0]
	routeTableId := items[1]

	instanceBindList, err := service.DescribeRouteTableAssociatedInstancesById(ctx, meta, ccnId, routeTableId)
	if err != nil {
		return err
	}

	if instanceBindList == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `tencentcloud_ccn_route_table_associate_instance_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("ccn_id", ccnId)
	_ = d.Set("route_table_id", routeTableId)

	tmpList := make([]map[string]interface{}, 0, len(instanceBindList))
	for _, instanceBind := range instanceBindList {
		instanceMap := map[string]interface{}{}
		if instanceBind.instanceId != "" {
			instanceMap["instance_id"] = instanceBind.instanceId
		}

		if instanceBind.instanceType != "" {
			instanceMap["instance_type"] = instanceBind.instanceType
		}

		tmpList = append(tmpList, instanceMap)
	}

	_ = d.Set("instances", tmpList)

	return nil
}

func resourceTencentCloudCcnRouteTableAssociateInstanceConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn_route_table_associate_instance_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = vpc.NewAssociateInstancesToCcnRouteTableRequest()
	)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	ccnId := items[0]
	routeTableId := items[1]

	request.CcnId = helper.String(ccnId)
	request.RouteTableId = helper.String(routeTableId)

	if v, ok := d.GetOk("instances"); ok {
		tmpV := v.(*schema.Set).List()
		for _, item := range tmpV {
			instanceMap := item.(map[string]interface{})
			instanceInfo := vpc.CcnInstanceWithoutRegion{}
			if v, ok := instanceMap["instance_id"]; ok {
				instanceInfo.InstanceId = helper.String(v.(string))
			}

			if v, ok := instanceMap["instance_type"]; ok {
				instanceInfo.InstanceType = helper.String(v.(string))
			}

			request.Instances = append(request.Instances, &instanceInfo)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().AssociateInstancesToCcnRouteTable(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update dc AssociateInstancesToCcnRouteTable failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCcnRouteTableAssociateInstanceConfigRead(d, meta)
}

func resourceTencentCloudCcnRouteTableAssociateInstanceConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn_route_table_associate_instance_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
