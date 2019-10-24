package tencentcloud

import (
	"context"
	"errors"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

type imageSorter []*cvm.Image

func (a imageSorter) Len() int {
	return len(a)
}

func (a imageSorter) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a imageSorter) Less(i, j int) bool {
	if a[i].CreatedTime == nil || a[j].CreatedTime == nil {
		return false
	}

	itime, _ := time.Parse(time.RFC3339, *a[i].CreatedTime)
	jtime, _ := time.Parse(time.RFC3339, *a[j].CreatedTime)

	return itime.Unix() < jtime.Unix()
}

// Sort images by creation date, in descending order.
func sortImages(images imageSorter) imageSorter {
	sortedImages := images
	sort.Sort(sort.Reverse(sortedImages))
	return sortedImages
}

// Returns the most recent image out of a slice of images.
func mostRecentImages(images imageSorter) imageSorter {
	return sortImages(images)
}

func dataSourceTencentCloudSourceImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudImagesRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceTencentCloudFiltersSchema(),
			"image_name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameRegex,
			},
			"os_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNotEmpty,
			},
			"image_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTencentCloudImagesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_image.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	cvmService := CvmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
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
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		images, errRet = cvmService.DescribeImagesByFilter(ctx, filter)
		if errRet != nil {
			return retryError(errRet, "InternalError")
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
		imageNameRegex = regexp.MustCompile(regImageName)
	}

	var resultImageId string
	images = mostRecentImages(images)
	for _, image := range images {
		if osName != "" {
			if strings.Contains(strings.ToLower(*image.OsName), strings.ToLower(osName)) {
				resultImageId = *image.ImageId
				d.Set("image_name", *image.ImageName)
				break
			}
			continue
		}

		if regImageName != "" {
			if imageNameRegex.MatchString(*image.ImageName) {
				resultImageId = *image.ImageId
				d.Set("image_name", *image.ImageName)
				break
			}
			continue
		}

		resultImageId = *image.ImageId
		d.Set("image_name", *image.ImageName)
		break
	}

	if resultImageId == "" {
		return errors.New("No image found")
	}

	id := dataResourceIdHash(resultImageId)
	d.SetId(id)

	if err := d.Set("image_id", resultImageId); err != nil {
		return err
	}

	return nil
}
