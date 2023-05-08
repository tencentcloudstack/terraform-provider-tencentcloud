/*
Use this data source to query list information of api_gateway api_doc
Example Usage
```hcl
data "tencentcloud_api_gateway_api_docs" "my_api_doc" {
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudAPIGatewayAPIDocs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAPIGatewayAPIDocsRead,
		Schema: map[string]*schema.Schema{
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			"api_doc_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of ApiDocs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_doc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Api Doc ID.",
						},
						"api_doc_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Api Doc Name.",
						},
						"api_doc_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Api Doc Status.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudAPIGatewayAPIDocsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_api_gateway_api_docs.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		apiDoc            []*apigateway.APIDoc
	)

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := apiGatewayService.DescribeApiDocList(ctx)
		if e != nil {
			return retryError(e)
		}

		apiDoc = results
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read api_gateway apiDocs failed, reason:%+v", logId, err)
		return err
	}

	apiDocList := []interface{}{}
	ids := make([]string, 0, len(apiDoc))
	if apiDoc != nil {
		for _, item := range apiDoc {
			docMap := map[string]interface{}{}
			if item.ApiDocId != nil {
				docMap["api_doc_id"] = item.ApiDocId
			}
			if item.ApiDocName != nil {
				docMap["api_doc_name"] = item.ApiDocName
			}
			if item.ApiDocStatus != nil {
				docMap["api_doc_status"] = item.ApiDocStatus
			}
			apiDocList = append(apiDocList, docMap)
			ids = append(ids, *item.ApiDocId)
		}
		_ = d.Set("api_doc_list", apiDocList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), apiDocList); e != nil {
			return e
		}
	}

	return nil
}
