package es

import (
	"context"
	"fmt"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	elasticsearch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es/v20180416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudElasticsearchLogstash() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudElasticsearchLogstashCreate,
		Read:   resourceTencentCloudElasticsearchLogstashRead,
		Update: resourceTencentCloudElasticsearchLogstashUpdate,
		Delete: resourceTencentCloudElasticsearchLogstashDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance name (compose of 1-50 letter, number, - or _).",
			},

			"zone": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Available zone.",
			},

			"logstash_version": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance version(6.8.13, 7.10.1).",
			},

			"vpc_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "VPC id.",
			},

			"subnet_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Subnet id.",
			},

			"node_num": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Node num(range 2-50).",
			},

			"charge_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Charge type. PREPAID: charged by months or years; POSTPAID_BY_HOUR: charged by hours; default vaule: POSTPAID_BY_HOUR.",
			},

			"charge_period": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Period when charged by months or years(unit depends on TimeUnit).",
			},

			"time_unit": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "charge time unit(set when ChargeType is PREPAID, default value: ms).",
			},

			"auto_voucher": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "whether to use voucher auto, 1 when use, else 0.",
			},

			"voucher_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Voucher list(only can use one voucher by now).",
			},

			"renew_flag": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Automatic renewal flag. RENEW_FLAG_AUTO: auto renewal; RENEW_FLAG_MANUAL: do not renew automatically, users renew manually. It needs to be set when ChargeType is PREPAID. If this parameter is not passed, ordinary users will not renew automatically by default, and SVIP users will renew automatically.",
			},

			"node_type": {
				Optional: true,
				Type:     schema.TypeString,
				Description: "Node type. Valid values:\n" +
					"- LOGSTASH.S1.SMALL2: 1 core 2G;\n" +
					"- LOGSTASH.S1.MEDIUM4:2 core 4G;\n" +
					"- LOGSTASH.S1.MEDIUM8:2 core 8G;\n" +
					"- LOGSTASH.S1.LARGE16:4 core 16G;\n" +
					"- LOGSTASH.S1.2XLARGE32:8 core 32G;\n" +
					"- LOGSTASH.S1.4XLARGE32:16 core 32G;\n" +
					"- LOGSTASH.S1.4XLARGE64:16 core 64G.",
			},

			"disk_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Disk type. CLOUD_SSD: SSD cloud disk; CLOUD_PREMIUM: high hard energy cloud disk; default: CLOUD_SSD.",
			},

			"disk_size": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "node disk size (unit GB).",
			},

			"license_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "License type. oss: open source version; xpack:xpack version; default: xpack.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Description: "Tag description list.",
			},

			"operation_duration": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "operation time by tencent clound.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"periods": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Required:    true,
							Description: "day of week, from Monday to Sunday, value range: [0, 6]notes: may return null when missing.",
						},
						"time_start": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "operation start time.",
						},
						"time_end": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "operation end time.",
						},
						"time_zone": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "time zone, for example: UTC+8.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudElasticsearchLogstashCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elasticsearch_logstash.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request    = elasticsearch.NewCreateLogstashInstanceRequest()
		response   = elasticsearch.NewCreateLogstashInstanceResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("zone"); ok {
		request.Zone = helper.String(v.(string))
	}

	if v, ok := d.GetOk("logstash_version"); ok {
		request.LogstashVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("node_num"); ok {
		request.NodeNum = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("charge_type"); ok {
		request.ChargeType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("charge_period"); ok {
		request.ChargePeriod = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("time_unit"); ok {
		request.TimeUnit = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("auto_voucher"); ok {
		request.AutoVoucher = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("voucher_ids"); ok {
		voucherIdsSet := v.(*schema.Set).List()
		for i := range voucherIdsSet {
			voucherIds := voucherIdsSet[i].(string)
			request.VoucherIds = append(request.VoucherIds, &voucherIds)
		}
	}

	if v, ok := d.GetOk("renew_flag"); ok {
		request.RenewFlag = helper.String(v.(string))
	}

	if v, ok := d.GetOk("node_type"); ok {
		request.NodeType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("disk_type"); ok {
		request.DiskType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("disk_size"); ok {
		request.DiskSize = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("license_type"); ok {
		request.LicenseType = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "operation_duration"); ok {
		operationDuration := elasticsearch.OperationDuration{}
		if v, ok := dMap["periods"]; ok {
			periodsSet := v.(*schema.Set).List()
			for i := range periodsSet {
				periods := periodsSet[i].(int)
				operationDuration.Periods = append(operationDuration.Periods, helper.IntUint64(periods))
			}
		}
		if v, ok := dMap["time_start"]; ok {
			operationDuration.TimeStart = helper.String(v.(string))
		}
		if v, ok := dMap["time_end"]; ok {
			operationDuration.TimeEnd = helper.String(v.(string))
		}
		if v, ok := dMap["time_zone"]; ok {
			operationDuration.TimeZone = helper.String(v.(string))
		}
		request.OperationDuration = &operationDuration
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEsClient().CreateLogstashInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create elasticsearch logstash failed, reason:%+v", logId, err)
		return err
	}
	instanceId = *response.Response.InstanceId

	service := ElasticsearchService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"1"}, 3*tccommon.ReadRetryTimeout, time.Second, service.ElasticsearchLogstashStateRefreshFunc(instanceId, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}
	d.SetId(instanceId)

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := fmt.Sprintf("qcs::es:%s:uin/:logstash/%s", region, instanceId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudElasticsearchLogstashRead(d, meta)
}

func resourceTencentCloudElasticsearchLogstashRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elasticsearch_logstash.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ElasticsearchService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()

	logstash, err := service.DescribeElasticsearchLogstashById(ctx, instanceId)
	if err != nil {
		return err
	}

	if logstash == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ElasticsearchLogstash` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if logstash.InstanceName != nil {
		_ = d.Set("instance_name", logstash.InstanceName)
	}

	if logstash.Zone != nil {
		_ = d.Set("zone", logstash.Zone)
	}

	if logstash.LogstashVersion != nil {
		_ = d.Set("logstash_version", logstash.LogstashVersion)
	}

	if logstash.VpcId != nil {
		_ = d.Set("vpc_id", logstash.VpcId)
	}

	if logstash.SubnetId != nil {
		_ = d.Set("subnet_id", logstash.SubnetId)
	}

	if logstash.NodeNum != nil {
		_ = d.Set("node_num", logstash.NodeNum)
	}

	if logstash.ChargeType != nil {
		_ = d.Set("charge_type", logstash.ChargeType)
	}

	if logstash.ChargePeriod != nil {
		_ = d.Set("charge_period", logstash.ChargePeriod)
	}

	if logstash.RenewFlag != nil {
		_ = d.Set("renew_flag", logstash.RenewFlag)
	}

	if logstash.NodeType != nil {
		_ = d.Set("node_type", logstash.NodeType)
	}

	if logstash.DiskType != nil {
		_ = d.Set("disk_type", logstash.DiskType)
	}

	if logstash.DiskSize != nil {
		_ = d.Set("disk_size", logstash.DiskSize)
	}

	if logstash.LicenseType != nil {
		_ = d.Set("license_type", logstash.LicenseType)
	}

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "es", "logstash", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	if logstash.OperationDuration != nil {
		operationDurationMap := map[string]interface{}{}

		if logstash.OperationDuration.Periods != nil {
			operationDurationMap["periods"] = logstash.OperationDuration.Periods
		}

		if logstash.OperationDuration.TimeStart != nil {
			operationDurationMap["time_start"] = logstash.OperationDuration.TimeStart
		}

		if logstash.OperationDuration.TimeEnd != nil {
			operationDurationMap["time_end"] = logstash.OperationDuration.TimeEnd
		}

		if logstash.OperationDuration.TimeZone != nil {
			operationDurationMap["time_zone"] = logstash.OperationDuration.TimeZone
		}

		_ = d.Set("operation_duration", []interface{}{operationDurationMap})
	}

	return nil
}

func resourceTencentCloudElasticsearchLogstashUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elasticsearch_logstash.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	if d.HasChange("tags") {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := tccommon.BuildTagResourceName("es", "logstash", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	request := elasticsearch.NewUpdateLogstashInstanceRequest()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"zone", "logstash_version", "vpc_id", "subnet_id", "charge_type", "charge_period", "time_unit", "auto_voucher", "voucher_ids", "renew_flag", "disk_type", "license_type"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	service := ElasticsearchService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	if d.HasChange("instance_name") {
		if v, ok := d.GetOk("instance_name"); ok {
			err := service.UpdateLogstashInstance(ctx, instanceId, map[string]interface{}{"instance_name": v.(string)})
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("node_num") {
		if v, ok := d.GetOkExists("node_num"); ok {
			err := service.UpdateLogstashInstance(ctx, instanceId, map[string]interface{}{"node_num": v.(int)})
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("node_type") {
		if v, ok := d.GetOk("node_type"); ok {
			request.NodeType = helper.String(v.(string))
			err := service.UpdateLogstashInstance(ctx, instanceId, map[string]interface{}{"node_type": v.(string)})
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("disk_size") {
		if v, ok := d.GetOkExists("disk_size"); ok {
			err := service.UpdateLogstashInstance(ctx, instanceId, map[string]interface{}{"disk_size": v.(int)})
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("operation_duration") {
		if dMap, ok := helper.InterfacesHeadMap(d, "operation_duration"); ok {
			operationDuration := elasticsearch.OperationDurationUpdated{}
			if v, ok := dMap["periods"]; ok {
				periodsSet := v.(*schema.Set).List()
				for i := range periodsSet {
					periods := periodsSet[i].(int)
					operationDuration.Periods = append(operationDuration.Periods, helper.IntUint64(periods))
				}
			}
			if v, ok := dMap["time_start"]; ok {
				operationDuration.TimeStart = helper.String(v.(string))
			}
			if v, ok := dMap["time_end"]; ok {
				operationDuration.TimeEnd = helper.String(v.(string))
			}
			if v, ok := dMap["time_zone"]; ok {
				operationDuration.TimeZone = helper.String(v.(string))
			}
			err := service.UpdateLogstashInstance(ctx, instanceId, map[string]interface{}{"operation_duration": operationDuration})
			if err != nil {
				return err
			}
		}
	}

	return resourceTencentCloudElasticsearchLogstashRead(d, meta)
}

func resourceTencentCloudElasticsearchLogstashDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elasticsearch_logstash.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ElasticsearchService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	instanceId := d.Id()

	if err := service.DeleteElasticsearchLogstashById(ctx, instanceId); err != nil {
		return err
	}
	conf := tccommon.BuildStateChangeConf([]string{}, []string{"-99"}, 3*tccommon.ReadRetryTimeout, time.Second, service.ElasticsearchLogstashStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}
	return nil
}
