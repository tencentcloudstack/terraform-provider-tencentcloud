/*
Use this data source to query detailed information of ssl managers

Example Usage

```hcl
data "tencentcloud_ssl_managers" "managers" {
  company_id = 1
  manager_name = "leader"
  manager_mail = "xx@x.com"
  status = "none"
  search_key = "xxx"
}
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudSslManagers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSslManagersRead,
		Schema: map[string]*schema.Schema{
			"company_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Company ID.",
			},

			"manager_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Manager&amp;#39;s name (will be abandoned), please use Searchkey.",
			},

			"manager_mail": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Vague query manager email (will be abandoned), please use Searchkey.",
			},

			"status": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filter according to the status of the manager, and the value is available&amp;#39;none&amp;#39; Unable to submit review&amp;#39;audit&amp;#39; Asian Credit Review&amp;#39;CAaudit&amp;#39; CA review&amp;#39;ok&amp;#39; has been reviewed&amp;#39;invalid&amp;#39;  审核失败&amp;#39;expiring&amp;#39;  is about to expire&amp;#39;expired&amp;#39; expired.",
			},

			"search_key": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Manager&amp;#39;s surname/Manager name/mailbox/department precise matching.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudSslManagersRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ssl_managers.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("company_id"); v != nil {
		paramMap["CompanyId"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("manager_name"); ok {
		paramMap["ManagerName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("manager_mail"); ok {
		paramMap["ManagerMail"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		paramMap["Status"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("search_key"); ok {
		paramMap["SearchKey"] = helper.String(v.(string))
	}

	service := SslService{client: meta.(*TencentCloudClient).apiV3Conn}

	var managers []*ssl.ManagerInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSslManagersByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		managers = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(managers))
	tmpList := make([]map[string]interface{}, 0, len(managers))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
