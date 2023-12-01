/*
Use this data source to query detailed information of lighthouse instance_disk_num

Example Usage

```hcl
data "tencentcloud_lighthouse_instance_disk_num" "instance_disk_num" {
  instance_ids = ["lhins-xxxxxx"]
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudLighthouseInstanceDiskNum() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudLighthouseInstanceDiskNumRead,
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of instance IDs.",
			},

			"attach_detail_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Mount information list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Id.",
						},
						"attached_disk_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of elastic cloud disks mounted to the instance.",
						},
						"max_attach_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of elastic cloud disks that can be mounted.",
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

func dataSourceTencentCloudLighthouseInstanceDiskNumRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_lighthouse_instance_disk_num.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	instanceIds := make([]string, 0)
	for _, instanceId := range d.Get("instance_ids").(*schema.Set).List() {
		instanceIds = append(instanceIds, instanceId.(string))
	}
	service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	var attachDetailSet []*lighthouse.AttachDetail

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeLighthouseInstanceDiskNum(ctx, instanceIds)
		if e != nil {
			return retryError(e)
		}
		attachDetailSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(attachDetailSet))
	tmpList := make([]map[string]interface{}, 0, len(attachDetailSet))

	if attachDetailSet != nil {
		for _, attachDetail := range attachDetailSet {
			attachDetailMap := map[string]interface{}{}

			if attachDetail.InstanceId != nil {
				attachDetailMap["instance_id"] = attachDetail.InstanceId
			}

			if attachDetail.AttachedDiskCount != nil {
				attachDetailMap["attached_disk_count"] = attachDetail.AttachedDiskCount
			}

			if attachDetail.MaxAttachCount != nil {
				attachDetailMap["max_attach_count"] = attachDetail.MaxAttachCount
			}

			ids = append(ids, *attachDetail.InstanceId)
			tmpList = append(tmpList, attachDetailMap)
		}

		_ = d.Set("attach_detail_set", tmpList)
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
