package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMongodbAuditService() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMongodbAuditServiceCreate,
		Read:   resourceTencentCloudMongodbAuditServiceRead,
		Update: resourceTencentCloudMongodbAuditServiceUpdate,
		Delete: resourceTencentCloudMongodbAuditServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance ID, for example: cmgo-xfts****.",
			},
			"log_expire_day": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Audit log retention days. Valid values: 7, 30, 90, 180, 365, 1095, 1825.",
			},
			"audit_all": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to enable full audit. true: full audit, false: rule-based audit. When AuditAll is true, RuleFilters is not required.",
			},
			"rule_filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Audit filter rules. Only required when audit_all is false.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter condition name. Valid values: SrcIp, DB, Collection, User, SqlType.",
						},
						"compare": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter match type. Must be EQ.",
						},
						"value": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Filter match values.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance name.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time when audit was enabled.",
			},
			"log_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Audit log storage type.",
			},
			"is_closing": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Whether audit is being closed.",
			},
			"is_opening": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Whether audit is being opened.",
			},
		},
	}
}

func resourceTencentCloudMongodbAuditServiceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_audit_service.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request    = mongodb.NewOpenAuditServiceRequest()
		instanceId string
	)

	_ = ctx

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	if v, ok := d.GetOkExists("log_expire_day"); ok {
		request.LogExpireDay = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("audit_all"); ok {
		request.AuditAll = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("rule_filters"); ok {
		for _, item := range v.([]interface{}) {
			filterMap := item.(map[string]interface{})
			logFilter := mongodb.LogFilter{}
			if v, ok := filterMap["type"]; ok {
				logFilter.Type = helper.String(v.(string))
			}
			if v, ok := filterMap["compare"]; ok {
				logFilter.Compare = helper.String(v.(string))
			}
			if v, ok := filterMap["value"]; ok {
				valueList := v.([]interface{})
				for _, val := range valueList {
					logFilter.Value = append(logFilter.Value, helper.String(val.(string)))
				}
			}
			request.RuleFilters = append(request.RuleFilters, &logFilter)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().OpenAuditService(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Open mongodb audit service failed, Response is nil"))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mongodb audit service failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		describeRequest := mongodb.NewDescribeAuditConfigRequest()
		describeRequest.InstanceId = helper.String(instanceId)
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().DescribeAuditConfig(describeRequest)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe mongodb audit config failed, Response is nil"))
		}

		if result.Response.IsOpening != nil && *result.Response.IsOpening == "true" {
			return resource.RetryableError(fmt.Errorf("mongodb audit service is still opening"))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mongodb audit service poll failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMongodbAuditServiceRead(d, meta)
}

func resourceTencentCloudMongodbAuditServiceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_audit_service.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request    = mongodb.NewDescribeAuditConfigRequest()
		response   *mongodb.DescribeAuditConfigResponse
		instanceId = d.Id()
	)

	_ = ctx

	request.InstanceId = helper.String(instanceId)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().DescribeAuditConfig(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read mongodb audit service failed, reason:%+v", logId, err)
		return err
	}

	if response == nil || response.Response == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MongodbAuditService` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if response.Response.LogExpireDay != nil {
		_ = d.Set("log_expire_day", response.Response.LogExpireDay)
	}

	if response.Response.AuditAll != nil {
		_ = d.Set("audit_all", response.Response.AuditAll)
	}

	if response.Response.InstanceName != nil {
		_ = d.Set("instance_name", response.Response.InstanceName)
	}

	if response.Response.CreateTime != nil {
		_ = d.Set("create_time", response.Response.CreateTime)
	}

	if response.Response.LogType != nil {
		_ = d.Set("log_type", response.Response.LogType)
	}

	if response.Response.IsClosing != nil {
		_ = d.Set("is_closing", response.Response.IsClosing)
	}

	if response.Response.IsOpening != nil {
		_ = d.Set("is_opening", response.Response.IsOpening)
	}

	return nil
}

func resourceTencentCloudMongodbAuditServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_audit_service.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request    = mongodb.NewModifyAuditServiceRequest()
		instanceId = d.Id()
	)

	_ = ctx

	request.InstanceId = helper.String(instanceId)

	if v, ok := d.GetOkExists("log_expire_day"); ok {
		request.LogExpireDay = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOkExists("audit_all"); ok {
		request.AuditAll = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("rule_filters"); ok {
		for _, item := range v.([]interface{}) {
			filterMap := item.(map[string]interface{})
			logFilter := mongodb.LogFilter{}
			if v, ok := filterMap["type"]; ok {
				logFilter.Type = helper.String(v.(string))
			}
			if v, ok := filterMap["compare"]; ok {
				logFilter.Compare = helper.String(v.(string))
			}
			if v, ok := filterMap["value"]; ok {
				valueList := v.([]interface{})
				for _, val := range valueList {
					logFilter.Value = append(logFilter.Value, helper.String(val.(string)))
				}
			}
			request.RuleFilters = append(request.RuleFilters, &logFilter)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().ModifyAuditService(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update mongodb audit service failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMongodbAuditServiceRead(d, meta)
}

func resourceTencentCloudMongodbAuditServiceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_audit_service.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request    = mongodb.NewCloseAuditServiceRequest()
		instanceId = d.Id()
	)

	_ = ctx

	request.InstanceId = helper.String(instanceId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().CloseAuditService(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete mongodb audit service failed, reason:%+v", logId, err)
		return err
	}

	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		describeRequest := mongodb.NewDescribeAuditConfigRequest()
		describeRequest.InstanceId = helper.String(instanceId)
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().DescribeAuditConfig(describeRequest)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe mongodb audit config failed, Response is nil"))
		}

		if result.Response.IsClosing != nil && *result.Response.IsClosing == "true" {
			return resource.RetryableError(fmt.Errorf("mongodb audit service is still closing"))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete mongodb audit service poll failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
