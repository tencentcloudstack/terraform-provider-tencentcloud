/*
Provides a resource to create a vpc bandwidth_package

Example Usage

```hcl
resource "tencentcloud_vpc_bandwidth_package" "bandwidth_package" {
  network_type            = "BGP"
  charge_type             = "TOP5_POSTPAID_BY_MONTH"
  bandwidth_package_name  = "test-001"
  tags = {
    "createdBy" = "terraform"
  }
}

```
Import

vpc bandwidth_package can be imported using the id, e.g.
```
$ terraform import tencentcloud_vpc_bandwidth_package.bandwidth_package bandwidthPackage_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudVpcBandwidthPackage() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudVpcBandwidthPackageRead,
		Create: resourceTencentCloudVpcBandwidthPackageCreate,
		Update: resourceTencentCloudVpcBandwidthPackageUpdate,
		Delete: resourceTencentCloudVpcBandwidthPackageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"network_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Bandwidth packet type, default:BGP, optional:- `BGP`: common BGP shared bandwidth package- `HIGH_QUALITY_BGP`: Quality BGP Shared Bandwidth Package.",
			},

			"charge_type": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Bandwidth package billing type, default: TOP5_POSTPAID_BY_MONTH." +
					" Optional value: `TOP5_POSTPAID_BY_MONTH`: TOP5 billed by monthly postpaid; `PERCENT95_POSTPAID_BY_MONTH`: 95 billed monthly postpaid;" +
					" `FIXED_PREPAID_BY_MONTH`: Monthly prepaid billing (Type FIXED_PREPAID_BY_MONTH product API capability is under construction);" +
					" `BANDWIDTH_POSTPAID_BY_DAY`: bandwidth billed by daily postpaid; `ENHANCED95_POSTPAID_BY_MONTH`: enhanced 95 billed monthly postpaid.",
			},

			"bandwidth_package_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Bandwidth package name.",
			},

			"internet_max_bandwidth": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Bandwidth packet speed limit size. Unit: Mbps, -1 means no speed limit.",
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
	defer logElapsed("resource.tencentcloud_bwp_bandwidth_package.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = vpc.NewCreateBandwidthPackageRequest()
		response *vpc.CreateBandwidthPackageResponse
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

	if v, ok := d.GetOk("internet_max_bandwidth"); ok {
		request.InternetMaxBandwidth = helper.IntInt64(v.(int))
	}

	if v := helper.GetTags(d, "tags"); len(v) > 0 {
		for tagKey, tagValue := range v {
			tag := vpc.Tag{
				Key:   helper.String(tagKey),
				Value: helper.String(tagValue),
			}
			request.Tags = append(request.Tags, &tag)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().CreateBandwidthPackage(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create bwp bandwidthPackage failed, reason:%+v", logId, err)
		return err
	}

	bandwidthPackageId := *response.Response.BandwidthPackageId

	d.SetId(bandwidthPackageId)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::vpc:%s:uin/:bandwidthPackage/%s", region, bandwidthPackageId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	err = resource.Retry(3*readRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeVpcBandwidthPackage(ctx, bandwidthPackageId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if instance == nil {
			return resource.RetryableError(fmt.Errorf("vpc bandwidthPackage instance is being created, retry..."))
		}
		if *instance.Status == "CREATED" {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("vpc bandwidthPackage instance status is %v, retry...", *instance.Status))
	})
	if err != nil {
		return err
	}

	return resourceTencentCloudVpcBandwidthPackageRead(d, meta)
}

func resourceTencentCloudVpcBandwidthPackageRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_bwp_bandwidth_package.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	bandwidthPackageId := d.Id()

	bandwidthPackage, err := service.DescribeVpcBandwidthPackage(ctx, bandwidthPackageId)

	if err != nil {
		return err
	}

	if bandwidthPackage == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `tencentcloud_vpc_bandwidth_package` [%s] not found, please check if it has been deleted.",
			logId, bandwidthPackageId,
		)
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

	if bandwidthPackage.Bandwidth != nil {
		_ = d.Set("internet_max_bandwidth", bandwidthPackage.Bandwidth)
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
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := vpc.NewModifyBandwidthPackageAttributeRequest()

	bandwidthPackageId := d.Id()

	request.BandwidthPackageId = &bandwidthPackageId

	if d.HasChange("network_type") {
		return fmt.Errorf("`network_type` do not support change now.")
	}

	if d.HasChange("bandwidth_package_count") {
		return fmt.Errorf("`bandwidth_package_count` do not support change now.")
	}

	if d.HasChange("internet_max_bandwidth") {
		return fmt.Errorf("`internet_max_bandwidth` do not support change now.")
	}

	if d.HasChange("protocol") {
		return fmt.Errorf("`protocol` do not support change now.")
	}

	if d.HasChange("internet_max_bandwidth") {
		return fmt.Errorf("`internet_max_bandwidth` do not support change now.")
	}

	if d.HasChange("charge_type") {
		if v, ok := d.GetOk("charge_type"); ok {
			request.ChargeType = helper.String(v.(string))
		}
	}

	if v, ok := d.GetOk("bandwidth_package_name"); ok {
		request.BandwidthPackageName = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().ModifyBandwidthPackageAttribute(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create vpc bandwidthPackage failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
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
	defer logElapsed("resource.tencentcloud_bwp_bandwidth_package.delete")()
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
