package tencentcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTencentCloudCcnBandwidthLimit() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCcnBandwidthLimitCreate,
		Read:   resourceTencentCloudCcnBandwidthLimitRead,
		Update: resourceTencentCloudCcnBandwidthLimitUpdate,
		Delete: resourceTencentCloudCcnBandwidthLimitDelete,

		Schema: map[string]*schema.Schema{
			"ccn_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the CCN.",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Limitation of region.",
			},
			"bandwidth_limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Limitation of bandwidth. Default is `0`.",
			},
			"dst_region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "Destination area restriction. If the `CCN` rate limit type is `OUTER_REGION_LIMIT`, " +
					"this value does not need to be set.",
			},
		},
	}
}

func resourceTencentCloudCcnBandwidthLimitCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ccn_bandwidth_limit.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		ccnId  = d.Get("ccn_id").(string)
		region = d.Get("region").(string)
	)

	_, has, err := service.DescribeCcn(ctx, ccnId)
	if err != nil {
		return err
	}
	if has == 0 {
		return fmt.Errorf("ccn[%s] doesn't exist", ccnId)
	}
	id := fmt.Sprintf("%s#%s", ccnId, region)
	var (
		dstRegion string
		limit     int64
	)
	if v, ok := d.GetOkExists("dst_region"); ok {
		dstRegion = v.(string)
	}
	if v, ok := d.GetOk("bandwidth_limit"); ok {
		limit = int64(v.(int))
	}
	if err := service.SetCcnRegionBandwidthLimits(ctx, ccnId, region, dstRegion, limit, false); err != nil {
		return err
	}
	d.SetId(id)

	return resourceTencentCloudCcnBandwidthLimitRead(d, meta)
}

func resourceTencentCloudCcnBandwidthLimitUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ccn_bandwidth_limit.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		ccnId  = d.Get("ccn_id").(string)
		region = d.Get("region").(string)
	)
	_, has, err := service.DescribeCcn(ctx, ccnId)
	if err != nil {
		return err
	}
	if has == 0 {
		return fmt.Errorf("ccn[%s] doesn't exist", ccnId)
	}
	if d.HasChange("bandwidth_limit") {
		var limitTemp int64
		if v, ok := d.GetOk("bandwidth_limit"); ok {
			limitTemp = int64(v.(int))
		}
		_, dstRegion := d.GetChange("dst_region")
		if err := service.SetCcnRegionBandwidthLimits(ctx, ccnId, region, dstRegion.(string), limitTemp, false); err != nil {
			return err
		}
	}
	return resourceTencentCloudCcnBandwidthLimitRead(d, meta)
}

func resourceTencentCloudCcnBandwidthLimitRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ccn_bandwidth_limit.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		ccnId     = d.Get("ccn_id").(string)
		region    = d.Get("region").(string)
		dstRegion = d.Get("dst_region").(string)
		onlineHas = true
		info      CcnBasicInfo
	)
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		infoTmp, has, e := service.DescribeCcn(ctx, ccnId)
		if e != nil {
			return retryError(e)
		}

		if has == 0 {
			d.SetId("")
			onlineHas = false
			return nil
		}
		info = infoTmp
		return nil
	})
	if err != nil {
		return err
	}
	if !onlineHas {
		return nil
	}
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		bandwidth, e := service.GetCcnRegionBandwidthLimit(ctx, ccnId, region, dstRegion, info.bandWithLimitType)
		if e != nil {
			return retryError(e)
		}
		_ = d.Set("bandwidth_limit", bandwidth)
		_ = d.Set("dst_region", dstRegion)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudCcnBandwidthLimitDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ccn_bandwidth_limit.delete")()

	var (
		ccnId     = d.Get("ccn_id").(string)
		region    = d.Get("region").(string)
		dstRegion string
		limit     int64
	)

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	if v, ok := d.GetOk("dst_region"); ok {
		dstRegion = v.(string)
	}
	if v, ok := d.GetOk("bandwidth_limit"); ok {
		limit = int64(v.(int))
	}
	if err := service.SetCcnRegionBandwidthLimits(ctx, ccnId, region, dstRegion, limit, true); err != nil {
		return err
	}
	return nil
}
