package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs/v20190719"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCfsFileSystems() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCfsFileSystemsRead,

		Schema: map[string]*schema.Schema{
			"file_system_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A specified file system ID used to query.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A file system name used to query.",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The available zone that the file system locates at.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the vpc to be queried.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of a vpc subnet.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			"file_system_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of cloud file system. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_system_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the file system.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the file system.",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The available zone that the file system locates at.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Protocol of the file system.",
						},
						"access_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the access group.",
						},
						"storage_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Storage type of the file system.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the file system.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of the file system.",
						},
						"size_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Size limit of the file system.",
						},
						"size_used": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Size used of the file system.",
						},
						"mount_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP of the file system.",
						},
						"fs_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Mount root-directory.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCfsFileSystemsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cfs_file_systems.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cfsService := CfsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	var fileSystemId string
	var vpcId string
	var subnetId string
	var name string
	var zone string
	if v, ok := d.GetOk("file_system_id"); ok {
		fileSystemId = v.(string)
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		vpcId = v.(string)
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		subnetId = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}
	if v, ok := d.GetOk("availability_zone"); ok {
		zone = v.(string)
	}

	var fileSystems []*cfs.FileSystemInfo
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		fileSystems, errRet = cfsService.DescribeFileSystem(ctx, fileSystemId, vpcId, subnetId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	fileSystemList := make([]map[string]interface{}, 0, len(fileSystems))
	ids := make([]string, 0, len(fileSystems))
	for _, fileSystem := range fileSystems {
		if name != "" && name != *fileSystem.FsName {
			continue
		}
		if zone != "" && zone != *fileSystem.Zone {
			continue
		}
		mapping := map[string]interface{}{
			"file_system_id":    fileSystem.FileSystemId,
			"name":              fileSystem.FsName,
			"availability_zone": fileSystem.Zone,
			"protocol":          fileSystem.Protocol,
			"access_group_id":   fileSystem.PGroup.PGroupId,
			"storage_type":      fileSystem.StorageType,
			"status":            fileSystem.LifeCycleState,
			"create_time":       fileSystem.CreationTime,
			"size_limit":        fileSystem.SizeLimit,
			"size_used":         fileSystem.SizeByte,
		}
		targets, err := cfsService.DescribeMountTargets(ctx, *fileSystem.FileSystemId)
		if err != nil {
			return err
		}
		var mountTarget *cfs.MountInfo
		if len(targets) > 0 {
			mountTarget = targets[0]
		}
		if mountTarget != nil {
			mapping["mount_ip"] = mountTarget.IpAddress
			mapping["fs_id"] = mountTarget.FSID
		}
		fileSystemList = append(fileSystemList, mapping)
		ids = append(ids, *fileSystem.FileSystemId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("file_system_list", fileSystemList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set cfs file system list fail, reason:%s\n ", logId, err.Error())
		return err
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := writeToFile(output.(string), fileSystemList); err != nil {
			return err
		}
	}
	return nil
}
