package tcr

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTcrTagRetentionRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcrTagRetentionRuleCreate,
		Read:   resourceTencentCloudTcrTagRetentionRuleRead,
		Update: resourceTencentCloudTcrTagRetentionRuleUpdate,
		Delete: resourceTencentCloudTcrTagRetentionRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"registry_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The main instance ID.",
			},

			"namespace_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The Name of the namespace.",
			},

			"retention_rule": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Retention Policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The supported policies are latestPushedK (retain the latest `k` pushed versions) and nDaysSinceLastPush (retain pushed versions within the last `n` days).",
						},
						"value": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "corresponding values for rule settings.",
						},
					},
				},
			},

			"advanced_rule_items": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				Description: "The advanced retention policy takes precedence; when both the basic and advanced retention policies are configured, the advanced retention policy will be used.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"retention_policy": {
							Optional:    true,
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Version retention rules.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Supported strategies, with possible values: latestPushedK (retain the latest K pushed versions), nDaysSinceLastPush (retain versions pushed within the last n days).",
									},
									"value": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Corresponding values under the rule settings.",
									},
								},
							},
						},
						"tag_filter": {
							Optional:    true,
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Tag filter.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"decoration": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Filter rule types: In tag filtering, the available options are matches (match) and excludes (exclude). In repository filtering, the available options are repoMatches (repository match) and repoExcludes (repository exclude).",
									},
									"pattern": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Filter expression.",
									},
								},
							},
						},
						"repository_filter": {
							Optional:    true,
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Warehouse filter.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"decoration": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Filter rule types: In tag filtering, the available options are matches (match) and excludes (exclude). In repository filtering, the available options are repoMatches (repository match) and repoExcludes (repository exclude).",
									},
									"pattern": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Filter expression.",
									},
								},
							},
						},
					},
				},
			},

			"cron_setting": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Execution cycle, currently only available selections are: manual; daily; weekly; monthly.",
			},

			"disabled": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to disable the rule, with the default value of false.",
			},

			// computed
			"retention_id": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The ID of the retention task.",
			},
		},
	}
}

func resourceTencentCloudTcrTagRetentionRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcr_tag_retention_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		tcrService    = TCRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request       = tcr.NewCreateTagRetentionRuleRequest()
		registryId    string
		namespaceName string
	)

	if v, ok := d.GetOk("registry_id"); ok {
		registryId = v.(string)
		request.RegistryId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("namespace_id"); ok {
		request.NamespaceId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("namespace_name"); ok {
		namespaceName = v.(string)
		namespace, has, err := tcrService.DescribeTCRNameSpaceById(ctx, registryId, namespaceName)
		if !has || namespace == nil {
			return fmt.Errorf("TCR namespace not found.")
		}

		if err != nil {
			return err
		}

		request.NamespaceId = namespace.NamespaceId
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "retention_rule"); ok {
		retentionRule := tcr.RetentionRule{}
		if v, ok := dMap["key"]; ok {
			retentionRule.Key = helper.String(v.(string))
		}

		if v, ok := dMap["value"]; ok {
			retentionRule.Value = helper.IntInt64(v.(int))
		}

		request.RetentionRule = &retentionRule
	}

	if v, ok := d.GetOk("advanced_rule_items"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			advancedRuleItem := tcr.RetentionRuleItem{}
			if v, ok := dMap["retention_policy"].([]interface{}); ok && len(v) > 0 {
				for _, item := range v {
					dMap := item.(map[string]interface{})
					retentionPolicy := tcr.RetentionRule{}
					if v, ok := dMap["key"]; ok {
						retentionPolicy.Key = helper.String(v.(string))
					}

					if v, ok := dMap["value"]; ok {
						retentionPolicy.Value = helper.IntInt64(v.(int))
					}

					advancedRuleItem.RetentionPolicy = &retentionPolicy
				}
			}

			if v, ok := dMap["tag_filter"].([]interface{}); ok && len(v) > 0 {
				for _, item := range v {
					dMap := item.(map[string]interface{})
					tagFilter := tcr.FilterSelector{}
					if v, ok := dMap["decoration"]; ok {
						tagFilter.Decoration = helper.String(v.(string))
					}

					if v, ok := dMap["pattern"]; ok {
						tagFilter.Pattern = helper.String(v.(string))
					}

					advancedRuleItem.TagFilter = &tagFilter
				}
			}

			if v, ok := dMap["repository_filter"].([]interface{}); ok && len(v) > 0 {
				for _, item := range v {
					dMap := item.(map[string]interface{})
					repositoryFilter := tcr.FilterSelector{}
					if v, ok := dMap["decoration"]; ok {
						repositoryFilter.Decoration = helper.String(v.(string))
					}

					if v, ok := dMap["pattern"]; ok {
						repositoryFilter.Pattern = helper.String(v.(string))
					}

					advancedRuleItem.RepositoryFilter = &repositoryFilter
				}
			}

			request.AdvancedRuleItems = append(request.AdvancedRuleItems, &advancedRuleItem)
		}
	}

	if v, ok := d.GetOk("cron_setting"); ok {
		request.CronSetting = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("disabled"); ok {
		request.Disabled = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTCRClient().CreateTagRetentionRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tcr TagRetentionRule failed, reason:%+v", logId, err)
		return err
	}

	TagRetentionRule, err := tcrService.DescribeTcrTagRetentionRuleById(ctx, registryId, namespaceName, nil)
	if err != nil {
		return fmt.Errorf("Query retention rule by id failed, reason:[%s]", err.Error())
	}

	if TagRetentionRule != nil {
		retentionId := helper.Int64ToStr(*TagRetentionRule.RetentionId)
		d.SetId(strings.Join([]string{registryId, namespaceName, retentionId}, tccommon.FILED_SP))
	} else {
		log.Printf("[CRITAL]%s TagRetentionRule is nil! Set unique id as empty.", logId)
		d.SetId("")
	}

	return resourceTencentCloudTcrTagRetentionRuleRead(d, meta)
}

func resourceTencentCloudTcrTagRetentionRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcr_tag_retention_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = TCRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	registryId := idSplit[0]
	namespaceName := idSplit[1]
	retentionId := idSplit[2]

	TagRetentionRule, err := service.DescribeTcrTagRetentionRuleById(ctx, registryId, namespaceName, &retentionId)
	if err != nil {
		return err
	}

	if TagRetentionRule == nil {
		log.Printf("[WARN]%s resource `tencentcloud_tcr_tag_retention_rule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("registry_id", registryId)

	if TagRetentionRule.RetentionId != nil {
		_ = d.Set("retention_id", TagRetentionRule.RetentionId)
	}

	if TagRetentionRule.NamespaceName != nil {
		_ = d.Set("namespace_name", TagRetentionRule.NamespaceName)
	}

	if len(TagRetentionRule.RetentionRuleList) > 0 {
		retentionRuleMap := map[string]interface{}{}
		retentionRule := TagRetentionRule.RetentionRuleList[0]
		if retentionRule.Key != nil {
			retentionRuleMap["key"] = retentionRule.Key
		}

		if retentionRule.Value != nil {
			retentionRuleMap["value"] = retentionRule.Value
		}

		_ = d.Set("retention_rule", []interface{}{retentionRuleMap})
	}

	if len(TagRetentionRule.AdvancedRuleItems) > 0 {
		advancedRuleItems := make([]map[string]interface{}, 0, len(TagRetentionRule.AdvancedRuleItems))
		for _, item := range TagRetentionRule.AdvancedRuleItems {
			advancedRuleItem := map[string]interface{}{}
			if item.RetentionPolicy != nil {
				retentionPolicyMap := map[string]interface{}{}
				if item.RetentionPolicy.Key != nil {
					retentionPolicyMap["key"] = item.RetentionPolicy.Key
				}

				if item.RetentionPolicy.Value != nil {
					retentionPolicyMap["value"] = item.RetentionPolicy.Value
				}

				advancedRuleItem["retention_policy"] = []interface{}{retentionPolicyMap}
			}

			if item.TagFilter != nil {
				tagFilterMap := map[string]interface{}{}
				if item.TagFilter.Decoration != nil {
					tagFilterMap["decoration"] = item.TagFilter.Decoration
				}

				if item.TagFilter.Pattern != nil {
					tagFilterMap["pattern"] = item.TagFilter.Pattern
				}

				advancedRuleItem["tag_filter"] = []interface{}{tagFilterMap}
			}

			if item.RepositoryFilter != nil {
				repositoryFilterMap := map[string]interface{}{}
				if item.RepositoryFilter.Decoration != nil {
					repositoryFilterMap["decoration"] = item.RepositoryFilter.Decoration
				}

				if item.RepositoryFilter.Pattern != nil {
					repositoryFilterMap["pattern"] = item.RepositoryFilter.Pattern
				}

				advancedRuleItem["repository_filter"] = []interface{}{repositoryFilterMap}
			}

			advancedRuleItems = append(advancedRuleItems, advancedRuleItem)
		}

		_ = d.Set("advanced_rule_items", advancedRuleItems)
	}

	if TagRetentionRule.CronSetting != nil {
		_ = d.Set("cron_setting", TagRetentionRule.CronSetting)
	}

	if TagRetentionRule.Disabled != nil {
		_ = d.Set("disabled", TagRetentionRule.Disabled)
	}

	return nil
}

func resourceTencentCloudTcrTagRetentionRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcr_tag_retention_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		tcrService = TCRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request    = tcr.NewModifyTagRetentionRuleRequest()
	)

	immutableArgs := []string{"registry_id", "namespace_name"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	registryId := idSplit[0]
	namespaceName := idSplit[1]
	retentionId := idSplit[2]

	namespace, has, err := tcrService.DescribeTCRNameSpaceById(ctx, registryId, namespaceName)
	if !has || namespace == nil {
		return fmt.Errorf("TCR namespace not found.")
	}

	if err != nil {
		return err
	}

	if v, ok := d.GetOkExists("cron_setting"); ok {
		request.CronSetting = helper.String(v.(string))
	}

	if d.HasChange("retention_rule") || d.HasChange("advanced_rule_items") || d.HasChange("disabled") {
		if dMap, ok := helper.InterfacesHeadMap(d, "retention_rule"); ok {
			retentionRule := tcr.RetentionRule{}
			if v, ok := dMap["key"]; ok {
				retentionRule.Key = helper.String(v.(string))
			}

			if v, ok := dMap["value"]; ok {
				retentionRule.Value = helper.IntInt64(v.(int))
			}

			request.RetentionRule = &retentionRule
		}

		if v, ok := d.GetOk("advanced_rule_items"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				advancedRuleItem := tcr.RetentionRuleItem{}
				if v, ok := dMap["retention_policy"].([]interface{}); ok && len(v) > 0 {
					for _, item := range v {
						dMap := item.(map[string]interface{})
						retentionPolicy := tcr.RetentionRule{}
						if v, ok := dMap["key"]; ok {
							retentionPolicy.Key = helper.String(v.(string))
						}

						if v, ok := dMap["value"]; ok {
							retentionPolicy.Value = helper.IntInt64(v.(int))
						}

						advancedRuleItem.RetentionPolicy = &retentionPolicy
					}
				}

				if v, ok := dMap["tag_filter"].([]interface{}); ok && len(v) > 0 {
					for _, item := range v {
						dMap := item.(map[string]interface{})
						tagFilter := tcr.FilterSelector{}
						if v, ok := dMap["decoration"]; ok {
							tagFilter.Decoration = helper.String(v.(string))
						}

						if v, ok := dMap["pattern"]; ok {
							tagFilter.Pattern = helper.String(v.(string))
						}

						advancedRuleItem.TagFilter = &tagFilter
					}
				}

				if v, ok := dMap["repository_filter"].([]interface{}); ok && len(v) > 0 {
					for _, item := range v {
						dMap := item.(map[string]interface{})
						repositoryFilter := tcr.FilterSelector{}
						if v, ok := dMap["decoration"]; ok {
							repositoryFilter.Decoration = helper.String(v.(string))
						}

						if v, ok := dMap["pattern"]; ok {
							repositoryFilter.Pattern = helper.String(v.(string))
						}

						advancedRuleItem.RepositoryFilter = &repositoryFilter
					}
				}

				request.AdvancedRuleItems = append(request.AdvancedRuleItems, &advancedRuleItem)
			}
		}

		if v, ok := d.GetOkExists("disabled"); ok {
			request.Disabled = helper.Bool(v.(bool))
		}
	}

	request.RegistryId = &registryId
	request.NamespaceId = namespace.NamespaceId
	request.RetentionId = helper.StrToInt64Point(retentionId)
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTCRClient().ModifyTagRetentionRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update tcr TagRetentionRule failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTcrTagRetentionRuleRead(d, meta)
}

func resourceTencentCloudTcrTagRetentionRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcr_tag_retention_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = TCRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	registryId := idSplit[0]
	retentionId := idSplit[2]

	if err := service.DeleteTcrTagRetentionRuleById(ctx, registryId, retentionId); err != nil {
		return err
	}

	return nil
}
