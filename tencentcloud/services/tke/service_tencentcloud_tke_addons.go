package tke

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type AppType string
type RawValuesType string
type AppPhase string

type App struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              AddonSpec `json:"spec,omitempty" `
	Status            AppStatus `json:"status,omitempty"`
}

type AppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []App `json:"items"`
}

type AppStatus struct {
	Phase              AppPhase    `json:"phase"`
	ObservedGeneration int64       `json:"observedGeneration,omitempty"`
	ReleaseStatus      string      `json:"releaseStatus,omitempty"`
	ReleaseLastUpdated metav1.Time `json:"releaseLastUpdated,omitempty"`
	Revision           int64       `json:"revision,omitempty"`
	RollbackRevision   int64       `json:"rollbackRevision,omitempty"`
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	Reason             string      `json:"reason,omitempty"`
	Message            string      `json:"message,omitempty"`
	Manifest           string      `json:"manifest"`
}

type AddonSpecChart struct {
	ChartName    *string `json:"chartName,omitempty"`
	ChartVersion *string `json:"chartVersion,omitempty"`
}

type AddonSpecValues struct {
	RawValuesType *string   `json:"rawValuesType,omitempty"`
	RawValues     *string   `json:"rawValues,omitempty"`
	Values        []*string `json:"values,omitempty"`
}

type AddonSpec struct {
	Chart  *AddonSpecChart  `json:"chart,omitempty"`
	Values *AddonSpecValues `json:"values,omitempty"`
}

type AddonRequestBody struct {
	Kind       *string    `json:"kind,omitempty"`
	ApiVersion *string    `json:"apiVersion,omitempty"`
	Spec       *AddonSpec `json:"spec,omitempty"`
}

type AddonResponseMeta struct {
	Name              *string            `json:"name,omitempty"`
	GenerateName      *string            `json:"generateName,omitempty"`
	Namespace         *string            `json:"namespace,omitempty"`
	SelfLink          *string            `json:"selfLink,omitempty"`
	Uid               *string            `json:"uid,omitempty"`
	ResourceVersion   *string            `json:"resourceVersion,omitempty"`
	Generation        *int               `json:"generation,omitempty"`
	CreationTimestamp *string            `json:"creationTimestamp,omitempty"`
	Labels            map[string]*string `json:"labels,omitempty"`
}

type AddonResponseData struct {
	Kind       *string                `json:"kind,omitempty"`
	ApiVersion *string                `json:"apiVersion,omitempty"`
	Metadata   *AddonResponseMeta     `json:"metadata,omitempty"`
	Spec       *AddonSpec             `json:"spec,omitempty"`
	Status     map[string]interface{} `json:"status,omitempty"`
}

func (me *TkeService) GetTkeAppChartList(ctx context.Context, request *tke.GetTkeAppChartListRequest) (info []*tke.AppChart, errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().GetTkeAppChartList(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	info = response.Response.AppCharts

	return
}

func (me *TkeService) DescribeAddonList(ctx context.Context, clusterId, addonName string) (ret []*tke.Addon, errRet error) {
	request := tke.NewDescribeAddonRequest()
	response := tke.NewDescribeAddonResponse()
	request.ClusterId = &clusterId
	if addonName != "" {
		request.AddonName = &addonName
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseTkeClient().DescribeAddon(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.Addons) == 0 {
		return
	}

	ret = response.Response.Addons
	return
}

func (me *TkeService) DescribeAddonValuesList(ctx context.Context, clusterId, addonName string) (ret *tke.DescribeAddonValuesResponseParams, errRet error) {
	request := tke.NewDescribeAddonValuesRequest()
	response := tke.NewDescribeAddonValuesResponse()
	request.ClusterId = &clusterId
	request.AddonName = &addonName

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseTkeClient().DescribeAddonValues(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response
	return
}

func (me *TkeService) GetAddonReqBody(addon, version string, values []*string, rawValues, rawValuesType *string) (string, error) {
	var reqBody = &AddonRequestBody{}
	//reqBody.Kind = helper.String("App") // Optional
	//reqBody.ApiVersion = helper.String("application.tkestack.io/v1") // Optional
	reqBody.Spec = &AddonSpec{
		Chart: &AddonSpecChart{
			ChartName:    &addon,
			ChartVersion: &version,
		},
	}

	addonValues := &AddonSpecValues{}
	if len(values) > 0 {
		addonValues.RawValuesType = helper.String("yaml")
		addonValues.Values = values
	}

	if rawValuesType != nil && rawValues != nil {
		base64EncodeValues := base64.StdEncoding.EncodeToString([]byte(*rawValues))
		addonValues.RawValuesType = rawValuesType
		addonValues.RawValues = &base64EncodeValues
	}

	reqBody.Spec.Values = addonValues

	result, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func (me *TkeService) GetAddonsPath(cluster, addon string) string {
	addonPath := ""
	if addon != "" {
		addonPath = fmt.Sprintf("/%s", addon)
	}
	return fmt.Sprintf("/apis/application.tkestack.io/v1/namespaces/%s/apps%s", cluster, addonPath)
}

func (me *TkeService) GetAddonNameFromJson(reqJson string) (name string, err error) {
	reqBody := &AddonRequestBody{}
	err = json.Unmarshal([]byte(reqJson), reqBody)
	if err != nil {
		err = fmt.Errorf("error when reading chart name in addon param: %s", err.Error())
	}
	if reqBody.Spec != nil && reqBody.Spec.Chart != nil && reqBody.Spec.Chart.ChartName != nil {
		name = *reqBody.Spec.Chart.ChartName
	}
	return
}
