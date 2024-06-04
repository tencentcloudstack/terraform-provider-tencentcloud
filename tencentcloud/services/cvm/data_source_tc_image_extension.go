package cvm

import (
	"context"
	"errors"
	"fmt"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"regexp"
	"strings"

	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func dataSourceTencentCloudImageReadPostHandleResponse0(ctx context.Context, req map[string]interface{}, resp *[]*cvm.Image) error {
	d := tccommon.ResourceDataFromContext(ctx)
	images := *resp
	var err error
	if len(images) == 0 {
		return errors.New("No image found")
	}

	var osName string
	if v, ok := d.GetOk("os_name"); ok {
		osName = v.(string)
	}

	var regImageName string
	var imageNameRegex *regexp.Regexp
	if v, ok := d.GetOk("image_name_regex"); ok {
		regImageName = v.(string)
		imageNameRegex, err = regexp.Compile(regImageName)
		if err != nil {
			return fmt.Errorf("image_name_regex format error,%s", err.Error())
		}
	}

	var resultImageId string
	images = sortImages(images)
	for _, image := range images {
		if osName != "" {
			if strings.Contains(strings.ToLower(*image.OsName), strings.ToLower(osName)) {
				resultImageId = *image.ImageId
				_ = d.Set("image_name", *image.ImageName)
				break
			}
			continue
		}

		if imageNameRegex != nil {
			if imageNameRegex.MatchString(*image.ImageName) {
				resultImageId = *image.ImageId
				_ = d.Set("image_name", *image.ImageName)
				break
			}
			continue
		}

		resultImageId = *image.ImageId
		_ = d.Set("image_name", *image.ImageName)
		break
	}

	if resultImageId == "" {
		return errors.New("No image found")
	}

	d.SetId(helper.DataResourceIdHash(resultImageId))
	if err := d.Set("image_id", resultImageId); err != nil {
		return err
	}

	context.WithValue(ctx, "resultImageId", resultImageId)
	return nil
}

func dataSourceTencentCloudImageReadOutputContent(ctx context.Context) interface{} {
	resultImageId := ctx.Value("resultImageId")
	return resultImageId
}
