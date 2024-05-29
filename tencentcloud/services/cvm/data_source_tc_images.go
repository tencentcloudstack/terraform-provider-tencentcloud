package cvm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudImagesRead,
		Schema: map[string]*schema.Schema{
			"image_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the image to be queried.",
			},

			"image_name_regex": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"os_name"},
				Description:   "A regex string to apply to the image list returned by TencentCloud, conflict with 'os_name'. **NOTE**: it is not wildcard, should look like `image_name_regex = \"^CentOS\\s+6\\.8\\s+64\\w*\"`.",
			},

			"image_type": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of the image type to be queried. Valid values: 'PUBLIC_IMAGE', 'PRIVATE_IMAGE', 'SHARED_IMAGE', 'MARKET_IMAGE'.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"images": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of image. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"architecture": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Architecture of the image.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Created time of the image.",
						},
						"image_creator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image creator of the image.",
						},
						"image_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the image.",
						},
						"image_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the image.",
						},
						"image_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the image.",
						},
						"image_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Size of the image.",
						},
						"image_source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image source of the image.",
						},
						"image_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "State of the image.",
						},
						"image_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the image.",
						},
						"os_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "OS name of the image.",
						},
						"platform": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Platform of the image.",
						},
						"snapshots": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of snapshot details.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Size of the cloud disk used to create the snapshot; unit: GB.",
									},
									"disk_usage": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of the cloud disk used to create the snapshot.",
									},
									"snapshot_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Snapshot ID.",
									},
									"snapshot_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Snapshot name, the user-defined snapshot alias.",
									},
								},
							},
						},
						"support_cloud_init": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether support cloud-init.",
						},
						"sync_percent": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Sync percent of the image.",
						},
					},
				},
			},

			"instance_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Instance type, such as `S1.SMALL1`.",
			},

			"os_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"image_name_regex"},
				Description:   "A string to apply with fuzzy match to the os_name attribute on the image list returned by TencentCloud, conflict with 'image_name_regex'.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudImagesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_images.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	var filtersList []*cvm.Filter
	filtersMap := map[string]*cvm.Filter{}
	filter := cvm.Filter{}
	name := "image-id"
	filter.Name = &name
	if v, ok := d.GetOk("image_id"); ok {
		filter.Values = []*string{helper.String(v.(string))}
	}
	filtersMap["Temp0"] = &filter
	if v, ok := filtersMap["Temp0"]; ok {
		filtersList = append(filtersList, v)
	}
	filter2 := cvm.Filter{}
	name2 := "image_type"
	filter2.Name = &name2
	if v, ok := d.GetOk("image_type"); ok {
		filter2.Values = []*string{helper.String(v.(string))}
	}
	filtersMap["Temp1"] = &filter2
	if v, ok := filtersMap["Temp1"]; ok {
		filtersList = append(filtersList, v)
	}
	paramMap["Filters"] = filtersList

	if err := dataSourceTencentCloudImagesReadPostFillRequest0(ctx, paramMap); err != nil {
		return err
	}

	var respData *cvm.DescribeImagesResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeImagesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})
	if err != nil {
		return err
	}

	if err := dataSourceTencentCloudImagesReadPostHandleResponse0(ctx, paramMap, respData); err != nil {
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), dataSourceTencentCloudImagesReadOutputContent(ctx)); e != nil {
			return e
		}
	}

	return nil
}
