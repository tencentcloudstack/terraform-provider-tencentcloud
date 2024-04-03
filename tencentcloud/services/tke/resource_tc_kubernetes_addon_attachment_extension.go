package tke

import (
	"context"
	"fmt"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTencentCloudKubernetesAddonAttachmentCreatePreRequest0(d *schema.ResourceData, meta interface{}, req *tke.ForwardApplicationRequestV3Request) error {

	var (
		addonName     = d.Get("name").(string)
		version       = d.Get("version").(string)
		values        = d.Get("values").([]interface{})
		rawValues     *string
		rawValuesType *string
		reqBody       = d.Get("request_body").(string)
		service       = TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		ctx           = context.Background()
	)
	clusterName := *req.ClusterName
	if version == "" {
		request := tke.NewGetTkeAppChartListRequest()
		chartList, err := service.GetTkeAppChartList(ctx, request)
		if err != nil {
			return fmt.Errorf("error while fetching latest chart versions, %s", err.Error())
		}
		for i := range chartList {
			chart := chartList[i]
			if *chart.Name == addonName {
				version = *chart.LatestVersion
				break
			}
		}
	}

	if reqBody == "" {
		if v, ok := d.GetOk("raw_values"); ok {
			rawValues = helper.String(v.(string))
		}
		if v, ok := d.GetOk("raw_values_type"); ok {
			rawValuesType = helper.String(v.(string))
		}

		var reqErr error
		v := helper.InterfacesStringsPoint(values)
		reqBody, reqErr = service.GetAddonReqBody(addonName, version, v, rawValues, rawValuesType)
		if reqErr != nil {
			return reqErr
		}
	}
	req.RequestBody = &reqBody
	req.Path = helper.String(service.GetAddonsPath(clusterName, ""))
	req.Method = helper.String("POST")
	return nil
}

func resourceTencentCloudKubernetesAddonAttachmentCreatePostRequest0(ctx context.Context, d *schema.ResourceData, meta interface{}, req *tke.ForwardApplicationRequestV3Request, resp *tke.ForwardApplicationRequestV3Response) error {
	var (
		clusterId string
		name      string
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}
	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	addonName := name
	resData := &AddonResponseData{}
	reason := "unknown error"
	phase, has, _ := service.PollingAddonsPhase(ctx, clusterId, addonName, resData)

	if resData.Status != nil && resData.Status["reason"] != nil {
		reason = resData.Status["reason"].(string)
	}

	if !has {
		return fmt.Errorf("addon %s not exists", addonName)
	}

	if phase == "ChartFetchFailed" || phase == "Failed" || phase == "RollbackFailed" || phase == "SyncFailed" {
		msg := fmt.Sprintf("Unexpected chart phase `%s`, reason: %s", phase, reason)
		if err := resourceTencentCloudKubernetesAddonAttachmentDelete(d, meta); err != nil {
			return err
		}
		d.SetId("")
		return fmt.Errorf(msg)
	}
	return nil
}

func resourceTencentCloudKubernetesAddonAttachmentUpdatePreRequest0(ctx context.Context, d *schema.ResourceData, meta interface{}, req *tke.ForwardApplicationRequestV3Request) error {
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	var (
		version       = d.Get("version").(string)
		values        = d.Get("values").([]interface{})
		reqBody       = d.Get("request_body").(string)
		rawValues     *string
		rawValuesType *string
		err           error
	)
	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	clusterName := idSplit[0]
	addonName := idSplit[1]

	if d.HasChange("version") || d.HasChange("values") || d.HasChange("raw_values") || d.HasChange("raw_values_type") {
		if v, ok := d.GetOk("raw_values"); ok {
			rawValues = helper.String(v.(string))
		}
		if v, ok := d.GetOk("raw_values_type"); ok {
			rawValuesType = helper.String(v.(string))
		}
		reqBody, err = service.GetAddonReqBody(addonName, version, helper.InterfacesStringsPoint(values), rawValues, rawValuesType)
	}
	if err != nil {
		return err
	}
	req.Method = helper.String("PATCH")
	req.ContentType = helper.String("application/strategic-merge-patch+json")
	req.ClusterName = helper.String(clusterName)
	req.Path = helper.String(service.GetAddonsPath(clusterName, addonName))
	req.RequestBody = &reqBody
	return nil
}
func getFilteredValues(d *schema.ResourceData, values []*string) []string {
	rawValues := helper.InterfacesStrings(d.Get("values").([]interface{}))

	for _, value := range values {
		kv := strings.Split(*value, "=")
		key := kv[0]

		if tccommon.IsContains(TKE_ADDON_DEFAULT_VALUES_KEY, key) || tccommon.IsContains(rawValues, *value) {
			continue
		}
		rawValues = append(rawValues, *value)
	}
	return rawValues
}
