/*
Provides an available EMR for the user.

The EMR data source fetch proper EMR from user's EMR pool.

Example Usage

```hcl
data "tencentcloud_emr" "my_emr" {
  filter {
    name   = "address-status"
    values = ["UNBIND"]
  }
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	emr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudEmr() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudEmrRead,

		Schema: map[string]*schema.Schema{
			"display_strategy": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display strategy(e.g.:clusterList, monitorManage)",
			},
			"prefix_instance_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "fetch all instances with same prefix(e.g.:emr-xxxxxx).",
			},
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Default 0 meaning first page",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Default 10, max value 100",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Fetch all instances which owner same project. Default 0 meaning use default project id",
			},
			"order_field": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Fetch all instances which owner same project. Default 0 meaning use default project id",
			},
			"asc": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Default 0 => descending order, 1 => ascending order",
			},
			"clusters": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of clusters will be exported and its every element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Id of instance",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster id of instance",
						},
						"ftitle": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Title of instance",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster name of instance",
						},
						"region_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "region id of instance",
						},
						"zone_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "region id of instance",
						},
						"app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "app id of instance",
						},
						"uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "user uin of instance",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "project id of instance",
						},
						"vpc_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "vpc id of instance",
						},
						"subnet_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "subnet id of instance",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "status of instance",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudEmrRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_emr.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	emrServer := EMRService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	filters := map[string]interface{}{}
	if v, ok := d.GetOk("display_strategy"); ok {
		filters["display_strategy"] = v.(string)
	}
	if v, ok := d.GetOk("prefix_instance_ids"); ok {
		filters["prefix_instance_ids"] = v.(string)
	}
	if v, ok := d.GetOk("offset"); ok {
		filters["offset"] = v.(uint64)
	}
	if v, ok := d.GetOk("limit"); ok {
		filters["limit"] = v.(uint64)
	}
	if v, ok := d.GetOk("project_id"); ok {
		filters["project_id"] = v.(int64)
	}
	if v, ok := d.GetOk("order_field"); ok {
		filters["order_field"] = v.(string)
	}
	if v, ok := d.GetOk("asc"); ok {
		filters["asc"] = v.(int64)
	}

	var clusters []*emr.ClusterInstancesInfo
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		clusters, errRet = emrServer.DescribeInstances(ctx, filters)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	emr_instances := make([]map[string]interface{}, 0, len(clusters))
	ids := make([]string, 0, len(clusters))

	for _, cluster := range clusters {
		mapping := map[string]interface{}{
			"cluster_id":   cluster.ClusterId,
			"ftitle":       cluster.Ftitle,
			"cluster_name": cluster.ClusterName,
			"region_id":    cluster.RegionId,
			"zone_id":      cluster.ZoneId,
		}
		emr_instances = append(emr_instances, mapping)
		ids = append(ids, (string)(*cluster.Id))
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("clusters", emr_instances)
	if err != nil {
		log.Printf("[CRITAL]%s provider set zones list fail, reason:%s\n ", logId, err.Error())
		return err
	}
	return nil
}
