package tke

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var addonResponseData = &AddonResponseData{}

func resourceTencentCloudKubernetesAddonAttachmentCreatePostFillRequest0(ctx context.Context, req *tke.ForwardApplicationRequestV3Request) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)

	var (
		addonName     = d.Get("name").(string)
		version       = d.Get("version").(string)
		values        = d.Get("values").([]interface{})
		rawValues     *string
		rawValuesType *string
		reqBody       = d.Get("request_body").(string)
		service       = TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
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

func resourceTencentCloudKubernetesAddonAttachmentCreatePostHandleResponse0(ctx context.Context, resp *tke.ForwardApplicationRequestV3Response) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
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
func resourceTencentCloudKubernetesAddonAttachmentReadPreRequest0(ctx context.Context, req *tke.ForwardApplicationRequestV3Request) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	var (
		err error
	)
	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	has := false
	clusterName := d.Get("cluster_id").(string)
	addonName := d.Get("name").(string)

	_, has, err = service.PollingAddonsPhase(ctx, clusterName, addonName, addonResponseData)

	if err != nil || !has {
		d.SetId("")
		return err
	}

	req.Method = helper.String("GET")
	req.ClusterName = &clusterName
	req.Path = helper.String(service.GetAddonsPath(clusterName, addonName))

	return nil
}

func resourceTencentCloudKubernetesAddonAttachmentReadPostHandleResponse0(ctx context.Context, resp *tke.ForwardApplicationRequestV3ResponseParams) error {
	d := tccommon.ResourceDataFromContext(ctx)

	spec := addonResponseData.Spec
	statuses := addonResponseData.Status
	clusterId := d.Get("cluster_id").(string)

	if spec != nil {
		_ = d.Set("cluster_id", clusterId)
		_ = d.Set("name", spec.Chart.ChartName)
		_ = d.Set("version", spec.Chart.ChartVersion)
		if spec.Values != nil && len(spec.Values.Values) > 0 {

			// Filter auto-filled values from addon creation
			filteredValues := getFilteredValues(d, spec.Values.Values)
			_ = d.Set("values", filteredValues)
		}

		if spec.Values != nil && spec.Values.RawValues != nil {
			rawValues := spec.Values.RawValues
			rawValuesType := spec.Values.RawValuesType

			base64DecodeValues, _ := base64.StdEncoding.DecodeString(*rawValues)
			jsonValues := string(base64DecodeValues)

			_ = d.Set("raw_values", jsonValues)
			_ = d.Set("raw_values_type", rawValuesType)
		}
	}

	if statuses != nil || len(statuses) == 0 {
		strMap := helper.CovertInterfaceMapToStrPtr(statuses)
		err := d.Set("status", strMap)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceTencentCloudKubernetesAddonAttachmentUpdateOnExit(ctx context.Context) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)

	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var (
		id            = d.Id()
		split         = strings.Split(id, tccommon.FILED_SP)
		clusterId     = split[0]
		addonName     = split[1]
		version       = d.Get("version").(string)
		values        = d.Get("values").([]interface{})
		reqBody       = d.Get("request_body").(string)
		err           error
		rawValues     *string
		rawValuesType *string
	)

	if d.HasChange("request_body") && reqBody == "" || d.HasChange("version") || d.HasChange("values") || d.HasChange("raw_values") || d.HasChange("raw_values_type") {
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

	err = service.UpdateExtensionAddon(ctx, clusterId, addonName, reqBody)

	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudKubernetesAddonAttachmentDeletePostFillRequest0(ctx context.Context, req *tke.ForwardApplicationRequestV3Request) error {

	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	clusterName := d.Get("cluster_id").(string)
	addonName := d.Get("name").(string)
	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	req.Method = helper.String("DELETE")
	req.ClusterName = &clusterName
	req.Path = helper.String(service.GetAddonsPath(clusterName, addonName))
	return nil
}
func resourceTencentCloudKubernetesAddonAttachmentDeletePostHandleResponse0(ctx context.Context, resp *tke.ForwardApplicationRequestV3Response) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	var (
		id        = d.Id()
		split     = strings.Split(id, tccommon.FILED_SP)
		clusterId = split[0]
		addonName = split[1]
		has       bool
	)
	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	// check if addon terminating or still exists
	_, has, _ = service.PollingAddonsPhase(ctx, clusterId, addonName, nil)

	if has {
		return fmt.Errorf("addon %s still exists", addonName)
	}
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
