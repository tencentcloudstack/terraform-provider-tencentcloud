package cdb

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdbv20170320 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudMysqlAuditService() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlAuditServiceCreate,
		Read:   resourceTencentCloudMysqlAuditServiceRead,
		Update: resourceTencentCloudMysqlAuditServiceUpdate,
		Delete: resourceTencentCloudMysqlAuditServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "TencentDB for MySQL instance ID.",
			},

			"log_expire_day": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Retention period of the audit log. Valid values:  `7` (one week), `30` (one month), `90` (three months), `180` (six months), `365` (one year), `1095` (three years), `1825` (five years).",
			},

			"high_log_expire_day": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Retention period of high-frequency audit logs. Valid values:  `7` (one week), `30` (one month).",
			},

			"rule_template_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Rule template ID. If both this parameter and AuditRuleFilters are not specified, all SQL statements will be recorded.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"audit_all": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Audit type. Valid values: true: Record all; false: Record by rules (default value).",
			},
		},
	}
}

func resourceTencentCloudMysqlAuditServiceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_audit_service.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = cdbv20170320.NewOpenAuditServiceRequest()
		instanceId string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOkExists("log_expire_day"); ok {
		request.LogExpireDay = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("high_log_expire_day"); ok {
		request.HighLogExpireDay = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("rule_template_ids"); ok {
		ruleTemplateIdsSet := v.(*schema.Set).List()
		for i := range ruleTemplateIdsSet {
			if ruleTemplateId, ok := ruleTemplateIdsSet[i].(string); ok && ruleTemplateId != "" {
				request.RuleTemplateIds = append(request.RuleTemplateIds, helper.String(ruleTemplateId))
			}
		}
	}

	if v, ok := d.GetOkExists("audit_all"); ok {
		request.AuditAll = helper.Bool(v.(bool))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().OpenAuditServiceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create mysql audit service failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(instanceId)

	// wait
	waitRequest := cdbv20170320.NewDescribeAuditInstanceListRequest()
	waitRequest.Filters = []*cdbv20170320.AuditInstanceFilters{
		{
			Name:       helper.String("InstanceId"),
			ExactMatch: helper.Bool(true),
			Values:     helper.Strings([]string{instanceId}),
		},
	}

	reqErr = resource.Retry(tccommon.ReadRetryTimeout*5, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().DescribeAuditInstanceListWithContext(ctx, waitRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitRequest.GetAction(), waitRequest.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Items == nil || len(result.Response.Items) == 0 {
			return resource.RetryableError(fmt.Errorf("Describe audit instance list failed, Response is nil."))
		}

		if len(result.Response.Items) != 1 {
			return resource.RetryableError(fmt.Errorf("Describe audit instance list failed, more than one instance item found."))
		}

		item := result.Response.Items[0]
		if item.AuditStatus != nil && *item.AuditStatus == "ON" {
			if item.AuditTask != nil && *item.AuditTask == 0 {
				return nil
			}
		}

		return resource.RetryableError(fmt.Errorf("waiting for mysql [%s] audit service opening", instanceId))
	})

	if reqErr != nil {
		return reqErr
	}

	return resourceTencentCloudMysqlAuditServiceRead(d, meta)
}

func resourceTencentCloudMysqlAuditServiceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_audit_service.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	respData, err := service.DescribeMysqlAuditInstanceListById(ctx, instanceId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_mysql_audit_service` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.InstanceId != nil {
		_ = d.Set("instance_id", respData.InstanceId)
	}

	if respData.LogExpireDay != nil {
		_ = d.Set("log_expire_day", respData.LogExpireDay)
	}

	if respData.HighLogExpireDay != nil {
		_ = d.Set("high_log_expire_day", respData.HighLogExpireDay)
	}

	if respData.RuleTemplateIds != nil {
		_ = d.Set("rule_template_ids", respData.RuleTemplateIds)
	}

	if respData.AuditAll != nil {
		_ = d.Set("audit_all", respData.AuditAll)
	}

	return nil
}

func resourceTencentCloudMysqlAuditServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_audit_service.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		instanceId = d.Id()
	)

	needChange := false
	mutableArgs := []string{"log_expire_day", "high_log_expire_day", "rule_template_ids", "audit_all"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := cdbv20170320.NewModifyAuditServiceRequest()
		if v, ok := d.GetOkExists("log_expire_day"); ok {
			request.LogExpireDay = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("high_log_expire_day"); ok {
			request.HighLogExpireDay = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("rule_template_ids"); ok {
			ruleTemplateIdsSet := v.(*schema.Set).List()
			for i := range ruleTemplateIdsSet {
				if ruleTemplateId, ok := ruleTemplateIdsSet[i].(string); ok && ruleTemplateId != "" {
					request.RuleTemplateIds = append(request.RuleTemplateIds, helper.String(ruleTemplateId))
				}
			}
		}

		if v, ok := d.GetOkExists("audit_all"); ok {
			request.AuditAll = helper.Bool(v.(bool))
		}

		request.InstanceId = &instanceId
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().ModifyAuditServiceWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update mysql audit service failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudMysqlAuditServiceRead(d, meta)
}

func resourceTencentCloudMysqlAuditServiceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_audit_service.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = cdbv20170320.NewCloseAuditServiceRequest()
		instanceId = d.Id()
	)

	request.InstanceId = &instanceId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().CloseAuditServiceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete mysql audit service failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// wait
	waitRequest := cdbv20170320.NewDescribeAuditInstanceListRequest()
	waitRequest.Filters = []*cdbv20170320.AuditInstanceFilters{
		{
			Name:       helper.String("InstanceId"),
			ExactMatch: helper.Bool(true),
			Values:     helper.Strings([]string{instanceId}),
		},
	}

	reqErr = resource.Retry(tccommon.ReadRetryTimeout*5, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().DescribeAuditInstanceListWithContext(ctx, waitRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitRequest.GetAction(), waitRequest.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Items == nil || len(result.Response.Items) == 0 {
			return resource.RetryableError(fmt.Errorf("Describe audit instance list failed, Response is nil."))
		}

		if len(result.Response.Items) != 1 {
			return resource.RetryableError(fmt.Errorf("Describe audit instance list failed, more than one instance item found."))
		}

		item := result.Response.Items[0]
		if item.AuditStatus != nil && *item.AuditStatus == "OFF" {
			if item.AuditTask != nil && *item.AuditTask == 0 {
				return nil
			}
		}

		return resource.RetryableError(fmt.Errorf("waiting for mysql [%s] audit service closing", instanceId))
	})

	if reqErr != nil {
		return reqErr
	}

	return nil
}
