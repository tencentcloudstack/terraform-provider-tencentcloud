/*
Use this data source to query images.

Example Usage

```hcl
data "tencentcloud_images" "foo" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "centos 7.5"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudImagesRead,

		Schema: map[string]*schema.Schema{
			"image_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the image to be queried.",
			},
			"image_type": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of the image type to be queried. Available values include: 'PUBLIC_IMAGE', 'PRIVATE_IMAGE', 'SHARED_IMAGE', 'MARKET_IMAGE'.",
			},
			"image_name_regex": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"os_name"},
				ValidateFunc:  validateNameRegex,
				Description:   "A regex string to apply to the image list returned by TencentCloud, conflict with 'os_name'. **NOTE**: it is not wildcard, should look like `image_name_regex = \"^CentOS\\s+6\\.8\\s+64\\w*\"`.",
			},
			"os_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"image_name_regex"},
				ValidateFunc:  validateNotEmpty,
				Description:   "A string to apply with fuzzy match to the os_name attribute on the image list returned by TencentCloud, conflict with 'image_name_regex'.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"images": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of image. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the image.",
						},
						"os_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "OS name of the image.",
						},
						"image_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the image.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Created time of the image.",
						},
						"image_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the image.",
						},
						"image_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the image.",
						},
						"image_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Size of the image.",
						},
						"architecture": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Architecture of the image.",
						},
						"image_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "State of the image.",
						},
						"platform": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Platform of the image.",
						},
						"image_creator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image creator of the image.",
						},
						"image_source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image source of the image.",
						},
						"sync_percent": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Sync percent of the image.",
						},
						"support_cloud_init": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether support cloud-init.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudImagesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_images.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	cvmService := CvmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	var (
		imageId        string
		imageType      []string
		imageName      string
		osName         string
		imageNameRegex *regexp.Regexp
		err            error
	)

	filter := make(map[string][]string)

	if v, ok := d.GetOk("image_id"); ok {
		imageId = v.(string)
		if imageId != "" {
			filter["image-id"] = []string{imageId}
		}
	}

	if v, ok := d.GetOk("image_type"); ok {
		for _, vv := range v.([]interface{}) {
			if vv.(string) != "" {
				imageType = append(imageType, vv.(string))
			}
		}
		if len(imageType) > 0 {
			filter["image-type"] = imageType
		}
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

	if v, ok := d.GetOk("os_name"); ok {
		osName = v.(string)
	}

	var images []*cvm.Image
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		var e error
		images, e = cvmService.DescribeImagesByFilter(ctx, filter)
		if e != nil {
			return retryError(e, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
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

	imageList := make([]map[string]interface{}, 0, len(results))
	ids := make([]string, 0, len(results))
	for _, image := range results {
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

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := writeToFile(output.(string), imageList); err != nil {
			return err
		}
	}

	return nil
}
