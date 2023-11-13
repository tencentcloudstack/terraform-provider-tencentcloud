/*
Use this data source to query detailed information of cdb uploaded_files

Example Usage

```hcl
data "tencentcloud_cdb_uploaded_files" "uploaded_files" {
  path = ""
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCdbUploadedFiles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCdbUploadedFilesRead,
		Schema: map[string]*schema.Schema{
			"path": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "File path. This field should be filled with the OwnerUin information of the user&amp;amp;#39;s primary account.",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of SQL files returned.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"upload_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Upload time.",
						},
						"upload_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Upload progress.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"all_slice_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of all fragments of the file.",
									},
									"complete_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of shards completed.",
									},
								},
							},
						},
						"file_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "File name.",
						},
						"file_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "File size, the unit is Bytes.",
						},
						"is_upload_finished": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether the upload is completed or not, optional values: 0 - not completed, 1 - completed.",
						},
						"file_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "File ID.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCdbUploadedFilesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cdb_uploaded_files.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("path"); ok {
		paramMap["Path"] = helper.String(v.(string))
	}

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCdbUploadedFilesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		totalCount = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(totalCount))
	if items != nil {
		for _, sqlFileInfo := range items {
			sqlFileInfoMap := map[string]interface{}{}

			if sqlFileInfo.UploadTime != nil {
				sqlFileInfoMap["upload_time"] = sqlFileInfo.UploadTime
			}

			if sqlFileInfo.UploadInfo != nil {
				uploadInfoMap := map[string]interface{}{}

				if sqlFileInfo.UploadInfo.AllSliceNum != nil {
					uploadInfoMap["all_slice_num"] = sqlFileInfo.UploadInfo.AllSliceNum
				}

				if sqlFileInfo.UploadInfo.CompleteNum != nil {
					uploadInfoMap["complete_num"] = sqlFileInfo.UploadInfo.CompleteNum
				}

				sqlFileInfoMap["upload_info"] = []interface{}{uploadInfoMap}
			}

			if sqlFileInfo.FileName != nil {
				sqlFileInfoMap["file_name"] = sqlFileInfo.FileName
			}

			if sqlFileInfo.FileSize != nil {
				sqlFileInfoMap["file_size"] = sqlFileInfo.FileSize
			}

			if sqlFileInfo.IsUploadFinished != nil {
				sqlFileInfoMap["is_upload_finished"] = sqlFileInfo.IsUploadFinished
			}

			if sqlFileInfo.FileId != nil {
				sqlFileInfoMap["file_id"] = sqlFileInfo.FileId
			}

			ids = append(ids, *sqlFileInfo.IdsHash)
			tmpList = append(tmpList, sqlFileInfoMap)
		}

		_ = d.Set("items", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string)); e != nil {
			return e
		}
	}
	return nil
}
