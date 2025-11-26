package cynosdb

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdbv20190107 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCynosdbAuditService() *schema.Resource {
	return &schema.Resource{
		Create: ResourceTencentCloudCynosdbAuditServiceCreate,
		Read:   ResourceTencentCloudCynosdbAuditServiceRead,
		Update: ResourceTencentCloudCynosdbAuditServiceUpdate,
		Delete: ResourceTencentCloudCynosdbAuditServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance ID.",
			},

			"log_expire_day": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Log retention period.",
			},

			"high_log_expire_day": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Frequent log retention period.",
			},

			"rule_template_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Rule template ID set.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"audit_all": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Audit type. true - full audit; default false - rule-based audit.",
			},
		},
	}
}

func ResourceTencentCloudCynosdbAuditServiceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_audit_service.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = cynosdbv20190107.NewOpenAuditServiceRequest()
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
			ruleTemplateIds := ruleTemplateIdsSet[i].(string)
			request.RuleTemplateIds = append(request.RuleTemplateIds, helper.String(ruleTemplateIds))
		}
	}

	if v, ok := d.GetOkExists("audit_all"); ok {
		request.AuditAll = helper.Bool(v.(bool))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().OpenAuditServiceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create cynosdb audit service failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(instanceId)

	// wait
	waitReq := cynosdbv20190107.NewDescribeAuditInstanceListRequest()
	waitReq.Offset = helper.Uint64(0)
	waitReq.Limit = helper.Uint64(1)
	waitReq.Filters = []*cynosdbv20190107.AuditInstanceFilters{
		{
			Name:       helper.String("InstanceId"),
			ExactMatch: helper.Bool(true),
			Values:     helper.Strings([]string{instanceId}),
		},
	}

	reqErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().DescribeAuditInstanceListWithContext(ctx, waitReq)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Items == nil || len(result.Response.Items) == 0 {
			return resource.NonRetryableError(fmt.Errorf("Describe cynosdb audit service failed, Response is nil."))
		}

		item := result.Response.Items[0]
		if item.AuditStatus == nil {
			return resource.NonRetryableError(fmt.Errorf("AuditStatus is nil."))
		}

		if *item.AuditStatus == "ON" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("cynosdb audit service is still running, audit status is %s.", *item.AuditStatus))
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create cynosdb audit service failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return ResourceTencentCloudCynosdbAuditServiceRead(d, meta)
}

func ResourceTencentCloudCynosdbAuditServiceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_audit_service.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	respData, err := service.DescribeCynosdbAuditServiceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_cynosdb_audit_service` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
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

func ResourceTencentCloudCynosdbAuditServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_audit_service.update")()
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
		request := cynosdbv20190107.NewModifyAuditServiceRequest()
		if v, ok := d.GetOkExists("log_expire_day"); ok {
			request.LogExpireDay = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("high_log_expire_day"); ok {
			request.HighLogExpireDay = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("rule_template_ids"); ok {
			ruleTemplateIdsSet := v.(*schema.Set).List()
			for i := range ruleTemplateIdsSet {
				ruleTemplateIds := ruleTemplateIdsSet[i].(string)
				request.RuleTemplateIds = append(request.RuleTemplateIds, helper.String(ruleTemplateIds))
			}
		}

		if v, ok := d.GetOkExists("audit_all"); ok {
			request.AuditAll = helper.Bool(v.(bool))
		}

		request.InstanceId = &instanceId
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().ModifyAuditServiceWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update cynosdb audit service failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return ResourceTencentCloudCynosdbAuditServiceRead(d, meta)
}

func ResourceTencentCloudCynosdbAuditServiceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_audit_service.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = cynosdbv20190107.NewCloseAuditServiceRequest()
		instanceId = d.Id()
	)

	request.InstanceId = &instanceId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().CloseAuditServiceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete cynosdb audit service failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// wait
	waitReq := cynosdbv20190107.NewDescribeAuditInstanceListRequest()
	waitReq.Offset = helper.Uint64(0)
	waitReq.Limit = helper.Uint64(1)
	waitReq.Filters = []*cynosdbv20190107.AuditInstanceFilters{
		{
			Name:       helper.String("InstanceId"),
			ExactMatch: helper.Bool(true),
			Values:     helper.Strings([]string{instanceId}),
		},
	}

	reqErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().DescribeAuditInstanceListWithContext(ctx, waitReq)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Items == nil || len(result.Response.Items) == 0 {
			return resource.NonRetryableError(fmt.Errorf("Describe cynosdb audit service failed, Response is nil."))
		}

		item := result.Response.Items[0]
		if item.AuditStatus == nil {
			return resource.NonRetryableError(fmt.Errorf("AuditStatus is nil."))
		}

		if *item.AuditStatus == "OFF" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("cynosdb audit service is still running, audit status is %s.", *item.AuditStatus))
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete cynosdb audit service failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
