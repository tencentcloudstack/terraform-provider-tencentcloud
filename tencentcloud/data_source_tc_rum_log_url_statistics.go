/*
Use this data source to query detailed information of rum log_url_statistics

Example Usage

```hcl
data "tencentcloud_rum_log_url_statistics" "log_url_statistics" {
  start_time = 1625444040
  type = "analysis"
  end_time = 1625454840
  i_d = 1
  ext_second = "ext2"
  engine = "Blink(79.0)"
  isp = "中国电信"
  from = "https://user.qzone.qq.com/"
  level = "1"
  brand = "Apple"
  area = "广州市"
  version_num = "1.0"
  platform = "2"
  ext_third = "ext3"
  ext_first = "ext1"
  net_type = "2"
  device = "Apple - iPhone"
  is_abroad = "0"
  os = "Windows - 10"
  browser = "Chrome(79.0)"
  env = "production"
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudRumLogUrlStatistics() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudRumLogUrlStatisticsRead,
		Schema: map[string]*schema.Schema{
			"start_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Start time but is represented using a timestamp in seconds.",
			},

			"type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Query Data Type. `analysis`:&amp;amp;#39;query analysis data&amp;amp;#39;, `compare`:&amp;amp;#39;query compare data&amp;amp;#39;, `allcount`:&amp;amp;#39;query allcount&amp;amp;#39;, `condition`:&amp;amp;#39;query in condition&amp;amp;#39;, `nettype`: &amp;amp;#39;CostType sort by nettype&amp;amp;#39; , `version`: &amp;amp;#39;CostType sort by version&amp;amp;#39; , `platform`: &amp;amp;#39;CostType sort by platform&amp;amp;#39; , `isp`: &amp;amp;#39;CostType sort by isp&amp;amp;#39; , `region`: &amp;amp;#39;CostType sort by region&amp;amp;#39; , `device`: &amp;amp;#39;CostType sort by device&amp;amp;#39; , `browser`: &amp;amp;#39;CostType sort by browser&amp;amp;#39; , `ext1`: &amp;amp;#39;CostType sort by ext1&amp;amp;#39; , `ext2`: &amp;amp;#39;CostType sort by ext2&amp;amp;#39; , `ext3`: &amp;amp;#39;CostType sort by ext3&amp;amp;#39; , `ret`: &amp;amp;#39;CostType sort by ret&amp;amp;#39; , `status`: &amp;amp;#39;CostType sort by status&amp;amp;#39; , `from`: &amp;amp;#39;CostType sort by from&amp;amp;#39; , `url`: &amp;amp;#39;CostType sort by url&amp;amp;#39; , `env`: &amp;amp;#39;CostType sort by env&amp;amp;#39;.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "End time but is represented using a timestamp in seconds.",
			},

			"i_d": {
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
				Description: "Log level for data reporting(`1`: &amp;amp;#39;whitelist&amp;amp;#39;, `2`: &amp;amp;#39;normal&amp;amp;#39;, `4`: &amp;amp;#39;error&amp;amp;#39;, `8`: &amp;amp;#39;promise error&amp;amp;#39;, `16`: &amp;amp;#39;ajax request error&amp;amp;#39;, `32`: &amp;amp;#39;js resource load error&amp;amp;#39;, `64`: &amp;amp;#39;image resource load error&amp;amp;#39;, `128`: &amp;amp;#39;css resource load error&amp;amp;#39;, `256`: &amp;amp;#39;console.error&amp;amp;#39;, `512`: &amp;amp;#39;video resource load error&amp;amp;#39;, `1024`: &amp;amp;#39;request retcode error&amp;amp;#39;, `2048`: &amp;amp;#39;sdk self monitor error&amp;amp;#39;, `4096`: &amp;amp;#39;pv log&amp;amp;#39;, `8192`: &amp;amp;#39;event log&amp;amp;#39;).",
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
				Description: "The platform where the data reporting takes place.(`1`: &amp;amp;#39;Android&amp;amp;#39;, `2`: &amp;amp;#39;IOS&amp;amp;#39;, `3`: &amp;amp;#39;Windows&amp;amp;#39;, `4`: &amp;amp;#39;Mac&amp;amp;#39;, `5`: &amp;amp;#39;Linux&amp;amp;#39;, `100`: &amp;amp;#39;Other&amp;amp;#39;).",
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
				Description: "The network type used for data reporting.(`1`: &amp;amp;#39;Wifi&amp;amp;#39;, `2`: &amp;amp;#39;2G&amp;amp;#39;, `3`: &amp;amp;#39;3G&amp;amp;#39;, `4`: &amp;amp;#39;4G&amp;amp;#39;, `5`: &amp;amp;#39;5G&amp;amp;#39;, `6`: &amp;amp;#39;6G&amp;amp;#39;, `100`: &amp;amp;#39;Unknown&amp;amp;#39;).",
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

			"env": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The code environment where the data reporting takes place.(`production`: &amp;amp;#39;production env&amp;amp;#39;, `development`: &amp;amp;#39;development env&amp;amp;#39;, `gray`: &amp;amp;#39;gray env&amp;amp;#39;, `pre`: &amp;amp;#39;pre env&amp;amp;#39;, `daily`: &amp;amp;#39;daily env&amp;amp;#39;, `local`: &amp;amp;#39;local env&amp;amp;#39;, `others`: &amp;amp;#39;others env&amp;amp;#39;).",
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

func dataSourceTencentCloudRumLogUrlStatisticsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_rum_log_url_statistics.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("start_time"); v != nil {
		paramMap["StartTime"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("type"); ok {
		paramMap["Type"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("end_time"); v != nil {
		paramMap["EndTime"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("i_d"); v != nil {
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

	if v, ok := d.GetOk("env"); ok {
		paramMap["Env"] = helper.String(v.(string))
	}

	service := RumService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeRumLogUrlStatisticsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		result = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result))
	if result != nil {
		_ = d.Set("result", result)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string)); e != nil {
			return e
		}
	}
	return nil
}
