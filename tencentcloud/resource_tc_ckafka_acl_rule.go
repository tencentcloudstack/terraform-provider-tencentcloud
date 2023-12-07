package tencentcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCkafkaAclRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCkafkaAclRuleCreate,
		Read:   resourceTencentCloudCkafkaAclRuleRead,
		Update: resourceTencentCloudCkafkaAclRuleUpdate,
		Delete: resourceTencentCloudCkafkaAclRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},

			"resource_type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Acl resource type, currently only supports Topic, enumeration value list{Topic}.",
			},

			"pattern_type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Match type, currently supports prefix matching and preset strategy, enumeration value list{PREFIXED/PRESET}.",
			},

			"rule_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "rule name.",
			},

			"rule_list": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "List of configured ACL rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"operation": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Acl operation mode, enumeration value (all operations All, read Read, write Write).",
						},
						"permission_type": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "permission type, (Deny|Allow).",
						},
						"host": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The default is *, which means that any host can be accessed. Currently, ckafka does not support host and ip network segment.",
						},
						"principal": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "User list, the default is User:, which means that any user can access, and the current user can only be the user included in the user list. The input format needs to be prefixed with [User:]. For example, user A is passed in as User:A.",
						},
					},
				},
			},

			"pattern": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "A value representing the prefix that the prefix matches.",
			},

			"is_applied": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether the preset ACL rule is applied to the newly added topic.",
			},
		},
	}
}

func resourceTencentCloudCkafkaAclRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ckafka_acl_rule.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = ckafka.NewCreateAclRuleRequest()
		instanceId string
		ruleName   string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	if v, ok := d.GetOk("resource_type"); ok {
		request.ResourceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("pattern_type"); ok {
		request.PatternType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rule_name"); ok {
		ruleName = v.(string)
		request.RuleName = helper.String(ruleName)
	}

	if v, ok := d.GetOk("rule_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			aclRuleInfo := ckafka.AclRuleInfo{}
			if v, ok := dMap["operation"]; ok {
				aclRuleInfo.Operation = helper.String(v.(string))
			}
			if v, ok := dMap["permission_type"]; ok {
				aclRuleInfo.PermissionType = helper.String(v.(string))
			}
			if v, ok := dMap["host"]; ok {
				aclRuleInfo.Host = helper.String(v.(string))
			}
			if v, ok := dMap["principal"]; ok {
				aclRuleInfo.Principal = helper.String(v.(string))
			}
			request.RuleList = append(request.RuleList, &aclRuleInfo)
		}
	}

	if v, ok := d.GetOk("pattern"); ok {
		request.Pattern = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("is_applied"); ok {
		request.IsApplied = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCkafkaClient().CreateAclRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ckafka aclRule failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId + FILED_SP + ruleName)

	return resourceTencentCloudCkafkaAclRuleRead(d, meta)
}

func resourceTencentCloudCkafkaAclRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ckafka_acl_rule.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CkafkaService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	ruleName := idSplit[1]

	aclRule, err := service.DescribeCkafkaAclRuleById(ctx, instanceId, ruleName)
	if err != nil {
		return err
	}

	if aclRule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CkafkaAclRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if aclRule.InstanceId != nil {
		_ = d.Set("instance_id", aclRule.InstanceId)
	}

	if aclRule.ResourceType != nil {
		_ = d.Set("resource_type", aclRule.ResourceType)
	}

	if aclRule.PatternType != nil {
		_ = d.Set("pattern_type", aclRule.PatternType)
	}

	if aclRule.RuleName != nil {
		_ = d.Set("rule_name", aclRule.RuleName)
	}

	if aclRule.AclList != nil {
		aclMapList := []interface{}{}
		aclList := make([]*ckafka.AclRuleInfo, 0)
		err := json.Unmarshal([]byte(*aclRule.AclList), &aclList)
		if err != nil {
			return err
		}
		for _, acl := range aclList {
			AclListMap := map[string]interface{}{}

			if acl.Operation != nil {
				AclListMap["operation"] = acl.Operation
			}

			if acl.PermissionType != nil {
				AclListMap["permission_type"] = acl.PermissionType
			}

			if acl.Host != nil {
				AclListMap["host"] = acl.Host
			}

			if acl.Principal != nil {
				AclListMap["principal"] = acl.Principal
			}

			aclMapList = append(aclMapList, AclListMap)
		}

		_ = d.Set("rule_list", aclMapList)

	}

	if aclRule.Pattern != nil {
		_ = d.Set("pattern", aclRule.Pattern)
	}

	if aclRule.IsApplied != nil {
		_ = d.Set("is_applied", aclRule.IsApplied)
	}

	return nil
}

func resourceTencentCloudCkafkaAclRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ckafka_acl_rule.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := ckafka.NewModifyAclRuleRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	ruleName := idSplit[1]

	request.InstanceId = &instanceId
	request.RuleName = &ruleName

	var hasChange bool

	if d.HasChange("is_applied") {
		if v, ok := d.GetOkExists("is_applied"); ok {
			request.IsApplied = helper.IntInt64(v.(int))
			hasChange = true
		}
	}

	if !hasChange {
		return resourceTencentCloudCkafkaAclRuleRead(d, meta)
	}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCkafkaClient().ModifyAclRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update ckafka aclRule failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCkafkaAclRuleRead(d, meta)
}

func resourceTencentCloudCkafkaAclRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ckafka_acl_rule.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CkafkaService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	ruleName := idSplit[1]

	if err := service.DeleteCkafkaAclRuleById(ctx, instanceId, ruleName); err != nil {
		return err
	}

	return nil
}
