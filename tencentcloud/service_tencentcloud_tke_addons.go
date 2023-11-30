package tencentcloud

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	logId := getLogId(ctx)
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

func (me *TkeService) DescribeExtensionAddonList(ctx context.Context, clusterId string) (AppList, error) {
	var (
		err      error
		response string
		appList  AppList
	)

	err = resource.Retry(readRetryTimeout*5, func() *resource.RetryError {
		response, _, err = me.DescribeExtensionAddon(ctx, clusterId, "")
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if err := json.Unmarshal([]byte(response), &appList); err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	return appList, err
}

func (me *TkeService) PollingAddonsPhase(ctx context.Context, clusterId, addonName string, addonResponseData *AddonResponseData) (string, bool, error) {
	var (
		err      error
		phase    string
		response string
		has      bool
	)

	if addonResponseData == nil {
		addonResponseData = &AddonResponseData{}
	}

	err = resource.Retry(readRetryTimeout*5, func() *resource.RetryError {
		response, has, err = me.DescribeExtensionAddon(ctx, clusterId, addonName)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if err = json.Unmarshal([]byte(response), addonResponseData); err != nil {
			return resource.NonRetryableError(err)
		}

		if addonResponseData.Status == nil {
			return nil
		}

		if addonResponseData.Status["phase"] != nil {
			phase = addonResponseData.Status["phase"].(string)
		}

		if phase == "Upgrading" || phase == "Installing" || phase == "ChartFetched" || phase == "RollingBack" || phase == "Terminating" {
			return resource.RetryableError(fmt.Errorf("addon %s is %s, retrying", addonName, phase))
		}

		return nil
	})

	return response, has, err
}

func (me *TkeService) ProcessExtensionAddons(ctx context.Context, request *tke.ForwardApplicationRequestV3Request) (response *tke.ForwardApplicationRequestV3Response, err error) {
	logId := getLogId(ctx)
	defer func() {
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err = me.client.UseTkeClient().ForwardApplicationRequestV3(request)

	if err != nil {
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TkeService) DescribeExtensionAddon(ctx context.Context, clusterName, addon string) (result string, has bool, errRet error) {
	request := tke.NewForwardApplicationRequestV3Request()
	request.Method = helper.String("GET")
	request.ClusterName = &clusterName
	request.Path = helper.String(me.GetAddonsPath(clusterName, addon))
	response, err := me.ProcessExtensionAddons(ctx, request)
	if err != nil {
		errRet = err
		return
	}
	has = true
	result = *response.Response.ResponseBody
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

func (me *TkeService) CreateExtensionAddon(ctx context.Context, clusterName, reqBody string) (errRet error) {
	request := tke.NewForwardApplicationRequestV3Request()
	request.Method = helper.String("POST")
	request.ClusterName = &clusterName
	request.Path = helper.String(me.GetAddonsPath(clusterName, ""))
	request.RequestBody = &reqBody
	_, err := me.ProcessExtensionAddons(ctx, request)
	if err != nil {
		errRet = err
		return
	}
	return
}

func (me *TkeService) UpdateExtensionAddon(ctx context.Context, clusterName, addon, reqBody string) (errRet error) {
	request := tke.NewForwardApplicationRequestV3Request()
	request.Method = helper.String("PATCH")
	request.ContentType = helper.String("application/strategic-merge-patch+json")
	request.ClusterName = &clusterName
	request.Path = helper.String(me.GetAddonsPath(clusterName, addon))
	request.RequestBody = &reqBody
	_, err := me.ProcessExtensionAddons(ctx, request)
	if err != nil {
		errRet = err
		return
	}
	return
}

func (me *TkeService) DeleteExtensionAddon(ctx context.Context, clusterName, addon string) (errRet error) {
	request := tke.NewForwardApplicationRequestV3Request()
	request.Method = helper.String("DELETE")
	request.ClusterName = &clusterName
	request.Path = helper.String(me.GetAddonsPath(clusterName, addon))
	_, err := me.ProcessExtensionAddons(ctx, request)
	if err != nil {
		errRet = err
		return
	}
	return
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
