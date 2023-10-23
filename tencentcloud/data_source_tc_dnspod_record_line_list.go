/*
Use this data source to query detailed information of dnspod record_line_list

Example Usage

```hcl
data "tencentcloud_dnspod_record_line_list" "record_line_list" {
  domain = "iac-tf.cloud"
  domain_grade = "DP_FREE"
  domain_id = 123
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDnspodRecordLineList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDnspodRecordLineListRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Domain.",
			},

			"domain_grade": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Domain level. + Old packages: D_FREE, D_PLUS, D_EXTRA, D_EXPERT, D_ULTRA correspond to free package, personal luxury, enterprise 1, enterprise 2, enterprise 3. + New packages: DP_FREE, DP_PLUS, DP_EXTRA, DP_EXPERT, DP_ULTRA correspond to new free, personal professional, enterprise basic, enterprise standard, enterprise flagship.",
			},

			"domain_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Domain ID. The parameter DomainId has a higher priority than the parameter Domain. If the parameter DomainId is passed, the parameter Domain will be ignored. You can find all Domains and DomainIds through the DescribeDomainList interface.",
			},

			"line_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Line list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Line name.",
						},
						"line_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Line ID.",
						},
					},
				},
			},

			"line_group_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Line group list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"line_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Line group ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Line group name.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Group type.",
						},
						"line_list": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "List of lines included in the line group.",
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

func dataSourceTencentCloudDnspodRecordLineListRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dnspod_record_line_list.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		domain        string
		lineList      []*dnspod.LineInfo
		lineGroupList []*dnspod.LineGroupInfo
		e             error
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("domain"); ok {
		paramMap["Domain"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("domain_grade"); ok {
		paramMap["DomainGrade"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("domain_id"); ok {
		paramMap["DomainId"] = helper.IntUint64(v.(int))
	}

	service := DnspodService{client: meta.(*TencentCloudClient).apiV3Conn}

	// var lineList []*dnspod.LineInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		lineList, lineGroupList, e = service.DescribeDnspodRecordLineListByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		// lineList = result
		return nil
	})
	if err != nil {
		return err
	}

	// ids := make([]string, 0, len(lineList))
	if lineList != nil {
		tmpList := make([]map[string]interface{}, 0, len(lineList))
		for _, lineInfo := range lineList {
			lineInfoMap := map[string]interface{}{}

			if lineInfo.Name != nil {
				lineInfoMap["name"] = lineInfo.Name
			}

			if lineInfo.LineId != nil {
				lineInfoMap["line_id"] = lineInfo.LineId
			}

			// ids = append(ids, *lineInfo.Domain)
			tmpList = append(tmpList, lineInfoMap)
		}

		_ = d.Set("line_list", tmpList)
	}

	if lineGroupList != nil {
		tmpList := make([]map[string]interface{}, 0, len(lineGroupList))
		for _, lineGroupInfo := range lineGroupList {
			lineGroupInfoMap := map[string]interface{}{}

			if lineGroupInfo.LineId != nil {
				lineGroupInfoMap["line_id"] = lineGroupInfo.LineId
			}

			if lineGroupInfo.Name != nil {
				lineGroupInfoMap["name"] = lineGroupInfo.Name
			}

			if lineGroupInfo.Type != nil {
				lineGroupInfoMap["type"] = lineGroupInfo.Type
			}

			if lineGroupInfo.LineList != nil {
				lineGroupInfoMap["line_list"] = lineGroupInfo.LineList
			}

			// ids = append(ids, *lineGroupInfo.Domain)
			tmpList = append(tmpList, lineGroupInfoMap)
		}

		_ = d.Set("line_group_list", tmpList)
	}

	d.SetId(helper.DataResourceIdHash(domain))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		e = writeToFile(output.(string), map[string]interface{}{
			"line_list":       lineList,
			"line_group_list": lineGroupList,
		})
		if e != nil {
			return e
		}
	}
	return nil
}
