package ccn

import (
	"context"
	"fmt"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCcn() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCcnCreate,
		Read:   resourceTencentCloudCcnRead,
		Update: resourceTencentCloudCcnUpdate,
		Delete: resourceTencentCloudCcnDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
				Description:  "Name of the CCN to be queried, and maximum length does not exceed 60 bytes.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(0, 100),
				Description:  "Description of CCN, and maximum length does not exceed 100 bytes.",
			},
			"qos": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      CNN_QOS_AU,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{CNN_QOS_PT, CNN_QOS_AU, CNN_QOS_AG}),
				Description:  "Service quality of CCN. Valid values: `PT`, `AU`, `AG`. The default is `AU`.",
			},
			"charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      POSTPAID,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{POSTPAID, PREPAID}),
				Description: "Billing mode. Valid values: `PREPAID`, `POSTPAID`. " +
					"`PREPAID` means prepaid, which means annual and monthly subscription, " +
					"`POSTPAID` means post-payment, which means billing by volume. " +
					"The default is `POSTPAID`. The prepaid model only supports inter-regional speed limit, " +
					"and the post-paid model supports inter-regional speed limit and regional export speed limit.",
			},
			"bandwidth_limit_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      OuterRegionLimit,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{OuterRegionLimit, InterRegionLimit}),
				Description: "The speed limit type. Valid values: `INTER_REGION_LIMIT`, `OUTER_REGION_LIMIT`. " +
					"`OUTER_REGION_LIMIT` represents the regional export speed limit, " +
					"`INTER_REGION_LIMIT` is the inter-regional speed limit. " +
					"The default is `OUTER_REGION_LIMIT`.",
			},
			// Computed values
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "States of instance. Valid values: `ISOLATED`(arrears) and `AVAILABLE`.",
			},
			"instance_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of attached instances.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of resource.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Instance tag.",
			},
		},
	}
}

func resourceTencentCloudCcnCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		name               = d.Get("name").(string)
		description        = ""
		qos                = d.Get("qos").(string)
		chargeType         = d.Get("charge_type").(string)
		bandwidthLimitType = d.Get("bandwidth_limit_type").(string)
	)
	if temp, ok := d.GetOk("description"); ok {
		description = temp.(string)
	}
	info, err := service.CreateCcn(ctx, name, description, qos, chargeType, bandwidthLimitType)
	if err != nil {
		return err
	}
	d.SetId(info.ccnId)

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		resourceName := tccommon.BuildTagResourceName("vpc", "ccn", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudCcnRead(d, meta)
}

func resourceTencentCloudCcnRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		info, has, e := service.DescribeCcn(ctx, d.Id())
		if e != nil {
			return tccommon.RetryError(e)
		}

		if has == 0 {
			d.SetId("")
			return nil
		}

		_ = d.Set("name", info.name)
		_ = d.Set("description", info.description)
		_ = d.Set("qos", strings.ToUpper(info.qos))
		_ = d.Set("state", strings.ToUpper(info.state))
		_ = d.Set("instance_count", info.instanceCount)
		_ = d.Set("create_time", info.createTime)
		_ = d.Set("charge_type", info.chargeType)
		_ = d.Set("bandwidth_limit_type", info.bandWithLimitType)

		return nil
	})
	if err != nil {
		return err
	}
	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "vpc", "ccn", tcClient.Region, d.Id())
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)
	return nil
}

func resourceTencentCloudCcnUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		name        = ""
		description = ""
		change      = false
	)
	if d.HasChange("name") {
		name = d.Get("name").(string)
		change = true
	}

	if d.HasChange("description") {
		if temp, ok := d.GetOk("description"); ok {
			description = temp.(string)
		}
		if description == "" {
			return fmt.Errorf("can not set description='' ")
		}
		change = true
	}

	d.Partial(true)
	if change {
		if err := service.ModifyCcnAttribute(ctx, d.Id(), name, description); err != nil {
			return err
		}
	}
	// modify band width limit type
	if d.HasChange("bandwidth_limit_type") {
		_, news := d.GetChange("bandwidth_limit_type")
		if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			if err := service.ModifyCcnRegionBandwidthLimitsType(ctx, d.Id(), news.(string)); err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		}); err != nil {
			return err
		}
	}

	if d.HasChange("tags") {

		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))

		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		resourceName := tccommon.BuildTagResourceName("vpc", "ccn", tcClient.Region, d.Id())
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}

	}
	d.Partial(false)
	return resourceTencentCloudCcnRead(d, meta)
}

func resourceTencentCloudCcnDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		_, has, e := service.DescribeCcn(ctx, d.Id())
		if e != nil {
			return tccommon.RetryError(e)
		}
		if has == 0 {
			d.SetId("")
			return nil
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err = service.DeleteCcn(ctx, d.Id()); err != nil {
		return err
	}

	return resource.Retry(2*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		_, has, err := service.DescribeCcn(ctx, d.Id())
		if err != nil {
			return resource.RetryableError(err)
		}
		if has == 0 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("delete fail"))
	})
}
