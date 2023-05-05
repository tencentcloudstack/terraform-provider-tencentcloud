/*
Use this data source to query detailed information of dayu DDoS policy attachments

Example Usage

```hcl
data "tencentcloud_dayu_ddos_policy_attachments" "foo_type" {
  resource_type = tencentcloud_dayu_ddos_policy_attachment.dayu_ddos_policy_attachment.resource_type
}
data "tencentcloud_dayu_ddos_policy_attachments" "foo_resource" {
  resource_id   = tencentcloud_dayu_ddos_policy_attachment.dayu_ddos_policy_attachment.resource_id
  resource_type = tencentcloud_dayu_ddos_policy_attachment.dayu_ddos_policy_attachment.resource_type
}
data "tencentcloud_dayu_ddos_policy_attachments" "foo_policy" {
  resource_type = tencentcloud_dayu_ddos_policy_attachment.dayu_ddos_policy_attachment.resource_type
  policy_id     = tencentcloud_dayu_ddos_policy_attachment.dayu_ddos_policy_attachment.policy_id
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDayuDdosPolicyAttachments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDayuDdosPolicyAttachmentsRead,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the attached resource to be queried.",
			},
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(DAYU_RESOURCE_TYPE),
				Description:  "Type of the resource that the DDoS policy works for, valid values are `bgpip`, `bgp`, `bgp-multip` and `net`.",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Id of the policy to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"dayu_ddos_policy_attachment_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of dayu DDoS policy attachments. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the attached resource.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the resource that the DDoS policy works for.",
						},
						"policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the policy.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudDayuDdosPolicyAttachmentsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dayu_ddos_policy_attachments.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	resourceId := d.Get("resource_id").(string)
	resourceType := d.Get("resource_type").(string)
	policyId := d.Get("policy_id").(string)

	dayuService := DayuService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	attachments, _, err := dayuService.DescribeDdosPolicyAttachments(ctx, resourceId, resourceType, policyId)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			attachments, _, err = dayuService.DescribeDdosPolicyAttachments(ctx, resourceId, resourceType, policyId)
			if err != nil {
				return retryError(err, "ClientError.NetworkError")
			}
			return nil
		})
	}
	if err != nil {
		log.Printf("[CRITAL]%s read DDoS Policy attachment failed, reason:%s\n", logId, err)
		return err
	}

	ddosPolicyAttachmentList := make([]map[string]interface{}, 0, len(attachments))
	ids := make([]string, 0, len(attachments))
	for _, attachment := range attachments {
		ddosPolicyAttachmentList = append(ddosPolicyAttachmentList, attachment)
		ids = append(ids, resourceId+FILED_SP+resourceType+""+FILED_SP+policyId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("dayu_ddos_policy_attachment_list", ddosPolicyAttachmentList); e != nil {
		log.Printf("[CRITAL]%s provider set DDoS policy attachment list fail, reason:%s\n", logId, e)
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), ddosPolicyAttachmentList); e != nil {
			return e
		}
	}

	return nil
}
