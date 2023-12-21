package monitor

import (
	"fmt"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func DataSourceTencentCloudMonitorProductNamespace() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentMonitorProductNamespaceRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name for filter, eg:`Load Banlancer`.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to store results.",
			},
			// Computed values
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list product namespaces. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"product_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "English name of this product.",
						},
						"product_chinese_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Chinese name of this product.",
						},
						"namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Namespace of each cloud product in monitor system.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentMonitorProductNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_monitor_product_namespace.read")()

	var (
		monitorService = MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request        = monitor.NewDescribeProductListRequest()
		response       *monitor.DescribeProductListResponse
		products       []*monitor.ProductSimple
		offset         uint64 = 0
		limit          uint64 = 20
		err            error
		filterName     = d.Get("name").(string)
	)

	request.Offset = &offset
	request.Limit = &limit
	request.Module = helper.String("monitor")

	var finish = false
	for {

		if finish {
			break
		}

		if err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			if response, err = monitorService.client.UseMonitorClient().DescribeProductList(request); err != nil {
				return tccommon.RetryError(err, tccommon.InternalError)
			}
			products = append(products, response.Response.ProductList...)
			if len(response.Response.ProductList) < int(limit) {
				finish = true
			}
			return nil
		}); err != nil {
			return err
		}
		offset = offset + limit
	}

	var list = make([]interface{}, 0, len(products))

	for _, product := range products {
		var listItem = map[string]interface{}{}
		listItem["product_name"] = product.ProductEnName
		listItem["product_chinese_name"] = product.ProductName
		listItem["namespace"] = product.Namespace
		if filterName == "" {
			list = append(list, listItem)
			continue
		}
		if product.ProductEnName != nil && strings.Contains(*product.ProductEnName, filterName) {
			list = append(list, listItem)
			continue
		}
		if product.ProductName != nil && strings.Contains(*product.ProductName, filterName) {
			list = append(list, listItem)
			continue
		}

	}
	if err = d.Set("list", list); err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("product_namespace_%s", filterName))
	if output, ok := d.GetOk("result_output_file"); ok {
		return tccommon.WriteToFile(output.(string), list)
	}
	return nil
}
