/*
Use this data source to query detailed information of CLB redirections

Example Usage

```hcl
data "tencentcloud_clb_redirections" "foo" {
  clb_id                 = "lb-p7olt9e5"
  source_listener_id     = "lbl-jc1dx6ju#lb-p7olt9e5"
  target_listener_id     = "lbl-asj1hzuo#lb-p7olt9e5"
  rewrite_source_rule_id = "loc-ft8fmngv#lbl-jc1dx6ju#lb-p7olt9e5"
  rewrite_target_rule_id = "loc-4xxr2cy7#lbl-asj1hzuo#lb-p7olt9e5"
  result_output_file     = "mytestpath"
}
```
*/
package tencentcloud

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
			"rewrite_source_rule_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule Id of source listener to be queried.",
			},
			"rewrite_target_rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Rule Id of target listener to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"redirection_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of cloud load redirection configurations. Each element contains the following attributes:",
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
						"rewrite_source_rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule Id of source listener.",
						},
						"rewrite_target_rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule Id of target listener.",
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
	params["source_listener_id"] = strings.Split(d.Get("source_listener_id").(string), "#")[0]
	params["rewrite_source_rule_id"] = strings.Split(d.Get("rewrite_source_rule_id").(string), "#")[0]
	if v, ok := d.GetOk("target_listener_id"); ok {
		params["target_listener_id"] = strings.Split(v.(string), "#")[0]
	}
	if v, ok := d.GetOk("rewrite_target_rule_id"); ok {
		params["rewrite_target_rule_id"] = strings.Split(v.(string), "#")[0]
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
		log.Printf("[CRITAL]%s read clb redirections failed, reason:%s\n ", logId, err.Error())
		return err
	}
	redirectionList := make([]map[string]interface{}, 0, len(redirections))
	ids := make([]string, 0, len(redirections))
	for _, rewrite := range redirections {
		mapping := map[string]interface{}{
			"clb_id":                 (*rewrite)["clb_id"],
			"source_listener_id":     (*rewrite)["source_listener_id"] + "#" + (*rewrite)["clb_id"],
			"target_listener_id":     (*rewrite)["target_listener_id"] + "#" + (*rewrite)["clb_id"],
			"rewrite_source_rule_id": (*rewrite)["rewrite_source_rule_id"] + "#" + (*rewrite)["source_listener_id"] + "#" + (*rewrite)["clb_id"],
			"rewrite_target_rule_id": (*rewrite)["rewrite_target_rule_id"] + "#" + (*rewrite)["target_listener_id"] + "#" + (*rewrite)["clb_id"],
		}

		redirectionList = append(redirectionList, mapping)
		ids = append(ids, (*rewrite)["rewrite_source_rule_id"]+"#"+(*rewrite)["rewrite_target_rule_id"]+(*rewrite)["source_listener_id"]+"#"+(*rewrite)["target_listener_id"]+"#"+(*rewrite)["clb_id"])
	}

	d.SetId(dataResourceIdsHash(ids))
	if e := d.Set("redirection_list", redirectionList); e != nil {
		log.Printf("[CRITAL]%s provider set redirection list fail, reason:%s\n ", logId, e.Error())
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
