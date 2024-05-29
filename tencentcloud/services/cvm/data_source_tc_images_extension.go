package cvm

import (
	"context"
	"fmt"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svccbs "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cbs"
	"log"
	"regexp"
	"strings"

	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func dataSourceTencentCloudImagesReadOutputContent(ctx context.Context) interface{} {
	imageList := ctx.Value("imageList").([]interface{})
	return imageList
}

func dataSourceTencentCloudImagesReadPostHandleResponse0(ctx context.Context, req map[string]interface{}, resp *cvm.DescribeImagesResponseParams) error {
	d := tccommon.ResourceDataFromContext(ctx)
	images := resp.ImageSet
	osName := ctx.Value("osName").(string)
	imageName := ctx.Value("imageName").(string)
	imageNameRegex := ctx.Value("imageNameRegex").(*regexp.Regexp)
	var results []*cvm.Image
	var err error
	images = sortImages(images)
	if osName == "" && imageName == "" {
		results = images
	} else {
		for _, image := range images {
			if osName != "" {
				if strings.Contains(strings.ToLower(*image.OsName), strings.ToLower(osName)) {
					results = append(results, image)
					continue
				}
			}
			if imageNameRegex != nil {
				if imageNameRegex.MatchString(*image.ImageName) {
					results = append(results, image)
					continue
				}
			}
		}
	}

	logId := tccommon.GetLogId(tccommon.ContextNil)
	meta := tccommon.ProviderMetaFromContext(ctx)
	cbsService := svccbs.NewCbsService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	imageList := make([]map[string]interface{}, 0, len(results))
	ids := make([]string, 0, len(results))
	for _, image := range results {
		snapshots, err := imagesReadSnapshotByIds(ctx, cbsService, image)
		if err != nil {
			return err
		}

		mapping := map[string]interface{}{
			"image_id":           image.ImageId,
			"os_name":            image.OsName,
			"image_type":         image.ImageType,
			"created_time":       image.CreatedTime,
			"image_name":         image.ImageName,
			"image_description":  image.ImageDescription,
			"image_size":         image.ImageSize,
			"architecture":       image.Architecture,
			"image_state":        image.ImageState,
			"platform":           image.Platform,
			"image_creator":      image.ImageCreator,
			"image_source":       image.ImageSource,
			"sync_percent":       image.SyncPercent,
			"support_cloud_init": image.IsSupportCloudinit,
			"snapshots":          snapshots,
		}
		imageList = append(imageList, mapping)
		ids = append(ids, *image.ImageId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("images", imageList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set image list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	context.WithValue(ctx, "imageList", imageList)
	return nil
}

func dataSourceTencentCloudImagesReadPostFillRequest0(ctx context.Context, req map[string]interface{}) error {
	d := tccommon.ResourceDataFromContext(ctx)
	var (
		imageName      string
		osName         string
		instanceType   string
		imageNameRegex *regexp.Regexp
		err            error
	)
	if v, ok := d.GetOk("image_name_regex"); ok {
		imageName = v.(string)
		if imageName != "" {
			imageNameRegex, err = regexp.Compile(imageName)
			if err != nil {
				return fmt.Errorf("image_name_regex format error,%s", err.Error())
			}
		}
	}

	if v, ok := d.GetOk("os_name"); ok {
		osName = v.(string)
	}

	if v, ok := d.GetOk("instance_type"); ok {
		instanceType = v.(string)
	}

	context.WithValue(ctx, "imageNameRegex", imageNameRegex)
	context.WithValue(ctx, "osName", osName)
	context.WithValue(ctx, "instanceType", instanceType)
	return nil
}

func imagesReadSnapshotByIds(ctx context.Context, cbsService svccbs.CbsService, image *cvm.Image) (snapshotResults []map[string]interface{}, errRet error) {
	if len(image.SnapshotSet) == 0 {
		return
	}

	snapshotByIds := make([]*string, 0, len(image.SnapshotSet))
	for _, snapshot := range image.SnapshotSet {
		snapshotByIds = append(snapshotByIds, snapshot.SnapshotId)
	}

	snapshots, errRet := cbsService.DescribeSnapshotByIds(ctx, snapshotByIds)
	if errRet != nil {
		return
	}

	snapshotResults = make([]map[string]interface{}, 0, len(snapshots))
	for _, snapshot := range snapshots {
		snapshotMap := make(map[string]interface{}, 4)
		snapshotMap["snapshot_id"] = snapshot.SnapshotId
		snapshotMap["disk_usage"] = snapshot.DiskUsage
		snapshotMap["disk_size"] = snapshot.DiskSize
		snapshotMap["snapshot_name"] = snapshot.SnapshotName

		snapshotResults = append(snapshotResults, snapshotMap)
	}

	return
}
