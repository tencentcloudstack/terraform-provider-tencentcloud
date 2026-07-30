package ga2

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	ga2v20250115 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudGa2GlobalAcceleratorAclRuleSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGa2GlobalAcceleratorAclRuleSetCreate,
		Read:   resourceTencentCloudGa2GlobalAcceleratorAclRuleSetRead,
		Update: resourceTencentCloudGa2GlobalAcceleratorAclRuleSetUpdate,
		Delete: resourceTencentCloudGa2GlobalAcceleratorAclRuleSetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"global_accelerator_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Global accelerator instance ID.",
			},
			"global_accelerator_acl_policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ACL policy ID that owns the rule set.",
			},
			"acl_entries": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 0,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Protocol. Valid values: `TCP`, `UDP`, `ALL`.",
						},
						"port": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Port.",
						},
						"source_cidr_block": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Source CIDR block.",
						},
						"policy": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Action. Valid values: `ACCEPT` (allow), `DROP` (deny).",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description. Maximum length is 100 bytes.",
						},
						"global_accelerator_acl_rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ACL rule ID.",
						},
					},
				},
				Description: "The desired full set of ACL rules under the policy.",
			},
			"task_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Async task ID from the latest write operation.",
			},
		},
	}
}

func resourceTencentCloudGa2GlobalAcceleratorAclRuleSetCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_global_accelerator_acl_rule_set.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service  = Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		gaId     string
		policyId string
	)

	if v, ok := d.GetOk("global_accelerator_id"); ok {
		gaId = v.(string)
	}

	if v, ok := d.GetOk("global_accelerator_acl_policy_id"); ok {
		policyId = v.(string)
	}

	entries := buildAclEntriesFromSchema(d)
	if len(entries) == 0 {
		d.SetId(strings.Join([]string{gaId, policyId}, tccommon.FILED_SP))
		return resourceTencentCloudGa2GlobalAcceleratorAclRuleSetRead(d, meta)
	}

	request := ga2v20250115.NewCreateGlobalAcceleratorAclRuleRequest()
	request.GlobalAcceleratorId = helper.String(gaId)
	request.GlobalAcceleratorAclPolicyId = helper.String(policyId)
	request.AclEntries = entries

	var (
		response       = ga2v20250115.NewCreateGlobalAcceleratorAclRuleResponse()
		createdRuleIds []string
	)

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().CreateGlobalAcceleratorAclRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create ga2 global_accelerator_acl_rule_set failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create ga2 global_accelerator_acl_rule_set failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[DEBUG]%s create ga2 global_accelerator_acl_rule_set, d.Id()=%s", logId, d.Id())

	if response.Response.GlobalAcceleratorAclRuleIds == nil || len(response.Response.GlobalAcceleratorAclRuleIds) == 0 {
		return fmt.Errorf("Create ga2 global_accelerator_acl_rule_set failed, GlobalAcceleratorAclRuleIds is nil or empty.")
	}

	if len(response.Response.GlobalAcceleratorAclRuleIds) != len(entries) {
		return fmt.Errorf("Create ga2 global_accelerator_acl_rule_set failed, the count of GlobalAcceleratorAclRuleIds [%d] does not match the count of input AclEntries [%d].", len(response.Response.GlobalAcceleratorAclRuleIds), len(entries))
	}

	for _, idPtr := range response.Response.GlobalAcceleratorAclRuleIds {
		if idPtr == nil {
			return fmt.Errorf("Create ga2 global_accelerator_acl_rule_set failed, one of GlobalAcceleratorAclRuleIds is nil.")
		}
		createdRuleIds = append(createdRuleIds, *idPtr)
	}

	if response.Response.TaskId == nil {
		return fmt.Errorf("Create ga2 global_accelerator_acl_rule_set failed, TaskId is nil.")
	}

	taskId := *response.Response.TaskId
	if err := service.WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return err
	}

	// Map the returned rule ids back onto the acl_entries list in input order.
	rawEntries := d.Get("acl_entries").([]interface{})
	for i := range rawEntries {
		if i < len(createdRuleIds) {
			if entry, ok := rawEntries[i].(map[string]interface{}); ok {
				entry["global_accelerator_acl_rule_id"] = createdRuleIds[i]
			}
		}
	}
	_ = d.Set("acl_entries", rawEntries)
	_ = d.Set("task_id", taskId)

	d.SetId(strings.Join([]string{gaId, policyId}, tccommon.FILED_SP))

	return resourceTencentCloudGa2GlobalAcceleratorAclRuleSetRead(d, meta)
}

func resourceTencentCloudGa2GlobalAcceleratorAclRuleSetRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_global_accelerator_acl_rule_set.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	gaId, policyId, err := parseGa2GlobalAcceleratorAclRuleSetId(d.Id())
	if err != nil {
		return err
	}

	respSet, err := service.DescribeGa2GlobalAcceleratorAclRulesByPolicyId(ctx, policyId)
	if err != nil {
		return err
	}

	// Preserve context for diagnosis before any empty-result handling.
	log.Printf("[CRUD] ga2 global_accelerator_acl_rule_set id=%s", d.Id())

	_ = d.Set("global_accelerator_id", gaId)
	_ = d.Set("global_accelerator_acl_policy_id", policyId)

	if len(respSet) == 0 {
		// An empty rule set is a valid desired state; keep the resource in state.
		_ = d.Set("acl_entries", []interface{}{})
		return nil
	}

	// Sort by rule id for deterministic state.
	sort.Slice(respSet, func(i, j int) bool {
		var a, b string
		if respSet[i] != nil && respSet[i].GlobalAcceleratorAclRuleId != nil {
			a = *respSet[i].GlobalAcceleratorAclRuleId
		}
		if respSet[j] != nil && respSet[j].GlobalAcceleratorAclRuleId != nil {
			b = *respSet[j].GlobalAcceleratorAclRuleId
		}
		return a < b
	})

	_ = d.Set("acl_entries", flattenAclRuleSetToSchema(respSet))

	return nil
}

func resourceTencentCloudGa2GlobalAcceleratorAclRuleSetUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_global_accelerator_acl_rule_set.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	gaId, policyId, err := parseGa2GlobalAcceleratorAclRuleSetId(d.Id())
	if err != nil {
		return err
	}

	// global_accelerator_id and global_accelerator_acl_policy_id are ForceNew,
	// so changing them never reaches Update.
	if d.HasChange("acl_entries") {
		oldVal, newVal := d.GetChange("acl_entries")
		oldEntries := expandAclEntryMaps(oldVal.([]interface{}))
		newEntries := expandAclEntryMaps(newVal.([]interface{}))

		oldByKey := indexAclEntriesByKey(oldEntries)

		var (
			toCreate   []map[string]interface{}
			toRemove   []string
			toModify   []map[string]interface{}
			lastTaskId string
		)

		// Detect new and changed entries keyed by global_accelerator_acl_rule_id.
		newSeen := make(map[string]bool)
		for _, e := range newEntries {
			id := aclEntryRuleId(e)
			if id == "" {
				// No server-assigned id yet: treat as a brand-new entry to create.
				toCreate = append(toCreate, e)
				continue
			}
			newSeen[id] = true
			if old, ok := oldByKey[id]; ok {
				if !aclEntryEqual(old, e) {
					toModify = append(toModify, e)
				}
			} else {
				toCreate = append(toCreate, e)
			}
		}

		for _, e := range oldEntries {
			id := aclEntryRuleId(e)
			if id == "" {
				continue
			}
			if !newSeen[id] {
				toRemove = append(toRemove, id)
			}
		}

		// (a) batch-create new entries.
		if len(toCreate) > 0 {
			aclEntries := make([]*ga2v20250115.AclEntries, 0, len(toCreate))
			for _, e := range toCreate {
				aclEntries = append(aclEntries, buildAclEntryFromMap(e))
			}

			request := ga2v20250115.NewCreateGlobalAcceleratorAclRuleRequest()
			request.GlobalAcceleratorId = helper.String(gaId)
			request.GlobalAcceleratorAclPolicyId = helper.String(policyId)
			request.AclEntries = aclEntries

			var createdRuleIds []string
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().CreateGlobalAcceleratorAclRuleWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				if result == nil || result.Response == nil {
					return resource.NonRetryableError(fmt.Errorf("Create ga2 global_accelerator_acl_rule_set failed, Response is nil."))
				}

				if result.Response.GlobalAcceleratorAclRuleIds == nil || len(result.Response.GlobalAcceleratorAclRuleIds) != len(toCreate) {
					return resource.NonRetryableError(fmt.Errorf("Create ga2 global_accelerator_acl_rule_set failed, the count of GlobalAcceleratorAclRuleIds does not match the input."))
				}

				for _, idPtr := range result.Response.GlobalAcceleratorAclRuleIds {
					if idPtr == nil {
						return resource.NonRetryableError(fmt.Errorf("Create ga2 global_accelerator_acl_rule_set failed, one of GlobalAcceleratorAclRuleIds is nil."))
					}
					createdRuleIds = append(createdRuleIds, *idPtr)
				}

				if result.Response.TaskId == nil {
					return resource.NonRetryableError(fmt.Errorf("Create ga2 global_accelerator_acl_rule_set failed, TaskId is nil."))
				}

				lastTaskId = *result.Response.TaskId
				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s update ga2 global_accelerator_acl_rule_set (create) failed, reason:%+v", logId, reqErr)
				return reqErr
			}

			if err := service.WaitForGa2TaskFinish(ctx, lastTaskId, d.Timeout(schema.TimeoutUpdate)); err != nil {
				return err
			}

			// Map created ids back onto the toCreate entries in input order.
			for i := range toCreate {
				if i < len(createdRuleIds) {
					toCreate[i]["global_accelerator_acl_rule_id"] = createdRuleIds[i]
				}
			}
		}

		// (b) batch-delete removed entries.
		if len(toRemove) > 0 {
			ruleIds := make([]*string, 0, len(toRemove))
			for i := range toRemove {
				ruleIds = append(ruleIds, helper.String(toRemove[i]))
			}

			request := ga2v20250115.NewDeleteGlobalAcceleratorAclRuleRequest()
			request.GlobalAcceleratorId = helper.String(gaId)
			request.GlobalAcceleratorAclPolicyId = helper.String(policyId)
			request.GlobalAcceleratorAclRuleIds = ruleIds

			deleted := false
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().DeleteGlobalAcceleratorAclRuleWithContext(ctx, request)
				if e != nil {
					if sdkerr, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
						if sdkerr.Code == "ResourceNotFound" {
							// Rules already absent on the cloud side; nothing to wait for.
							return nil
						}
					}
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				if result == nil || result.Response == nil || result.Response.TaskId == nil {
					return resource.NonRetryableError(fmt.Errorf("Delete ga2 global_accelerator_acl_rule_set failed, Response is nil."))
				}

				lastTaskId = *result.Response.TaskId
				deleted = true
				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s update ga2 global_accelerator_acl_rule_set (delete) failed, reason:%+v", logId, reqErr)
				return reqErr
			}

			if deleted {
				if err := service.WaitForGa2TaskFinish(ctx, lastTaskId, d.Timeout(schema.TimeoutUpdate)); err != nil {
					return err
				}
			}
		}

		// (c) modify changed entries one-by-one.
		for _, e := range toModify {
			modRequest := ga2v20250115.NewModifyGlobalAcceleratorAclRuleRequest()
			modRequest.GlobalAcceleratorId = helper.String(gaId)
			modRequest.GlobalAcceleratorAclPolicyId = helper.String(policyId)
			modRequest.GlobalAcceleratorAclRuleId = helper.String(aclEntryRuleId(e))
			modRequest.Protocol = helper.String(aclEntryStr(e, "protocol"))
			modRequest.Port = helper.String(aclEntryStr(e, "port"))
			modRequest.SourceCidrBlock = helper.String(aclEntryStr(e, "source_cidr_block"))
			modRequest.Policy = helper.String(aclEntryStr(e, "policy"))
			modRequest.Description = helper.String(aclEntryStr(e, "description"))

			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().ModifyGlobalAcceleratorAclRuleWithContext(ctx, modRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modRequest.GetAction(), modRequest.ToJsonString(), result.ToJsonString())
				}

				if result == nil || result.Response == nil || result.Response.TaskId == nil {
					return resource.NonRetryableError(fmt.Errorf("Modify ga2 global_accelerator_acl_rule_set failed, Response is nil."))
				}

				lastTaskId = *result.Response.TaskId
				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s update ga2 global_accelerator_acl_rule_set (modify) failed, reason:%+v", logId, reqErr)
				return reqErr
			}

			if err := service.WaitForGa2TaskFinish(ctx, lastTaskId, d.Timeout(schema.TimeoutUpdate)); err != nil {
				return err
			}
		}

		if lastTaskId != "" {
			_ = d.Set("task_id", lastTaskId)
		}
	}

	return resourceTencentCloudGa2GlobalAcceleratorAclRuleSetRead(d, meta)
}

func resourceTencentCloudGa2GlobalAcceleratorAclRuleSetDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_global_accelerator_acl_rule_set.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	gaId, policyId, err := parseGa2GlobalAcceleratorAclRuleSetId(d.Id())
	if err != nil {
		return err
	}

	entries := d.Get("acl_entries").([]interface{})
	ruleIds := make([]*string, 0, len(entries))
	for _, raw := range entries {
		entry, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if v, ok := entry["global_accelerator_acl_rule_id"].(string); ok && v != "" {
			ruleIds = append(ruleIds, helper.String(v))
		}
	}

	if len(ruleIds) == 0 {
		// Nothing to delete on the cloud side.
		return nil
	}

	request := ga2v20250115.NewDeleteGlobalAcceleratorAclRuleRequest()
	request.GlobalAcceleratorId = helper.String(gaId)
	request.GlobalAcceleratorAclPolicyId = helper.String(policyId)
	request.GlobalAcceleratorAclRuleIds = ruleIds

	var taskId string
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().DeleteGlobalAcceleratorAclRuleWithContext(ctx, request)
		if e != nil {
			if sdkerr, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkerr.Code == "ResourceNotFound" {
					return nil
				}
			}
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.TaskId == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete ga2 global_accelerator_acl_rule_set failed, Response is nil."))
		}

		taskId = *result.Response.TaskId
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete ga2 global_accelerator_acl_rule_set failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if taskId != "" {
		if err := service.WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutDelete)); err != nil {
			return err
		}
	}

	return nil
}

func parseGa2GlobalAcceleratorAclRuleSetId(id string) (gaId, policyId string, err error) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 2 {
		err = fmt.Errorf("invalid resource id %q, expected format <global_accelerator_id>%s<global_accelerator_acl_policy_id>", id, tccommon.FILED_SP)
		return
	}
	gaId, policyId = parts[0], parts[1]
	if gaId == "" || policyId == "" {
		err = fmt.Errorf("invalid resource id %q, components must all be non-empty", id)
		return
	}
	return
}

// buildAclEntriesFromSchema reads acl_entries from the ResourceData and builds an SDK AclEntries slice.
func buildAclEntriesFromSchema(d *schema.ResourceData) []*ga2v20250115.AclEntries {
	raw := d.Get("acl_entries").([]interface{})
	result := make([]*ga2v20250115.AclEntries, 0, len(raw))
	for _, r := range raw {
		entry, ok := r.(map[string]interface{})
		if !ok {
			continue
		}
		result = append(result, buildAclEntryFromMap(entry))
	}
	return result
}

// buildAclEntryFromMap builds a single SDK AclEntries element from a schema map.
func buildAclEntryFromMap(entry map[string]interface{}) *ga2v20250115.AclEntries {
	return &ga2v20250115.AclEntries{
		Protocol:        helper.String(aclEntryStr(entry, "protocol")),
		Port:            helper.String(aclEntryStr(entry, "port")),
		SourceCidrBlock: helper.String(aclEntryStr(entry, "source_cidr_block")),
		Policy:          helper.String(aclEntryStr(entry, "policy")),
		Description:     helper.String(aclEntryStr(entry, "description")),
	}
}

// flattenAclRuleSetToSchema converts the API response set into the acl_entries list value.
func flattenAclRuleSetToSchema(set []*ga2v20250115.GlobalAcceleratorAclRuleSet) []interface{} {
	result := make([]interface{}, 0, len(set))
	for _, item := range set {
		if item == nil {
			continue
		}
		entry := map[string]interface{}{}
		if item.GlobalAcceleratorAclRuleId != nil {
			entry["global_accelerator_acl_rule_id"] = *item.GlobalAcceleratorAclRuleId
		}
		if item.Protocol != nil {
			entry["protocol"] = *item.Protocol
		}
		if item.Port != nil {
			entry["port"] = *item.Port
		}
		if item.SourceCidrBlock != nil {
			entry["source_cidr_block"] = *item.SourceCidrBlock
		}
		if item.Policy != nil {
			entry["policy"] = *item.Policy
		}
		if item.Description != nil {
			entry["description"] = *item.Description
		}
		result = append(result, entry)
	}
	return result
}

// expandAclEntryMaps converts the raw schema list into a slice of maps for diffing.
func expandAclEntryMaps(raw []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(raw))
	for _, r := range raw {
		entry, ok := r.(map[string]interface{})
		if !ok {
			continue
		}
		result = append(result, entry)
	}
	return result
}

// indexAclEntriesByKey indexes the entry maps by global_accelerator_acl_rule_id (empty ids are skipped).
func indexAclEntriesByKey(entries []map[string]interface{}) map[string]map[string]interface{} {
	result := make(map[string]map[string]interface{})
	for _, e := range entries {
		id := aclEntryRuleId(e)
		if id == "" {
			continue
		}
		result[id] = e
	}
	return result
}

// aclEntryRuleId returns the global_accelerator_acl_rule_id value of an entry map.
func aclEntryRuleId(entry map[string]interface{}) string {
	if v, ok := entry["global_accelerator_acl_rule_id"].(string); ok {
		return v
	}
	return ""
}

// aclEntryStr safely returns a string field of an entry map.
func aclEntryStr(entry map[string]interface{}, key string) string {
	if v, ok := entry[key].(string); ok {
		return v
	}
	return ""
}

// aclEntryEqual compares two entry maps by their editable fields (ignoring the computed rule id).
func aclEntryEqual(a, b map[string]interface{}) bool {
	return aclEntryStr(a, "protocol") == aclEntryStr(b, "protocol") &&
		aclEntryStr(a, "port") == aclEntryStr(b, "port") &&
		aclEntryStr(a, "source_cidr_block") == aclEntryStr(b, "source_cidr_block") &&
		aclEntryStr(a, "policy") == aclEntryStr(b, "policy") &&
		aclEntryStr(a, "description") == aclEntryStr(b, "description")
}
