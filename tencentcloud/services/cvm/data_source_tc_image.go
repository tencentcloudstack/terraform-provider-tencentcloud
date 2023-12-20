package cvm

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudImageRead,

		Schema: map[string]*schema.Schema{
			"filter": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "One or more name/value pairs to filter.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Key of the filter, valid keys: `image-id`, `image-type`, `image-name`.",
						},
						"values": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Values of the filter.",
						},
					},
				},
			},
			"image_name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateNameRegex,
				Description:  "A regex string to apply to the image list returned by TencentCloud. **NOTE**: it is not wildcard, should look like `image_name_regex = \"^CentOS\\s+6\\.8\\s+64\\w*\"`.",
			},
			"os_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateNotEmpty,
				Description:  "A string to apply with fuzzy match to the os_name attribute on the image list returned by TencentCloud. **NOTE**: when os_name is provided, highest priority is applied in this field instead of `image_name_regex`.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"image_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An image id indicate the uniqueness of a certain image,  which can be used for instance creation or resetting.",
			},
			"image_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of this image.",
			},
		},
	}
}

func dataSourceTencentCloudImageRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_image.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	cvmService := CvmService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	filter := make(map[string][]string)
	filters, ok := d.GetOk("filter")
	if ok {
		for _, v := range filters.(*schema.Set).List() {
			vv := v.(map[string]interface{})
			name := vv["name"].(string)
			filter[name] = []string{}
			for _, vvv := range vv["values"].([]interface{}) {
				filter[name] = append(filter[name], vvv.(string))
			}
		}
	}

	var images []*cvm.Image
	var errRet error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		images, errRet = cvmService.DescribeImagesByFilter(ctx, filter, "")
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}

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

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err = tccommon.WriteToFile(output.(string), resultImageId); err != nil {
			return err
		}
	}

	return nil
}
