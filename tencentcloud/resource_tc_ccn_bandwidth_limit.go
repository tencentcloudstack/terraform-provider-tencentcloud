/*
Provides a resource to limit CCN bandwidth.

Example Usage

Set the upper limit of regional outbound bandwidth

```hcl
variable "other_region1" {
  default = "ap-shanghai"
}

resource "tencentcloud_ccn" "main" {
  name        = "ci-temp-test-ccn"
  description = "ci-temp-test-ccn-des"
  qos         = "AG"
}

resource "tencentcloud_ccn_bandwidth_limit" "limit1" {
  ccn_id          = tencentcloud_ccn.main.id
  region          = var.other_region1
  bandwidth_limit = 500
}
```

Set the upper limit between regions

```hcl
variable "other_region1" {
  default = "ap-shanghai"
}

variable "other_region2" {
  default = "ap-nanjing"
}

resource tencentcloud_ccn main {
  name                 = "ci-temp-test-ccn"
  description          = "ci-temp-test-ccn-des"
  qos                  = "AG"
  bandwidth_limit_type = "INTER_REGION_LIMIT"
}

resource tencentcloud_ccn_bandwidth_limit limit1 {
  ccn_id          = tencentcloud_ccn.main.id
  region          = var.other_region1
  dst_region      = var.other_region2
  bandwidth_limit = 100
}
```

*/
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
				Description: "Limitation of bandwidth.",
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
	if err := service.SetCcnRegionBandwidthLimits(ctx, ccnId, region, dstRegion, limit); err != nil {
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
		if err := service.SetCcnRegionBandwidthLimits(ctx, ccnId, region, dstRegion.(string), limitTemp); err != nil {
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

	_, _ = d, meta

	return nil
}
