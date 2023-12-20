package cfs

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs/v20190719"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudCfsAccessRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfsAccessRuleCreate,
		Read:   resourceTencentCloudCfsAccessRuleRead,
		Update: resourceTencentCloudCfsAccessRuleUpdate,
		Delete: resourceTencentCloudCfsAccessRuleDelete,

		Schema: map[string]*schema.Schema{
			"access_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of a access group.",
			},
			"auth_client_ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A single IP or a single IP address range such as 10.1.10.11 or 10.10.1.0/24 indicates that all IPs are allowed. Please note that the IP entered should be CVM's private IP.",
			},
			"priority": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 100),
				Description:  "The priority level of rule. Valid value ranges: (1~100). `1` indicates the highest priority.",
			},
			"rw_permission": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      CFS_RW_PERMISSION_RO,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CFS_RW_PERMISSION),
				Description:  "Read and write permissions. Valid values are `RO` and `RW`. and default is `RO`.",
			},
			"user_permission": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      CFS_USER_PERMISSION_ROOT_SQUASH,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CFS_USER_PERMISSION),
				Description:  "The permissions of accessing users. Valid values are `all_squash`, `no_all_squash`, `root_squash` and `no_root_squash`. and default is `root_squash`. `all_squash` indicates that all access users are mapped as anonymous users or user groups; `no_all_squash` indicates that access users will match local users first and be mapped to anonymous users or user groups after matching failed; `root_squash` indicates that map access root users to anonymous users or user groups; `no_root_squash` indicates that access root users keep root account permission.",
			},
		},
	}
}

func resourceTencentCloudCfsAccessRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfs_access_rule.create")()
	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := cfs.NewCreateCfsRuleRequest()
	request.PGroupId = helper.String(d.Get("access_group_id").(string))
	request.AuthClientIp = helper.String(d.Get("auth_client_ip").(string))
	request.Priority = helper.IntInt64(d.Get("priority").(int))
	request.RWPermission = helper.String(d.Get("rw_permission").(string))
	request.UserPermission = helper.String(d.Get("user_permission").(string))
	ruleId := ""
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfsClient().CreateCfsRule(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return tccommon.RetryError(err)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response.Response.RuleId == nil {
			return resource.NonRetryableError(fmt.Errorf("access rule id is nil"))
		}
		ruleId = *response.Response.RuleId
		return nil
	})
	if err != nil {
		return err
	}
	d.SetId(ruleId)

	return resourceTencentCloudCfsAccessRuleRead(d, meta)
}

func resourceTencentCloudCfsAccessRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfs_access_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	ruleId := d.Id()
	groupId := d.Get("access_group_id").(string)
	cfsService := CfsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	var accessRule *cfs.PGroupRuleInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		rules, errRet := cfsService.DescribeAccessRule(ctx, groupId, ruleId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if len(rules) > 0 {
			accessRule = rules[0]
		}
		return nil
	})
	if err != nil {
		return err
	}
	if accessRule == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("auth_client_ip", accessRule.AuthClientIp)
	_ = d.Set("user_permission", accessRule.UserPermission)
	_ = d.Set("priority", accessRule.Priority)
	if accessRule.RWPermission != nil {
		_ = d.Set("rw_permission", strings.ToUpper(*accessRule.RWPermission))
	}

	return nil
}

func resourceTencentCloudCfsAccessRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfs_access_rule.update")()
	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := cfs.NewUpdateCfsRuleRequest()
	request.RuleId = helper.String(d.Id())
	request.PGroupId = helper.String(d.Get("access_group_id").(string))
	if d.HasChange("auth_client_ip") {
		request.AuthClientIp = helper.String(d.Get("auth_client_ip").(string))
	}
	if d.HasChange("rw_permission") {
		request.RWPermission = helper.String(d.Get("rw_permission").(string))
	}
	if d.HasChange("user_permission") {
		request.UserPermission = helper.String(d.Get("user_permission").(string))
	}
	if d.HasChange("priority") {
		request.Priority = helper.IntInt64(d.Get("priority").(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfsClient().UpdateCfsRule(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return tccommon.RetryError(err)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		return nil
	})
	if err != nil {
		return err
	}

	return resourceTencentCloudCfsAccessRuleRead(d, meta)
}

func resourceTencentCloudCfsAccessRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfs_access_rule.delete")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cfsService := CfsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	ruleId := d.Id()
	groupId := d.Get("access_group_id").(string)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		errRet := cfsService.DeleteAccessRule(ctx, groupId, ruleId)
		if errRet != nil {
			return tccommon.RetryError(errRet)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
