package gwlb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gwlbv20240906 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gwlb/v20240906"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudGwlbInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGwlbInstanceCreate,
		Read:   resourceTencentCloudGwlbInstanceRead,
		Update: resourceTencentCloudGwlbInstanceUpdate,
		Delete: resourceTencentCloudGwlbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the VPC to which the backend target device of the GWLB belongs, such as vpc-12345678. It can be obtained through the DescribeVpcEx interface. If left blank, it defaults to DefaultVPC. This parameter is required when a private network CLB instance is created.",
			},

			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subnet ID of the VPC to which the backend target device of the GWLB belongs.",
			},

			"load_balancer_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "GWLB instance name. It supports input of 1 to 60 characters. If not filled in, it will be generated automatically by default.",
			},

			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "While the GWLB is purchased, it is tagged, with a maximum of 20 tag key-value pairs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag value.",
						},
					},
				},
			},

			"lb_charge_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "GWLB instance billing type, which currently supports POSTPAID_BY_HOUR only. The default is POSTPAID_BY_HOUR.",
			},
			"vips": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Gateway Load Balancer provides virtual IP services.",
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Gateway Load Balancer instance status. 0: Creating, 1: Running normally, 3: Removing.",
			},
			"target_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the associated target group.",
			},
			"delete_protect": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to turn on the deletion protection function.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time.",
			},
			"isolation": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "0: means not quarantined, 1: means quarantined.",
			},
			"isolated_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time when the Gateway Load Balancer instance was isolated.",
			},
			"operate_protect": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to enable the configuration modification protection function.",
			},
		},
	}
}

func resourceTencentCloudGwlbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gwlb_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	service := GwlbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var (
		instanceId string
	)
	var (
		request  = gwlbv20240906.NewCreateGatewayLoadBalancerRequest()
		response = gwlbv20240906.NewCreateGatewayLoadBalancerResponse()
	)

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("load_balancer_name"); ok {
		request.LoadBalancerName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		for _, item := range v.([]interface{}) {
			tagsMap := item.(map[string]interface{})
			tagInfo := gwlbv20240906.TagInfo{}
			if v, ok := tagsMap["tag_key"]; ok {
				tagInfo.TagKey = helper.String(v.(string))
			}
			if v, ok := tagsMap["tag_value"]; ok {
				tagInfo.TagValue = helper.String(v.(string))
			}
			request.Tags = append(request.Tags, &tagInfo)
		}
	}

	if v, ok := d.GetOk("lb_charge_type"); ok {
		request.LBChargeType = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGwlbV20240906Client().CreateGatewayLoadBalancerWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create gwlb instance failed, reason:%+v", logId, err)
		return err
	}

	if len(response.Response.LoadBalancerIds) > 0 {
		instanceId = *response.Response.LoadBalancerIds[0]
	}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"0"}, 10*tccommon.ReadRetryTimeout, time.Second, service.TaskStatusRefreshFunc(ctx, *response.Response.RequestId, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}
	_ = response

	d.SetId(instanceId)

	return resourceTencentCloudGwlbInstanceRead(d, meta)
}

func resourceTencentCloudGwlbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gwlb_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := GwlbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()

	respData, err := service.DescribeGwlbInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `gwlb_instance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if respData.VpcId != nil {
		_ = d.Set("vpc_id", respData.VpcId)
	}
	if respData.SubnetId != nil {
		_ = d.Set("subnet_id", respData.SubnetId)
	}
	if respData.LoadBalancerName != nil {
		_ = d.Set("load_balancer_name", respData.LoadBalancerName)
	}
	if len(respData.Tags) > 0 {
		tags := make([]interface{}, 0)
		for _, tag := range respData.Tags {
			tags = append(tags, map[string]interface{}{
				"tag_key":   tag.TagKey,
				"tag_value": tag.TagValue,
			})

		}
		_ = d.Set("tags", tags)
	}
	if respData.ChargeType != nil {
		_ = d.Set("lb_charge_type", respData.ChargeType)
	}
	if len(respData.Vips) > 0 {
		vips := make([]string, 0)
		for _, vip := range respData.Vips {
			vips = append(vips, *vip)
		}
		_ = d.Set("vips", vips)
	}
	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}
	if respData.TargetGroupId != nil {
		_ = d.Set("target_group_id", respData.TargetGroupId)
	}
	if respData.DeleteProtect != nil {
		_ = d.Set("delete_protect", respData.DeleteProtect)
	}
	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}
	if respData.Isolation != nil {
		_ = d.Set("isolation", respData.Isolation)
	}
	if respData.IsolatedTime != nil {
		_ = d.Set("isolated_time", respData.IsolatedTime)
	}
	if respData.OperateProtect != nil {
		_ = d.Set("operate_protect", respData.OperateProtect)
	}
	_ = instanceId
	return nil
}

func resourceTencentCloudGwlbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gwlb_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	immutableArgs := []string{"vpc_id", "subnet_id", "number", "tags", "lb_charge_type"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	instanceId := d.Id()

	needChange := false
	mutableArgs := []string{"load_balancer_name"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := gwlbv20240906.NewModifyGatewayLoadBalancerAttributeRequest()
		request.LoadBalancerId = helper.String(instanceId)
		if v, ok := d.GetOk("load_balancer_name"); ok {
			request.LoadBalancerName = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGwlbV20240906Client().ModifyGatewayLoadBalancerAttributeWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update gwlb instance failed, reason:%+v", logId, err)
			return err
		}
	}

	_ = instanceId
	return resourceTencentCloudGwlbInstanceRead(d, meta)
}

func resourceTencentCloudGwlbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gwlb_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	instanceId := d.Id()

	var (
		request  = gwlbv20240906.NewDeleteGatewayLoadBalancerRequest()
		response = gwlbv20240906.NewDeleteGatewayLoadBalancerResponse()
	)

	request.LoadBalancerIds = helper.Strings([]string{instanceId})
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGwlbV20240906Client().DeleteGatewayLoadBalancerWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete gwlb instance failed, reason:%+v", logId, err)
		return err
	}

	service := GwlbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	conf := tccommon.BuildStateChangeConf([]string{}, []string{"0"}, 10*tccommon.ReadRetryTimeout, time.Second, service.TaskStatusRefreshFunc(ctx, *response.Response.RequestId, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}
	_ = response
	_ = instanceId
	return nil
}
