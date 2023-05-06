/*
Use this data source to query detailed information of cfs file_system_clients

Example Usage

```hcl
data "tencentcloud_cfs_file_system_clients" "file_system_clients" {
  file_system_id = "cfs-iobiaxtj"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs/v20190719"
)

func dataSourceTencentCloudCfsFileSystemClients() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCfsFileSystemClientsRead,
		Schema: map[string]*schema.Schema{
			"file_system_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "File system ID.",
			},

			"client_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Client list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cfs_vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP address of the file system.",
						},
						"client_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Client IP.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "File system VPCID.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the availability zone, e.g. ap-beijing-1. For more information, see regions and availability zones in the Overview document.",
						},
						"zone_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AZ name.",
						},
						"mount_directory": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Path in which the file system is mounted to the client.",
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

func dataSourceTencentCloudCfsFileSystemClientsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cfs_file_system_clients.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	fsId := d.Get("file_system_id").(string)

	service := CfsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var clientList []*cfs.FileSystemClient

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCfsFileSystemClientsById(ctx, fsId)
		if e != nil {
			return retryError(e)
		}
		clientList = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(clientList))

	if clientList != nil {
		for _, fileSystemClient := range clientList {
			fileSystemClientMap := map[string]interface{}{}

			if fileSystemClient.CfsVip != nil {
				fileSystemClientMap["cfs_vip"] = fileSystemClient.CfsVip
			}

			if fileSystemClient.ClientIp != nil {
				fileSystemClientMap["client_ip"] = fileSystemClient.ClientIp
			}

			if fileSystemClient.VpcId != nil {
				fileSystemClientMap["vpc_id"] = fileSystemClient.VpcId
			}

			if fileSystemClient.Zone != nil {
				fileSystemClientMap["zone"] = fileSystemClient.Zone
			}

			if fileSystemClient.ZoneName != nil {
				fileSystemClientMap["zone_name"] = fileSystemClient.ZoneName
			}

			if fileSystemClient.MountDirectory != nil {
				fileSystemClientMap["mount_directory"] = fileSystemClient.MountDirectory
			}

			tmpList = append(tmpList, fileSystemClientMap)
		}

		_ = d.Set("client_list", tmpList)
	}

	d.SetId(fsId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
