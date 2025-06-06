package mqtt

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mqttv20240516 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mqtt/v20240516"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMqttInstance() *schema.Resource {
	return &schema.Resource{
		Create: ResourceTencentCloudMqttInstanceCreate,
		Read:   ResourceTencentCloudMqttInstanceRead,
		Update: ResourceTencentCloudMqttInstanceUpdate,
		Delete: ResourceTencentCloudMqttInstanceDelete,
		Schema: map[string]*schema.Schema{
			"instance_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance type. PRO for Professional Edition; PLATINUM for Platinum Edition.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance name.",
			},

			"sku_code": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Product SKU, available SKUs can be queried via the DescribeProductSKUList API.",
			},

			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Remarks.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of the MQTT instance.",
			},

			"vpc_list": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "VPC information bound to the instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "VPC ID.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Subnet ID.",
						},
					},
				},
			},

			"renew_flag": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable auto-renewal (0: Disabled; 1: Enabled).",
			},

			"time_span": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Purchase duration (unit: months).",
			},

			"pay_mode": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Payment mode (0: Postpaid; 1: Prepaid).",
			},

			"device_certificate_provision_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Client certificate registration method: JITP: Automatic registration; API: Manually register through the API.",
			},

			"automatic_activation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Is the automatic registration certificate automatically activated. Default is false.",
			},

			"authorization_policy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Authorization policy switch. Default is false.",
			},

			"force_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicate whether to force delete the instance. Default is `false`. If set true, the instance will be permanently deleted instead of being moved into the recycle bin. Note: only works for `PREPAID` instance.",
			},
		},
	}
}

func ResourceTencentCloudMqttInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = mqttv20240516.NewCreateInstanceRequest()
		response   = mqttv20240516.NewCreateInstanceResponse()
		instanceId string
	)

	if v, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sku_code"); ok {
		request.SkuCode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_list"); ok {
		for _, item := range v.([]interface{}) {
			vpcListMap := item.(map[string]interface{})
			vpcInfo := mqttv20240516.VpcInfo{}
			if v, ok := vpcListMap["vpc_id"].(string); ok && v != "" {
				vpcInfo.VpcId = helper.String(v)
			}

			if v, ok := vpcListMap["subnet_id"].(string); ok && v != "" {
				vpcInfo.SubnetId = helper.String(v)
			}

			request.VpcList = append(request.VpcList, &vpcInfo)
		}
	}

	if v, ok := d.GetOkExists("renew_flag"); ok {
		request.RenewFlag = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("time_span"); ok {
		request.TimeSpan = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("pay_mode"); ok {
		request.PayMode = helper.IntInt64(v.(int))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().CreateInstanceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create mqtt failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create mqtt failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.InstanceId == nil {
		return fmt.Errorf("InstanceId is nil.")
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	// wait
	waitReq := mqttv20240516.NewDescribeInstanceRequest()
	waitReq.InstanceId = &instanceId
	err := resource.Retry(6*tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().DescribeInstanceWithContext(ctx, waitReq)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.InstanceStatus == nil {
			return resource.NonRetryableError(fmt.Errorf("Wait mqtt failed, Response is nil."))
		}

		if *result.Response.InstanceStatus == "RUNNING" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Instance is still creating, status is %s", *result.Response.InstanceStatus))
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mqtt failed, reason:%+v", logId, err)
		return reqErr
	}

	var (
		isAutomaticActivation bool
		isAuthorizationPolicy bool
	)

	if v, ok := d.GetOkExists("automatic_activation"); ok {
		isAutomaticActivation = v.(bool)
	}

	if v, ok := d.GetOkExists("authorization_policy"); ok {
		isAuthorizationPolicy = v.(bool)
	}

	// open automatic_activation or authorization_policy
	if isAutomaticActivation || isAuthorizationPolicy {
		modifyRequest := mqttv20240516.NewModifyInstanceRequest()
		modifyRequest.InstanceId = &instanceId
		modifyRequest.AutomaticActivation = helper.Bool(isAutomaticActivation)
		modifyRequest.AuthorizationPolicy = helper.Bool(isAuthorizationPolicy)
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().ModifyInstanceWithContext(ctx, modifyRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyRequest.GetAction(), modifyRequest.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update mqtt failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := fmt.Sprintf("qcs::mqtt:%s:uin/:instance/%s", region, instanceId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return ResourceTencentCloudMqttInstanceRead(d, meta)
}

func ResourceTencentCloudMqttInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = MqttService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		tagService = svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region     = meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		instanceId = d.Id()
	)

	respData, err := service.DescribeMqttById(ctx, instanceId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `mqtt` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if respData.InstanceType != nil {
		_ = d.Set("instance_type", respData.InstanceType)
	}

	if respData.InstanceName != nil {
		_ = d.Set("name", respData.InstanceName)
	}

	if respData.Remark != nil {
		_ = d.Set("remark", respData.Remark)
	}

	if respData.SkuCode != nil {
		_ = d.Set("sku_code", respData.SkuCode)
	}

	if respData.RenewFlag != nil {
		_ = d.Set("renew_flag", respData.RenewFlag)
	}

	if respData.PayMode != nil {
		if *respData.PayMode == "POSTPAID" {
			_ = d.Set("pay_mode", 0)
		} else {
			_ = d.Set("pay_mode", 1)
		}
	}

	if respData.DeviceCertificateProvisionType != nil {
		_ = d.Set("device_certificate_provision_type", respData.DeviceCertificateProvisionType)
	}

	if respData.AutomaticActivation != nil {
		_ = d.Set("automatic_activation", respData.AutomaticActivation)
	}

	if respData.AuthorizationPolicy != nil {
		_ = d.Set("authorization_policy", respData.AuthorizationPolicy)
	}

	forceDelete := false
	if v, ok := d.GetOkExists("force_delete"); ok {
		forceDelete = v.(bool)
		_ = d.Set("force_delete", forceDelete)
	}

	tags, err := tagService.DescribeResourceTags(ctx, "mqtt", "instance", region, instanceId)
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)

	return nil
}

func ResourceTencentCloudMqttInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		instanceId = d.Id()
	)

	immutableArgs := []string{"instance_type", "vpc_list", "renew_flag", "time_span", "pay_mode"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	needChange := false
	mutableArgs := []string{"name", "remark", "sku_code", "device_certificate_provision_type", "automatic_activation", "authorization_policy"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := mqttv20240516.NewModifyInstanceRequest()
		request.InstanceId = helper.String(instanceId)
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}

		if v, ok := d.GetOk("sku_code"); ok {
			request.SkuCode = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("automatic_activation"); ok {
			request.AutomaticActivation = helper.Bool(v.(bool))
		}

		if v, ok := d.GetOkExists("authorization_policy"); ok {
			request.AuthorizationPolicy = helper.Bool(v.(bool))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().ModifyInstanceWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update mqtt failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := fmt.Sprintf("qcs::mqtt:%s:uin/:instance/%s", region, instanceId)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return ResourceTencentCloudMqttInstanceRead(d, meta)
}

func ResourceTencentCloudMqttInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request     = mqttv20240516.NewDeleteInstanceRequest()
		instanceId  = d.Id()
		payMode     int
		forceDelete bool
	)

	if v, ok := d.GetOkExists("pay_mode"); ok {
		payMode = v.(int)
	}

	if v, ok := d.GetOkExists("force_delete"); ok {
		forceDelete = v.(bool)
	}

	request.InstanceId = helper.String(instanceId)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().DeleteInstanceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete mqtt failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// wait
	waitReq := mqttv20240516.NewDescribeInstanceRequest()
	waitReq.InstanceId = &instanceId
	err := resource.Retry(4*tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().DescribeInstanceWithContext(ctx, waitReq)
		if e != nil {
			if sdkError, ok := e.(*errors.TencentCloudSDKError); ok {
				if sdkError.Code == "ResourceNotFound.Instance" {
					return nil
				}
			}

			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.InstanceStatus == nil {
			return resource.NonRetryableError(fmt.Errorf("Wait mqtt failed, Response is nil."))
		}

		if payMode == 1 && *result.Response.InstanceStatus == "OVERDUE" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Instance is still destroying, status is %s", *result.Response.InstanceStatus))
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete mqtt failed, reason:%+v", logId, err)
		return reqErr
	}

	// PREPAID need delete again
	if payMode == 1 && forceDelete == true {
		// delete again
		reqErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().DeleteInstanceWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s delete mqtt failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		// wait
		err = resource.Retry(4*tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().DescribeInstanceWithContext(ctx, waitReq)
			if e != nil {
				if sdkError, ok := e.(*errors.TencentCloudSDKError); ok {
					if sdkError.Code == "ResourceNotFound.Instance" {
						return nil
					}
				}

				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), result.ToJsonString())
			}

			return resource.RetryableError(fmt.Errorf("Instance is still destroying, status is %s", *result.Response.InstanceStatus))
		})

		if err != nil {
			log.Printf("[CRITAL]%s delete mqtt failed, reason:%+v", logId, err)
			return reqErr
		}
	}

	return nil
}
