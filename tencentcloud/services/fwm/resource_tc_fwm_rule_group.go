package fwm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudFwmRuleGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudFwmRuleGroupCreate,
		Read:   resourceTencentCloudFwmRuleGroupRead,
		Update: resourceTencentCloudFwmRuleGroupUpdate,
		Delete: resourceTencentCloudFwmRuleGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule group name, cannot be empty, up to 50 characters.",
			},

			"product": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Product the rule group belongs to. Valid values: `enterprise_sg` (enterprise security group), `cfw_vpc_acl` (VPC border firewall), `cfw_nat_acl` (NAT border firewall), `cfw_edge_acl` (internet border firewall).",
			},

			// computed
			"group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule group ID.",
			},
		},
	}
}

func resourceTencentCloudFwmRuleGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_rule_group.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = NewFwmService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		request = fwmv20250611.NewCreateRuleGroupRequest()
		groupId string
	)

	if v, ok := d.GetOk("group_name"); ok {
		request.GroupName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		request.Product = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().CreateRuleGroupWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create tencentcloud_fwm_rule_group failed, Response is nil."))
		}

		if result.Response.GroupId == nil || *result.Response.GroupId == "" {
			return resource.NonRetryableError(fmt.Errorf("Create tencentcloud_fwm_rule_group failed, GroupId is empty."))
		}

		groupId = *result.Response.GroupId
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create tencentcloud_fwm_rule_group failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[DEBUG]%s create tencentcloud_fwm_rule_group success, groupId: %s", logId, groupId)
	d.SetId(groupId)

	pollErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ruleGroup, e := service.DescribeFwmRuleGroupById(ctx, groupId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if ruleGroup == nil {
			return resource.RetryableError(fmt.Errorf("tencentcloud_fwm_rule_group %s is not visible yet", groupId))
		}
		return nil
	})

	if pollErr != nil {
		log.Printf("[CRITAL]%s wait tencentcloud_fwm_rule_group %s visible failed, reason:%+v", logId, groupId, pollErr)
		return pollErr
	}

	return resourceTencentCloudFwmRuleGroupRead(d, meta)
}

func resourceTencentCloudFwmRuleGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_rule_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = NewFwmService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		groupId = d.Id()
	)

	respData, err := service.DescribeFwmRuleGroupById(ctx, groupId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_fwm_rule_group` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.GroupName != nil {
		_ = d.Set("group_name", respData.GroupName)
	}

	if respData.Product != nil {
		_ = d.Set("product", respData.Product)
	}

	if respData.GroupId != nil {
		_ = d.Set("group_id", respData.GroupId)
	}

	return nil
}

func resourceTencentCloudFwmRuleGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_rule_group.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	if d.HasChange("group_name") {
		request := fwmv20250611.NewModifyRuleGroupRequest()
		request.GroupId = helper.String(d.Id())

		if v, ok := d.GetOk("group_name"); ok {
			request.GroupName = helper.String(v.(string))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().ModifyRuleGroupWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify tencentcloud_fwm_rule_group failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update tencentcloud_fwm_rule_group failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudFwmRuleGroupRead(d, meta)
}

func resourceTencentCloudFwmRuleGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_rule_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = fwmv20250611.NewDeleteRuleGroupRequest()
		groupId = d.Id()
	)

	request.GroupIds = []*string{&groupId}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().DeleteRuleGroupWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete tencentcloud_fwm_rule_group failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete tencentcloud_fwm_rule_group failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
