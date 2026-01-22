package gwlb

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gwlbv20240906 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gwlb/v20240906"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudGwlbTargetGroupRegisterInstances() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGwlbTargetGroupRegisterInstancesCreate,
		Read:   resourceTencentCloudGwlbTargetGroupRegisterInstancesRead,
		Delete: resourceTencentCloudGwlbTargetGroupRegisterInstancesDelete,
		Update: resourceTencentCloudGwlbTargetGroupRegisterInstancesUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"target_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Target group ID.",
			},

			"target_group_instances": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Server instance array.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bind_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Private network IP of target group instance.",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Port of target group instance. Only 6081 is supported.",
						},
						"weight": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Weight of target group instance. Only 0 or 16 is supported, and non-0 is uniformly treated as 16.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudGwlbTargetGroupRegisterInstancesCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gwlb_target_group_register_instances.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := GwlbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var (
		targetGroupId        string
		targetGroupInstances []*gwlbv20240906.TargetGroupInstance
	)

	if v, ok := d.GetOk("target_group_id"); ok {
		targetGroupId = v.(string)
	}

	if v, ok := d.GetOk("target_group_instances"); ok {
		targetGroupInstances = make([]*gwlbv20240906.TargetGroupInstance, 0)
		for _, item := range v.(*schema.Set).List() {
			targetGroupInstancesMap := item.(map[string]interface{})
			targetGroupInstance := gwlbv20240906.TargetGroupInstance{}
			if v, ok := targetGroupInstancesMap["bind_ip"]; ok {
				targetGroupInstance.BindIP = helper.String(v.(string))
			}
			if v, ok := targetGroupInstancesMap["port"]; ok {
				targetGroupInstance.Port = helper.IntUint64(v.(int))
			}
			if v, ok := targetGroupInstancesMap["weight"]; ok {
				targetGroupInstance.Weight = helper.IntUint64(v.(int))
			}
			targetGroupInstances = append(targetGroupInstances, &targetGroupInstance)
		}
	}

	err := service.RegisterTargetGroupInstances(ctx, targetGroupId, targetGroupInstances)
	if err != nil {
		log.Printf("[CRITAL]%s create gwlb target group register instances failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(targetGroupId)

	return resourceTencentCloudGwlbTargetGroupRegisterInstancesRead(d, meta)
}

func resourceTencentCloudGwlbTargetGroupRegisterInstancesUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gwlb_target_group_register_instances.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := GwlbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	targetGroupId := d.Id()

	if d.HasChange("target_group_instances") {
		o, n := d.GetChange("target_group_instances")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		add := ns.Difference(os).List()
		remove := os.Difference(ns).List()

		if len(remove) > 0 {
			removeTargetGroupInstances := make([]*gwlbv20240906.TargetGroupInstance, 0)
			for _, item := range remove {
				targetGroupInstancesMap := item.(map[string]interface{})
				targetGroupInstance := gwlbv20240906.TargetGroupInstance{}
				if v, ok := targetGroupInstancesMap["bind_ip"]; ok {
					targetGroupInstance.BindIP = helper.String(v.(string))
				}
				if v, ok := targetGroupInstancesMap["port"]; ok {
					targetGroupInstance.Port = helper.IntUint64(v.(int))
				}
				if v, ok := targetGroupInstancesMap["weight"]; ok {
					targetGroupInstance.Weight = helper.IntUint64(v.(int))
				}
				removeTargetGroupInstances = append(removeTargetGroupInstances, &targetGroupInstance)
			}
			err := service.DeregisterTargetGroupInstances(ctx, targetGroupId, removeTargetGroupInstances)
			if err != nil {
				log.Printf("[CRITAL]%s delete gwlb target group register instances failed, reason:%+v", logId, err)
				return err
			}
		}
		if len(add) > 0 {
			addTargetGroupInstances := make([]*gwlbv20240906.TargetGroupInstance, 0)
			for _, item := range add {
				targetGroupInstancesMap := item.(map[string]interface{})
				targetGroupInstance := gwlbv20240906.TargetGroupInstance{}
				if v, ok := targetGroupInstancesMap["bind_ip"]; ok {
					targetGroupInstance.BindIP = helper.String(v.(string))
				}
				if v, ok := targetGroupInstancesMap["port"]; ok {
					targetGroupInstance.Port = helper.IntUint64(v.(int))
				}
				if v, ok := targetGroupInstancesMap["weight"]; ok {
					targetGroupInstance.Weight = helper.IntUint64(v.(int))
				}
				addTargetGroupInstances = append(addTargetGroupInstances, &targetGroupInstance)
			}
			err := service.RegisterTargetGroupInstances(ctx, targetGroupId, addTargetGroupInstances)
			if err != nil {
				log.Printf("[CRITAL]%s create gwlb target group register instances failed, reason:%+v", logId, err)
				return err
			}
		}
	}
	return resourceTencentCloudGwlbTargetGroupRegisterInstancesRead(d, meta)
}

func resourceTencentCloudGwlbTargetGroupRegisterInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gwlb_target_group_register_instances.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := GwlbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	targetGroupId := d.Id()

	respData, err := service.DescribeGwlbTargetGroupRegisterInstancesById(ctx, targetGroupId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `gwlb_target_group_register_instances` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("target_group_id", targetGroupId)
	targetGroupInstances := make([]interface{}, 0)
	for _, item := range respData.TargetGroupInstanceSet {
		targetGroupInstanceMap := make(map[string]interface{})
		if len(item.PrivateIpAddresses) > 0 {
			targetGroupInstanceMap["bind_ip"] = item.PrivateIpAddresses[0]
		}
		if item.Port != nil {
			targetGroupInstanceMap["port"] = item.Port
		}
		if item.Weight != nil {
			targetGroupInstanceMap["weight"] = item.Weight
		}
		targetGroupInstances = append(targetGroupInstances, targetGroupInstanceMap)
	}

	_ = d.Set("target_group_instances", targetGroupInstances)
	_ = targetGroupId
	return nil
}

func resourceTencentCloudGwlbTargetGroupRegisterInstancesDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gwlb_target_group_register_instances.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	service := GwlbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	targetGroupId := d.Id()

	targetGroupInstances := make([]*gwlbv20240906.TargetGroupInstance, 0)
	if v, ok := d.GetOk("target_group_instances"); ok {
		for _, item := range v.(*schema.Set).List() {
			targetGroupInstancesMap := item.(map[string]interface{})
			targetGroupInstance := gwlbv20240906.TargetGroupInstance{}
			if v, ok := targetGroupInstancesMap["bind_ip"]; ok {
				targetGroupInstance.BindIP = helper.String(v.(string))
			}
			if v, ok := targetGroupInstancesMap["port"]; ok {
				targetGroupInstance.Port = helper.IntUint64(v.(int))
			}
			if v, ok := targetGroupInstancesMap["weight"]; ok {
				targetGroupInstance.Weight = helper.IntUint64(v.(int))
			}
			targetGroupInstances = append(targetGroupInstances, &targetGroupInstance)
		}
	}
	err := service.DeregisterTargetGroupInstances(ctx, targetGroupId, targetGroupInstances)
	if err != nil {
		log.Printf("[CRITAL]%s delete gwlb target group register instances failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
