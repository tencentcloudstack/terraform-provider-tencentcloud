/*
Provides a resource to create a tcr immutable_tag_rules

Example Usage

```hcl
resource "tencentcloud_tcr_immutable_tag_rules" "immutable_tag_rules" {
  registry_id = "tcr-xxx"
  namespace_name = "ns"
  rule {
		repository_pattern = "**"
		tag_pattern = "**"
		repository_decoration = "repoMatches"
		tag_decoration = "matches"
		disabled = false
		rule_id = 1
		ns_name = "ns"

  }
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

tcr immutable_tag_rules can be imported using the id, e.g.

```
terraform import tencentcloud_tcr_immutable_tag_rules.immutable_tag_rules immutable_tag_rules_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudTcrImmutableTagRules() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcrImmutableTagRulesCreate,
		Read:   resourceTencentCloudTcrImmutableTagRulesRead,
		Update: resourceTencentCloudTcrImmutableTagRulesUpdate,
		Delete: resourceTencentCloudTcrImmutableTagRulesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"registry_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"namespace_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Namespace name.",
			},

			"rule": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Rule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"repository_pattern": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Repository matching rules.",
						},
						"tag_pattern": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag matching rules.",
						},
						"repository_decoration": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Repository decoration type:repoMatches or repoExcludes.",
						},
						"tag_decoration": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag decoration type: matches or excludes.",
						},
						"disabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Disable rule.",
						},
						"rule_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Rule id.",
						},
						"ns_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Namespace name.",
						},
					},
				},
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudTcrImmutableTagRulesCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_immutable_tag_rules.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request       = tcr.NewCreateImmutableTagRulesRequest()
		response      = tcr.NewCreateImmutableTagRulesResponse()
		registryId    string
		namespaceName string
		ruleId        int
	)
	if v, ok := d.GetOk("registry_id"); ok {
		registryId = v.(string)
		request.RegistryId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace_name"); ok {
		namespaceName = v.(string)
		request.NamespaceName = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "rule"); ok {
		immutableTagRule := tcr.ImmutableTagRule{}
		if v, ok := dMap["repository_pattern"]; ok {
			immutableTagRule.RepositoryPattern = helper.String(v.(string))
		}
		if v, ok := dMap["tag_pattern"]; ok {
			immutableTagRule.TagPattern = helper.String(v.(string))
		}
		if v, ok := dMap["repository_decoration"]; ok {
			immutableTagRule.RepositoryDecoration = helper.String(v.(string))
		}
		if v, ok := dMap["tag_decoration"]; ok {
			immutableTagRule.TagDecoration = helper.String(v.(string))
		}
		if v, ok := dMap["disabled"]; ok {
			immutableTagRule.Disabled = helper.Bool(v.(bool))
		}
		if v, ok := dMap["rule_id"]; ok {
			immutableTagRule.RuleId = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["ns_name"]; ok {
			immutableTagRule.NsName = helper.String(v.(string))
		}
		request.Rule = &immutableTagRule
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTcrClient().CreateImmutableTagRules(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tcr ImmutableTagRules failed, reason:%+v", logId, err)
		return err
	}

	registryId = *response.Response.RegistryId
	d.SetId(strings.Join([]string{registryId, namespaceName, helper.Int64ToStr(ruleId)}, FILED_SP))

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::tcr:%s:uin/:instance/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudTcrImmutableTagRulesRead(d, meta)
}

func resourceTencentCloudTcrImmutableTagRulesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_immutable_tag_rules.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TcrService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	registryId := idSplit[0]
	namespaceName := idSplit[1]
	ruleId := idSplit[2]

	ImmutableTagRules, err := service.DescribeTcrImmutableTagRulesById(ctx, registryId, namespaceName, ruleId)
	if err != nil {
		return err
	}

	if ImmutableTagRules == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TcrImmutableTagRules` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if ImmutableTagRules.RegistryId != nil {
		_ = d.Set("registry_id", ImmutableTagRules.RegistryId)
	}

	if ImmutableTagRules.NamespaceName != nil {
		_ = d.Set("namespace_name", ImmutableTagRules.NamespaceName)
	}

	if ImmutableTagRules.Rule != nil {
		ruleMap := map[string]interface{}{}

		if ImmutableTagRules.Rule.RepositoryPattern != nil {
			ruleMap["repository_pattern"] = ImmutableTagRules.Rule.RepositoryPattern
		}

		if ImmutableTagRules.Rule.TagPattern != nil {
			ruleMap["tag_pattern"] = ImmutableTagRules.Rule.TagPattern
		}

		if ImmutableTagRules.Rule.RepositoryDecoration != nil {
			ruleMap["repository_decoration"] = ImmutableTagRules.Rule.RepositoryDecoration
		}

		if ImmutableTagRules.Rule.TagDecoration != nil {
			ruleMap["tag_decoration"] = ImmutableTagRules.Rule.TagDecoration
		}

		if ImmutableTagRules.Rule.Disabled != nil {
			ruleMap["disabled"] = ImmutableTagRules.Rule.Disabled
		}

		if ImmutableTagRules.Rule.RuleId != nil {
			ruleMap["rule_id"] = ImmutableTagRules.Rule.RuleId
		}

		if ImmutableTagRules.Rule.NsName != nil {
			ruleMap["ns_name"] = ImmutableTagRules.Rule.NsName
		}

		_ = d.Set("rule", []interface{}{ruleMap})
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "tcr", "instance", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTcrImmutableTagRulesUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_immutable_tag_rules.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tcr.NewModifyImmutableTagRulesRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	registryId := idSplit[0]
	namespaceName := idSplit[1]
	ruleId := idSplit[2]

	request.RegistryId = &registryId
	request.NamespaceName = &namespaceName
	request.RuleId = &ruleId

	immutableArgs := []string{"registry_id", "namespace_name", "rule"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("registry_id") {
		if v, ok := d.GetOk("registry_id"); ok {
			request.RegistryId = helper.String(v.(string))
		}
	}

	if d.HasChange("namespace_name") {
		if v, ok := d.GetOk("namespace_name"); ok {
			request.NamespaceName = helper.String(v.(string))
		}
	}

	if d.HasChange("rule") {
		if dMap, ok := helper.InterfacesHeadMap(d, "rule"); ok {
			immutableTagRule := tcr.ImmutableTagRule{}
			if v, ok := dMap["repository_pattern"]; ok {
				immutableTagRule.RepositoryPattern = helper.String(v.(string))
			}
			if v, ok := dMap["tag_pattern"]; ok {
				immutableTagRule.TagPattern = helper.String(v.(string))
			}
			if v, ok := dMap["repository_decoration"]; ok {
				immutableTagRule.RepositoryDecoration = helper.String(v.(string))
			}
			if v, ok := dMap["tag_decoration"]; ok {
				immutableTagRule.TagDecoration = helper.String(v.(string))
			}
			if v, ok := dMap["disabled"]; ok {
				immutableTagRule.Disabled = helper.Bool(v.(bool))
			}
			if v, ok := dMap["rule_id"]; ok {
				immutableTagRule.RuleId = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["ns_name"]; ok {
				immutableTagRule.NsName = helper.String(v.(string))
			}
			request.Rule = &immutableTagRule
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTcrClient().ModifyImmutableTagRules(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tcr ImmutableTagRules failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("tcr", "instance", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudTcrImmutableTagRulesRead(d, meta)
}

func resourceTencentCloudTcrImmutableTagRulesDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_immutable_tag_rules.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TcrService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	registryId := idSplit[0]
	namespaceName := idSplit[1]
	ruleId := idSplit[2]

	if err := service.DeleteTcrImmutableTagRulesById(ctx, registryId, namespaceName, ruleId); err != nil {
		return err
	}

	return nil
}
