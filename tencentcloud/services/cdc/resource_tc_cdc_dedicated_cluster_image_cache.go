package cdc

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdc/v20201214"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDedicatedClusterImageCache() *schema.Resource {
	return &schema.Resource{
		Create: ResourceTencentCloudCdcDedicatedClusterImageCacheCreate,
		Read:   ResourceTencentCloudCdcDedicatedClusterImageCacheRead,
		Delete: ResourceTencentCloudCdcDedicatedClusterImageCacheDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"dedicated_cluster_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},
			"image_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Image ID.",
			},
		},
	}
}

func ResourceTencentCloudCdcDedicatedClusterImageCacheCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdc_dedicated_cluster_image_cache.create")()

	var (
		logId              = tccommon.GetLogId(tccommon.ContextNil)
		request            = cdc.NewCreateDedicatedClusterImageCacheRequest()
		dedicatedClusterId string
		image              string
	)

	if v, ok := d.GetOk("dedicated_cluster_id"); ok {
		dedicatedClusterId = v.(string)
		request.DedicatedClusterId = helper.String(dedicatedClusterId)
	}

	if v, ok := d.GetOk("image_id"); ok {
		image = v.(string)
		request.ImageId = helper.String(image)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCdcClient().CreateDedicatedClusterImageCache(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cdc dedicatedClusterImageCache failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{dedicatedClusterId, image}, tccommon.FILED_SP))

	service := CdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	conf := tccommon.BuildStateChangeConf([]string{}, []string{CDC_CACHE_STATUS_CACHED, CDC_CACHE_STATUS_CACHE_FAILED}, 20*tccommon.ReadRetryTimeout, time.Second, service.DedicatedClusterImageCacheStateRefreshFunc(dedicatedClusterId, image, CDC_CACHE_STATUS_CACHED_ALL, []string{}))
	if object, e := conf.WaitForState(); e != nil {
		return e
	} else {
		imageCacheState := object.(*cvm.Image)
		if imageCacheState.CdcCacheStatus != nil && *imageCacheState.CdcCacheStatus == CDC_CACHE_STATUS_CACHE_FAILED {
			return fmt.Errorf("cache failed")
		}
	}

	return ResourceTencentCloudCdcDedicatedClusterImageCacheRead(d, meta)
}

func ResourceTencentCloudCdcDedicatedClusterImageCacheRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdc_dedicated_cluster_image_cache.read")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = CdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	images, err := service.DescribeImages(ctx, idSplit[0], idSplit[1], CDC_CACHE_STATUS_CACHED_ALL)
	if err != nil {
		return err
	}

	if len(images) == 0 {
		d.SetId("")
		log.Printf("[WARN]%s resource `CdcDedicatedClusterImageCache` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if images[0].CdcCacheStatus != nil && *images[0].CdcCacheStatus == CDC_CACHE_STATUS_CACHED {
		_ = d.Set("dedicated_cluster_id", idSplit[0])
		_ = d.Set("image_id", idSplit[1])
	} else {
		d.SetId("")
		log.Printf("[WARN]%s resource `CdcDedicatedClusterImageCache` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	return nil
}

func ResourceTencentCloudCdcDedicatedClusterImageCacheDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdc_dedicated_cluster_image_cache.delete")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = cdc.NewDeleteDedicatedClusterImageCacheRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	request.DedicatedClusterId = helper.String(idSplit[0])
	request.ImageId = helper.String(idSplit[1])

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCdcClient().DeleteDedicatedClusterImageCache(request)
		if e != nil {
			if sdkerr, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkerr.Code == "InvalidParameter.ImageNotCacheInCdc" {
					return nil
				}
			}
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete cdc dedicatedClusterImageCache failed, reason:%+v", logId, err)
		return err
	}

	service := CdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	conf := tccommon.BuildStateChangeConf([]string{}, []string{CDC_CACHE_STATUS_NO_CACHE}, 20*tccommon.ReadRetryTimeout, time.Second, service.DedicatedClusterImageCacheStateRefreshFunc(idSplit[0], idSplit[1], CDC_CACHE_STATUS_NO_CACHE, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
