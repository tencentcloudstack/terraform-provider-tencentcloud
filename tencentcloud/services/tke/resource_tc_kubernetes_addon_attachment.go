package tke

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTkeAddonAttachment() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of cluster.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of addon.",
			},
			"version": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Description:   "Addon version, default latest version. Conflict with `request_body`.",
				ConflictsWith: []string{"request_body"},
			},
			"values": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				Description:   "Values the addon passthroughs. Conflict with `request_body`.",
				ConflictsWith: []string{"request_body"},
				Elem:          &schema.Schema{Type: schema.TypeString},
			},
			"raw_values": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Description:   "Raw Values. Conflict with `request_body`. Required with `raw_values_type`.",
				ConflictsWith: []string{"request_body"},
				RequiredWith:  []string{"raw_values_type"},
			},
			"raw_values_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The type of raw Values. Required with `raw_values`.",
				RequiredWith: []string{"raw_values"},
			},
			"request_body": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Serialized json string as request body of addon spec. If set, will ignore `version` and `values`.",
				ConflictsWith: []string{"version", "values"},
			},
			"response_body": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Addon response body.",
			},
			"status": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Addon current status.",
			},
		},
		Create: resourceTencentCloudTkeAddonAttachmentCreate,
		Update: resourceTencentCloudTkeAddonAttachmentUpdate,
		Read:   resourceTencentCloudTkeAddonAttachmentRead,
		Delete: resourceTencentCloudTkeAddonAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceTencentCloudTkeAddonAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.resource_tc_kubernetes_addon_attachment.create")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := TkeService{client: client}
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		clusterId     = d.Get("cluster_id").(string)
		addonName     = d.Get("name").(string)
		version       = d.Get("version").(string)
		values        = d.Get("values").([]interface{})
		rawValues     *string
		rawValuesType *string
		reqBody       = d.Get("request_body").(string)
	)

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

	err := service.CreateExtensionAddon(ctx, clusterId, reqBody)

	if err != nil {
		return err
	}

	d.SetId(clusterId + tccommon.FILED_SP + addonName)

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
		if err := resourceTencentCloudTkeAddonAttachmentDelete(d, meta); err != nil {
			return err
		}
		d.SetId("")
		return fmt.Errorf(msg)
	}

	return resourceTencentCloudTkeAddonAttachmentRead(d, meta)
}

func resourceTencentCloudTkeAddonAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.resource_tc_kubernetes_addon_attachment.read")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	id := d.Id()
	has := false
	split := strings.Split(id, tccommon.FILED_SP)
	if len(split) < 2 {
		return fmt.Errorf("id expected format: cluster_id#addon_name but no addon_name provided")
	}
	clusterId := split[0]
	addonName := split[1]

	var (
		err               error
		response          string
		addonResponseData = &AddonResponseData{}
	)

	_, has, err = service.PollingAddonsPhase(ctx, clusterId, addonName, addonResponseData)

	if err != nil || !has {
		d.SetId("")
		return err
	}

	response, _, err = service.DescribeExtensionAddon(ctx, clusterId, addonName)

	if err != nil {
		d.SetId("")
		return err
	}

	_ = d.Set("response_body", response)

	spec := addonResponseData.Spec
	statuses := addonResponseData.Status

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

	d.SetId(id)

	return nil
}

func resourceTencentCloudTkeAddonAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.resource_tc_kubernetes_addon_attachment.update")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
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

	return resourceTencentCloudTkeAddonAttachmentRead(d, meta)
}

func resourceTencentCloudTkeAddonAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.resource_tc_kubernetes_addon_attachment.delete")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		id        = d.Id()
		split     = strings.Split(id, tccommon.FILED_SP)
		clusterId = split[0]
		addonName = split[1]
		has       bool
	)

	if err := service.DeleteExtensionAddon(ctx, clusterId, addonName); err != nil {
		return err
	}

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
