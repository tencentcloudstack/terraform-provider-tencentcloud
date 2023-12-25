package tcr

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTcrImmutableTagRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcrImmutableTagRuleCreate,
		Read:   resourceTencentCloudTcrImmutableTagRuleRead,
		Update: resourceTencentCloudTcrImmutableTagRuleUpdate,
		Delete: resourceTencentCloudTcrImmutableTagRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"registry_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},

			"namespace_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "namespace name.",
			},

			"rule": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "rule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"repository_pattern": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "repository matching rules.",
						},
						"tag_pattern": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "tag matching rules.",
						},
						"repository_decoration": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "repository decoration type:repoMatches or repoExcludes.",
						},
						"tag_decoration": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "tag decoration type: matches or excludes.",
						},
						"disabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "disable rule.",
						},
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "rule id.",
						},
						"ns_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "namespace name.",
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

func resourceTencentCloudTcrImmutableTagRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcr_immutable_tag_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request       = tcr.NewCreateImmutableTagRulesRequest()
		registryId    string
		namespaceName string
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service       = TCRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)
	if v, ok := d.GetOk("registry_id"); ok {
		request.RegistryId = helper.String(v.(string))
		registryId = v.(string)
	}

	if v, ok := d.GetOk("namespace_name"); ok {
		request.NamespaceName = helper.String(v.(string))
		namespaceName = v.(string)
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
		request.Rule = &immutableTagRule
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTCRClient().CreateImmutableTagRules(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tcr ImmutableTagRule failed, reason:%+v", logId, err)
		return err
	}

	ImmutableTagRules, err := service.DescribeTcrImmutableTagRuleById(ctx, registryId, &namespaceName, nil)
	if err != nil {
		return err
	}

	ruleId := helper.Int64ToStr(*ImmutableTagRules[0].RuleId)

	d.SetId(strings.Join([]string{registryId, namespaceName, ruleId}, tccommon.FILED_SP))

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := fmt.Sprintf("qcs::tcr:%s:uin/:instance/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudTcrImmutableTagRuleRead(d, meta)
}

func resourceTencentCloudTcrImmutableTagRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcr_immutable_tag_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	service := TCRService{client: tcClient}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	registryId := idSplit[0]
	namespaceName := idSplit[1]
	ruleId := idSplit[2]

	ImmutableTagRules, err := service.DescribeTcrImmutableTagRuleById(ctx, registryId, &namespaceName, &ruleId)
	if err != nil {
		return err
	}

	if ImmutableTagRules == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TcrImmutableTagRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	ImmutableTagRule := ImmutableTagRules[0]

	ruleMap := map[string]interface{}{}
	_ = d.Set("registry_id", registryId)
	_ = d.Set("namespace_name", namespaceName)

	if ImmutableTagRule.RepositoryPattern != nil {
		ruleMap["repository_pattern"] = ImmutableTagRule.RepositoryPattern
	}

	if ImmutableTagRule.TagPattern != nil {
		ruleMap["tag_pattern"] = ImmutableTagRule.TagPattern
	}

	if ImmutableTagRule.RepositoryDecoration != nil {
		ruleMap["repository_decoration"] = ImmutableTagRule.RepositoryDecoration
	}

	if ImmutableTagRule.TagDecoration != nil {
		ruleMap["tag_decoration"] = ImmutableTagRule.TagDecoration
	}

	if ImmutableTagRule.Disabled != nil {
		ruleMap["disabled"] = ImmutableTagRule.Disabled
	}

	if ImmutableTagRule.RuleId != nil {
		ruleMap["id"] = ImmutableTagRule.RuleId
	}

	if ImmutableTagRule.NsName != nil {
		ruleMap["ns_name"] = ImmutableTagRule.NsName
	}

	_ = d.Set("rule", []interface{}{ruleMap})

	tags, err := tagService.DescribeResourceTags(ctx, "tcr", "instance", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTcrImmutableTagRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcr_immutable_tag_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := tcr.NewModifyImmutableTagRulesRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	registryId := idSplit[0]
	namespaceName := idSplit[1]
	ruleId := idSplit[2]

	request.RegistryId = &registryId
	request.NamespaceName = &namespaceName
	request.RuleId = helper.StrToInt64Point(ruleId)

	immutableArgs := []string{"registry_id", "namespace_name"}

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
			request.Rule = &immutableTagRule
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTCRClient().ModifyImmutableTagRules(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tcr ImmutableTagRule failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := tccommon.BuildTagResourceName("tcr", "instance", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudTcrImmutableTagRuleRead(d, meta)
}

func resourceTencentCloudTcrImmutableTagRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcr_immutable_tag_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TCRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	registryId := idSplit[0]
	namespaceName := idSplit[1]
	ruleId := idSplit[2]

	if err := service.DeleteTcrImmutableTagRuleById(ctx, registryId, namespaceName, ruleId); err != nil {
		return err
	}

	return nil
}
