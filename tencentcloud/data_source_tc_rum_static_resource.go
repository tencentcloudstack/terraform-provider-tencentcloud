/*
Use this data source to query detailed information of rum static_resource

Example Usage

```hcl
data "tencentcloud_rum_static_resource" "static_resource" {
  start_time = 1625444040
  type       = "top"
  end_time   = 1625454840
  project_id = 1
}
```
*/
package tencentcloud

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudRumStaticResource() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudRumStaticResourceRead,
		Schema: map[string]*schema.Schema{
			"start_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Start time but is represented using a timestamp in seconds.",
			},

			"type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Query Data Type. `top`: CostType Sorting order by top, `count40x`: CostType Sorting order by statusCode(40x) count, `nettype`: CostType Sorting order by nettype, `allcount`: CostType allcount, `day`: CostType group by day, `condition`: CostType Sorting in condition, `pagepv`: CostType Sorting order by pagepv, `area`: CostType query in area, `version`: CostType sort by version, `platform`: CostType sort by platform, `isp`: CostType sort by isp, `region`: region, `device`: CostType sort by device, `browser`: CostType sort by browser, `ext1`: CostType sort by ext1, `ext2`: CostType sort by ext2, `ext3`: CostType sort by ext3, `ret`: CostType sort by ret, `status`: CostType sort by status, `from`: CostType sort by from, `url`: CostType sort by url, `env`: CostType sort by env.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "End time but is represented using a timestamp in seconds.",
			},

			"project_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Project ID.",
			},

			"ext_second": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Second Expansion parameter.",
			},

			"engine": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The browser engine used for data reporting.",
			},

			"isp": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The internet service provider used for data reporting.",
			},

			"from": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The source page of the data reporting.",
			},

			"level": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Log level for data reporting(`1`: whitelist, `2`: normal, `4`: error, `8`: promise error, `16`: ajax request error, `32`: js resource load error, `64`: image resource load error, `128`: css resource load error, `256`: console.error, `512`: video resource load error, `1024`: request retcode error, `2048`: sdk self monitor error, `4096`: pv log, `8192`: event log).",
			},

			"brand": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The mobile phone brand used for data reporting.",
			},

			"area": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The region where the data reporting takes place.",
			},

			"version_num": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The SDK version used for data reporting.",
			},

			"platform": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The platform where the data reporting takes place.(`1`: Android, `2`: IOS, `3`: Windows, `4`: Mac, `5`: Linux, `100`: Other).",
			},

			"ext_third": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Third Expansion parameter.",
			},

			"ext_first": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "First Expansion parameter.",
			},

			"net_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The network type used for data reporting.(`1`: Wifi, `2`: 2G, `3`: 3G, `4`: 4G, `5`: 5G, `6`: 6G, `100`: Unknown).",
			},

			"device": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The device used for data reporting.",
			},

			"is_abroad": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Whether it is non-China region.`1`: yes; `0`: no.",
			},

			"os": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The operating system used for data reporting.",
			},

			"browser": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The browser type used for data reporting.",
			},

			"cost_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The method used for calculating the elapsed time `50`: 50th percentile, `75`: 75th percentile., `90`: 90th percentile., `95`: 95th percentile., `99`: 99th percentile., `99.5`: 99.5th percentile., `avg`: Mean.",
			},

			"url": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The URL Key where the data reporting takes place.",
			},

			"env": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The code environment where the data reporting takes place.(`production`: production env, `development`: development env, `gray`: gray env, `pre`: pre env, `daily`: daily env, `local`: local env, `others`: others env).",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Return value.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudRumStaticResourceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_rum_static_resource.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		startTime int
		endTime   int
	)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("start_time"); v != nil {
		startTime = v.(int)
		paramMap["StartTime"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("end_time"); v != nil {
		endTime = v.(int)
		paramMap["EndTime"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("type"); ok {
		paramMap["Type"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("project_id"); v != nil {
		paramMap["ID"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("ext_second"); ok {
		paramMap["ExtSecond"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("engine"); ok {
		paramMap["Engine"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("isp"); ok {
		paramMap["Isp"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("from"); ok {
		paramMap["From"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("level"); ok {
		paramMap["Level"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("brand"); ok {
		paramMap["Brand"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("area"); ok {
		paramMap["Area"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("version_num"); ok {
		paramMap["VersionNum"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("platform"); ok {
		paramMap["Platform"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ext_third"); ok {
		paramMap["ExtThird"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ext_first"); ok {
		paramMap["ExtFirst"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("net_type"); ok {
		paramMap["NetType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("device"); ok {
		paramMap["Device"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("is_abroad"); ok {
		paramMap["IsAbroad"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("os"); ok {
		paramMap["Os"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("browser"); ok {
		paramMap["Browser"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cost_type"); ok {
		paramMap["CostType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("url"); ok {
		paramMap["Url"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("env"); ok {
		paramMap["Env"] = helper.String(v.(string))
	}

	service := RumService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result *string
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeRumStaticResourceByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		result = response
		return nil
	})
	if err != nil {
		return err
	}

	var ids string
	if result != nil {
		ids = *result
		_ = d.Set("result", result)
	}

	d.SetId(helper.DataResourceIdsHash([]string{strconv.Itoa(startTime), strconv.Itoa(endTime), ids}))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
