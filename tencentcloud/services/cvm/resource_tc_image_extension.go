package cvm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func resourceTencentCloudImageCreatePostFillRequest0(ctx context.Context, req *cvm.CreateImageRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)

	if d.Get("force_poweroff").(bool) {
		req.ForcePoweroff = helper.String(TRUE)
	} else {
		req.ForcePoweroff = helper.String(FALSE)
	}

	if v, ok := d.GetOkExists("sysprep"); ok {
		value := v.(bool)
		if value {
			req.Sysprep = helper.String(TRUE)
		} else {
			req.Sysprep = helper.String(FALSE)
		}
	}

	if len(req.SnapshotIds) > 0 && len(req.DataDiskIds) > 0 {
		return fmt.Errorf("`%s` and `%s` Can't appear in the profile China at the same time,The parameter `%s` depends on the pre_parameter `%s`",
			"snapshot_ids", "data_disk_ids", "data_disk_ids", "instance_id")
	}

	if v := helper.GetTags(d, "tags"); len(v) > 0 {
		tags := make([]*cvm.Tag, 0)
		for tagKey, tagValue := range v {
			tag := cvm.Tag{
				Key:   helper.String(tagKey),
				Value: helper.String(tagValue),
			}
			tags = append(tags, &tag)
		}
		tagSpecification := cvm.TagSpecification{
			ResourceType: helper.String("image"),
			Tags:         tags,
		}
		req.TagSpecification = append(req.TagSpecification, &tagSpecification)
	}

	return nil
}

func resourceTencentCloudImageCreatePostHandleResponse0(ctx context.Context, resp *cvm.CreateImageResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	logId := ctx.Value(tccommon.LogIdKey)

	cvmService := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	imageId := d.Id()

	// wait for status
	image, errRet := cvmService.DescribeImageById(ctx, imageId)
	if errRet != nil {
		return errRet
	}
	if image == nil {
		return fmt.Errorf("[CRITAL]%s creating cvm image failed, image doesn't exist", logId)
	}
	return nil
}

func resourceTencentCloudImageReadRequestOnSuccess0(ctx context.Context, resp *cvm.Image) *resource.RetryError {
	logId := ctx.Value(tccommon.LogIdKey)

	if resp != nil {
		if *resp.ImageState == "CREATEFAILED" {
			return resource.NonRetryableError(fmt.Errorf("[CRITAL]%s Create Image is failed", logId))
		}
		if *resp.ImageState != "NORMAL" {
			return resource.RetryableError(fmt.Errorf("iamge instance status is processing"))
		}
	}
	return nil
}

func resourceTencentCloudImageReadPostHandleResponse0(ctx context.Context, resp *cvm.Image) error {
	d := tccommon.ResourceDataFromContext(ctx)

	// Use the resource value when the instance_id in the resource is not empty.
	// the instance ID is not returned in the query response body.
	instanceId := ""
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	snapShotSysDisk := make([]interface{}, 0, len(resp.SnapshotSet))
	for _, v := range resp.SnapshotSet {
		snapShotSysDisk = append(snapShotSysDisk, v.SnapshotId)
	}

	if instanceId != "" {
		_ = d.Set("instance_id", helper.String(instanceId))
	} else {
		_ = d.Set("snapshot_ids", snapShotSysDisk)
	}

	return nil
}

func resourceTencentCloudImageDeletePostHandleResponse0(ctx context.Context, resp *cvm.DeleteImagesResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	logId := ctx.Value(tccommon.LogIdKey)

	cvmService := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	imageId := d.Id()

	//check image
	err := resource.Retry(3*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := cvmService.DescribeImageById(ctx, imageId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if result != nil {
			return resource.RetryableError(fmt.Errorf("image exits error,image_id = %s", imageId))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read image failed, reason:%+v", logId, err)
		return err
	}
	return nil
}
