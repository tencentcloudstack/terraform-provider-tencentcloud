/*
Use this data source to query monitor data. for complex queries, use (https://github.com/tencentyun/tencentcloud-exporter)

Example Usage

```hcl
data "tencentcloud_instances" "instances" {
}

#cvm
data "tencentcloud_monitor_data" "cvm_monitor_data" {
  namespace   = "QCE/CVM"
  metric_name = "CPUUsage"
  dimensions {
    name  = "InstanceId"
    value = data.tencentcloud_instances.instances.instance_list[0].instance_id
  }
  period     = 300
  start_time = "2020-04-28T18:45:00+08:00"
  end_time   = "2020-04-28T19:00:00+08:00"
}

#cos
data "tencentcloud_monitor_data" "cos_monitor_data" {
  namespace   = "QCE/COS"
  metric_name = "InternetTraffic"
  dimensions {
    name  = "appid"
    value = "1258798060"
  }
  dimensions {
    name  = "bucket"
    value = "test-1258798060"
  }

  period     = 300
  start_time = "2020-04-28T18:30:00+08:00"
  end_time   = "2020-04-28T19:00:00+08:00"
}
```

*/
package tencentcloud

import (
	"crypto/md5"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func dataSourceTencentMonitorData() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentMonitorDataRead,
		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Namespace of each cloud product in monitor system, refer to `data.tencentcloud_monitor_product_namespace`.",
			},
			"metric_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Metric name, please refer to the documentation of monitor interface of each product.",
			},
			"dimensions": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instance dimension name, eg: `InstanceId` for cvm.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instance dimension value, eg: `ins-j0hk02zo` for cvm.",
						},
					},
				},
				Description: "Dimensional composition of instance objects.",
			},
			"period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: "Statistical period.",
			},
			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Start time for this query, eg:`2018-09-22T19:51:23+08:00`.",
			},
			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "End time for this query, eg:`2018-09-22T20:00:00+08:00`.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to store results.",
			},
			// Computed values
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list data point. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timestamp": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Statistical timestamp.",
						},
						"value": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Statistical value.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentMonitorDataRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_monitor_data.read")()

	var (
		monitorService = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		request        = monitor.NewGetMonitorDataRequest()
		response       *monitor.GetMonitorDataResponse
		err            error
		dimensions     = d.Get("dimensions").([]interface{})
		instance       monitor.Instance
		list           []interface{}
	)
	request.Namespace = helper.String(d.Get("namespace").(string))
	request.MetricName = helper.String(d.Get("metric_name").(string))
	request.Period = helper.IntUint64(d.Get("period").(int))
	request.StartTime = helper.String(d.Get("start_time").(string))
	request.EndTime = helper.String(d.Get("end_time").(string))
	request.Instances = []*monitor.Instance{&instance}

	for _, dimension := range dimensions {
		kv := dimension.(map[string]interface{})
		instance.Dimensions = append(instance.Dimensions, &monitor.Dimension{
			Name:  helper.String(kv["name"].(string)),
			Value: helper.String(kv["value"].(string)),
		})
	}

	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		if response, err = monitorService.client.UseMonitorClient().GetMonitorData(request); err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	if len(response.Response.DataPoints) > 0 {
		data := response.Response.DataPoints[0]
		min := len(data.Values)
		if min > len(data.Timestamps) {
			min = len(data.Timestamps)
		}
		for i := 0; i < min; i++ {
			kv := make(map[string]interface{})
			kv["timestamp"] = int64(*data.Timestamps[i])
			kv["value"] = data.Values[i]
			list = append(list, kv)
		}
	}

	if err = d.Set("list", list); err != nil {
		return err
	}

	md := md5.New()
	_, _ = md.Write([]byte(request.ToJsonString()))
	id := fmt.Sprintf("%x", md.Sum(nil))
	d.SetId(id)
	if output, ok := d.GetOk("result_output_file"); ok {
		return writeToFile(output.(string), list)
	}
	return nil
}
