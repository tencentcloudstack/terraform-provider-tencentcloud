/*
Provides a resource to create a vpc bandwidth_package

Example Usage

```hcl
resource "tencentcloud_vpc_bandwidth_package" "bandwidth_package" {
  network_type = &lt;nil&gt;
  charge_type = &lt;nil&gt;
  bandwidth_package_name = &lt;nil&gt;
  bandwidth_package_count = &lt;nil&gt;
  internet_max_bandwidth = &lt;nil&gt;
  protocol = &lt;nil&gt;
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

vpc bandwidth_package can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_bandwidth_package.bandwidth_package bandwidth_package_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudVpcBandwidthPackage() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcBandwidthPackageCreate,
		Read:   resourceTencentCloudVpcBandwidthPackageRead,
		Update: resourceTencentCloudVpcBandwidthPackageUpdate,
		Delete: resourceTencentCloudVpcBandwidthPackageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"network_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Bandwidth packet type, default:BGP, optional:- `BGP`: common BGP shared bandwidth package- `HIGH_QUALITY_BGP`: Quality BGP Shared Bandwidth Package.",
			},

			"charge_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Bandwidth package billing type, default: TOP5_POSTPAID_BY_MONTH, optional value:- `TOP5_POSTPAID_BY_MONTH`: TOP5 billed by monthly postpaid- `PERCENT95_POSTPAID_BY_MONTH`: 95 billed monthly postpaid- `FIXED_PREPAID_BY_MONTH`: Monthly prepaid billing.",
			},

			"bandwidth_package_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Bandwidth package name.",
			},

			"bandwidth_package_count": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The number of bandwidth packages (traditional account type can only fill in 1), the value range of standard account type is 1~20.",
			},

			"internet_max_bandwidth": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Bandwidth packet rate limit size. Unit: Mbps, `-1` means unlimited speed.",
			},

			"protocol": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Bandwidth packet protocol type. Currently supports `ipv4` and `ipv6` protocol bandwidth packets, the default value is `ipv4`.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudVpcBandwidthPackageCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_bandwidth_package.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request            = vpc.NewCreateBandwidthPackageRequest()
		response           = vpc.NewCreateBandwidthPackageResponse()
		bandwidthPackageId string
	)
	if v, ok := d.GetOk("network_type"); ok {
		request.NetworkType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("charge_type"); ok {
		request.ChargeType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("bandwidth_package_name"); ok {
		request.BandwidthPackageName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("bandwidth_package_count"); ok {
		request.BandwidthPackageCount = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("internet_max_bandwidth"); ok {
		request.InternetMaxBandwidth = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("protocol"); ok {
		request.Protocol = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().CreateBandwidthPackage(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vpc bandwidthPackage failed, reason:%+v", logId, err)
		return err
	}

	bandwidthPackageId = *response.Response.BandwidthPackageId
	d.SetId(bandwidthPackageId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"CREATED"}, 3*readRetryTimeout, time.Second, service.VpcBandwidthPackageStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::vpc:%s:uin/:bandwidthPackage/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudVpcBandwidthPackageRead(d, meta)
}

func resourceTencentCloudVpcBandwidthPackageRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_bandwidth_package.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	bandwidthPackageId := d.Id()

	bandwidthPackage, err := service.DescribeVpcBandwidthPackageById(ctx, bandwidthPackageId)
	if err != nil {
		return err
	}

	if bandwidthPackage == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcBandwidthPackage` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if bandwidthPackage.NetworkType != nil {
		_ = d.Set("network_type", bandwidthPackage.NetworkType)
	}

	if bandwidthPackage.ChargeType != nil {
		_ = d.Set("charge_type", bandwidthPackage.ChargeType)
	}

	if bandwidthPackage.BandwidthPackageName != nil {
		_ = d.Set("bandwidth_package_name", bandwidthPackage.BandwidthPackageName)
	}

	if bandwidthPackage.BandwidthPackageCount != nil {
		_ = d.Set("bandwidth_package_count", bandwidthPackage.BandwidthPackageCount)
	}

	if bandwidthPackage.InternetMaxBandwidth != nil {
		_ = d.Set("internet_max_bandwidth", bandwidthPackage.InternetMaxBandwidth)
	}

	if bandwidthPackage.Protocol != nil {
		_ = d.Set("protocol", bandwidthPackage.Protocol)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "vpc", "bandwidthPackage", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudVpcBandwidthPackageUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_bandwidth_package.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := vpc.NewModifyBandwidthPackageAttributeRequest()

	bandwidthPackageId := d.Id()

	request.BandwidthPackageId = &bandwidthPackageId

	immutableArgs := []string{"network_type", "charge_type", "bandwidth_package_name", "bandwidth_package_count", "internet_max_bandwidth", "protocol"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("charge_type") {
		if v, ok := d.GetOk("charge_type"); ok {
			request.ChargeType = helper.String(v.(string))
		}
	}

	if d.HasChange("bandwidth_package_name") {
		if v, ok := d.GetOk("bandwidth_package_name"); ok {
			request.BandwidthPackageName = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().ModifyBandwidthPackageAttribute(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update vpc bandwidthPackage failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("vpc", "bandwidthPackage", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudVpcBandwidthPackageRead(d, meta)
}

func resourceTencentCloudVpcBandwidthPackageDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_bandwidth_package.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	bandwidthPackageId := d.Id()

	if err := service.DeleteVpcBandwidthPackageById(ctx, bandwidthPackageId); err != nil {
		return err
	}

	return nil
}
