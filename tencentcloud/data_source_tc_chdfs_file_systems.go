/*
Use this data source to query detailed information of chdfs file_systems

Example Usage

```hcl
data "tencentcloud_chdfs_file_systems" "file_systems" {
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	chdfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/chdfs/v20201112"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudChdfsFileSystems() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudChdfsFileSystemsRead,
		Schema: map[string]*schema.Schema{
			"file_systems": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "File system list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Appid of the user.",
						},
						"file_system_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "File system name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Desc of the file system.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region of the file system.",
						},
						"file_system_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "File system id.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time.",
						},
						"block_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Block size of the file system(byte).",
						},
						"capacity_quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Capacity of the file system(byte).",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Status of the file system（1: creating create success 3: create failed).",
						},
						"super_users": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Super users of the file system.",
						},
						"posix_acl": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Check POSIX ACL or not.",
						},
						"enable_ranger": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Check the ranger address or not.",
						},
						"ranger_service_addresses": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Ranger adress list.",
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

func dataSourceTencentCloudChdfsFileSystemsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_chdfs_file_systems.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	service := ChdfsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var fileSystems []*chdfs.FileSystem

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeChdfsFileSystemsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		fileSystems = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(fileSystems))
	tmpList := make([]map[string]interface{}, 0, len(fileSystems))

	if fileSystems != nil {
		for _, fileSystem := range fileSystems {
			fileSystemMap := map[string]interface{}{}

			if fileSystem.AppId != nil {
				fileSystemMap["app_id"] = fileSystem.AppId
			}

			if fileSystem.FileSystemName != nil {
				fileSystemMap["file_system_name"] = fileSystem.FileSystemName
			}

			if fileSystem.Description != nil {
				fileSystemMap["description"] = fileSystem.Description
			}

			if fileSystem.Region != nil {
				fileSystemMap["region"] = fileSystem.Region
			}

			if fileSystem.FileSystemId != nil {
				fileSystemMap["file_system_id"] = fileSystem.FileSystemId
			}

			if fileSystem.CreateTime != nil {
				fileSystemMap["create_time"] = fileSystem.CreateTime
			}

			if fileSystem.BlockSize != nil {
				fileSystemMap["block_size"] = fileSystem.BlockSize
			}

			if fileSystem.CapacityQuota != nil {
				fileSystemMap["capacity_quota"] = fileSystem.CapacityQuota
			}

			if fileSystem.Status != nil {
				fileSystemMap["status"] = fileSystem.Status
			}

			if fileSystem.SuperUsers != nil {
				fileSystemMap["super_users"] = fileSystem.SuperUsers
			}

			if fileSystem.PosixAcl != nil {
				fileSystemMap["posix_acl"] = fileSystem.PosixAcl
			}

			if fileSystem.EnableRanger != nil {
				fileSystemMap["enable_ranger"] = fileSystem.EnableRanger
			}

			if fileSystem.RangerServiceAddresses != nil {
				fileSystemMap["ranger_service_addresses"] = fileSystem.RangerServiceAddresses
			}

			ids = append(ids, *fileSystem.FileSystemId)
			tmpList = append(tmpList, fileSystemMap)
		}

		_ = d.Set("file_systems", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
