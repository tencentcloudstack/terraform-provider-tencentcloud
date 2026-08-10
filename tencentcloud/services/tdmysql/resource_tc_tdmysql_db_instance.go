package tdmysql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tdmysqlv20211122 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmysql/v20211122"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudTdmysqlDbInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmysqlDbInstanceCreate,
		Read:   resourceTencentCloudTdmysqlDbInstanceRead,
		Update: resourceTencentCloudTdmysqlDbInstanceUpdate,
		Delete: resourceTencentCloudTdmysqlDbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance zone.",
			},

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

			"spec_code": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specification code.",
			},

			"disk": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Storage node disk capacity, unit GB.",
			},

			"storage_node_num": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Storage node number.",
			},

			"replications": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Storage node replica number, max 5, must be odd.",
			},

			"full_replications": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Full replica number.",
			},

			"create_version": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Create version, defaults to latest.",
			},

			"instance_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance name, length 1-60.",
			},

			"resource_tags": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Resource tag list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Tag value.",
						},
					},
				},
			},

			"init_params": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "Init instance params.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"param": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Param key.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Param value.",
						},
					},
				},
			},

			"time_unit": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Time unit, m: month.",
			},

			"time_span": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Time span.",
			},

			"storage_node_cpu": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Storage node CPU cores.",
			},

			"storage_node_mem": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Storage node memory.",
			},

			"pay_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Pay mode, 0: postpaid, 1: prepaid.",
			},

			"mc_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Control node number.",
			},

			"vport": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Custom port.",
			},

			"zones": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Multi AZ zone list.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"auto_voucher": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to use voucher.",
			},

			"voucher_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Voucher ID list.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"instance_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Instance architecture type, separate or hybrid.",
			},

			"storage_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Disk type, CLOUD_HSSD or CLOUD_TCS.",
			},

			"az_mode": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "AZ mode, 1: single AZ, 2: multi AZ non-master, 3: multi AZ master.",
			},

			"instance_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Instance mode.",
			},

			"template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Parameter template ID.",
			},

			"sql_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Compatible mode, MySQL or HBase.",
			},

			"auto_scale_config": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Auto scaling config for svls instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"range_min": {
							Type:        schema.TypeFloat,
							Required:    true,
							ForceNew:    true,
							Description: "CCU min value.",
						},
						"range_max": {
							Type:        schema.TypeFloat,
							Required:    true,
							ForceNew:    true,
							Description: "CCU max value.",
						},
					},
				},
			},

			"security_group_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Security group ID list.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Root username, defaults to dbaadmin.",
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "dbaadmin password.",
			},

			"encryption_enable": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Transparent encryption, 0: disable, 1: enable.",
			},

			"auto_renew_flag": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Auto renew flag, 1 indicates enabling auto-renewal; 0 indicates disabling auto-renewal.",
			},

			"enable_ssl": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable SSL.",
			},

			// computed
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance ID.",
			},
		},
	}
}

func resourceTencentCloudTdmysqlDbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmysql_db_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = tdmysqlv20211122.NewCreateDBInstancesRequest()
		response = tdmysqlv20211122.NewCreateDBInstancesResponse()
	)

	if v, ok := d.GetOk("zone"); ok {
		request.Zone = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("spec_code"); ok {
		request.SpecCode = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("disk"); ok {
		request.Disk = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOkExists("storage_node_num"); ok {
		request.StorageNodeNum = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOkExists("replications"); ok {
		request.Replications = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOkExists("full_replications"); ok {
		request.FullReplications = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("create_version"); ok {
		request.CreateVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource_tags"); ok {
		for _, item := range v.([]interface{}) {
			resourceTagMap := item.(map[string]interface{})
			resourceTag := tdmysqlv20211122.ResourceTag{}
			if v, ok := resourceTagMap["tag_key"].(string); ok && v != "" {
				resourceTag.TagKey = helper.String(v)
			}
			if v, ok := resourceTagMap["tag_value"].(string); ok && v != "" {
				resourceTag.TagValue = helper.String(v)
			}
			request.ResourceTags = append(request.ResourceTags, &resourceTag)
		}
	}

	if v, ok := d.GetOk("init_params"); ok {
		initParamsSet := v.(*schema.Set).List()
		for _, item := range initParamsSet {
			initParamsMap := item.(map[string]interface{})
			instanceParam := tdmysqlv20211122.InstanceParam{}
			if v, ok := initParamsMap["param"].(string); ok && v != "" {
				instanceParam.Param = helper.String(v)
			}
			if v, ok := initParamsMap["value"].(string); ok && v != "" {
				instanceParam.Value = helper.String(v)
			}
			request.InitParams = append(request.InitParams, &instanceParam)
		}
	}

	if v, ok := d.GetOk("time_unit"); ok {
		request.TimeUnit = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("time_span"); ok {
		request.TimeSpan = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOkExists("storage_node_cpu"); ok {
		request.StorageNodeCpu = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOkExists("storage_node_mem"); ok {
		request.StorageNodeMem = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("pay_mode"); ok {
		request.PayMode = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("mc_num"); ok {
		request.MCNum = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOkExists("vport"); ok {
		request.Vport = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("zones"); ok {
		zonesList := v.([]interface{})
		for _, item := range zonesList {
			request.Zones = append(request.Zones, helper.String(item.(string)))
		}
	}

	if v, ok := d.GetOkExists("auto_voucher"); ok {
		request.AutoVoucher = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("voucher_ids"); ok {
		voucherIdsList := v.([]interface{})
		for _, item := range voucherIdsList {
			request.VoucherIds = append(request.VoucherIds, helper.String(item.(string)))
		}
	}

	if v, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("storage_type"); ok {
		request.StorageType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("az_mode"); ok {
		request.AZMode = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("instance_mode"); ok {
		request.InstanceMode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("template_id"); ok {
		request.TemplateId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sql_mode"); ok {
		request.SQLMode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("auto_scale_config"); ok {
		autoScaleConfigList := v.([]interface{})
		if len(autoScaleConfigList) > 0 {
			autoScaleConfigMap := autoScaleConfigList[0].(map[string]interface{})
			autoScaleConfig := tdmysqlv20211122.AutoScalingConfig{}
			if v, ok := autoScaleConfigMap["range_min"].(float64); ok {
				autoScaleConfig.RangeMin = helper.Float64(v)
			}
			if v, ok := autoScaleConfigMap["range_max"].(float64); ok {
				autoScaleConfig.RangeMax = helper.Float64(v)
			}
			request.AutoScaleConfig = &autoScaleConfig
		}
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIdsList := v.([]interface{})
		for _, item := range securityGroupIdsList {
			request.SecurityGroupIds = append(request.SecurityGroupIds, helper.String(item.(string)))
		}
	}

	if v, ok := d.GetOk("user_name"); ok {
		request.UserName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("password"); ok {
		request.Password = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("encryption_enable"); ok {
		request.EncryptionEnable = helper.Int64(int64(v.(int)))
	}

	request.InstanceCount = helper.Int64(int64(1))
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmysqlV20211122Client().CreateDBInstancesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.InstanceIds == nil || len(result.Response.InstanceIds) == 0 {
			return resource.NonRetryableError(fmt.Errorf("Create tdmysql db instance failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create tdmysql db instance failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	instanceId := *response.Response.InstanceIds[0]
	d.SetId(instanceId)

	if response.Response.FlowId == nil {
		return fmt.Errorf("Create tdmysql db instance failed, FlowId is nil.")
	}

	// wait DescribeFlow until Status=success
	flowId := *response.Response.FlowId
	flowRequest := tdmysqlv20211122.NewDescribeFlowRequest()
	flowRequest.FlowId = helper.Int64(flowId)
	flowErr := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		flowResult, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmysqlV20211122Client().DescribeFlowWithContext(ctx, flowRequest)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, flowRequest.GetAction(), flowRequest.ToJsonString(), flowResult.ToJsonString())
		if flowResult == nil || flowResult.Response == nil || flowResult.Response.Status == nil {
			return resource.NonRetryableError(fmt.Errorf("DescribeFlow failed, Response is nil."))
		}

		status := *flowResult.Response.Status
		if status == FLOW_STATUS_SUCCESS {
			return nil
		}

		if status == FLOW_STATUS_FAILED || status == FLOW_STATUS_PAUSED {
			return resource.NonRetryableError(fmt.Errorf("Create tdmysql db instance async flow failed, status is %s.", status))
		}

		return resource.RetryableError(fmt.Errorf("Create tdmysql db instance async flow is running, status is %s.", status))
	})

	if flowErr != nil {
		log.Printf("[CRITAL]%s create tdmysql db instance async flow polling failed, reason:%+v", logId, flowErr)
		return flowErr
	}

	// enable auto renew flag
	if v, ok := d.GetOkExists("auto_renew_flag"); ok {
		if v.(int) == 1 {
			request := tdmysqlv20211122.NewModifyAutoRenewFlagRequest()
			request.InstanceIds = helper.Strings([]string{instanceId})
			request.AutoRenewFlag = helper.Int64(int64(v.(int)))
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmysqlV20211122Client().ModifyAutoRenewFlagWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				if result == nil || result.Response == nil {
					return resource.NonRetryableError(fmt.Errorf("Modify tdmysql db instance auto renew failed, Response is nil."))
				}

				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s update tdmysql db instance auto renew failed failed, reason:%+v", logId, reqErr)
				return reqErr
			}
		}
	}

	// enable ssl
	if v, ok := d.GetOkExists("enable_ssl"); ok {
		if v.(bool) {
			request := tdmysqlv20211122.NewModifyInstanceSSLStatusRequest()
			response := tdmysqlv20211122.NewModifyInstanceSSLStatusResponse()
			request.InstanceId = helper.String(instanceId)
			request.Enabled = helper.Bool(v.(bool))
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmysqlV20211122Client().ModifyInstanceSSLStatusWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				if result == nil || result.Response == nil {
					return resource.NonRetryableError(fmt.Errorf("Modify tdmysql db instance ssl failed, Response is nil."))
				}

				response = result
				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s update tdmysql db instance ssl failed, reason:%+v", logId, reqErr)
				return reqErr
			}

			if response.Response.FlowId == nil {
				return fmt.Errorf("Update tdmysql db instance ssl failed, FlowId is nil.")
			}

			// wait
			flowId := *response.Response.FlowId
			flowRequest := tdmysqlv20211122.NewDescribeFlowRequest()
			flowRequest.FlowId = helper.Int64(flowId)
			flowErr := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
				flowResult, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmysqlV20211122Client().DescribeFlowWithContext(ctx, flowRequest)
				if e != nil {
					return tccommon.RetryError(e)
				}

				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, flowRequest.GetAction(), flowRequest.ToJsonString(), flowResult.ToJsonString())
				if flowResult == nil || flowResult.Response == nil || flowResult.Response.Status == nil {
					return resource.NonRetryableError(fmt.Errorf("DescribeFlow failed, Response is nil."))
				}

				status := *flowResult.Response.Status
				if status == FLOW_STATUS_SUCCESS {
					return nil
				}

				if status == FLOW_STATUS_FAILED || status == FLOW_STATUS_PAUSED {
					return resource.NonRetryableError(fmt.Errorf("Update tdmysql db instance ssl failed, status is %s.", status))
				}

				return resource.RetryableError(fmt.Errorf("Update tdmysql db instance ssl is running, status is %s.", status))
			})

			if flowErr != nil {
				log.Printf("[CRITAL]%s update tdmysql db instance ssl polling failed, reason:%+v", logId, flowErr)
				return flowErr
			}
		}
	}

	return resourceTencentCloudTdmysqlDbInstanceRead(d, meta)
}

func resourceTencentCloudTdmysqlDbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmysql_db_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = TdmysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	respData, err := service.DescribeTdmysqlDbInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_tdmysql_db_instance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.InstanceId != nil {
		_ = d.Set("instance_id", respData.InstanceId)
	}

	if respData.Zone != nil {
		_ = d.Set("zone", respData.Zone)
	}

	if respData.VpcId != nil {
		_ = d.Set("vpc_id", respData.VpcId)
	}

	if respData.SubnetId != nil {
		_ = d.Set("subnet_id", respData.SubnetId)
	}

	if respData.SpecCode != nil {
		_ = d.Set("spec_code", respData.SpecCode)
	}

	if respData.Disk != nil {
		_ = d.Set("disk", respData.Disk)
	}

	if respData.StorageNodeNum != nil {
		_ = d.Set("storage_node_num", respData.StorageNodeNum)
	}

	if respData.Replications != nil {
		_ = d.Set("replications", respData.Replications)
	}

	if respData.FullReplications != nil {
		_ = d.Set("full_replications", respData.FullReplications)
	}

	if respData.CreateVersion != nil {
		_ = d.Set("create_version", respData.CreateVersion)
	}

	if respData.InstanceName != nil {
		_ = d.Set("instance_name", respData.InstanceName)
	}

	if respData.ResourceTags != nil {
		resourceTagsList := make([]map[string]interface{}, 0, len(respData.ResourceTags))
		for _, resourceTag := range respData.ResourceTags {
			resourceTagMap := map[string]interface{}{}
			if resourceTag.TagKey != nil {
				resourceTagMap["tag_key"] = resourceTag.TagKey
			}
			if resourceTag.TagValue != nil {
				resourceTagMap["tag_value"] = resourceTag.TagValue
			}
			resourceTagsList = append(resourceTagsList, resourceTagMap)
		}
		_ = d.Set("resource_tags", resourceTagsList)
	}

	if respData.InitParams != nil {
		initParamsList := make([]map[string]interface{}, 0, len(respData.InitParams))
		for _, initParam := range respData.InitParams {
			initParamsMap := map[string]interface{}{}
			if initParam.Param != nil {
				initParamsMap["param"] = initParam.Param
			}
			if initParam.Value != nil {
				initParamsMap["value"] = initParam.Value
			}
			initParamsList = append(initParamsList, initParamsMap)
		}
		_ = d.Set("init_params", initParamsList)
	}

	if respData.StorageNodeCpu != nil {
		_ = d.Set("storage_node_cpu", respData.StorageNodeCpu)
	}

	if respData.StorageNodeMem != nil {
		_ = d.Set("storage_node_mem", respData.StorageNodeMem)
	}

	if respData.PayMode != nil {
		if *respData.PayMode == PAY_MODE_POSTPAY {
			_ = d.Set("pay_mode", "0")
		} else if *respData.PayMode == PAY_MODE_PREPAY {
			_ = d.Set("pay_mode", "1")
		}
	}

	if respData.Vport != nil {
		_ = d.Set("vport", respData.Vport)
	}

	if respData.Zones != nil {
		_ = d.Set("zones", respData.Zones)
	}

	if respData.InstanceType != nil {
		_ = d.Set("instance_type", respData.InstanceType)
	}

	if respData.StorageType != nil {
		_ = d.Set("storage_type", respData.StorageType)
	}

	if respData.AZMode != nil {
		_ = d.Set("az_mode", respData.AZMode)
	}

	if respData.InstanceMode != nil {
		_ = d.Set("instance_mode", respData.InstanceMode)
	}

	if respData.TemplateId != nil {
		_ = d.Set("template_id", respData.TemplateId)
	}

	if respData.SQLMode != nil {
		_ = d.Set("sql_mode", respData.SQLMode)
	}

	if respData.AutoScaleConfig != nil {
		autoScaleConfigList := make([]map[string]interface{}, 0, 1)
		autoScaleConfigMap := map[string]interface{}{}
		if respData.AutoScaleConfig.RangeMin != nil {
			autoScaleConfigMap["range_min"] = respData.AutoScaleConfig.RangeMin
		}
		if respData.AutoScaleConfig.RangeMax != nil {
			autoScaleConfigMap["range_max"] = respData.AutoScaleConfig.RangeMax
		}
		autoScaleConfigList = append(autoScaleConfigList, autoScaleConfigMap)
		_ = d.Set("auto_scale_config", autoScaleConfigList)
	}

	if respData.EncryptionEnable != nil {
		_ = d.Set("encryption_enable", respData.EncryptionEnable)
	}

	sgResp, err := service.DescribeTdmysqlDBSecurityGroupsById(ctx, instanceId)
	if err == nil && sgResp != nil {
		if sgResp.Groups != nil && len(sgResp.Groups) > 0 {
			securityGroupIds := make([]string, 0, len(sgResp.Groups))
			for _, group := range sgResp.Groups {
				if group.SecurityGroupId != nil {
					securityGroupIds = append(securityGroupIds, *group.SecurityGroupId)
				}
			}

			_ = d.Set("security_group_ids", securityGroupIds)
		}
	}

	return nil
}

func resourceTencentCloudTdmysqlDbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmysql_db_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		instanceId = d.Id()
	)

	immutableArgs := []string{
		"zone",
		"replications",
		"create_version",
		"resource_tags",
		"time_unit",
		"time_span",
		"pay_mode",
		"mc_num",
		"instance_type",
		"instance_mode",
		"template_id",
		"sql_mode",
		"auto_scale_config",
		"user_name",
		"encryption_enable",
	}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("tencentcloud_tdmysql_db_instance argument `%s` cannot be changed, please recreate the resource if you need to change it.", v)
		}
	}

	if d.HasChange("vpc_id") || d.HasChange("subnet_id") {
		request := tdmysqlv20211122.NewModifyInstanceNetworkRequest()
		request.InstanceId = helper.String(instanceId)
		if v, ok := d.GetOk("vpc_id"); ok {
			request.VpcId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("subnet_id"); ok {
			request.SubnetId = helper.String(v.(string))
		}

		request.VipReleaseDelay = helper.IntUint64(0)
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmysqlV20211122Client().ModifyInstanceNetworkWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify tdmysql db instance network failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update tdmysql db instance network failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		// wait
		time.Sleep(10 * time.Second)
		waitReq := tdmysqlv20211122.NewDescribeDBInstanceDetailRequest()
		waitReq.InstanceId = helper.String(instanceId)
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			ratelimit.Check(waitReq.GetAction())
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmysqlV20211122Client().DescribeDBInstanceDetail(waitReq)
			if e != nil {
				return tccommon.RetryError(e)
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe tdmysql db instance failed, Response is nil."))
			}

			if result.Response.Status == nil {
				return resource.NonRetryableError(fmt.Errorf("Status is nil."))
			}

			if *result.Response.Status == DB_INSTANCE_STATUS_RUNNING {
				return nil
			}

			return resource.RetryableError(fmt.Errorf("waiting for tdmysql db instance update network, status is %s.", *result.Response.Status))
		})

		if err != nil {
			return err
		}
	}

	if d.HasChange("spec_code") || d.HasChange("disk") ||
		d.HasChange("storage_node_cpu") || d.HasChange("storage_node_mem") ||
		d.HasChange("storage_type") {
		request := tdmysqlv20211122.NewUpgradeInstanceRequest()
		response := tdmysqlv20211122.NewUpgradeInstanceResponse()
		request.InstanceId = helper.String(instanceId)
		if v, ok := d.GetOk("spec_code"); ok {
			request.SpecCode = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("disk"); ok {
			request.Disk = helper.Int64(int64(v.(int)))
		}

		if v, ok := d.GetOkExists("storage_node_cpu"); ok {
			request.StorageNodeCpu = helper.Int64(int64(v.(int)))
		}

		if v, ok := d.GetOkExists("storage_node_mem"); ok {
			request.StorageNodeMem = helper.Int64(int64(v.(int)))
		}

		if v, ok := d.GetOk("storage_type"); ok {
			request.StorageType = helper.String(v.(string))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmysqlV20211122Client().UpgradeInstanceWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Upgrade tdmysql db instance failed, Response is nil."))
			}

			response = result
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s upgrade tdmysql db instance failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		// wait
		flowId := *response.Response.FlowId
		flowRequest := tdmysqlv20211122.NewDescribeFlowRequest()
		flowRequest.FlowId = helper.Int64(flowId)
		flowErr := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			flowResult, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmysqlV20211122Client().DescribeFlowWithContext(ctx, flowRequest)
			if e != nil {
				return tccommon.RetryError(e)
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, flowRequest.GetAction(), flowRequest.ToJsonString(), flowResult.ToJsonString())
			if flowResult == nil || flowResult.Response == nil || flowResult.Response.Status == nil {
				return resource.NonRetryableError(fmt.Errorf("DescribeFlow failed, Response is nil."))
			}

			status := *flowResult.Response.Status
			if status == FLOW_STATUS_SUCCESS {
				return nil
			}

			if status == FLOW_STATUS_FAILED || status == FLOW_STATUS_PAUSED {
				return resource.NonRetryableError(fmt.Errorf("Upgrade tdmysql db instance failed, status is %s.", status))
			}

			return resource.RetryableError(fmt.Errorf("Upgrade tdmysql db instance is running, status is %s.", status))
		})

		if flowErr != nil {
			log.Printf("[CRITAL]%s upgrade tdmysql db instance polling failed, reason:%+v", logId, flowErr)
			return flowErr
		}
	}

	if d.HasChange("zones") || d.HasChange("storage_node_num") ||
		d.HasChange("az_mode") || d.HasChange("full_replications") {
		request := tdmysqlv20211122.NewExpandInstanceRequest()
		response := tdmysqlv20211122.NewExpandInstanceResponse()
		request.InstanceId = helper.String(instanceId)
		if v, ok := d.GetOk("zones"); ok {
			zonesList := v.([]interface{})
			for _, item := range zonesList {
				request.Zones = append(request.Zones, helper.String(item.(string)))
			}
		}

		if v, ok := d.GetOkExists("storage_node_num"); ok {
			request.StorageNodeNum = helper.Int64(int64(v.(int)))
		}

		if v, ok := d.GetOkExists("az_mode"); ok {
			request.AZMode = helper.Int64(int64(v.(int)))
		}

		if v, ok := d.GetOkExists("full_replications"); ok {
			request.FullReplications = helper.Int64(int64(v.(int)))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmysqlV20211122Client().ExpandInstanceWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Expand tdmysql db instance failed, Response is nil."))
			}

			response = result
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s expand tdmysql db instance failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		// wait
		flowId := *response.Response.FlowId
		flowRequest := tdmysqlv20211122.NewDescribeFlowRequest()
		flowRequest.FlowId = helper.Int64(flowId)
		flowErr := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			flowResult, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmysqlV20211122Client().DescribeFlowWithContext(ctx, flowRequest)
			if e != nil {
				return tccommon.RetryError(e)
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, flowRequest.GetAction(), flowRequest.ToJsonString(), flowResult.ToJsonString())
			if flowResult == nil || flowResult.Response == nil || flowResult.Response.Status == nil {
				return resource.NonRetryableError(fmt.Errorf("DescribeFlow failed, Response is nil."))
			}

			status := *flowResult.Response.Status
			if status == FLOW_STATUS_SUCCESS {
				return nil
			}

			if status == FLOW_STATUS_FAILED || status == FLOW_STATUS_PAUSED {
				return resource.NonRetryableError(fmt.Errorf("Expand tdmysql db instance failed, status is %s.", status))
			}

			return resource.RetryableError(fmt.Errorf("Expand tdmysql db instance is running, status is %s.", status))
		})

		if flowErr != nil {
			log.Printf("[CRITAL]%s expand tdmysql db instance polling failed, reason:%+v", logId, flowErr)
			return flowErr
		}
	}

	if d.HasChange("vport") {
		request := tdmysqlv20211122.NewModifyDBInstanceVPortRequest()
		request.InstanceId = helper.String(instanceId)
		if v, ok := d.GetOkExists("vport"); ok {
			request.Vport = helper.Int64(int64(v.(int)))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmysqlV20211122Client().ModifyDBInstanceVPortWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify tdmysql db instance vport failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update tdmysql db instance vport failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	if d.HasChange("instance_name") {
		request := tdmysqlv20211122.NewModifyInstanceNameRequest()
		request.InstanceId = helper.String(instanceId)
		if v, ok := d.GetOk("instance_name"); ok {
			request.InstanceName = helper.String(v.(string))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmysqlV20211122Client().ModifyInstanceNameWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify tdmysql db instance name failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update tdmysql db instance name failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	if d.HasChange("init_params") {
		request := tdmysqlv20211122.NewModifyDBParametersRequest()
		response := tdmysqlv20211122.NewModifyDBParametersResponse()
		request.InstanceId = helper.String(instanceId)
		if v, ok := d.GetOk("init_params"); ok {
			initParamsSet := v.(*schema.Set).List()
			for _, item := range initParamsSet {
				paramsMap := item.(map[string]interface{})
				instanceParam := tdmysqlv20211122.DBParamValue{}
				if v, ok := paramsMap["param"].(string); ok && v != "" {
					instanceParam.Param = helper.String(v)
				}
				if v, ok := paramsMap["value"].(string); ok && v != "" {
					instanceParam.Value = helper.String(v)
				}
				request.Params = append(request.Params, &instanceParam)
			}
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmysqlV20211122Client().ModifyDBParametersWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify tdmysql db instance params failed, Response is nil."))
			}

			response = result
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update tdmysql db instance params failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		// wait
		flowId := *response.Response.TaskID
		flowRequest := tdmysqlv20211122.NewDescribeFlowRequest()
		flowRequest.FlowId = helper.Int64(flowId)
		flowErr := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			flowResult, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmysqlV20211122Client().DescribeFlowWithContext(ctx, flowRequest)
			if e != nil {
				return tccommon.RetryError(e)
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, flowRequest.GetAction(), flowRequest.ToJsonString(), flowResult.ToJsonString())
			if flowResult == nil || flowResult.Response == nil || flowResult.Response.Status == nil {
				return resource.NonRetryableError(fmt.Errorf("DescribeFlow failed, Response is nil."))
			}

			status := *flowResult.Response.Status
			if status == FLOW_STATUS_SUCCESS {
				return nil
			}

			if status == FLOW_STATUS_FAILED || status == FLOW_STATUS_PAUSED {
				return resource.NonRetryableError(fmt.Errorf("Update tdmysql db instance params failed, status is %s.", status))
			}

			return resource.RetryableError(fmt.Errorf("Update tdmysql db instance params is running, status is %s.", status))
		})

		if flowErr != nil {
			log.Printf("[CRITAL]%s update tdmysql db instance params polling failed, reason:%+v", logId, flowErr)
			return flowErr
		}
	}

	if d.HasChange("auto_renew_flag") {
		request := tdmysqlv20211122.NewModifyAutoRenewFlagRequest()
		request.InstanceIds = helper.Strings([]string{instanceId})
		if v, ok := d.GetOkExists("auto_renew_flag"); ok {
			request.AutoRenewFlag = helper.Int64(int64(v.(int)))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmysqlV20211122Client().ModifyAutoRenewFlagWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify tdmysql db instance auto renew failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update tdmysql db instance auto renew failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	if d.HasChange("security_group_ids") {
		request := tdmysqlv20211122.NewModifyDBInstanceSecurityGroupsRequest()
		request.InstanceId = helper.String(instanceId)
		if v, ok := d.GetOk("security_group_ids"); ok {
			securityGroupIdsList := v.([]interface{})
			for _, item := range securityGroupIdsList {
				request.SecurityGroupIds = append(request.SecurityGroupIds, helper.String(item.(string)))
			}
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmysqlV20211122Client().ModifyDBInstanceSecurityGroupsWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify tdmysql db instance security groups failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update tdmysql db instance security groups failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	if d.HasChange("password") {
		request := tdmysqlv20211122.NewResetUsersPasswordRequest()
		response := tdmysqlv20211122.NewResetUsersPasswordResponse()
		request.InstanceId = helper.String(instanceId)

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmysqlV20211122Client().ResetUsersPasswordWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Reset tdmysql db instance password failed, Response is nil."))
			}

			response = result
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s Reset tdmysql db instance password failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		if response.Response.FlowId == nil {
			return fmt.Errorf("Reset tdmysql db instance password failed, FlowId is nil.")
		}

		// wait
		flowId := *response.Response.FlowId
		flowRequest := tdmysqlv20211122.NewDescribeFlowRequest()
		flowRequest.FlowId = helper.Int64(flowId)
		flowErr := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			flowResult, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmysqlV20211122Client().DescribeFlowWithContext(ctx, flowRequest)
			if e != nil {
				return tccommon.RetryError(e)
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, flowRequest.GetAction(), flowRequest.ToJsonString(), flowResult.ToJsonString())
			if flowResult == nil || flowResult.Response == nil || flowResult.Response.Status == nil {
				return resource.NonRetryableError(fmt.Errorf("DescribeFlow failed, Response is nil."))
			}

			status := *flowResult.Response.Status
			if status == FLOW_STATUS_SUCCESS {
				return nil
			}

			if status == FLOW_STATUS_FAILED || status == FLOW_STATUS_PAUSED {
				return resource.NonRetryableError(fmt.Errorf("Reset tdmysql db instance password failed, status is %s.", status))
			}

			return resource.RetryableError(fmt.Errorf("Reset tdmysql db instance password is running, status is %s.", status))
		})

		if flowErr != nil {
			log.Printf("[CRITAL]%s reset tdmysql db instance password polling failed, reason:%+v", logId, flowErr)
			return flowErr
		}
	}

	if d.HasChange("enable_ssl") {
		request := tdmysqlv20211122.NewModifyInstanceSSLStatusRequest()
		response := tdmysqlv20211122.NewModifyInstanceSSLStatusResponse()
		request.InstanceId = helper.String(instanceId)
		if v, ok := d.GetOkExists("enable_ssl"); ok {
			request.Enabled = helper.Bool(v.(bool))
		}
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmysqlV20211122Client().ModifyInstanceSSLStatusWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify tdmysql db instance ssl failed, Response is nil."))
			}

			response = result
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update tdmysql db instance ssl failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		if response.Response.FlowId == nil {
			return fmt.Errorf("Update tdmysql db instance ssl failed, FlowId is nil.")
		}

		// wait
		flowId := *response.Response.FlowId
		flowRequest := tdmysqlv20211122.NewDescribeFlowRequest()
		flowRequest.FlowId = helper.Int64(flowId)
		flowErr := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			flowResult, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmysqlV20211122Client().DescribeFlowWithContext(ctx, flowRequest)
			if e != nil {
				return tccommon.RetryError(e)
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, flowRequest.GetAction(), flowRequest.ToJsonString(), flowResult.ToJsonString())
			if flowResult == nil || flowResult.Response == nil || flowResult.Response.Status == nil {
				return resource.NonRetryableError(fmt.Errorf("DescribeFlow failed, Response is nil."))
			}

			status := *flowResult.Response.Status
			if status == FLOW_STATUS_SUCCESS {
				return nil
			}

			if status == FLOW_STATUS_FAILED || status == FLOW_STATUS_PAUSED {
				return resource.NonRetryableError(fmt.Errorf("Update tdmysql db instance ssl failed, status is %s.", status))
			}

			return resource.RetryableError(fmt.Errorf("Update tdmysql db instance ssl is running, status is %s.", status))
		})

		if flowErr != nil {
			log.Printf("[CRITAL]%s update tdmysql db instance ssl polling failed, reason:%+v", logId, flowErr)
			return flowErr
		}
	}

	return resourceTencentCloudTdmysqlDbInstanceRead(d, meta)
}

func resourceTencentCloudTdmysqlDbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmysql_db_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		instanceId = d.Id()
	)

	// isolate first
	isolateRequest := tdmysqlv20211122.NewIsolateDBInstanceRequest()
	isolateRequest.InstanceIds = []*string{helper.String(instanceId)}
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmysqlV20211122Client().IsolateDBInstanceWithContext(ctx, isolateRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, isolateRequest.GetAction(), isolateRequest.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Isolate tdmysql db instance failed, Response is nil."))
		}

		if result.Response.SuccessInstanceIds != nil {
			for _, id := range result.Response.SuccessInstanceIds {
				if id != nil && *id == instanceId {
					return nil
				}
			}
		}

		return resource.NonRetryableError(fmt.Errorf("Isolate tdmysql db instance failed failed, instanceId %s not in SuccessInstanceIds.", instanceId))
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s isolate tdmysql db instance failed failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// wait
	waitReq := tdmysqlv20211122.NewDescribeDBInstanceDetailRequest()
	waitReq.InstanceId = helper.String(instanceId)
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		ratelimit.Check(waitReq.GetAction())
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmysqlV20211122Client().DescribeDBInstanceDetail(waitReq)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe tdmysql db instance failed, Response is nil."))
		}

		if result.Response.Status == nil {
			return resource.NonRetryableError(fmt.Errorf("Status is nil."))
		}

		if *result.Response.Status == DB_INSTANCE_STATUS_ISOLATED {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("waiting for tdmysql db instance to be isolated, status is %s.", *result.Response.Status))
	})

	if err != nil {
		return err
	}

	// destroy
	destroyRequest := tdmysqlv20211122.NewDestroyInstancesRequest()
	destroyRequest.InstanceIds = []*string{helper.String(instanceId)}
	reqErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmysqlV20211122Client().DestroyInstancesWithContext(ctx, destroyRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, destroyRequest.GetAction(), destroyRequest.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("destroy tdmysql db instance failed, Response is nil."))
		}

		if result.Response.SuccessInstanceIds != nil {
			for _, id := range result.Response.SuccessInstanceIds {
				if id != nil && *id == instanceId {
					return nil
				}
			}
		}

		return resource.NonRetryableError(fmt.Errorf("destroy tdmysql db instance failed failed, instanceId %s not in SuccessInstanceIds.", instanceId))
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s destroy tdmysql db instance failed failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// wait
	waitReq = tdmysqlv20211122.NewDescribeDBInstanceDetailRequest()
	waitReq.InstanceId = helper.String(instanceId)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		ratelimit.Check(waitReq.GetAction())
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmysqlV20211122Client().DescribeDBInstanceDetail(waitReq)
		if e != nil {
			if sdkErr, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkErr.Code == DESTROY_DB_INSTANCE_SUCCESS_ERROR_CODE {
					return nil
				}
			}

			return tccommon.RetryError(e)
		}

		return resource.RetryableError(fmt.Errorf("waiting for tdmysql db instance to be destroy, status is %s.", *result.Response.Status))
	})

	if err != nil {
		return err
	}

	return nil
}
