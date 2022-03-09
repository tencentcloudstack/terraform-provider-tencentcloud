/*
Provides an available EMR for the user.

The EMR data source fetch proper EMR from user's EMR pool.

Example Usage

```hcl
data "tencentcloud_emr" "my_emr" {
  display_strategy="clusterList"
  instance_ids=["emr-rnzqrleq"]
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
				Description: "Display strategy(e.g.:clusterList, monitorManage).",
			},
			"instance_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "fetch all instances with same prefix(e.g.:emr-xxxxxx).",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Fetch all instances which owner same project. Default 0 meaning use default project id.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
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
							Description: "Id of instance.",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster id of instance.",
						},
						"ftitle": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Title of instance.",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster name of instance.",
						},
						"region_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Region id of instance.",
						},
						"zone_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Zone id of instance.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Zone of instance.",
						},
						"master_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Master ip of instance.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Project id of instance.",
						},
						"charge_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Charge type of instance.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Status of instance.",
						},
						"add_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Add time of instance.",
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
	if v, ok := d.GetOk("instance_ids"); ok {
		filters["instance_ids"] = v.([]interface{})
	}
	if v, ok := d.GetOk("project_id"); ok {
		filters["project_id"] = v.(int64)
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
			"cluster_name": cluster.ClusterName,
			"ftitle":       cluster.Ftitle,
			"status":       cluster.Status,
			"region_id":    cluster.RegionId,
			"zone_id":      cluster.ZoneId,
			"zone":         cluster.Zone,
			"charge_type":  cluster.ChargeType,
			"master_ip":    cluster.MasterIp,
			"add_time":     cluster.AddTime,
		}
		emr_instances = append(emr_instances, mapping)
		ids = append(ids, (string)(*cluster.Id))
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("clusters", emr_instances)
	if err != nil {
		log.Printf("[CRITAL]%s provider set cluster list fail, reason:%s\n ", logId, err.Error())
		return err
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := writeToFile(output.(string), emr_instances); err != nil {
			return err
		}
	}
	return nil
}
