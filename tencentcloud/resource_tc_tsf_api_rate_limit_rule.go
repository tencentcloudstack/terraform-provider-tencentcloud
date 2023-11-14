/*
Provides a resource to create a tsf api_rate_limit_rule

Example Usage

```hcl
resource "tencentcloud_tsf_api_rate_limit_rule" "api_rate_limit_rule" {
  api_id = ""
  max_qps =
  usable_status = ""
  }
```

Import

tsf api_rate_limit_rule can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_api_rate_limit_rule.api_rate_limit_rule api_rate_limit_rule_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudTsfApiRateLimitRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfApiRateLimitRuleCreate,
		Read:   resourceTencentCloudTsfApiRateLimitRuleRead,
		Update: resourceTencentCloudTsfApiRateLimitRuleUpdate,
		Delete: resourceTencentCloudTsfApiRateLimitRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"api_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Api Id.",
			},

			"max_qps": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Qps value.",
			},

			"usable_status": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Enable/disable, enabled/disabled, if not passed, it is enabled by default.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Current limiting result.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule Id.",
						},
						"api_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API ID.",
						},
						"rule_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Current limit name.",
						},
						"max_qps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum current limit qps.",
						},
						"usable_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Effective/disabled, enabled/disabled.",
						},
						"rule_content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule content.",
						},
						"tsf_rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tsf Rule ID.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Describe.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
						"updated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Update time.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTsfApiRateLimitRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_api_rate_limit_rule.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = tsf.NewCreateApiRateLimitRuleRequest()
		response = tsf.NewCreateApiRateLimitRuleResponse()
		apiId    string
	)
	if v, ok := d.GetOk("api_id"); ok {
		apiId = v.(string)
		request.ApiId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("max_qps"); ok {
		request.MaxQps = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("usable_status"); ok {
		request.UsableStatus = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().CreateApiRateLimitRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf apiRateLimitRule failed, reason:%+v", logId, err)
		return err
	}

	apiId = *response.Response.ApiId
	d.SetId(apiId)

	return resourceTencentCloudTsfApiRateLimitRuleRead(d, meta)
}

func resourceTencentCloudTsfApiRateLimitRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_api_rate_limit_rule.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	apiRateLimitRuleId := d.Id()

	apiRateLimitRule, err := service.DescribeTsfApiRateLimitRuleById(ctx, apiId)
	if err != nil {
		return err
	}

	if apiRateLimitRule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfApiRateLimitRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if apiRateLimitRule.ApiId != nil {
		_ = d.Set("api_id", apiRateLimitRule.ApiId)
	}

	if apiRateLimitRule.MaxQps != nil {
		_ = d.Set("max_qps", apiRateLimitRule.MaxQps)
	}

	if apiRateLimitRule.UsableStatus != nil {
		_ = d.Set("usable_status", apiRateLimitRule.UsableStatus)
	}

	if apiRateLimitRule.Result != nil {
		resultList := []interface{}{}
		for _, result := range apiRateLimitRule.Result {
			resultMap := map[string]interface{}{}

			if apiRateLimitRule.Result.RuleId != nil {
				resultMap["rule_id"] = apiRateLimitRule.Result.RuleId
			}

			if apiRateLimitRule.Result.ApiId != nil {
				resultMap["api_id"] = apiRateLimitRule.Result.ApiId
			}

			if apiRateLimitRule.Result.RuleName != nil {
				resultMap["rule_name"] = apiRateLimitRule.Result.RuleName
			}

			if apiRateLimitRule.Result.MaxQps != nil {
				resultMap["max_qps"] = apiRateLimitRule.Result.MaxQps
			}

			if apiRateLimitRule.Result.UsableStatus != nil {
				resultMap["usable_status"] = apiRateLimitRule.Result.UsableStatus
			}

			if apiRateLimitRule.Result.RuleContent != nil {
				resultMap["rule_content"] = apiRateLimitRule.Result.RuleContent
			}

			if apiRateLimitRule.Result.TsfRuleId != nil {
				resultMap["tsf_rule_id"] = apiRateLimitRule.Result.TsfRuleId
			}

			if apiRateLimitRule.Result.Description != nil {
				resultMap["description"] = apiRateLimitRule.Result.Description
			}

			if apiRateLimitRule.Result.CreatedTime != nil {
				resultMap["created_time"] = apiRateLimitRule.Result.CreatedTime
			}

			if apiRateLimitRule.Result.UpdatedTime != nil {
				resultMap["updated_time"] = apiRateLimitRule.Result.UpdatedTime
			}

			resultList = append(resultList, resultMap)
		}

		_ = d.Set("result", resultList)

	}

	return nil
}

func resourceTencentCloudTsfApiRateLimitRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_api_rate_limit_rule.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tsf.NewUpdateApiRateLimitRuleRequest()

	apiRateLimitRuleId := d.Id()

	request.ApiId = &apiId

	immutableArgs := []string{"api_id", "max_qps", "usable_status", "result"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("max_qps") {
		if v, ok := d.GetOkExists("max_qps"); ok {
			request.MaxQps = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("usable_status") {
		if v, ok := d.GetOk("usable_status"); ok {
			request.UsableStatus = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().UpdateApiRateLimitRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tsf apiRateLimitRule failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTsfApiRateLimitRuleRead(d, meta)
}

func resourceTencentCloudTsfApiRateLimitRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_api_rate_limit_rule.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	apiRateLimitRuleId := d.Id()

	if err := service.DeleteTsfApiRateLimitRuleById(ctx, apiId); err != nil {
		return err
	}

	return nil
}
