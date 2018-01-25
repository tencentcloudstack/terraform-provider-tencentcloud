package tencentcloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"sort"
	"time"

	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

const (
	tencentCloudApiDescribeImagesParaLimitMaxFiltersNumber       = 10
	tencentCloudApiDescribeImagesParaLimitMaxFilterValuessNumber = 5
)

type imageSorter []struct {
	ImageId     string `json:tag"ImageId"`
	OsName      string `json:tag"OsName"`
	ImageName   string `json:tag"ImageName"`
	CreatedTime string `json:tag"CreatedTime"`
}

func (a imageSorter) Len() int {
	return len(a)
}

func (a imageSorter) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a imageSorter) Less(i, j int) bool {
	itime, _ := time.Parse(time.RFC3339, a[i].CreatedTime)
	jtime, _ := time.Parse(time.RFC3339, a[j].CreatedTime)
	return itime.Unix() < jtime.Unix()
}

// Sort images by creation date, in descending order.
func sortImages(images imageSorter) imageSorter {
	sortedImages := images
	sort.Sort(sort.Reverse(imageSorter(sortedImages)))
	return sortedImages
}

func dataSourceTencentCloudSourceImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudImagesRead,

		Schema: map[string]*schema.Schema{
			"image_name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateNameRegex,
			},

			"os_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateNotEmpty,
			},

			//"most_recent": {
			//	Type:     schema.TypeBool,
			//	Optional: true,
			//	Default:  false,
			//	ForceNew: true,
			//},

			"filter": dataSourceTencentCloudFiltersSchema(),

			// Computed values.
			"image_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTencentCloudImagesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*TencentCloudClient).commonConn
	filters, filtersOk := d.GetOk("filter")
	imageNameRegex, nameRegexOk := d.GetOk("image_name_regex")
	osName, osNameOk := d.GetOk("os_name")

	if !nameRegexOk && !filtersOk && !osNameOk {
		return fmt.Errorf("One of image_name_regex, os_name, filter must be assigned")
	}
	imageNameRegexStr := imageNameRegex.(string)
	osNameStr := osName.(string)

	params := map[string]string{
		"Version": "2017-03-12",
		"Action":  "DescribeImages",
		"Limit":   "100",
	}

	if filtersOk {
		filterList := filters.(*schema.Set)
		err := buildFiltersParam(
			params,
			filterList,
			tencentCloudApiDescribeImagesParaLimitMaxFiltersNumber,
			tencentCloudApiDescribeImagesParaLimitMaxFilterValuessNumber,
		)
		if err != nil {
			return err
		}

	}

	log.Printf("[DEBUG] tencentcloud_image - param: %v", params)
	response, err := client.SendRequest("image", params)
	if err != nil {
		return err
	}

	var jsonresp struct {
		Response struct {
			Error struct {
				Code    string `json:tag"Code"`
				Message string `json:tag"Message"`
			}
			ImageSet []struct {
				ImageId     string `json:tag"ImageId"`
				OsName      string `json:tag"OsName"`
				ImageName   string `json:tag"ImageName"`
				CreatedTime string `json:tag"CreatedTime"`
			}
		}
	}

	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		return err
	}
	if jsonresp.Response.Error.Code != "" {
		return fmt.Errorf(
			"tencentcloud_image got error, code:%v, message:%v",
			jsonresp.Response.Error.Code,
			jsonresp.Response.Error.Message,
		)
	}

	var (
		resultImageId string
		regImageName  = regexp.MustCompile(imageNameRegexStr)
	)
	imageList := jsonresp.Response.ImageSet
	if len(imageList) == 0 {
		return errors.New("No image found")
	}

	// Dont expose most_recent as a outlet for the sake of eaiser usage, let's sort by recent by default
	//recent := d.Get("most_recent").(bool)
	recent := true
	if recent && len(imageList) > 1 {
		imageList = mostRecentImages(imageList)
	}

	for _, image := range imageList {
		imageId := image.ImageId
		imageName := image.ImageName
		osName := image.OsName
		log.Printf(
			"[DEBUG] tencentcloud_image - Images found imageId: %v, osName:% v, imageName: %v, imageNameRegexStr: %v",
			imageId,
			osName,
			imageName,
			imageNameRegexStr,
		)

		// osName full match is in the highest priority
		if osNameOk {
			if strings.Contains(strings.ToLower(osName), strings.ToLower(osNameStr)) {
				resultImageId = imageId
				break
			} else {
				continue
			}
		}

		if nameRegexOk {
			if regImageName.MatchString(imageName) {
				resultImageId = imageId
				break
			} else {
				continue
			}
		}

		if filtersOk {
			resultImageId = imageId
			break
		}
	}

	if nameRegexOk && resultImageId == "" {
		err = errors.New("osName not matched")
		return err
	}

	id := dataResourceIdHash(resultImageId)
	d.SetId(id)
	log.Printf("[DEBUG] tencentcloud_image - imageId: %#v", resultImageId)
	if err := d.Set("image_id", resultImageId); err != nil {
		return err
	}
	return nil
}

// Returns the most recent image out of a slice of images.
func mostRecentImages(images imageSorter) imageSorter {
	return sortImages(images)
}
