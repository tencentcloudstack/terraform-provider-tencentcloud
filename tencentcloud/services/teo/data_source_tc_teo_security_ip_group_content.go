package teo

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func DataSourceTencentCloudTeoSecurityIPGroupContent() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTeoSecurityIPGroupContentRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site ID.",
			},

			"group_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "IP group ID.",
			},

			"ip_total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total count of IPs or CIDR blocks in the IP group.",
			},

			"ip_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of IPs or CIDR blocks in the IP group.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTeoSecurityIPGroupContentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_teo_security_ip_group_content.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient()

	request := teo.NewDescribeSecurityIPGroupContentRequest()

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("group_id"); ok {
		request.GroupId = helper.IntInt64(v.(int))
	}

	var ipList []*string
	var ipTotalCount int64
	var offset int64 = 0
	var limit int64 = 100000

	for {
		request.Offset = &offset
		request.Limit = &limit
		ratelimit.Check(request.GetAction())

		var response *teo.DescribeSecurityIPGroupContentResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			resp, e := client.DescribeSecurityIPGroupContent(request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			response = resp
			return nil
		})
		if err != nil {
			return err
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil {
			break
		}

		if response.Response.IPTotalCount != nil {
			ipTotalCount = *response.Response.IPTotalCount
		}

		if response.Response.IPList != nil {
			ipList = append(ipList, response.Response.IPList...)
		}

		if len(response.Response.IPList) < int(limit) {
			break
		}
		offset += limit
	}

	if ipTotalCount > 0 {
		_ = d.Set("ip_total_count", ipTotalCount)
	}

	if ipList != nil {
		_ = d.Set("ip_list", ipList)
	}

	d.SetId(helper.BuildToken())

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
