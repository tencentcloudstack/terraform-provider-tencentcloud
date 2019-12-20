/*
Use this data source to query detailed information of CLB redirections

Example Usage

```hcl
data "tencentcloud_clb_redirections" "foo" {
  clb_id             = "lb-p7olt9e5"
  source_listener_id = "lbl-jc1dx6ju"
  target_listener_id = "lbl-asj1hzuo"
  source_rule_id     = "loc-ft8fmngv"
  target_rule_id     = "loc-4xxr2cy7"
  result_output_file = "mytestpath"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceTencentCloudClbRedirections() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbRedirectionsRead,

		Schema: map[string]*schema.Schema{
			"clb_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of the CLB to be queried.",
			},
			"source_listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of source listener to be queried.",
			},
			"target_listener_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Id of target listener to be queried.",
			},
			"source_rule_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule id of source listener to be queried.",
			},
			"target_rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Rule id of target listener to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"redirection_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of cloud load balancer redirection configurations. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"clb_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the CLB.",
						},
						"source_listener_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of source listener.",
						},
						"target_listener_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of target listener.",
						},
						"source_rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule id of source listener.",
						},
						"target_rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule id of target listener.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudClbRedirectionsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_clb_redirections.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	params := make(map[string]string)
	params["clb_id"] = d.Get("clb_id").(string)
	params["source_listener_id"] = d.Get("source_listener_id").(string)
	params["source_rule_id"] = d.Get("source_rule_id").(string)
	if v, ok := d.GetOk("target_listener_id"); ok {
		params["target_listener_id"] = v.(string)
	}
	if v, ok := d.GetOk("target_rule_id"); ok {
		params["target_rule_id"] = v.(string)
	}

	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var redirections []*map[string]string
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := clbService.DescribeRedirectionsByFilter(ctx, params)
		if e != nil {
			return retryError(e)
		}
		redirections = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CLB redirections failed, reason:%+v", logId, err)
		return err
	}
	redirectionList := make([]map[string]interface{}, 0, len(redirections))
	ids := make([]string, 0, len(redirections))
	for _, r := range redirections {
		mapping := map[string]interface{}{
			"clb_id":             (*r)["clb_id"],
			"source_listener_id": (*r)["source_listener_id"],
			"target_listener_id": (*r)["target_listener_id"],
			"source_rule_id":     (*r)["source_rule_id"],
			"target_rule_id":     (*r)["target_rule_id"],
		}

		redirectionList = append(redirectionList, mapping)
		ids = append(ids, (*r)["source_rule_id"]+"#"+(*r)["target_rule_id"]+(*r)["source_listener_id"]+"#"+(*r)["target_listener_id"]+"#"+(*r)["clb_id"])
	}

	d.SetId(dataResourceIdsHash(ids))
	if e := d.Set("redirection_list", redirectionList); e != nil {
		log.Printf("[CRITAL]%s provider set CLB redirection list fail, reason:%+v", logId, e)
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), redirectionList); e != nil {
			return e
		}
	}

	return nil
}
