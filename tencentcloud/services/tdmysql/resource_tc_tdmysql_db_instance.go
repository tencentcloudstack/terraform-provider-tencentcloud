package tdmysql

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmysqlv20211122 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmysql/v20211122"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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
		Schema: map[string]*schema.Schema{
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Instance region.",
			},

			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "VPC ID.",
			},

			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "VPC subnet ID.",
			},

			"spec_code": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Purchase specification.",
			},

			"disk": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Storage node disk capacity, unit: GB.",
			},

			"storage_node_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Number of storage nodes.",
			},

			"replications": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Number of storage node replicas, up to 5, must be odd.",
			},

			"instance_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Number of instances to create, up to 10.",
			},

			"full_replications": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Number of full replicas.",
			},

			"create_version": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Instance version, use the latest version by default.",
			},

			"instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Instance name, length 1-60, supports Chinese/English/numbers/-/_.",
			},

			"resource_tags": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Tag key-value pair array.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag value.",
						},
					},
				},
			},

			"init_params": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Instance init params, such as character_set_server, lower_case_table_names.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"param": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Config param key.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Config param value.",
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
				Description: "Time span of the product.",
			},

			"storage_node_cpu": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Storage node CPU cores.",
			},

			"storage_node_mem": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Storage node memory size.",
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
				Description: "Management control node number.",
			},

			"vport": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Custom port.",
			},

			"zones": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Multi-AZ list.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"auto_voucher": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Whether to use voucher.",
			},

			"voucher_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Voucher list.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"instance_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Instance architecture type, separate: separate architecture; hybrid: peer architecture.",
			},

			"storage_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Disk type, CLOUD_HSSD: enhanced SSD, CLOUD_TCS: local SSD.",
			},

			"az_mode": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "AZ mode, 1: single AZ, 2: multi-AZ non-primary AZ, 3: multi-AZ primary AZ.",
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
				Description: "Compatibility mode, enum: MySQL, HBase.",
			},

			"auto_scale_config": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "CCU scaling config of svls instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"range_min": {
							Type:        schema.TypeFloat,
							Optional:    true,
							Description: "CCU min value.",
						},
						"range_max": {
							Type:        schema.TypeFloat,
							Optional:    true,
							Description: "CCU max value.",
						},
					},
				},
			},

			"security_group_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Security group list to bind.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Root username, default dbaadmin, will be reset to dbaadmin.",
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "Dbaadmin password.",
			},

			"encryption_enable": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Whether to enable transparent encryption, 0: no, 1: yes.",
			},

			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance ID.",
			},

			"vip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Subnet IP.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance status.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Update time.",
			},

			"char_set": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Character set.",
			},

			"node": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Node info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node IP.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node type, such as sqlengine, tdstore, mc.",
						},
						"node_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node unique ID.",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Node port.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node zone.",
						},
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node host IP.",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Node CPU count.",
						},
						"mem": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Node memory size.",
						},
						"data_disk": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Node disk size.",
						},
						"binlog_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Node binlog service info.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sid": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Binlog service unique ID.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Binlog service type.",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Binlog instance unique ID.",
									},
								},
							},
						},
					},
				},
			},

			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance region.",
			},

			"status_desc": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance status description in Chinese.",
			},

			"renew_flag": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Renew flag.",
			},

			"expire_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expire time.",
			},

			"isolated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Isolate time.",
			},

			"disk_usage": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Max node disk usage.",
			},

			"binlog_status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Whether binlog is enabled.",
			},

			"standby_flag": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Standby flag, 1: no standby, 2: standby primary, 3: standby secondary.",
			},

			"binlog_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CDC node type.",
			},

			"timing_modify_instance_flag": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "1 means supported, 0 means not supported.",
			},

			"columnar_node_cpu": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Columnar node CPU cores.",
			},

			"columnar_node_mem": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Columnar node memory size.",
			},

			"columnar_node_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Columnar node count.",
			},

			"columnar_node_disk": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Columnar node disk size.",
			},

			"columnar_node_storage_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Columnar node disk type.",
			},

			"columnar_node_spec_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Columnar node spec code.",
			},

			"columnar_vip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Columnar vip.",
			},

			"columnar_vport": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Columnar vport.",
			},

			"is_support_columnar": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the instance supports columnar.",
			},

			"instance_category": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Instance category.",
			},

			"is_switch_full_replications_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether support modifying full replica count.",
			},

			"dumper_vip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Dumper vip.",
			},

			"dumper_vport": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Dumper vport.",
			},

			"template_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Parameter template name.",
			},

			"analysis_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance analysis engine mode, enum: libra.",
			},

			"analysis_relation_infos": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Analysis engine instance relation.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"primary_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Source instance ID.",
						},
						"analysis_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Analysis engine instance ID.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Analysis engine relation status.",
						},
						"create_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time.",
						},
						"update_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Update time.",
						},
					},
				},
			},

			"analysis_instance_info": {
				Type:        schema.TypeList,
				Computed:    true,
				MaxItems:    1,
				Description: "Analysis engine instance info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"replicas_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Replicas num.",
						},
					},
				},
			},

			"maintenance_window": {
				Type:        schema.TypeList,
				Computed:    true,
				MaxItems:    1,
				Description: "Maintenance window config.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Start time.",
						},
						"duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Duration.",
						},
						"week_days": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Week days.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"encryption_kms_region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Real kms region used for kms service.",
			},

			"instance_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Instance ID list returned by CreateDBInstances.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"flow_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Flow ID returned by CreateDBInstances.",
			},
		},
	}
}

func resourceTencentCloudTdmysqlDbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmysql_db_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = tdmysqlv20211122.NewCreateDBInstancesRequest()
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

	if v, ok := d.GetOkExists("instance_count"); ok {
		request.InstanceCount = helper.Int64(int64(v.(int)))
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
			resourceTagsMap := item.(map[string]interface{})
			resourceTag := tdmysqlv20211122.ResourceTag{}
			if v, ok := resourceTagsMap["tag_key"].(string); ok && v != "" {
				resourceTag.TagKey = helper.String(v)
			}

			if v, ok := resourceTagsMap["tag_value"].(string); ok && v != "" {
				resourceTag.TagValue = helper.String(v)
			}

			request.ResourceTags = append(request.ResourceTags, &resourceTag)
		}
	}

	if v, ok := d.GetOk("init_params"); ok {
		for _, item := range v.([]interface{}) {
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
		zonesSet := v.([]interface{})
		for i := range zonesSet {
			zone := zonesSet[i].(string)
			request.Zones = append(request.Zones, helper.String(zone))
		}
	}

	if v, ok := d.GetOkExists("auto_voucher"); ok {
		request.AutoVoucher = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("voucher_ids"); ok {
		voucherIdsSet := v.([]interface{})
		for i := range voucherIdsSet {
			voucherId := voucherIdsSet[i].(string)
			request.VoucherIds = append(request.VoucherIds, helper.String(voucherId))
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
		autoScaleConfigs := v.([]interface{})
		if len(autoScaleConfigs) > 0 {
			autoScaleConfigMap := autoScaleConfigs[0].(map[string]interface{})
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
		securityGroupIdsSet := v.([]interface{})
		for i := range securityGroupIdsSet {
			securityGroupId := securityGroupIdsSet[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, helper.String(securityGroupId))
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

	var instanceId string
	var flowId int64
	var response *tdmysqlv20211122.CreateDBInstancesResponse
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmysqlV20211122Client().CreateDBInstancesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create tdmysql db instance failed, Response is nil."))
		}

		if result.Response.InstanceIds == nil || len(result.Response.InstanceIds) == 0 {
			return resource.NonRetryableError(fmt.Errorf("Create tdmysql db instance failed, InstanceIds is empty."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create tdmysql db instance failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	instanceId = *response.Response.InstanceIds[0]
	if response.Response.FlowId != nil {
		flowId = *response.Response.FlowId
	}
	log.Printf("[DEBUG]%s tdmysql_db_instance create success, logId=%s, id=%s, flowId=%d", logId, logId, instanceId, flowId)

	if response.Response.FlowId != nil {
		tdmysqlService := NewTdmysqlService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			status, e := tdmysqlService.DescribeTdmysqlFlow(ctx, flowId)
			if e != nil {
				return tccommon.RetryError(e)
			}

			if status == "success" {
				return nil
			}

			if status == "failed" || status == "paused" {
				return resource.NonRetryableError(fmt.Errorf("Create tdmysql db instance flow failed, status is %s.", status))
			}

			return resource.RetryableError(fmt.Errorf("Create tdmysql db instance flow is running, status is %s.", status))
		})

		if err != nil {
			log.Printf("[CRITAL]%s create tdmysql db instance flow polling failed, id=%s, flowId=%d, reason:%+v", logId, instanceId, flowId, err)
			return err
		}
	}

	d.SetId(instanceId)
	if response.Response.InstanceIds != nil {
		instanceIdsList := make([]string, 0, len(response.Response.InstanceIds))
		for _, id := range response.Response.InstanceIds {
			if id != nil {
				instanceIdsList = append(instanceIdsList, *id)
			}
		}
		_ = d.Set("instance_ids", instanceIdsList)
	}

	_ = d.Set("flow_id", flowId)

	return resourceTencentCloudTdmysqlDbInstanceRead(d, meta)
}

func resourceTencentCloudTdmysqlDbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmysql_db_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = NewTdmysqlService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	respData, err := service.DescribeTdmysqlDbInstanceById(ctx, d.Id())
	if err != nil {
		return err
	}

	if respData == nil || respData.InstanceId == nil {
		log.Printf("[CRUD] tdmysql_db_instance id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if respData.InstanceName != nil {
		_ = d.Set("instance_name", respData.InstanceName)
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

	if respData.CreateVersion != nil {
		_ = d.Set("create_version", respData.CreateVersion)
	}

	if respData.Vip != nil {
		_ = d.Set("vip", respData.Vip)
	}

	if respData.Vport != nil {
		_ = d.Set("vport", respData.Vport)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.Disk != nil {
		_ = d.Set("disk", respData.Disk)
	}

	if respData.StorageNodeNum != nil {
		_ = d.Set("storage_node_num", respData.StorageNodeNum)
	}

	if respData.InitParams != nil && len(respData.InitParams) > 0 {
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

	if respData.ResourceTags != nil && len(respData.ResourceTags) > 0 {
		resourceTagsList := make([]map[string]interface{}, 0, len(respData.ResourceTags))
		for _, resourceTag := range respData.ResourceTags {
			resourceTagsMap := map[string]interface{}{}
			if resourceTag.TagKey != nil {
				resourceTagsMap["tag_key"] = resourceTag.TagKey
			}

			if resourceTag.TagValue != nil {
				resourceTagsMap["tag_value"] = resourceTag.TagValue
			}

			resourceTagsList = append(resourceTagsList, resourceTagsMap)
		}

		_ = d.Set("resource_tags", resourceTagsList)
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.UpdateTime != nil {
		_ = d.Set("update_time", respData.UpdateTime)
	}

	if respData.Replications != nil {
		_ = d.Set("replications", respData.Replications)
	}

	if respData.FullReplications != nil {
		_ = d.Set("full_replications", respData.FullReplications)
	}

	if respData.CharSet != nil {
		_ = d.Set("char_set", respData.CharSet)
	}

	if respData.Node != nil && len(respData.Node) > 0 {
		nodeList := make([]map[string]interface{}, 0, len(respData.Node))
		for _, node := range respData.Node {
			nodeMap := map[string]interface{}{}
			if node.IP != nil {
				nodeMap["ip"] = node.IP
			}

			if node.Type != nil {
				nodeMap["type"] = node.Type
			}

			if node.NodeId != nil {
				nodeMap["node_id"] = node.NodeId
			}

			if node.Port != nil {
				nodeMap["port"] = node.Port
			}

			if node.Zone != nil {
				nodeMap["zone"] = node.Zone
			}

			if node.Host != nil {
				nodeMap["host"] = node.Host
			}

			if node.Cpu != nil {
				nodeMap["cpu"] = node.Cpu
			}

			if node.Mem != nil {
				nodeMap["mem"] = node.Mem
			}

			if node.DataDisk != nil {
				nodeMap["data_disk"] = node.DataDisk
			}

			if node.BinlogInfo != nil && len(node.BinlogInfo) > 0 {
				binlogInfoList := make([]map[string]interface{}, 0, len(node.BinlogInfo))
				for _, binlogInfo := range node.BinlogInfo {
					binlogInfoMap := map[string]interface{}{}
					if binlogInfo.Sid != nil {
						binlogInfoMap["sid"] = binlogInfo.Sid
					}

					if binlogInfo.Type != nil {
						binlogInfoMap["type"] = binlogInfo.Type
					}

					if binlogInfo.InstanceId != nil {
						binlogInfoMap["instance_id"] = binlogInfo.InstanceId
					}

					binlogInfoList = append(binlogInfoList, binlogInfoMap)
				}

				nodeMap["binlog_info"] = binlogInfoList
			}

			nodeList = append(nodeList, nodeMap)
		}

		_ = d.Set("node", nodeList)
	}

	if respData.Region != nil {
		_ = d.Set("region", respData.Region)
	}

	if respData.SpecCode != nil {
		_ = d.Set("spec_code", respData.SpecCode)
	}

	if respData.InstanceId != nil {
		_ = d.Set("instance_id", respData.InstanceId)
	}

	if respData.StatusDesc != nil {
		_ = d.Set("status_desc", respData.StatusDesc)
	}

	if respData.StorageNodeCpu != nil {
		_ = d.Set("storage_node_cpu", respData.StorageNodeCpu)
	}

	if respData.StorageNodeMem != nil {
		_ = d.Set("storage_node_mem", respData.StorageNodeMem)
	}

	if respData.RenewFlag != nil {
		_ = d.Set("renew_flag", respData.RenewFlag)
	}

	if respData.PayMode != nil {
		_ = d.Set("pay_mode", respData.PayMode)
	}

	if respData.ExpireAt != nil {
		_ = d.Set("expire_at", respData.ExpireAt)
	}

	if respData.IsolatedAt != nil {
		_ = d.Set("isolated_at", respData.IsolatedAt)
	}

	if respData.InstanceType != nil {
		_ = d.Set("instance_type", respData.InstanceType)
	}

	if respData.StorageType != nil {
		_ = d.Set("storage_type", respData.StorageType)
	}

	if respData.Zones != nil {
		zonesList := make([]string, 0, len(respData.Zones))
		for _, zone := range respData.Zones {
			if zone != nil {
				zonesList = append(zonesList, *zone)
			}
		}
		_ = d.Set("zones", zonesList)
	}

	if respData.DiskUsage != nil {
		_ = d.Set("disk_usage", respData.DiskUsage)
	}

	if respData.BinlogStatus != nil {
		_ = d.Set("binlog_status", respData.BinlogStatus)
	}

	if respData.AZMode != nil {
		_ = d.Set("az_mode", respData.AZMode)
	}

	if respData.StandbyFlag != nil {
		_ = d.Set("standby_flag", respData.StandbyFlag)
	}

	if respData.BinlogType != nil {
		_ = d.Set("binlog_type", respData.BinlogType)
	}

	if respData.TimingModifyInstanceFlag != nil {
		_ = d.Set("timing_modify_instance_flag", respData.TimingModifyInstanceFlag)
	}

	if respData.ColumnarNodeCpu != nil {
		_ = d.Set("columnar_node_cpu", respData.ColumnarNodeCpu)
	}

	if respData.ColumnarNodeMem != nil {
		_ = d.Set("columnar_node_mem", respData.ColumnarNodeMem)
	}

	if respData.ColumnarNodeNum != nil {
		_ = d.Set("columnar_node_num", respData.ColumnarNodeNum)
	}

	if respData.ColumnarNodeDisk != nil {
		_ = d.Set("columnar_node_disk", respData.ColumnarNodeDisk)
	}

	if respData.ColumnarNodeStorageType != nil {
		_ = d.Set("columnar_node_storage_type", respData.ColumnarNodeStorageType)
	}

	if respData.ColumnarNodeSpecCode != nil {
		_ = d.Set("columnar_node_spec_code", respData.ColumnarNodeSpecCode)
	}

	if respData.ColumnarVip != nil {
		_ = d.Set("columnar_vip", respData.ColumnarVip)
	}

	if respData.ColumnarVport != nil {
		_ = d.Set("columnar_vport", respData.ColumnarVport)
	}

	if respData.IsSupportColumnar != nil {
		_ = d.Set("is_support_columnar", respData.IsSupportColumnar)
	}

	if respData.InstanceCategory != nil {
		_ = d.Set("instance_category", respData.InstanceCategory)
	}

	if respData.SQLMode != nil {
		_ = d.Set("sql_mode", respData.SQLMode)
	}

	if respData.IsSwitchFullReplicationsEnable != nil {
		_ = d.Set("is_switch_full_replications_enable", respData.IsSwitchFullReplicationsEnable)
	}

	if respData.InstanceMode != nil {
		_ = d.Set("instance_mode", respData.InstanceMode)
	}

	if respData.DumperVip != nil {
		_ = d.Set("dumper_vip", respData.DumperVip)
	}

	if respData.DumperVport != nil {
		_ = d.Set("dumper_vport", respData.DumperVport)
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

	if respData.TemplateId != nil {
		_ = d.Set("template_id", respData.TemplateId)
	}

	if respData.TemplateName != nil {
		_ = d.Set("template_name", respData.TemplateName)
	}

	if respData.AnalysisMode != nil {
		_ = d.Set("analysis_mode", respData.AnalysisMode)
	}

	if respData.AnalysisRelationInfos != nil && len(respData.AnalysisRelationInfos) > 0 {
		analysisRelationInfosList := make([]map[string]interface{}, 0, len(respData.AnalysisRelationInfos))
		for _, analysisRelationInfo := range respData.AnalysisRelationInfos {
			analysisRelationInfoMap := map[string]interface{}{}
			if analysisRelationInfo.PrimaryInstanceId != nil {
				analysisRelationInfoMap["primary_instance_id"] = analysisRelationInfo.PrimaryInstanceId
			}

			if analysisRelationInfo.AnalysisInstanceId != nil {
				analysisRelationInfoMap["analysis_instance_id"] = analysisRelationInfo.AnalysisInstanceId
			}

			if analysisRelationInfo.Status != nil {
				analysisRelationInfoMap["status"] = analysisRelationInfo.Status
			}

			if analysisRelationInfo.CreateAt != nil {
				analysisRelationInfoMap["create_at"] = analysisRelationInfo.CreateAt
			}

			if analysisRelationInfo.UpdateAt != nil {
				analysisRelationInfoMap["update_at"] = analysisRelationInfo.UpdateAt
			}

			analysisRelationInfosList = append(analysisRelationInfosList, analysisRelationInfoMap)
		}

		_ = d.Set("analysis_relation_infos", analysisRelationInfosList)
	}

	if respData.AnalysisInstanceInfo != nil {
		analysisInstanceInfoList := make([]map[string]interface{}, 0, 1)
		analysisInstanceInfoMap := map[string]interface{}{}
		if respData.AnalysisInstanceInfo.ReplicasNum != nil {
			analysisInstanceInfoMap["replicas_num"] = respData.AnalysisInstanceInfo.ReplicasNum
		}

		analysisInstanceInfoList = append(analysisInstanceInfoList, analysisInstanceInfoMap)
		_ = d.Set("analysis_instance_info", analysisInstanceInfoList)
	}

	if respData.MaintenanceWindow != nil {
		maintenanceWindowList := make([]map[string]interface{}, 0, 1)
		maintenanceWindowMap := map[string]interface{}{}
		if respData.MaintenanceWindow.StartTime != nil {
			maintenanceWindowMap["start_time"] = respData.MaintenanceWindow.StartTime
		}

		if respData.MaintenanceWindow.Duration != nil {
			maintenanceWindowMap["duration"] = respData.MaintenanceWindow.Duration
		}

		if respData.MaintenanceWindow.WeekDays != nil {
			weekDaysList := make([]string, 0, len(respData.MaintenanceWindow.WeekDays))
			for _, weekDay := range respData.MaintenanceWindow.WeekDays {
				if weekDay != nil {
					weekDaysList = append(weekDaysList, *weekDay)
				}
			}
			maintenanceWindowMap["week_days"] = weekDaysList
		}

		maintenanceWindowList = append(maintenanceWindowList, maintenanceWindowMap)
		_ = d.Set("maintenance_window", maintenanceWindowList)
	}

	if respData.EncryptionEnable != nil {
		_ = d.Set("encryption_enable", respData.EncryptionEnable)
	}

	if respData.EncryptionKmsRegion != nil {
		_ = d.Set("encryption_kms_region", respData.EncryptionKmsRegion)
	}

	return nil
}

func resourceTencentCloudTdmysqlDbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmysql_db_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	immutableArgs := []string{
		"zone", "vpc_id", "subnet_id", "spec_code", "disk", "storage_node_num",
		"replications", "instance_count", "full_replications", "create_version",
		"resource_tags", "init_params", "time_unit", "time_span", "storage_node_cpu",
		"storage_node_mem", "pay_mode", "mc_num", "vport", "zones", "auto_voucher",
		"voucher_ids", "instance_type", "storage_type", "az_mode", "instance_mode",
		"template_id", "sql_mode", "auto_scale_config", "security_group_ids",
		"user_name", "password", "encryption_enable",
	}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("instance_name") {
		request := tdmysqlv20211122.NewModifyInstanceNameRequest()
		request.InstanceId = helper.String(d.Id())
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

	return resourceTencentCloudTdmysqlDbInstanceRead(d, meta)
}

func resourceTencentCloudTdmysqlDbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmysql_db_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = NewTdmysqlService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	err := service.IsolateTdmysqlDbInstance(ctx, d.Id())
	if err != nil {
		log.Printf("[CRITAL]%s delete tdmysql db instance failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
