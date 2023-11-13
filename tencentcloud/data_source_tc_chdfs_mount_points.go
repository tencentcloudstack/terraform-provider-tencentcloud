/*
Use this data source to query detailed information of chdfs mount_points

Example Usage

```hcl
data "tencentcloud_chdfs_mount_points" "mount_points" {
  file_system_id = &lt;nil&gt;
  access_group_id = &lt;nil&gt;
  owner_uin = &lt;nil&gt;
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

func dataSourceTencentCloudChdfsMountPoints() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudChdfsMountPointsRead,
		Schema: map[string]*schema.Schema{
			"file_system_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Get mount points belongs to file system id, only can use one of the AccessGroupId,FileSystemId,OwnerUin parameters.",
			},

			"access_group_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Get mount points belongs to access group id, only can use one of the AccessGroupId,FileSystemId,OwnerUin parameters.",
			},

			"owner_uin": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Get mount points belongs to owner uin, only can use one of the AccessGroupId,FileSystemId,OwnerUin parameters.",
			},

			"mount_points": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Mount point list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mount_point_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Mount point id.",
						},
						"mount_point_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Mount point name.",
						},
						"file_system_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Mounted file system id.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Mount point status.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time.",
						},
						"access_group_ids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Associated group ids.",
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

func dataSourceTencentCloudChdfsMountPointsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_chdfs_mount_points.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("file_system_id"); ok {
		paramMap["FileSystemId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("access_group_id"); ok {
		paramMap["AccessGroupId"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("owner_uin"); v != nil {
		paramMap["OwnerUin"] = helper.IntUint64(v.(int))
	}

	service := ChdfsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var mountPoints []*chdfs.MountPoint

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeChdfsMountPointsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		mountPoints = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(mountPoints))
	tmpList := make([]map[string]interface{}, 0, len(mountPoints))

	if mountPoints != nil {
		for _, mountPoint := range mountPoints {
			mountPointMap := map[string]interface{}{}

			if mountPoint.MountPointId != nil {
				mountPointMap["mount_point_id"] = mountPoint.MountPointId
			}

			if mountPoint.MountPointName != nil {
				mountPointMap["mount_point_name"] = mountPoint.MountPointName
			}

			if mountPoint.FileSystemId != nil {
				mountPointMap["file_system_id"] = mountPoint.FileSystemId
			}

			if mountPoint.Status != nil {
				mountPointMap["status"] = mountPoint.Status
			}

			if mountPoint.CreateTime != nil {
				mountPointMap["create_time"] = mountPoint.CreateTime
			}

			if mountPoint.AccessGroupIds != nil {
				mountPointMap["access_group_ids"] = mountPoint.AccessGroupIds
			}

			ids = append(ids, *mountPoint.MountPointId)
			tmpList = append(tmpList, mountPointMap)
		}

		_ = d.Set("mount_points", tmpList)
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
