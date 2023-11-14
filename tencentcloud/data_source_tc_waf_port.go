/*
Use this data source to query detailed information of waf port

Example Usage

```hcl
data "tencentcloud_waf_port" "port" {
  edition = ""
  instance_i_d = ""
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

func dataSourceTencentCloudWafPort() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWafPortRead,
		Schema: map[string]*schema.Schema{
			"edition": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Instance type, sparta-waf represents SAAS WAF, clb-waf represents CLB WAF.",
			},

			"instance_i_d": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Instance unique ID.",
			},

			"http_ports": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Http port list for instance.",
			},

			"https_ports": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Https port list for instance.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudWafPortRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_waf_port.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("edition"); ok {
		paramMap["Edition"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_i_d"); ok {
		paramMap["InstanceID"] = helper.String(v.(string))
	}

	service := WafService{client: meta.(*TencentCloudClient).apiV3Conn}

	var httpPorts []*string

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWafPortByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		httpPorts = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(httpPorts))
	if httpPorts != nil {
		_ = d.Set("http_ports", httpPorts)
	}

	if httpsPorts != nil {
		_ = d.Set("https_ports", httpsPorts)
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
