package cvm

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svccbs "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cbs"

	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func dataSourceTencentCloudImagesReadOutputContent(ctx context.Context) interface{} {
	imageList := ctx.Value("imageList")
	return imageList
}

func dataSourceTencentCloudImagesReadPreRequest0(ctx context.Context, req *cvm.DescribeImagesRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)
	if v, ok := d.GetOk("instance_type"); ok {
		req.InstanceType = helper.String(v.(string))
	}

	return nil
}

func dataSourceTencentCloudImagesReadPostHandleResponse0(ctx context.Context, req map[string]interface{}, resp *[]*cvm.Image) error {
	d := tccommon.ResourceDataFromContext(ctx)
	images := *resp
	var (
		imageName      string
		osName         string
		imageNameRegex *regexp.Regexp
		err            error
	)

	if v, ok := d.GetOk("os_name"); ok {
		osName = v.(string)
	}

	if v, ok := d.GetOk("image_name_regex"); ok {
		imageName = v.(string)
		if imageName != "" {
			imageNameRegex, err = regexp.Compile(imageName)
			if err != nil {
				return fmt.Errorf("image_name_regex format error,%s", err.Error())
			}
		}
	}

	var results []*cvm.Image
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

	ctx = context.WithValue(ctx, "imageList", imageList)
	return nil
}

func dataSourceTencentCloudImagesReadPostFillRequest0(ctx context.Context, req map[string]interface{}) error {
	d := tccommon.ResourceDataFromContext(ctx)
	var (
		imageName string
		imageType []string
		err       error
	)

	if v, ok := d.GetOk("image_name_regex"); ok {
		imageName = v.(string)
		if imageName != "" {
			_, err = regexp.Compile(imageName)
			if err != nil {
				return fmt.Errorf("image_name_regex format error,%s", err.Error())
			}
		}
	}

	if v, ok := d.GetOk("image_type"); ok {
		for _, vv := range v.([]interface{}) {
			if vv.(string) != "" {
				imageType = append(imageType, vv.(string))
			}
		}
		if len(imageType) > 0 {
			req["image-type"] = imageType
		}
	}

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
