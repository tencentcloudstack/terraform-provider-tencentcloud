/*
Use this data source to query detailed information of CLB attachments

Example Usage

```hcl
data "tencentcloud_clb_attachments" "clblab" {
  listener_id = "lbl-hh141sn9"
  clb_id      = "lb-k2zjp9lv"
  rule_id     = "loc-4xxr2cy7"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
)

func dataSourceTencentCloudClbServerAttachments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbServerAttachmentsRead,

		Schema: map[string]*schema.Schema{
			"clb_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of the CLB to be queried.",
			},
			"listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of the CLB listener to be queried.",
			},
			"rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Id of the CLB listener rule. If the protocol of listener is HTTP/HTTPS, this para is required.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"attachment_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of cloud load balancer attachment configurations. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"clb_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the CLB.",
						},
						"listener_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the CLB listener.",
						},
						"protocol_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of protocol within the listener, and available values include 'TCP', 'UDP', 'HTTP', 'HTTPS' and 'TCP_SSL'.NOTES: TCP_SSL is testing internally, please apply if you need to use.",
						},
						"rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the CLB listener rule.",
						},
						"targets": {
							Type:        schema.TypeSet,
							Computed:    true,
							MinItems:    1,
							MaxItems:    100,
							Description: "Information of the backends to be attached.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Id of the backend server.",
									},
									"port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Port of the backend server.",
									},
									"weight": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Forwarding weight of the backend service, the range of [0, 100], defaults to 10.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudClbServerAttachmentsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_clb_attachments.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	params := make(map[string]string)
	clbId := d.Get("clb_id").(string)
	listenerId := d.Get("listener_id").(string)
	locationId := ""
	if v, ok := d.GetOk("rule_id"); ok {
		locationId = v.(string)
		params["rule_id"] = locationId
	}
	params["clb_id"] = clbId
	params["listener_id"] = listenerId

	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var attachments []*clb.ListenerBackend
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := clbService.DescribeAttachmentsByFilter(ctx, params)
		if e != nil {
			return retryError(e)
		}
		attachments = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CLB attachments failed, reason:%+v", logId, err)
		return err
	}
	attachmentList := make([]map[string]interface{}, 0, len(attachments))
	ids := make([]string, 0, len(attachments))
	for _, attachment := range attachments {
		mapping := map[string]interface{}{
			"clb_id":        clbId,
			"listener_id":   listenerId,
			"protocol_type": attachment.Protocol,
		}
		if locationId != "" {
			mapping["rule_id"] = locationId
		}
		if *attachment.Protocol == CLB_LISTENER_PROTOCOL_HTTP || *attachment.Protocol == CLB_LISTENER_PROTOCOL_HTTPS {
			if len(attachment.Rules) > 0 {
				for _, loc := range attachment.Rules {
					if locationId == "" || locationId == *loc.LocationId {
						mapping["targets"] = flattenBackendList(loc.Targets)
					}
				}
			}
		} else if *attachment.Protocol == CLB_LISTENER_PROTOCOL_TCP || *attachment.Protocol == CLB_LISTENER_PROTOCOL_UDP {
			mapping["targets"] = flattenBackendList(attachment.Targets)
		}
		attachmentList = append(attachmentList, mapping)
		ids = append(ids, locationId+"#"+listenerId+"#"+clbId)
	}

	d.SetId(dataResourceIdsHash(ids))
	if e := d.Set("attachment_list", attachmentList); e != nil {
		log.Printf("[CRITAL]%s provider set attachment list fail, reason:%+v", logId, e)
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), attachmentList); e != nil {
			return e
		}
	}

	return nil
}
