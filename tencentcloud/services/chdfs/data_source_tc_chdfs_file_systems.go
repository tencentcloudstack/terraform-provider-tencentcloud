package chdfs

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	chdfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/chdfs/v20201112"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudChdfsFileSystems() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudChdfsFileSystemsRead,
		Schema: map[string]*schema.Schema{
			"file_systems": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "file system list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "appid of the user.",
						},
						"file_system_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "file system name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "desc of the file system.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "region of the file system.",
						},
						"file_system_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "file system id.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "create time.",
						},
						"block_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "block size of the file system(byte).",
						},
						"capacity_quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "capacity of the file system(byte).",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "status of the file system(1: creating create success 3: create failed).",
						},
						"super_users": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "super users of the file system.",
						},
						"posix_acl": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "check POSIX ACL or not.",
						},
						"enable_ranger": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "check the ranger address or not.",
						},
						"ranger_service_addresses": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "ranger address list.",
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
	defer tccommon.LogElapsed("data_source.tencentcloud_chdfs_file_systems.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ChdfsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var fileSystems []*chdfs.FileSystem

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeChdfsFileSystems(ctx)
		if e != nil {
			return tccommon.RetryError(e)
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
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
