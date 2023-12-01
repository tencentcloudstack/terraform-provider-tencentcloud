/*
Provides a resource to create a chdfs life_cycle_rule

Example Usage

```hcl
resource "tencentcloud_chdfs_life_cycle_rule" "life_cycle_rule" {
  file_system_id = "f14mpfy5lh4e"

  life_cycle_rule {
    life_cycle_rule_name = "terraform-test"
    path                 = "/test"
    status               = 1

    transitions {
      days = 30
      type = 1
    }
  }
}
```

Import

chdfs life_cycle_rule can be imported using the id, e.g.

```
terraform import tencentcloud_chdfs_life_cycle_rule.life_cycle_rule file_system_id#life_cycle_rule_id
```
*/
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
				ForceNew:    true,
				Description: "file system id.",
			},

			"life_cycle_rule": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "life cycle rule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"life_cycle_rule_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "single rule id.",
						},
						"life_cycle_rule_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "rule name.",
						},
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "rule op path.",
						},
						"transitions": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "life cycle rule transition list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "trigger days(n day).",
									},
									"type": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "transition type, 1: archive, 2: delete, 3: low rate.",
									},
								},
							},
						},
						"status": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "rule status, 1:open, 2:close.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "rule create time.",
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
		fileSystemId string
		path         string
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
			path = v.(string)
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
		request.LifeCycleRules = append(request.LifeCycleRules, &lifeCycleRule)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseChdfsClient().CreateLifeCycleRules(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		//response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create chdfs lifeCycleRule failed, reason:%+v", logId, err)
		return err
	}

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := ChdfsService{client: meta.(*TencentCloudClient).apiV3Conn}
	lifeCycleRule, err := service.DescribeChdfsLifeCycleRuleByPath(ctx, fileSystemId, path)
	if err != nil {
		return err
	}

	d.SetId(fileSystemId + FILED_SP + helper.UInt64ToStr(*lifeCycleRule.LifeCycleRuleId))

	return resourceTencentCloudChdfsLifeCycleRuleRead(d, meta)
}

func resourceTencentCloudChdfsLifeCycleRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_chdfs_life_cycle_rule.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ChdfsService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	fileSystemId := idSplit[0]
	lifeCycleRuleId := idSplit[1]

	lifeCycleRule, err := service.DescribeChdfsLifeCycleRuleById(ctx, fileSystemId, lifeCycleRuleId)
	if err != nil {
		return err
	}

	if lifeCycleRule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ChdfsLifeCycleRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("file_system_id", fileSystemId)

	if lifeCycleRule != nil {
		lifeCycleRuleMap := map[string]interface{}{}

		if lifeCycleRule.LifeCycleRuleId != nil {
			lifeCycleRuleMap["life_cycle_rule_id"] = lifeCycleRule.LifeCycleRuleId
		}

		if lifeCycleRule.LifeCycleRuleName != nil {
			lifeCycleRuleMap["life_cycle_rule_name"] = lifeCycleRule.LifeCycleRuleName
		}

		if lifeCycleRule.Path != nil {
			lifeCycleRuleMap["path"] = lifeCycleRule.Path
		}

		if lifeCycleRule.Transitions != nil {
			transitionsList := []interface{}{}
			for _, transitions := range lifeCycleRule.Transitions {
				transitionsMap := map[string]interface{}{}

				if transitions.Days != nil {
					transitionsMap["days"] = transitions.Days
				}

				if transitions.Type != nil {
					transitionsMap["type"] = transitions.Type
				}

				transitionsList = append(transitionsList, transitionsMap)
			}

			lifeCycleRuleMap["transitions"] = transitionsList
		}

		if lifeCycleRule.Status != nil {
			lifeCycleRuleMap["status"] = lifeCycleRule.Status
		}

		if lifeCycleRule.CreateTime != nil {
			lifeCycleRuleMap["create_time"] = lifeCycleRule.CreateTime
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

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	lifeCycleRuleId := idSplit[1]

	if d.HasChange("life_cycle_rule") {
		if dMap, ok := helper.InterfacesHeadMap(d, "life_cycle_rule"); ok {
			lifeCycleRule := chdfs.LifeCycleRule{}

			lifeCycleRule.LifeCycleRuleId = helper.StrToUint64Point(lifeCycleRuleId)
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
			request.LifeCycleRules = append(request.LifeCycleRules, &lifeCycleRule)
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
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	lifeCycleRuleId := idSplit[1]

	if err := service.DeleteChdfsLifeCycleRuleById(ctx, lifeCycleRuleId); err != nil {
		return err
	}

	return nil
}
