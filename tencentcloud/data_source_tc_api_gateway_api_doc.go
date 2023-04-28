/*
Use this data source to query detailed information of api_gateway api_doc

Example Usage

```hcl
data "tencentcloud_api_gateway_api_doc" "my_api_doc" {
  api_doc_id = "apidoc-a7ragqam"
}
```
*/

package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceTencentCloudAPIGatewayAPIDoc() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAPIGatewayAPIDocRead,
		Schema: map[string]*schema.Schema{
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudAPIGatewayAPIDocRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tat_command.read")()
	defer inconsistentCheck(d, meta)()

	//var (
	//	logId := getLogId(contextNil)
	//	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	//)
	//err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
	//	results, e := tatService.DescribeTatCommandByFilter(ctx, paramMap)
	//	if e != nil {
	//		return retryError(e)
	//	}
	//	commandSet = results
	//	return nil
	//})
	//if err != nil {
	//	log.Printf("[CRITAL]%s read Tat commandSet failed, reason:%+v", logId, err)
	//	return err
	//}
	return nil
}
