/*
Use this data source to query detailed information of HA VIP EIP attachments

Example Usage

```hcl
data "tencentcloud_ha_vip_eip_attachments" "foo" {
  havip_id     = "havip-kjqwe4ba"
  address_ip   = "1.1.1.1"
}
```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudHaVipEipAttachments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudHaVipEipAttachmentsRead,

		Schema: map[string]*schema.Schema{
			"havip_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of the attached HA VIP to be queried.",
			},
			"address_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Public IP address of EIP to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"ha_vip_eip_attachment_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of HA VIP EIP attachments. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"havip_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the attached HA VIP.",
						},
						"address_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public IP address of EIP.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudHaVipEipAttachmentsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ha_vip_eip_attachments.read")()

	logId := getLogId(contextNil)

	haVipId := d.Get("havip_id").(string)
	eip := d.Get("address_ip").(string)

	request := vpc.NewDescribeHaVipsRequest()
	request.HaVipIds = []*string{&haVipId}
	var response *vpc.DescribeHaVipsResponse
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DescribeHaVips(request)
		if e != nil {
			return retryError(errors.WithStack(e))
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read HA VIP failed, reason:%+v", logId, err)
		return err
	}

	haVipEipAttachmentList := make([]map[string]interface{}, 0, len(response.Response.HaVipSet))
	ids := make([]string, 0, len(response.Response.HaVipSet))
	for _, haVip := range response.Response.HaVipSet {
		if eip != "" {
			if *haVip.AddressIp != eip {
				continue
			}
		}
		mapping := map[string]interface{}{
			"havip_id":   haVipId,
			"address_ip": *haVip.AddressIp,
		}
		haVipEipAttachmentList = append(haVipEipAttachmentList, mapping)
		ids = append(ids, haVipId+"#"+*haVip.AddressIp)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("ha_vip_eip_attachment_list", haVipEipAttachmentList); e != nil {
		log.Printf("[CRITAL]%s provider set HA VIP EIP attachment list fail, reason:%s\n", logId, e)
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), haVipEipAttachmentList); e != nil {
			return e
		}
	}

	return nil
}
