package dbdc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbdcv20201029 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDbdcDbCustomDisasterRecoverGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDbdcDbCustomDisasterRecoverGroupCreate,
		Read:   resourceTencentCloudDbdcDbCustomDisasterRecoverGroupRead,
		Update: resourceTencentCloudDbdcDbCustomDisasterRecoverGroupUpdate,
		Delete: resourceTencentCloudDbdcDbCustomDisasterRecoverGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Placement group name. Up to 60 characters, only Chinese and English are allowed.",
			},

			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Placement group type. Valid values: `HOST` (physical machine). Default value: `HOST`.",
			},

			"strategy": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Placement group strategy. Valid values: `SPREAD` (spread placement group). Default value: `SPREAD`.",
			},

			"affinity": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Affinity of the placement group. Instances in the group will be distributed according to this affinity. Valid values: `[1, 10]`. Default value: `1`.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			// computed
			"disaster_recover_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Placement group ID.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Placement group status. Valid values: `Creating`, `Available`, `CreateFailed`, `Deleting`, `Modifying`.",
			},

			"node_quota_total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum number of nodes that the placement group can hold.",
			},

			"current_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Current number of nodes in the placement group.",
			},

			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time.",
			},

			"node_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of DB Custom node IDs in the placement group.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceTencentCloudDbdcDbCustomDisasterRecoverGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dbdc_db_custom_disaster_recover_group.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                  = tccommon.GetLogId(tccommon.ContextNil)
		ctx                    = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service                = DbdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request                = dbdcv20201029.NewCreateDBCustomDisasterRecoverGroupRequest()
		disasterRecoverGroupId string
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("strategy"); ok {
		request.Strategy = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("affinity"); ok {
		request.Affinity = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("tags"); ok {
		for tagKey, tagValue := range v.(map[string]interface{}) {
			tag := dbdcv20201029.Tag{
				Key:   helper.String(tagKey),
				Value: helper.String(tagValue.(string)),
			}

			request.Tags = append(request.Tags, &tag)
		}
	}

	response := dbdcv20201029.NewCreateDBCustomDisasterRecoverGroupResponse()
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDbdcV20201029Client().CreateDBCustomDisasterRecoverGroupWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create dbdc db custom disaster recover group failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create dbdc db custom disaster recover group failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	logId = tccommon.GetLogId(ctx)
	log.Printf("[DEBUG]%s create dbdc db custom disaster recover group, logId=%s, d.Id()=%s", logId, logId, d.Id())
	if response.Response.DisasterRecoverGroupId == nil || *response.Response.DisasterRecoverGroupId == "" {
		return fmt.Errorf("Create dbdc db custom disaster recover group failed, DisasterRecoverGroupId is nil or empty.")
	}

	disasterRecoverGroupId = *response.Response.DisasterRecoverGroupId
	d.SetId(disasterRecoverGroupId)

	// Create finalization: poll Describe until Status == "Available".
	if err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		group, e := service.DescribeDBCustomDisasterRecoverGroupById(ctx, disasterRecoverGroupId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if group == nil || group.Status == nil {
			return resource.RetryableError(fmt.Errorf("dbdc db custom disaster recover group [%s] is still creating.", disasterRecoverGroupId))
		}

		status := *group.Status
		switch status {
		case "Available":
			return nil
		case "CreateFailed":
			return resource.NonRetryableError(fmt.Errorf("Create dbdc db custom disaster recover group [%s] failed, status is `CreateFailed`.", disasterRecoverGroupId))
		default:
			return resource.RetryableError(fmt.Errorf("dbdc db custom disaster recover group [%s] is still creating, status is `%s`.", disasterRecoverGroupId, status))
		}
	}); err != nil {
		return err
	}

	return resourceTencentCloudDbdcDbCustomDisasterRecoverGroupRead(d, meta)
}

func resourceTencentCloudDbdcDbCustomDisasterRecoverGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dbdc_db_custom_disaster_recover_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = DbdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	respData, err := service.DescribeDBCustomDisasterRecoverGroupById(ctx, d.Id())
	if err != nil {
		return err
	}

	if respData == nil || respData.DisasterRecoverGroupId == nil {
		log.Printf("[CRUD] tencentcloud_dbdc_db_custom_disaster_recover_group id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.Type != nil {
		_ = d.Set("type", respData.Type)
	}

	if respData.Strategy != nil {
		_ = d.Set("strategy", respData.Strategy)
	}

	if respData.Affinity != nil {
		_ = d.Set("affinity", respData.Affinity)
	}

	if respData.Tags != nil {
		tags := make(map[string]interface{}, len(respData.Tags))
		for _, tag := range respData.Tags {
			if tag == nil || tag.Key == nil {
				continue
			}

			if tag.Value != nil {
				tags[*tag.Key] = *tag.Value
			} else {
				tags[*tag.Key] = ""
			}
		}

		_ = d.Set("tags", tags)
	}

	if respData.DisasterRecoverGroupId != nil {
		_ = d.Set("disaster_recover_group_id", respData.DisasterRecoverGroupId)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.NodeQuotaTotal != nil {
		_ = d.Set("node_quota_total", respData.NodeQuotaTotal)
	}

	if respData.CurrentNum != nil {
		_ = d.Set("current_num", respData.CurrentNum)
	}

	if respData.CreatedTime != nil {
		_ = d.Set("created_time", respData.CreatedTime)
	}

	if respData.NodeIds != nil {
		_ = d.Set("node_ids", helper.StringsInterfaces(respData.NodeIds))
	}

	return nil
}

func resourceTencentCloudDbdcDbCustomDisasterRecoverGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dbdc_db_custom_disaster_recover_group.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                  = tccommon.GetLogId(tccommon.ContextNil)
		ctx                    = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service                = DbdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		disasterRecoverGroupId = d.Id()
	)

	immutableArgs := []string{"type", "strategy"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed in-place, please recreate the dbdc db custom disaster recover group resource.", v)
		}
	}

	if d.HasChange("name") || d.HasChange("affinity") {
		request := dbdcv20201029.NewModifyDBCustomDisasterRecoverGroupAttributeRequest()
		response := dbdcv20201029.NewModifyDBCustomDisasterRecoverGroupAttributeResponse()
		request.DisasterRecoverGroupId = helper.String(disasterRecoverGroupId)

		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("affinity"); ok {
			request.Affinity = helper.IntInt64(v.(int))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDbdcV20201029Client().ModifyDBCustomDisasterRecoverGroupAttributeWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify dbdc db custom disaster recover group attribute failed, Response is nil."))
			}

			response = result
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s modify dbdc db custom disaster recover group attribute failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		// wait
		if response.Response.TaskId != nil {
			if err := waitDBCustomTaskSucceeded(ctx, &service, *response.Response.TaskId, d.Timeout(schema.TimeoutUpdate)); err != nil {
				return err
			}
		}
	}

	if d.HasChange("tags") {
		oldRaw, newRaw := d.GetChange("tags")
		oldTags := oldRaw.(map[string]interface{})
		newTags := newRaw.(map[string]interface{})

		request := dbdcv20201029.NewModifyDBCustomDisasterRecoverGroupTagsRequest()
		request.DisasterRecoverGroupId = helper.String(disasterRecoverGroupId)

		for tagKey, tagValue := range newTags {
			tag := dbdcv20201029.Tag{
				Key:   helper.String(tagKey),
				Value: helper.String(tagValue.(string)),
			}

			request.AddTags = append(request.AddTags, &tag)
		}

		for tagKey := range oldTags {
			if _, ok := newTags[tagKey]; !ok {
				request.DeleteTagKeys = append(request.DeleteTagKeys, helper.String(tagKey))
			}
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDbdcV20201029Client().ModifyDBCustomDisasterRecoverGroupTagsWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify dbdc db custom disaster recover group attribute failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update dbdc db custom disaster recover group attribute failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudDbdcDbCustomDisasterRecoverGroupRead(d, meta)
}

func resourceTencentCloudDbdcDbCustomDisasterRecoverGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dbdc_db_custom_disaster_recover_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                  = tccommon.GetLogId(tccommon.ContextNil)
		ctx                    = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service                = DbdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request                = dbdcv20201029.NewDeleteDBCustomDisasterRecoverGroupsRequest()
		disasterRecoverGroupId = d.Id()
	)

	request.DisasterRecoverGroupIds = []*string{helper.String(disasterRecoverGroupId)}
	response := dbdcv20201029.NewDeleteDBCustomDisasterRecoverGroupsResponse()
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDbdcV20201029Client().DeleteDBCustomDisasterRecoverGroupsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete dbdc db custom disaster recover group failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete dbdc db custom disaster recover group failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// Delete is async, wait for the task to succeed.
	if response.Response.TaskId != nil {
		if err := waitDBCustomTaskSucceeded(ctx, &service, *response.Response.TaskId, d.Timeout(schema.TimeoutDelete)); err != nil {
			return err
		}
	}

	return nil
}
