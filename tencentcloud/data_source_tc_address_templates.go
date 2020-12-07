/*
Use this data source to query detailed information of address templates.

Example Usage

```hcl
data "tencentcloud_address_templates" "name" {
  name       = "test"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudAddressTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAddressTemplatesRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the address template to query.",
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Id of the address template to query.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// Computed values
			"template_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information list of the dedicated address templates.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the address template.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of address template.",
						},
						"addresses": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Set of the addresses.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudAddressTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_address_templates.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var name, templateId string
	var filters = make([]*vpc.Filter, 0)
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
		filters = append(filters, &vpc.Filter{Name: helper.String("address-template-name"), Values: []*string{&name}})
	}

	if v, ok := d.GetOk("id"); ok {
		templateId = v.(string)
		filters = append(filters, &vpc.Filter{Name: helper.String("address-template-id"), Values: []*string{&templateId}})
	}

	vpcService := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	var outErr, inErr error
	templates, outErr := vpcService.DescribeAddressTemplates(ctx, filters)
	if outErr != nil {
		outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			templates, inErr = vpcService.DescribeAddressTemplates(ctx, filters)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
	}

	if outErr != nil {
		return outErr
	}

	ids := make([]string, 0, len(templates))
	templateList := make([]map[string]interface{}, 0, len(templates))
	for _, ins := range templates {
		mapping := map[string]interface{}{
			"id":        ins.AddressTemplateId,
			"name":      ins.AddressTemplateName,
			"addresses": ins.AddressSet,
		}
		templateList = append(templateList, mapping)
		ids = append(ids, *ins.AddressTemplateId)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("template_list", templateList); e != nil {
		log.Printf("[CRITAL]%s provider set address template list fail, reason:%s\n", logId, e)
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), templateList); e != nil {
			return e
		}
	}

	return nil

}
