package tencentcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
	"log"
)

type AddonSpecChart struct {
	ChartName    *string `json:"chartName,omitempty"`
	ChartVersion *string `json:"chartVersion,omitempty"`
}

type AddonSpecValues struct {
	RawValuesType *string   `json:"rawValuesType,omitempty"`
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
	Kind       *string            `json:"kind,omitempty"`
	ApiVersion *string            `json:"apiVersion,omitempty"`
	Metadata   *AddonResponseMeta `json:"metadata,omitempty"`
	Spec       *AddonSpec         `json:"spec,omitempty"`
	Status     map[string]*string `json:"status,omitempty"`
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

func (me *TkeService) PollingAddonsPhase(ctx context.Context, clusterId, addonName string, addonResponseData *AddonResponseData) (string, bool, error) {
	var (
		err error
		phase string
		response string
		has bool
	)

	if addonResponseData == nil {
		addonResponseData = &AddonResponseData{}
	}

	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
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

		reason := addonResponseData.Status["reason"]
		if addonResponseData.Status["phase"] != nil {
			phase = *addonResponseData.Status["phase"]
		}
		if reason == nil {
			reason = helper.String("unknown error")
		}

		if phase == "Upgrading" || phase == "Installing" || phase == "ChartFetched" || phase == "RollingBack" || phase == "Terminating"{
			return resource.RetryableError(fmt.Errorf("addon %s is %s, retrying", addonName, phase))
		}

		return nil
	})

	return phase, has, err
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

func (me *TkeService) GetAddonReqBody(addon, version string, values []*string) (string, error) {
	var reqBody = &AddonRequestBody{}
	//reqBody.Kind = helper.String("App") // Optional
	//reqBody.ApiVersion = helper.String("application.tkestack.io/v1") // Optional
	reqBody.Spec = &AddonSpec{
		Chart: &AddonSpecChart{
			ChartName:    &addon,
			ChartVersion: &version,
		},
	}

	if len(values) > 0 {
		reqBody.Spec.Values = &AddonSpecValues{
			RawValuesType: helper.String("yaml"),
			Values:        values,
		}
	}

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
