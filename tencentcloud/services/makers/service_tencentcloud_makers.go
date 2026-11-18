package makers

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	makerssdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/makers/v20261113"
	sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/makers/v20261113"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewMakersService(client *connectivity.TencentCloudClient) MakersService {
	return MakersService{client: client}
}

type MakersService struct {
	client *connectivity.TencentCloudClient
}

// DescribeDeploymentLog returns the build log of the deployment identified by
// zoneId, projectId and deploymentId. The three keys form an exact query, so
// no pagination is involved. A nil LogUrl with a nil error is a valid outcome:
// the cloud side documents that LogUrl may be null when no log is available.
func (me *MakersService) DescribeDeploymentLog(ctx context.Context, request *makerssdk.DescribeDeploymentLogRequest) (logUrl *string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *makerssdk.DescribeDeploymentLogResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseMakersClient().DescribeDeploymentLogWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("DescribeDeploymentLog returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	logUrl = response.Response.LogUrl
	return
}

// CreateZoneCustomDomain binds a custom domain to a Makers zone and returns
// the domain echoed back by the cloud side. The cloud-side id is not used as
// a resource identity: all follow-up APIs address the domain by the composite
// key (zone id, project id, domain name).
func (me *MakersService) CreateZoneCustomDomain(ctx context.Context, request *makerssdk.CreateZoneCustomDomainRequest) (domain string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *makerssdk.CreateZoneCustomDomainResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseMakersClient().CreateZoneCustomDomainWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("CreateZoneCustomDomain returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.Domain == nil || *response.Response.Domain == "" {
		errRet = fmt.Errorf("CreateZoneCustomDomain returned an empty domain")
		return
	}

	domain = *response.Response.Domain
	return
}

// DescribeZoneCustomDomain queries a single custom domain by the composite
// key (zone id, project id, domain name). Only `Custom` type domains are
// considered managed; `Preset` domains are platform built-ins.
//
// A nil result with a nil error means "queried successfully, no such domain";
// an API failure is always reported as an error so the caller never confuses
// a failed request with a deleted resource.
func (me *MakersService) DescribeZoneCustomDomain(ctx context.Context, zoneId, projectId, domain string) (zoneDomain *makerssdk.ZoneDomain, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := makerssdk.NewDescribeZoneCustomDomainsRequest()
	request.ZoneId = &zoneId
	request.ProjectId = &projectId
	request.CustomDomains = []*makerssdk.CustomDomainInput{
		{
			Domain: &domain,
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *makerssdk.DescribeZoneCustomDomainsResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseMakersClient().DescribeZoneCustomDomainsWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("DescribeZoneCustomDomains returned an empty response for domain %s", domain)
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	for _, item := range response.Response.PagesDomains {
		if item == nil || item.Domain == nil || *item.Domain != domain {
			continue
		}
		if item.Type != nil && *item.Type != "Custom" {
			continue
		}
		zoneDomain = item
		return
	}

	return
}

// ModifyZoneCustomDomainEnv switches the environment a custom domain is
// bound to. This is the only in-place update the cloud API supports; any
// other change is handled by replacement in the Terraform layer.
func (me *MakersService) ModifyZoneCustomDomainEnv(ctx context.Context, request *makerssdk.ModifyZoneCustomDomainEnvRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *makerssdk.ModifyZoneCustomDomainEnvResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseMakersClient().ModifyZoneCustomDomainEnvWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("ModifyZoneCustomDomainEnv returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

// DeleteProjectCustomDomains unbinds the given custom domain from the
// project. The cloud API is batch-shaped; the Terraform layer always sends
// exactly one element so destroy only affects the managed domain.
func (me *MakersService) DeleteProjectCustomDomain(ctx context.Context, zoneId, projectId, domain string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := makerssdk.NewDeleteProjectCustomDomainsRequest()
	request.ZoneId = &zoneId
	request.ProjectId = &projectId
	request.CustomDomains = []*string{&domain}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *makerssdk.DeleteProjectCustomDomainsResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseMakersClient().DeleteProjectCustomDomainsWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("DeleteProjectCustomDomains returned an empty response for domain %s", domain)
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

// Zone domain status values reported by DescribeZoneCustomDomains.
const (
	// ZoneDomainStatusOnline is the terminal success status of a domain:
	// the domain is bound and serving traffic.
	ZoneDomainStatusOnline = "online"
	// ZoneDomainStatusProcess means an async operation (binding, unbinding or
	// environment switch) is still in progress.
	ZoneDomainStatusProcess = "process"
)

// zoneCustomDomainInProgress reports whether the reported status means an
// async operation is still running. A missing status is treated as in
// progress: the domain row may be returned before all fields settle.
func zoneCustomDomainInProgress(status *string) bool {
	return status == nil || strings.EqualFold(*status, ZoneDomainStatusProcess)
}

// zoneCustomDomainStatusOf safely formats a possibly nil status value.
func zoneCustomDomainStatusOf(status *string) string {
	if status == nil {
		return ""
	}
	return *status
}

// WaitForZoneCustomDomain polls DescribeZoneCustomDomains until the domain
// is visible with a settled (non-process) status. The binding operation is
// async on the cloud side: right after CreateZoneCustomDomain succeeds the
// domain may not be readable yet, or may be readable with a `process`
// status while the binding is still converging.
func (me *MakersService) WaitForZoneCustomDomain(ctx context.Context, zoneId, projectId, domain string, timeout time.Duration) error {
	logId := tccommon.GetLogId(ctx)

	err := resource.Retry(timeout, func() *resource.RetryError {
		item, e := me.DescribeZoneCustomDomain(ctx, zoneId, projectId, domain)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if item == nil {
			return resource.RetryableError(fmt.Errorf("makers zone custom domain [%s] is not ready yet", domain))
		}
		if zoneCustomDomainInProgress(item.Status) {
			return resource.RetryableError(fmt.Errorf("makers zone custom domain [%s] is still processing, status [%s]", domain, zoneCustomDomainStatusOf(item.Status)))
		}
		if !strings.EqualFold(*item.Status, ZoneDomainStatusOnline) {
			return resource.NonRetryableError(fmt.Errorf("makers zone custom domain [%s] reached unexpected status [%s]", domain, *item.Status))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s wait for makers zone custom domain [%s] failed, reason:%+v", logId, domain, err)
		return err
	}
	return nil
}

// WaitForZoneCustomDomainDeleted polls DescribeZoneCustomDomains until the
// domain is no longer visible. The unbinding operation is async on the cloud
// side, so the domain may still be reported (possibly with a `process`
// status) right after the delete call returns.
func (me *MakersService) WaitForZoneCustomDomainDeleted(ctx context.Context, zoneId, projectId, domain string, timeout time.Duration) error {
	logId := tccommon.GetLogId(ctx)

	err := resource.Retry(timeout, func() *resource.RetryError {
		item, e := me.DescribeZoneCustomDomain(ctx, zoneId, projectId, domain)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if item != nil {
			return resource.RetryableError(fmt.Errorf("makers zone custom domain [%s] is still being deleted, status [%s]", domain, zoneCustomDomainStatusOf(item.Status)))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s wait for makers zone custom domain [%s] deletion failed, reason:%+v", logId, domain, err)
		return err
	}
	return nil
}

// WaitForZoneCustomDomainEnv polls DescribeZoneCustomDomains until the
// domain reports the expected environment slug with a settled status. The
// environment switch is async on the cloud side.
func (me *MakersService) WaitForZoneCustomDomainEnv(ctx context.Context, zoneId, projectId, domain, env string, timeout time.Duration) error {
	logId := tccommon.GetLogId(ctx)

	err := resource.Retry(timeout, func() *resource.RetryError {
		item, e := me.DescribeZoneCustomDomain(ctx, zoneId, projectId, domain)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if item == nil {
			return resource.RetryableError(fmt.Errorf("makers zone custom domain [%s] is not readable while switching env", domain))
		}
		if zoneCustomDomainInProgress(item.Status) {
			return resource.RetryableError(fmt.Errorf("makers zone custom domain [%s] env switch is still processing, status [%s]", domain, zoneCustomDomainStatusOf(item.Status)))
		}
		if !strings.EqualFold(*item.Status, ZoneDomainStatusOnline) {
			return resource.NonRetryableError(fmt.Errorf("makers zone custom domain [%s] reached unexpected status [%s] while switching env", domain, *item.Status))
		}
		if item.Env == nil || *item.Env != env {
			return resource.RetryableError(fmt.Errorf("makers zone custom domain [%s] env is not ready yet, expect [%s]", domain, env))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s wait for makers zone custom domain [%s] env [%s] failed, reason:%+v", logId, domain, env, err)
		return err
	}
	return nil
}

// CreateProject creates a Makers project and returns the generated project id.
func (me *MakersService) CreateProject(ctx context.Context, request *makerssdk.CreateProjectRequest) (projectId *string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *makerssdk.CreateProjectResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseMakersClient().CreateProjectWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("CreateProject returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	projectId = response.Response.ProjectId
	return
}

// DescribeProject returns the project identified by zoneId and projectId.
// The two keys form an exact query, so no pagination is involved. A nil
// project with a nil error means the project does not exist on the cloud
// side.
func (me *MakersService) DescribeProject(ctx context.Context, zoneId, projectId string) (project *makerssdk.Project, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := sdk.NewDescribeProjectsRequest()
	request.ZoneId = helper.String(zoneId)
	request.ProjectIds = []*string{helper.String(projectId)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *makerssdk.DescribeProjectsResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseMakersClient().DescribeProjectsWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("DescribeProjects returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Projects) > 0 {
		for _, item := range response.Response.Projects {
			if item.ProjectId != nil && *item.ProjectId == projectId && item.ZoneId != nil && *item.ZoneId == zoneId {
				return item, nil
			}
		}
	}
	return nil, nil
}

// ModifyProject updates the basic information of a Makers project.
func (me *MakersService) ModifyProject(ctx context.Context, request *makerssdk.ModifyProjectRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *makerssdk.ModifyProjectResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseMakersClient().ModifyProjectWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("ModifyProject returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

// DeleteProject deletes a Makers project.
func (me *MakersService) DeleteProject(ctx context.Context, request *makerssdk.DeleteProjectRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *makerssdk.DeleteProjectResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseMakersClient().DeleteProjectWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("DeleteProject returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

// Deployment terminal statuses reported by the DescribeDeployments API.
const (
	// DeploymentStatusSuccess means the build finished successfully.
	DeploymentStatusSuccess = "Success"
	// DeploymentStatusFailed means the build finished with an error.
	DeploymentStatusFailed = "Failed"
	// DeploymentStatusProcess means the build is still running.
	DeploymentStatusProcess = "Process"
	// DeploymentStatusInvalid means the deployment is invalid.
	DeploymentStatusInvalid = "Invalid"
	// DeploymentStatusPending means the build is queued and not started yet.
	DeploymentStatusPending = "Pending"
	// DeploymentStatusTimeout means the build exceeded the cloud-side time limit.
	DeploymentStatusTimeout = "Timeout"
	// DeploymentStatusCancelled means the build was cancelled.
	DeploymentStatusCancelled = "Cancelled"
)

// CreateDeployment triggers a new Makers build via the CreateDeployment API.
// The API is asynchronous: it returns a DeploymentId as soon as the build is
// accepted, and the build result must be confirmed separately through
// DescribeDeployments. DeploymentId may legitimately be nil on the wire; the
// caller decides whether that is acceptable.
func (me *MakersService) CreateDeployment(ctx context.Context, request *makerssdk.CreateDeploymentRequest) (deploymentId *string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *makerssdk.CreateDeploymentResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseMakersClient().CreateDeploymentWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("CreateDeployment returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	deploymentId = response.Response.DeploymentId
	return
}

// DescribeDeploymentById queries a single deployment by its exact DeploymentId.
// ZoneId and ProjectId scope the query. It returns the deployment object, or
// (nil, nil) when no deployment matches the id.
func (me *MakersService) DescribeDeploymentById(ctx context.Context, zoneId, projectId, deploymentId string) (deployment *makerssdk.Deployment, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := makerssdk.NewDescribeDeploymentsRequest()
	request.ZoneId = helper.String(zoneId)
	request.ProjectId = helper.String(projectId)
	request.DeploymentIds = []*string{helper.String(deploymentId)}
	request.Limit = helper.Int64(1)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *makerssdk.DescribeDeploymentsResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseMakersClient().DescribeDeploymentsWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("DescribeDeployments returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	// The query filters by an exact deployment id, so at most one entry can
	// match. An empty list is a legitimate "not found" outcome, not an error.
	for _, item := range response.Response.Deployments {
		if item != nil && item.DeploymentId != nil && *item.DeploymentId == deploymentId {
			deployment = item
			break
		}
	}
	return
}

// WaitMakersDeploymentComplete polls DescribeDeployments until the deployment
// reaches a terminal state (Success / Failed / Invalid / Timeout / Cancelled),
// or the timeout budget is exhausted. It returns the last observed deployment
// snapshot so the caller can attach it to a timeout diagnostic. A nil
// deployment is treated as "not visible yet" and keeps the loop polling.
func (me *MakersService) WaitMakersDeploymentComplete(ctx context.Context, zoneId, projectId, deploymentId string, timeout time.Duration, pollInterval time.Duration) (last *makerssdk.Deployment, err error) {
	deadline := time.Now().Add(timeout)
	for {
		latest, e := me.DescribeDeploymentById(ctx, zoneId, projectId, deploymentId)
		if e != nil {
			err = e
			return
		}
		last = latest
		if latest != nil && latest.Status != nil {
			switch *latest.Status {
			case DeploymentStatusSuccess:
				return
			case DeploymentStatusFailed, DeploymentStatusInvalid, DeploymentStatusTimeout, DeploymentStatusCancelled:
				err = fmt.Errorf("deployment %s ended in status %s", deploymentId, *latest.Status)
				return
			}
		}

		if time.Now().After(deadline) {
			err = fmt.Errorf("deployment %s did not converge within %s", deploymentId, timeout)
			return
		}
		time.Sleep(pollInterval)
	}
}

// DescribeCosTempToken queries a temporary cos token for uploading files.
// ZoneId and ProjectId scope the query. It returns the token object.
func (me *MakersService) DescribeCosTempToken(ctx context.Context, zoneId, projectId string) (cosTempParam *makerssdk.DescribeCosTempTokenResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	project, err := me.DescribeProject(ctx, zoneId, projectId)
	if err != nil {
		errRet = fmt.Errorf("fetching the project: %w", err)
		return
	}
	if project == nil || project.Name == nil {
		errRet = fmt.Errorf("project %s does not exist or has no name", projectId)
		return
	}

	request := makerssdk.NewDescribeCosTempTokenRequest()
	request.ZoneId = helper.String(zoneId)
	request.ProjectId = helper.String(projectId)
	request.ProjectName = project.Name

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *makerssdk.DescribeCosTempTokenResponse
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseMakersClient().DescribeCosTempTokenWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("DescribeCosTempToken returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	cosTempParam = response.Response
	return
}

// UploadMakersLocalArtifact uploads a local build artifact to the Makers temp
// bucket using the temporary credential from DescribeCosTempToken, and
// returns the object path to pass to CreateDeployment as TempBucketPath.
//
// A folder given as localPath is compressed into a zip archive first. The
// cloud side scopes the credential to the TargetPath prefix returned by
// DescribeCosTempToken (format: <appId>/<projectName>/<timestamp>); the
// artifact is stored under that prefix using its base name.
func (me *MakersService) UploadMakersLocalArtifact(ctx context.Context, zoneId, projectId, localPath string, progress func(string)) (string, error) {
	token, err := me.DescribeCosTempToken(ctx, zoneId, projectId)
	if err != nil {
		return "", fmt.Errorf("fetching the COS temporary credential: %w", err)
	}
	if token.Bucket == nil || token.Region == nil || token.TargetPath == nil || token.Credentials == nil ||
		token.Credentials.TmpSecretId == nil || token.Credentials.TmpSecretKey == nil || token.Credentials.Token == nil {
		return "", fmt.Errorf("DescribeCosTempToken returned an incomplete credential: bucket, region, target path and the temporary key triplet must all be set")
	}

	// A folder is packaged into a zip first (the build expects an archive,
	// typically with dist_type = Folder); a plain file is uploaded as-is.
	uploadPath := localPath
	objectName := filepath.Base(localPath)
	info, err := os.Stat(localPath)
	if err != nil {
		return "", fmt.Errorf("reading the local artifact %s: %w", localPath, err)
	}
	if info.IsDir() {
		progress("compressing the local folder into a zip archive")
		zipPath, err := zipMakersDirectory(localPath)
		if err != nil {
			return "", err
		}
		defer func() {
			_ = os.Remove(zipPath)
		}()
		uploadPath = zipPath
		objectName += ".zip"
	}

	cosService := NewMakersCosService(
		me.client,
		*token.Region,
		*token.Credentials.TmpSecretId,
		*token.Credentials.TmpSecretKey,
		*token.Credentials.Token)
	objectKey := fmt.Sprintf("%s/%s", strings.Trim(*token.TargetPath, "/"), objectName)
	if _, err := cosService.PutFile(ctx, *token.Bucket, objectKey, uploadPath); err != nil {
		return "", fmt.Errorf("uploading the artifact to COS: %w", err)
	}
	return objectKey, nil
}

// zipMakersDirectory packages a local folder into a zip archive stored in a
// temporary file and returns its path. The archive stores files with paths
// relative to the folder root; directory entries are omitted (their paths
// are implied by the file entries).
func zipMakersDirectory(dir string) (string, error) {
	tmp, err := os.CreateTemp("", "makers-artifact-*.zip")
	if err != nil {
		return "", fmt.Errorf("creating a temporary zip file: %w", err)
	}
	defer func() {
		_ = tmp.Close()
	}()

	writer := zip.NewWriter(tmp)
	err = filepath.WalkDir(dir, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		entryWriter, err := writer.Create(filepath.ToSlash(rel))
		if err != nil {
			return err
		}
		src, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func() {
			_ = src.Close()
		}()
		_, err = io.Copy(entryWriter, src)
		return err
	})
	if closeErr := writer.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		_ = os.Remove(tmp.Name())
		return "", fmt.Errorf("compressing the folder %s: %w", dir, err)
	}
	return tmp.Name(), nil
}
