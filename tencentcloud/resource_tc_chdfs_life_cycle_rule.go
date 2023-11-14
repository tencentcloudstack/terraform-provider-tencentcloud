/*
Provides a resource to create a chdfs life_cycle_rule

Example Usage

```hcl
resource "tencentcloud_chdfs_life_cycle_rule" "life_cycle_rule" {
  file_system_id = &lt;nil&gt;
  life_cycle_rule {
		life_cycle_rule_name = &lt;nil&gt;
		path = &lt;nil&gt;
		transitions {
			days = &lt;nil&gt;
			type = &lt;nil&gt;
		}
		status =

  }
}
```

Import

chdfs life_cycle_rule can be imported using the id, e.g.

```
terraform import tencentcloud_chdfs_life_cycle_rule.life_cycle_rule life_cycle_rule_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	chdfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/chdfs/v20201112"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudChdfsLifeCycleRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudChdfsLifeCycleRuleCreate,
		Read:   resourceTencentCloudChdfsLifeCycleRuleRead,
		Update: resourceTencentCloudChdfsLifeCycleRuleUpdate,
		Delete: resourceTencentCloudChdfsLifeCycleRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"file_system_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "File system id.",
			},

			"life_cycle_rule": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Life cycle rule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"life_cycle_rule_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Single rule id.",
						},
						"life_cycle_rule_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Rule name.",
						},
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Rule op path.",
						},
						"transitions": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Life cycle rule transition list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Trigger days(n day).",
									},
									"type": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Transition type, 1: archive, 2: delete, 3: low rate.",
									},
								},
							},
						},
						"status": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Rule status, 1:open, 2:close.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule create time.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudChdfsLifeCycleRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_chdfs_life_cycle_rule.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = chdfs.NewCreateLifeCycleRulesRequest()
		response     = chdfs.NewCreateLifeCycleRulesResponse()
		fileSystemId string
	)
	if v, ok := d.GetOk("file_system_id"); ok {
		fileSystemId = v.(string)
		request.FileSystemId = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "life_cycle_rule"); ok {
		lifeCycleRule := chdfs.LifeCycleRule{}
		if v, ok := dMap["life_cycle_rule_name"]; ok {
			lifeCycleRule.LifeCycleRuleName = helper.String(v.(string))
		}
		if v, ok := dMap["path"]; ok {
			lifeCycleRule.Path = helper.String(v.(string))
		}
		if v, ok := dMap["transitions"]; ok {
			for _, item := range v.([]interface{}) {
				transitionsMap := item.(map[string]interface{})
				transition := chdfs.Transition{}
				if v, ok := transitionsMap["days"]; ok {
					transition.Days = helper.IntUint64(v.(int))
				}
				if v, ok := transitionsMap["type"]; ok {
					transition.Type = helper.IntUint64(v.(int))
				}
				lifeCycleRule.Transitions = append(lifeCycleRule.Transitions, &transition)
			}
		}
		if v, ok := dMap["status"]; ok {
			lifeCycleRule.Status = helper.IntUint64(v.(int))
		}
		request.LifeCycleRule = &lifeCycleRule
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseChdfsClient().CreateLifeCycleRules(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create chdfs lifeCycleRule failed, reason:%+v", logId, err)
		return err
	}

	fileSystemId = *response.Response.FileSystemId
	d.SetId(fileSystemId)

	return resourceTencentCloudChdfsLifeCycleRuleRead(d, meta)
}

func resourceTencentCloudChdfsLifeCycleRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_chdfs_life_cycle_rule.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ChdfsService{client: meta.(*TencentCloudClient).apiV3Conn}

	lifeCycleRuleId := d.Id()

	lifeCycleRule, err := service.DescribeChdfsLifeCycleRuleById(ctx, fileSystemId)
	if err != nil {
		return err
	}

	if lifeCycleRule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ChdfsLifeCycleRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if lifeCycleRule.FileSystemId != nil {
		_ = d.Set("file_system_id", lifeCycleRule.FileSystemId)
	}

	if lifeCycleRule.LifeCycleRule != nil {
		lifeCycleRuleMap := map[string]interface{}{}

		if lifeCycleRule.LifeCycleRule.LifeCycleRuleId != nil {
			lifeCycleRuleMap["life_cycle_rule_id"] = lifeCycleRule.LifeCycleRule.LifeCycleRuleId
		}

		if lifeCycleRule.LifeCycleRule.LifeCycleRuleName != nil {
			lifeCycleRuleMap["life_cycle_rule_name"] = lifeCycleRule.LifeCycleRule.LifeCycleRuleName
		}

		if lifeCycleRule.LifeCycleRule.Path != nil {
			lifeCycleRuleMap["path"] = lifeCycleRule.LifeCycleRule.Path
		}

		if lifeCycleRule.LifeCycleRule.Transitions != nil {
			transitionsList := []interface{}{}
			for _, transitions := range lifeCycleRule.LifeCycleRule.Transitions {
				transitionsMap := map[string]interface{}{}

				if transitions.Days != nil {
					transitionsMap["days"] = transitions.Days
				}

				if transitions.Type != nil {
					transitionsMap["type"] = transitions.Type
				}

				transitionsList = append(transitionsList, transitionsMap)
			}

			lifeCycleRuleMap["transitions"] = []interface{}{transitionsList}
		}

		if lifeCycleRule.LifeCycleRule.Status != nil {
			lifeCycleRuleMap["status"] = lifeCycleRule.LifeCycleRule.Status
		}

		if lifeCycleRule.LifeCycleRule.CreateTime != nil {
			lifeCycleRuleMap["create_time"] = lifeCycleRule.LifeCycleRule.CreateTime
		}

		_ = d.Set("life_cycle_rule", []interface{}{lifeCycleRuleMap})
	}

	return nil
}

func resourceTencentCloudChdfsLifeCycleRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_chdfs_life_cycle_rule.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := chdfs.NewModifyLifeCycleRulesRequest()

	lifeCycleRuleId := d.Id()

	request.FileSystemId = &fileSystemId

	immutableArgs := []string{"file_system_id", "life_cycle_rule"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("life_cycle_rule") {
		if dMap, ok := helper.InterfacesHeadMap(d, "life_cycle_rule"); ok {
			lifeCycleRule := chdfs.LifeCycleRule{}
			if v, ok := dMap["life_cycle_rule_name"]; ok {
				lifeCycleRule.LifeCycleRuleName = helper.String(v.(string))
			}
			if v, ok := dMap["path"]; ok {
				lifeCycleRule.Path = helper.String(v.(string))
			}
			if v, ok := dMap["transitions"]; ok {
				for _, item := range v.([]interface{}) {
					transitionsMap := item.(map[string]interface{})
					transition := chdfs.Transition{}
					if v, ok := transitionsMap["days"]; ok {
						transition.Days = helper.IntUint64(v.(int))
					}
					if v, ok := transitionsMap["type"]; ok {
						transition.Type = helper.IntUint64(v.(int))
					}
					lifeCycleRule.Transitions = append(lifeCycleRule.Transitions, &transition)
				}
			}
			if v, ok := dMap["status"]; ok {
				lifeCycleRule.Status = helper.IntUint64(v.(int))
			}
			request.LifeCycleRule = &lifeCycleRule
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseChdfsClient().ModifyLifeCycleRules(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update chdfs lifeCycleRule failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudChdfsLifeCycleRuleRead(d, meta)
}

func resourceTencentCloudChdfsLifeCycleRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_chdfs_life_cycle_rule.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ChdfsService{client: meta.(*TencentCloudClient).apiV3Conn}
	lifeCycleRuleId := d.Id()

	if err := service.DeleteChdfsLifeCycleRuleById(ctx, fileSystemId); err != nil {
		return err
	}

	return nil
}
