package tke

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func DataSourceTencentCloudTkeAddonValues() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTkeAddonValuesRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "ID of the TKE cluster, which can be obtained from the cluster list in the TKE console (https://console.cloud.tencent.com/tke2/cluster).",
			},

			"addon_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Name of the addon, which can be obtained from the component management page in the cluster details page of the TKE console.",
			},

			"values": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "Parameter list in JSON format. The value is the rendering result of the configured parameters merged with the default parameters of the chart. " +
					"Note: it is NOT empty for an addon that is not installed in the cluster but whose chart exists in the addon repository - " +
					"in that case the value is exactly the same as `default_values` and carries no cluster information " +
					"(e.g. `global.cluster.id` is an empty string). Use the `tencentcloud_kubernetes_addons` data source to check whether the addon is installed.",
			},

			"default_values": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Parameter list supported by the addon in JSON format, whose value is the default value of the chart.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTkeAddonValuesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tke_addon_values.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	var clusterId string
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}

	var addonName string
	if v, ok := d.GetOk("addon_name"); ok {
		addonName = v.(string)
	}

	var resp *tke.DescribeAddonValuesResponseParams
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeAddonValuesList(ctx, clusterId, addonName)
		if e != nil {
			// `not found chart` / `resource name may not be empty` are wrapped into
			// `ResourceUnavailable`, which is a *generic* retryable code (common.go:93).
			// They are deterministic failures (the chart does not exist, retrying will
			// never help), so they must NOT be retried - otherwise a simple typo in
			// `addon_name` wastes the whole 3-minute read timeout.
			if tccommon.IsDeterministicError(e) {
				log.Printf("[CRITAL] NonRetryable addon values error: %v", e)
				return resource.NonRetryableError(e)
			}
			return tccommon.RetryError(e)
		}

		if result == nil {
			log.Printf("[DATASOURCE] DescribeAddonValues returned empty result")
			return resource.NonRetryableError(fmt.Errorf("tencentcloud_tke_addon_values DescribeAddonValuesResponse is nil"))
		}

		resp = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	// A not-installed addon whose chart exists still yields a rendered payload here
	// (the chart defaults, byte-identical to DefaultValues and carrying no cluster
	// information). Pass it through as-is - see the schema description.
	if resp.Values != nil {
		_ = d.Set("values", *resp.Values)
	}

	if resp.DefaultValues != nil {
		_ = d.Set("default_values", *resp.DefaultValues)
	}

	d.SetId(strings.Join([]string{clusterId, addonName}, tccommon.FILED_SP))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		// `tccommon.WriteToFile` marshals the given value with `json.Marshal`.
		// Passing `d` (*schema.ResourceData) produces an empty object `{}` because
		// ResourceData has no exported fields and does not implement `MarshalJSON`.
		// Pass the actual payload instead - same style as
		// `data_source_tc_kubernetes_addons.go`.
		if e := tccommon.WriteToFile(output.(string), map[string]interface{}{
			"values":         d.Get("values"),
			"default_values": d.Get("default_values"),
		}); e != nil {
			return e
		}
	}

	return nil
}
