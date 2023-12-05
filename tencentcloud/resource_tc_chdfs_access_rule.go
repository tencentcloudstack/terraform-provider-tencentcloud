package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	chdfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/chdfs/v20201112"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudChdfsAccessRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudChdfsAccessRuleCreate,
		Read:   resourceTencentCloudChdfsAccessRuleRead,
		Update: resourceTencentCloudChdfsAccessRuleUpdate,
		Delete: resourceTencentCloudChdfsAccessRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"access_rule": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "rule detail.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_rule_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "single rule id.",
						},
						"address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "rule address, IP OR IP SEG.",
						},
						"access_mode": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "rule access mode, 1: read only, 2: read &amp; wirte.",
						},
						"priority": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "rule priority, range 1 - 100, value less higher priority.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "rule create time.",
						},
					},
				},
			},

			"access_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "access group id.",
			},
		},
	}
}

func resourceTencentCloudChdfsAccessRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_chdfs_access_rule.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request       = chdfs.NewCreateAccessRulesRequest()
		response      = chdfs.NewCreateAccessRulesResponse()
		accessGroupId string
		accessRuleId  uint64
	)
	if dMap, ok := helper.InterfacesHeadMap(d, "access_rule"); ok {
		accessRule := chdfs.AccessRule{}
		if v, ok := dMap["address"]; ok {
			accessRule.Address = helper.String(v.(string))
		}
		if v, ok := dMap["access_mode"]; ok {
			accessRule.AccessMode = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["priority"]; ok {
			accessRule.Priority = helper.IntUint64(v.(int))
		}
		request.AccessRules = append(request.AccessRules, &accessRule)
	}

	if v, ok := d.GetOk("access_group_id"); ok {
		accessGroupId = v.(string)
		request.AccessGroupId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseChdfsClient().CreateAccessRules(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create chdfs accessRule failed, reason:%+v", logId, err)
		return err
	}

	if len(response.Response.AccessRules) < 1 {
		return fmt.Errorf("create chdfs accessRules failed")
	}

	accessRuleId = *response.Response.AccessRules[0].AccessRuleId
	d.SetId(accessGroupId + FILED_SP + helper.UInt64ToStr(accessRuleId))

	return resourceTencentCloudChdfsAccessRuleRead(d, meta)
}

func resourceTencentCloudChdfsAccessRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_chdfs_access_rule.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ChdfsService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	accessGroupId := idSplit[0]
	accessRuleId := idSplit[1]

	accessRule, err := service.DescribeChdfsAccessRulesById(ctx, accessGroupId, accessRuleId)
	if err != nil {
		return err
	}

	if accessRule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ChdfsAccessRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if accessRule != nil {
		accessRuleMap := map[string]interface{}{}

		if accessRule.AccessRuleId != nil {
			accessRuleMap["access_rule_id"] = accessRule.AccessRuleId
		}

		if accessRule.Address != nil {
			accessRuleMap["address"] = accessRule.Address
		}

		if accessRule.AccessMode != nil {
			accessRuleMap["access_mode"] = accessRule.AccessMode
		}

		if accessRule.Priority != nil {
			accessRuleMap["priority"] = accessRule.Priority
		}

		if accessRule.CreateTime != nil {
			accessRuleMap["create_time"] = accessRule.CreateTime
		}

		_ = d.Set("access_rule", []interface{}{accessRuleMap})

	}

	_ = d.Set("access_group_id", accessGroupId)
	return nil
}

func resourceTencentCloudChdfsAccessRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_chdfs_access_rule.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := chdfs.NewModifyAccessRulesRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	accessRuleId := idSplit[1]

	if d.HasChange("access_rule") {
		if dMap, ok := helper.InterfacesHeadMap(d, "access_rule"); ok {
			accessRule := chdfs.AccessRule{}

			accessRule.AccessRuleId = helper.StrToUint64Point(accessRuleId)

			if v, ok := dMap["address"]; ok {
				accessRule.Address = helper.String(v.(string))
			}
			if v, ok := dMap["access_mode"]; ok {
				accessRule.AccessMode = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["priority"]; ok {
				accessRule.Priority = helper.IntUint64(v.(int))
			}
			request.AccessRules = append(request.AccessRules, &accessRule)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseChdfsClient().ModifyAccessRules(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update chdfs accessRule failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudChdfsAccessRuleRead(d, meta)
}

func resourceTencentCloudChdfsAccessRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_chdfs_access_rule.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ChdfsService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	accessRuleId := idSplit[1]

	if err := service.DeleteChdfsAccessRulesById(ctx, accessRuleId); err != nil {
		return err
	}

	return nil
}
