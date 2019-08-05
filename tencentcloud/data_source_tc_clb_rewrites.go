/*
Use this data source to query detailed information of CLB rewrite

Example Usage

```hcl
data "tencentcloud_clb_rewrites" "clblab" {
  clb_id                = "lb-p7olt9e5"
  source_listener_id    = "lbl-jc1dx6ju#lb-p7olt9e5"
  target_listener_id    = "lbl-asj1hzuo#lb-p7olt9e5"
  rewrite_source_loc_id = "loc-ft8fmngv#lbl-jc1dx6ju#lb-p7olt9e5"
  rewrite_target_loc_id = "loc-4xxr2cy7#lbl-asj1hzuo#lb-p7olt9e5"
  result_output_file    = "mytestpath"
}
```
*/
package tencentcloud

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudClbRewrites() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbRewritesRead,

		Schema: map[string]*schema.Schema{
			"clb_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: " ID of the CLB to be queried.",
			},
			"source_listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of source listener.",
			},
			"target_listener_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Id of source listener.",
			},
			"rewrite_source_loc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of rule id of source listener.",
			},
			"rewrite_target_loc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Id of rule id of target listener.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"rewrite_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of cloud load redirection configurations. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"clb_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: " ID of the CLB to be queried.",
						},
						"source_listener_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of source listener.",
						},
						"target_listener_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of source listener.",
						},
						"rewrite_source_loc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of rule id of source listener.",
						},
						"rewrite_target_loc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of rule id of target listener.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudClbRewritesRead(d *schema.ResourceData, meta interface{}) error {
	defer LogElapsed("data_source.tencentcloud_clb_rewrites.read")()

	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	params := make(map[string]string)
	params["clb_id"] = d.Get("clb_id").(string)
	params["source_listener_id"] = strings.Split(d.Get("source_listener_id").(string), "#")[0]
	params["rewrite_source_loc_id"] = strings.Split(d.Get("rewrite_source_loc_id").(string), "#")[0]
	if v, ok := d.GetOk("target_listener_id"); ok {
		params["target_listener_id"] = strings.Split(v.(string), "#")[0]
	}
	if v, ok := d.GetOk("rewrite_target_loc_id"); ok {
		params["rewrite_target_loc_id"] = strings.Split(v.(string), "#")[0]
	}

	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	rewrites, err := clbService.DescribeRewriteInfosByFilter(ctx, params)
	if err != nil {
		return err
	}

	rewriteList := make([]map[string]interface{}, 0, len(rewrites))
	ids := make([]string, 0, len(rewrites))
	for _, rewrite := range rewrites {
		mapping := map[string]interface{}{
			"clb_id":                (*rewrite)["clb_id"],
			"source_listener_id":    (*rewrite)["source_listener_id"] + "#" + (*rewrite)["clb_id"],
			"target_listener_id":    (*rewrite)["target_listener_id"] + "#" + (*rewrite)["clb_id"],
			"rewrite_source_loc_id": (*rewrite)["rewrite_source_loc_id"] + "#" + (*rewrite)["source_listener_id"] + "#" + (*rewrite)["clb_id"],
			"rewrite_target_loc_id": (*rewrite)["rewrite_target_loc_id"] + "#" + (*rewrite)["target_listener_id"] + "#" + (*rewrite)["clb_id"],
		}

		rewriteList = append(rewriteList, mapping)
		ids = append(ids, (*rewrite)["rewrite_source_loc_id"]+"#"+(*rewrite)["rewrite_target_loc_id"]+(*rewrite)["source_listener_id"]+"#"+(*rewrite)["target_listener_id"]+"#"+(*rewrite)["clb_id"])
	}

	d.SetId(dataResourceIdsHash(ids))
	if err = d.Set("rewrite_list", rewriteList); err != nil {
		log.Printf("[CRITAL]%s provider set rewrite list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := writeToFile(output.(string), rewriteList); err != nil {
			return err
		}
	}

	return nil
}
