/*
Provides a resource to create a CFS access rule.

Example Usage

```hcl
resource "tencentcloud_cfs_access_rule" "foo" {
  access_group_id = "pgroup-7nx89k7l"
  auth_client_ip = "10.10.1.0/24"
  priority = 1
  rw_permission = "RO"
  user_permission = "root_squash"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs/v20190719"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func resourceTencentCloudCfsAccessRule() *schema.Resource {
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
				ValidateFunc: validateIntegerInRange(1, 100),
				Description:  "The priority level of rule. The range is 1-100, and 1 indicates the highest priority.",
			},
			"rw_permission": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      CFS_RW_PERMISSION_RO,
				ValidateFunc: validateAllowedStringValue(CFS_RW_PERMISSION),
				Description:  "Read and write permissions. Valid values are `RO` and `RW`, and default is `RO`.",
			},
			"user_permission": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      CFS_USER_PERMISSION_ROOT_SQUASH,
				ValidateFunc: validateAllowedStringValue(CFS_USER_PERMISSION),
				Description:  "The permissions of accessing users. Valid values are `all_squash`, `no_all_squash`, `root_squash` and `no_root_squash`, and default is `root_squash`. `all_squash` indicates that all access users are mapped as anonymous users or user groups; `no_all_squash` indicates that access users will match local users first and be mapped to anonymous users or user groups after matching failed; `root_squash` indicates that map access root users to anonymous users or user groups; `no_root_squash` indicates that access root users keep root account permission.",
			},
		},
	}
}

func resourceTencentCloudCfsAccessRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfs_access_rule.create")()
	logId := getLogId(contextNil)

	request := cfs.NewCreateCfsRuleRequest()
	request.PGroupId = stringToPointer(d.Get("access_group_id").(string))
	request.AuthClientIp = stringToPointer(d.Get("auth_client_ip").(string))
	request.Priority = int64ToPointer(d.Get("priority").(int))
	request.RWPermission = stringToPointer(d.Get("rw_permission").(string))
	request.UserPermission = stringToPointer(d.Get("user_permission").(string))
	ruleId := ""
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := meta.(*TencentCloudClient).apiV3Conn.UseCfsClient().CreateCfsRule(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return retryError(err)
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
	defer logElapsed("resource.tencentcloud_cfs_access_rule.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	ruleId := d.Id()
	groupId := d.Get("access_group_id").(string)
	cfsService := CfsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var accessRule *cfs.PGroupRuleInfo
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		rules, errRet := cfsService.DescribeAccessRule(ctx, groupId, ruleId)
		if errRet != nil {
			return retryError(errRet, "InternalError")
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

	d.Set("auth_client_ip", accessRule.AuthClientIp)
	d.Set("user_permission", accessRule.UserPermission)
	d.Set("priority", accessRule.Priority)
	if accessRule.RWPermission != nil {
		d.Set("rw_permission", strings.ToUpper(*accessRule.RWPermission))
	}

	return nil
}

func resourceTencentCloudCfsAccessRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfs_access_rule.update")()
	logId := getLogId(contextNil)

	request := cfs.NewUpdateCfsRuleRequest()
	request.RuleId = stringToPointer(d.Id())
	request.PGroupId = stringToPointer(d.Get("access_group_id").(string))
	if d.HasChange("auth_client_ip") {
		request.AuthClientIp = stringToPointer(d.Get("auth_client_ip").(string))
	}
	if d.HasChange("rw_permission") {
		request.RWPermission = stringToPointer(d.Get("rw_permission").(string))
	}
	if d.HasChange("user_permission") {
		request.UserPermission = stringToPointer(d.Get("user_permission").(string))
	}
	if d.HasChange("priority") {
		request.Priority = int64ToPointer(d.Get("priority").(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := meta.(*TencentCloudClient).apiV3Conn.UseCfsClient().UpdateCfsRule(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return retryError(err)
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
	defer logElapsed("resource.tencentcloud_cfs_access_rule.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	cfsService := CfsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	ruleId := d.Id()
	groupId := d.Get("access_group_id").(string)
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		errRet := cfsService.DeleteAccessRule(ctx, groupId, ruleId)
		if errRet != nil {
			return retryError(errRet)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
