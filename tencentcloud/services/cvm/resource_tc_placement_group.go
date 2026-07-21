package cvm

import (
	"context"
	"fmt"
	"log"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func ResourceTencentCloudPlacementGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPlacementGroupCreate,
		Read:   resourceTencentCloudPlacementGroupRead,
		Update: resourceTencentCloudPlacementGroupUpdate,
		Delete: resourceTencentCloudPlacementGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
				Description:  "Name of the placement group, 1-60 characters in length.",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CVM_PLACEMENT_GROUP_TYPE),
				Description:  "Type of the placement group. Valid values: `HOST`, `SW` and `RACK`.",
			},
			"strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CVM_PLACEMENT_GROUP_STRATEGY),
				Description:  "Strategy of the placement group. Valid values: `SPREAD` and `PARTITION`. `SPREAD` is the default strategy. When strategy is `PARTITION`, `partition_count` must be set. This field cannot be modified after creation.",
			},
			"affinity": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 10),
				Description:  "Affinity of the placement group.Valid values: 1~10, default is 1.",
			},
			"partition_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(2, 30),
				Description:  "Partition count of the placement group. Valid values: 2~30. Only valid when `strategy` is set to `PARTITION`.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of the placement group.",
			},

			// computed
			"cvm_quota_total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum number of hosts in the placement group.",
			},
			"current_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of hosts in the placement group.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the placement group.",
			},
		},
	}
}

func resourceTencentCloudPlacementGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_placement_group.create")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	cvmService := CvmService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	placementName := d.Get("name").(string)
	placementType := d.Get("type").(string)

	var strategy string
	if v, ok := d.GetOk("strategy"); ok {
		strategy = v.(string)
	}

	var affinity int
	if v, ok := d.GetOkExists("affinity"); ok {
		affinity = v.(int)
	}

	var partitionCount int
	if v, ok := d.GetOkExists("partition_count"); ok {
		partitionCount = v.(int)
	}

	if partitionCount > 0 && strategy != CVM_PLACEMENT_GROUP_STRATEGY_PARTITION {
		return fmt.Errorf("`partition_count` is only valid when `strategy` is set to `PARTITION`")
	}

	tags := make([]*cvm.Tag, 0)
	if v := helper.GetTags(d, "tags"); len(v) > 0 {
		for tagKey, tagValue := range v {
			tag := cvm.Tag{
				Key:   helper.String(tagKey),
				Value: helper.String(tagValue),
			}
			tags = append(tags, &tag)
		}
	}
	var response *cvm.CreateDisasterRecoverGroupResponse
	var errRet error
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		response, errRet = cvmService.CreatePlacementGroup(ctx, placementName, placementType, strategy, affinity, partitionCount, tags)
		if errRet != nil {
			return tccommon.RetryError(errRet)
		}
		return nil
	})
	if err != nil {
		return err
	}
	d.SetId(*response.Response.DisasterRecoverGroupId)

	return resourceTencentCloudPlacementGroupRead(d, meta)
}

func resourceTencentCloudPlacementGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_placement_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	placementId := d.Id()
	cvmService := CvmService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	var placement *cvm.DisasterRecoverGroup
	var errRet error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		placement, errRet = cvmService.DescribePlacementGroupById(ctx, placementId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if placement == nil {
		log.Printf("[DEBUG] tencentcloud_placement_group read placement is nil, id=%s", d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("name", placement.Name)
	_ = d.Set("type", placement.Type)
	_ = d.Set("affinity", placement.Affinity)
	_ = d.Set("cvm_quota_total", placement.CvmQuotaTotal)
	_ = d.Set("current_num", placement.CurrentNum)
	_ = d.Set("create_time", placement.CreateTime)

	if placement.Strategy != nil {
		_ = d.Set("strategy", placement.Strategy)
	}
	if placement.PartitionCount != nil {
		_ = d.Set("partition_count", int(*placement.PartitionCount))
	}

	if len(placement.Tags) > 0 {
		_ = d.Set("tags", flattenInstanceTagsMapping(placement.Tags))
	}

	return nil
}

func resourceTencentCloudPlacementGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_placement_group.update")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	placementId := d.Id()
	cvmService := CvmService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	immutableArgs := []string{
		"type",
		"strategy",
		"affinity",
		"partition_count",
	}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("name") {
		placementName := d.Get("name").(string)
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			err := cvmService.ModifyPlacementGroup(ctx, placementId, placementName)
			if err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	if d.HasChange("tags") {
		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))

		client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(client)
		region := client.Region
		resourceName := tccommon.BuildTagResourceName("cvm", "ps", region, d.Id())
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudPlacementGroupRead(d, meta)
}

func resourceTencentCloudPlacementGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_placement_group.delete")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	placementId := d.Id()
	cvmService := CvmService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		err := cvmService.DeletePlacementGroup(ctx, placementId)
		if err != nil {
			return tccommon.RetryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
