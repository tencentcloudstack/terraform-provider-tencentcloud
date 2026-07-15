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
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
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
			"log_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Audit log storage type.",
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
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().OpenAuditServiceWithContext(ctx, request)
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

	// wait
	describeRequest := mongodb.NewDescribeAuditInstanceListRequest()
	describeRequest.Filters = []*mongodb.Filters{
		{
			Name:   helper.String("InstanceId"),
			Values: []*string{helper.String(instanceId)},
		},
	}
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().DescribeAuditInstanceListWithContext(ctx, describeRequest)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe mongodb audit instance list failed, Response is nil"))
		}

		var auditStatus string
		var auditTask int64
		for _, item := range result.Response.Items {
			if item != nil && item.InstanceId != nil && *item.InstanceId == instanceId {
				if item.AuditStatus != nil && item.AuditTask != nil {
					auditStatus = *item.AuditStatus
					auditTask = *item.AuditTask
				}

				break
			}
		}

		if auditStatus == "ON" && auditTask == 0 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("mongodb audit service is still opening, current status: %s, task: %d", auditStatus, auditTask))
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

	request.InstanceId = helper.String(instanceId)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().DescribeAuditConfigWithContext(ctx, request)
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
		log.Printf("[WARN]%s resource `tencentcloud_mongodb_audit_service` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
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

	if response.Response.LogType != nil {
		_ = d.Set("log_type", response.Response.LogType)
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
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().ModifyAuditServiceWithContext(ctx, request)
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

	request.InstanceId = helper.String(instanceId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().CloseAuditServiceWithContext(ctx, request)
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

	// wait
	describeRequest := mongodb.NewDescribeAuditInstanceListRequest()
	describeRequest.Filters = []*mongodb.Filters{
		{
			Name:   helper.String("InstanceId"),
			Values: []*string{helper.String(instanceId)},
		},
	}
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().DescribeAuditInstanceListWithContext(ctx, describeRequest)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe mongodb audit instance list failed, Response is nil"))
		}

		var auditStatus string
		var auditTask int64
		for _, item := range result.Response.Items {
			if item != nil && item.InstanceId != nil && *item.InstanceId == instanceId {
				if item.AuditStatus != nil && item.AuditTask != nil {
					auditStatus = *item.AuditStatus
					auditTask = *item.AuditTask
				}

				break
			}
		}

		if auditStatus == "OFF" && auditTask == 0 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("mongodb audit service is still closing, current status: %s, task: %d", auditStatus, auditTask))
	})

	if err != nil {
		log.Printf("[CRITAL]%s close mongodb audit service poll failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
