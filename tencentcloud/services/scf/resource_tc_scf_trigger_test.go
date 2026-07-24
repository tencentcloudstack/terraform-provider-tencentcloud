package scf_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	scfsvc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/scf"
)

// Run with: go test ./tencentcloud/services/scf/ -run "TestScfTrigger" -v -count=1 -gcflags="all=-l"

type mockMetaScfTrigger struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaScfTrigger) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaScfTrigger{}

func newMockMetaScfTrigger() *mockMetaScfTrigger {
	return &mockMetaScfTrigger{client: &connectivity.TencentCloudClient{}}
}

func ptrStringScfTrigger(s string) *string {
	return &s
}

func ptrUint64ScfTrigger(v uint64) *uint64 {
	return &v
}

// TestScfTriggerCreate verifies the create path sets the composite id and passes
// schema fields to the CreateTrigger request.
func TestScfTriggerCreate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	scfClient := &scf.Client{}
	patches.ApplyMethodReturn(newMockMetaScfTrigger().client, "UseScfClient", scfClient)

	var capturedRequest *scf.CreateTriggerRequest
	patches.ApplyMethodFunc(scfClient, "CreateTriggerWithContext", func(ctx interface{}, request *scf.CreateTriggerRequest) (*scf.CreateTriggerResponse, error) {
		capturedRequest = request
		resp := scf.NewCreateTriggerResponse()
		resp.Response = &scf.CreateTriggerResponseParams{
			TriggerInfo: &scf.Trigger{
				TriggerName: ptrStringScfTrigger("tf-trigger"),
				Type:        ptrStringScfTrigger("timer"),
			},
			RequestId: ptrStringScfTrigger("fake-request-id"),
		}
		return resp, nil
	})

	// Mock ListTriggers called by the read after create.
	patches.ApplyMethodFunc(scfClient, "ListTriggers", func(request *scf.ListTriggersRequest) (*scf.ListTriggersResponse, error) {
		resp := scf.NewListTriggersResponse()
		resp.Response = &scf.ListTriggersResponseParams{
			TotalCount: ptrUint64ScfTrigger(1),
			Triggers: []*scf.TriggerInfo{{
				TriggerName:     ptrStringScfTrigger("tf-trigger"),
				Type:            ptrStringScfTrigger("timer"),
				TriggerDesc:     ptrStringScfTrigger(`{"cron":"*/5 * * * * * *"}`),
				Qualifier:       ptrStringScfTrigger("$DEFAULT"),
				Enable:          ptrUint64ScfTrigger(1),
				AvailableStatus: ptrStringScfTrigger("AVAILABLE"),
				Description:     ptrStringScfTrigger("desc"),
				CustomArgument:  ptrStringScfTrigger("arg"),
				AddTime:         ptrStringScfTrigger("2025-01-01 00:00:00"),
				ModTime:         ptrStringScfTrigger("2025-01-01 00:00:00"),
			}},
			RequestId: ptrStringScfTrigger("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaScfTrigger()
	res := scfsvc.ResourceTencentCloudScfTrigger()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"function_name":   "tf-function",
		"trigger_name":    "tf-trigger",
		"type":            "timer",
		"namespace":       "default",
		"trigger_desc":    `{"cron":"*/5 * * * * *"}`,
		"qualifier":       "$DEFAULT",
		"enable":          "OPEN",
		"description":     "desc",
		"custom_argument": "arg",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)

	// Verify composite id.
	assert.Equal(t, "tf-function#default#tf-trigger", d.Id())

	// Verify required fields were passed to the create request.
	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "tf-function", *capturedRequest.FunctionName)
	assert.Equal(t, "tf-trigger", *capturedRequest.TriggerName)
	assert.Equal(t, "timer", *capturedRequest.Type)
	assert.Equal(t, "default", *capturedRequest.Namespace)
	assert.Equal(t, "OPEN", *capturedRequest.Enable)

	// Verify enable was converted from int64 to string on read.
	assert.Equal(t, "OPEN", d.Get("enable").(string))
	assert.Equal(t, "AVAILABLE", d.Get("available_status").(string))
}

// TestScfTriggerCreate_NilResponse verifies create fails when the API returns a nil TriggerInfo.
func TestScfTriggerCreate_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	scfClient := &scf.Client{}
	patches.ApplyMethodReturn(newMockMetaScfTrigger().client, "UseScfClient", scfClient)

	patches.ApplyMethodFunc(scfClient, "CreateTriggerWithContext", func(ctx interface{}, request *scf.CreateTriggerRequest) (*scf.CreateTriggerResponse, error) {
		resp := scf.NewCreateTriggerResponse()
		resp.Response = &scf.CreateTriggerResponseParams{
			RequestId: ptrStringScfTrigger("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaScfTrigger()
	res := scfsvc.ResourceTencentCloudScfTrigger()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"function_name": "tf-function",
		"trigger_name":  "tf-trigger",
		"type":          "timer",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	// Verify id was NOT set when TriggerInfo is nil.
	assert.Equal(t, "", d.Id())
}

// TestScfTriggerRead verifies the read path populates state from TriggerInfo.
func TestScfTriggerRead(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	scfClient := &scf.Client{}
	patches.ApplyMethodReturn(newMockMetaScfTrigger().client, "UseScfClient", scfClient)

	patches.ApplyMethodFunc(scfClient, "ListTriggers", func(request *scf.ListTriggersRequest) (*scf.ListTriggersResponse, error) {
		resp := scf.NewListTriggersResponse()
		resp.Response = &scf.ListTriggersResponseParams{
			TotalCount: ptrUint64ScfTrigger(1),
			Triggers: []*scf.TriggerInfo{{
				TriggerName:     ptrStringScfTrigger("tf-trigger"),
				Type:            ptrStringScfTrigger("timer"),
				TriggerDesc:     ptrStringScfTrigger(`{"cron":"*/5 * * * * *"}`),
				Qualifier:       ptrStringScfTrigger("$DEFAULT"),
				Enable:          ptrUint64ScfTrigger(1),
				AvailableStatus: ptrStringScfTrigger("AVAILABLE"),
				Description:     ptrStringScfTrigger("desc"),
				CustomArgument:  ptrStringScfTrigger("arg"),
				AddTime:         ptrStringScfTrigger("2025-01-01 00:00:00"),
				ModTime:         ptrStringScfTrigger("2025-01-01 00:00:00"),
			}},
			RequestId: ptrStringScfTrigger("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaScfTrigger()
	res := scfsvc.ResourceTencentCloudScfTrigger()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("tf-function#default#tf-trigger")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "tf-function#default#tf-trigger", d.Id())
	assert.Equal(t, "tf-function", d.Get("function_name").(string))
	assert.Equal(t, "default", d.Get("namespace").(string))
	assert.Equal(t, "tf-trigger", d.Get("trigger_name").(string))
	assert.Equal(t, "timer", d.Get("type").(string))
	assert.Equal(t, "OPEN", d.Get("enable").(string))
	assert.Equal(t, "AVAILABLE", d.Get("available_status").(string))
	assert.Equal(t, "2025-01-01 00:00:00", d.Get("add_time").(string))
	assert.Equal(t, "2025-01-01 00:00:00", d.Get("mod_time").(string))
}

// TestScfTriggerRead_NotFound verifies the read path clears the id when the trigger is gone.
func TestScfTriggerRead_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	scfClient := &scf.Client{}
	patches.ApplyMethodReturn(newMockMetaScfTrigger().client, "UseScfClient", scfClient)

	patches.ApplyMethodFunc(scfClient, "ListTriggers", func(request *scf.ListTriggersRequest) (*scf.ListTriggersResponse, error) {
		resp := scf.NewListTriggersResponse()
		resp.Response = &scf.ListTriggersResponseParams{
			TotalCount: ptrUint64ScfTrigger(0),
			Triggers:   []*scf.TriggerInfo{},
			RequestId:  ptrStringScfTrigger("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaScfTrigger()
	res := scfsvc.ResourceTencentCloudScfTrigger()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("tf-function#default#tf-trigger")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	// Verify the id was cleared because the trigger was not found.
	assert.Equal(t, "", d.Id())
}

// TestScfTriggerRead_BrokenId verifies the read path rejects a broken id.
func TestScfTriggerRead_BrokenId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	scfClient := &scf.Client{}
	patches.ApplyMethodReturn(newMockMetaScfTrigger().client, "UseScfClient", scfClient)

	meta := newMockMetaScfTrigger()
	res := scfsvc.ResourceTencentCloudScfTrigger()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("broken-id")

	err := res.Read(d, meta)
	assert.Error(t, err)
}

// TestScfTriggerUpdate verifies the update path passes mutable fields to UpdateTrigger.
func TestScfTriggerUpdate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	scfClient := &scf.Client{}
	patches.ApplyMethodReturn(newMockMetaScfTrigger().client, "UseScfClient", scfClient)

	var capturedRequest *scf.UpdateTriggerRequest
	patches.ApplyMethodFunc(scfClient, "UpdateTriggerWithContext", func(ctx interface{}, request *scf.UpdateTriggerRequest) (*scf.UpdateTriggerResponse, error) {
		capturedRequest = request
		resp := scf.NewUpdateTriggerResponse()
		resp.Response = &scf.UpdateTriggerResponseParams{
			RequestId: ptrStringScfTrigger("fake-request-id"),
		}
		return resp, nil
	})

	// Mock ListTriggers called by the read after update.
	patches.ApplyMethodFunc(scfClient, "ListTriggers", func(request *scf.ListTriggersRequest) (*scf.ListTriggersResponse, error) {
		resp := scf.NewListTriggersResponse()
		resp.Response = &scf.ListTriggersResponseParams{
			TotalCount: ptrUint64ScfTrigger(1),
			Triggers: []*scf.TriggerInfo{{
				TriggerName:     ptrStringScfTrigger("tf-trigger"),
				Type:            ptrStringScfTrigger("timer"),
				TriggerDesc:     ptrStringScfTrigger(`{"cron":"*/5 * * * * *"}`),
				Qualifier:       ptrStringScfTrigger("$DEFAULT"),
				Enable:          ptrUint64ScfTrigger(0),
				AvailableStatus: ptrStringScfTrigger("AVAILABLE"),
				Description:     ptrStringScfTrigger("desc-updated"),
				CustomArgument:  ptrStringScfTrigger("arg"),
				AddTime:         ptrStringScfTrigger("2025-01-01 00:00:00"),
				ModTime:         ptrStringScfTrigger("2025-01-02 00:00:00"),
			}},
			RequestId: ptrStringScfTrigger("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaScfTrigger()
	res := scfsvc.ResourceTencentCloudScfTrigger()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"function_name":   "tf-function",
		"trigger_name":    "tf-trigger",
		"type":            "timer",
		"namespace":       "default",
		"trigger_desc":    `{"cron":"*/5 * * * * *"}`,
		"qualifier":       "$DEFAULT",
		"enable":          "CLOSE",
		"description":     "desc-updated",
		"custom_argument": "arg",
	})
	d.SetId("tf-function#default#tf-trigger")

	// Patch HasChange to simulate changes in mutable fields only.
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		switch key {
		case "enable", "qualifier", "trigger_desc", "description", "custom_argument":
			return true
		}
		return false
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)

	// Verify mutable fields were passed to the update request.
	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "tf-function", *capturedRequest.FunctionName)
	assert.Equal(t, "default", *capturedRequest.Namespace)
	assert.Equal(t, "tf-trigger", *capturedRequest.TriggerName)
	assert.Equal(t, "CLOSE", *capturedRequest.Enable)
	assert.Equal(t, "desc-updated", *capturedRequest.Description)

	// Verify enable was read back as CLOSE (int64 0 -> CLOSE).
	assert.Equal(t, "CLOSE", d.Get("enable").(string))
	assert.Equal(t, "desc-updated", d.Get("description").(string))
}

// TestScfTriggerDelete verifies the delete path builds the DeleteTrigger request.
func TestScfTriggerDelete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	scfClient := &scf.Client{}
	patches.ApplyMethodReturn(newMockMetaScfTrigger().client, "UseScfClient", scfClient)

	var capturedRequest *scf.DeleteTriggerRequest
	patches.ApplyMethodFunc(scfClient, "DeleteTriggerWithContext", func(ctx interface{}, request *scf.DeleteTriggerRequest) (*scf.DeleteTriggerResponse, error) {
		capturedRequest = request
		resp := scf.NewDeleteTriggerResponse()
		resp.Response = &scf.DeleteTriggerResponseParams{
			RequestId: ptrStringScfTrigger("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaScfTrigger()
	res := scfsvc.ResourceTencentCloudScfTrigger()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"function_name": "tf-function",
		"trigger_name":  "tf-trigger",
		"type":          "timer",
		"namespace":     "default",
		"trigger_desc":  `{"cron":"*/5 * * * * *"}`,
		"qualifier":     "$DEFAULT",
	})
	d.SetId("tf-function#default#tf-trigger")

	err := res.Delete(d, meta)
	assert.NoError(t, err)

	// Verify the delete request was built correctly from state.
	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "tf-function", *capturedRequest.FunctionName)
	assert.Equal(t, "tf-trigger", *capturedRequest.TriggerName)
	assert.Equal(t, "timer", *capturedRequest.Type)
	assert.Equal(t, "default", *capturedRequest.Namespace)
	assert.Equal(t, "$DEFAULT", *capturedRequest.Qualifier)
}

// TestScfTriggerSchema validates the schema definition.
func TestScfTriggerSchema(t *testing.T) {
	res := scfsvc.ResourceTencentCloudScfTrigger()

	// Required + ForceNew identity fields.
	functionName, ok := res.Schema["function_name"]
	assert.True(t, ok, "function_name should exist in schema")
	assert.True(t, functionName.Required)
	assert.True(t, functionName.ForceNew)

	triggerName, ok := res.Schema["trigger_name"]
	assert.True(t, ok, "trigger_name should exist in schema")
	assert.True(t, triggerName.Required)
	assert.True(t, triggerName.ForceNew)

	typeField, ok := res.Schema["type"]
	assert.True(t, ok, "type should exist in schema")
	assert.True(t, typeField.Required)
	assert.True(t, typeField.ForceNew)

	// Namespace defaults to default and is ForceNew.
	namespace, ok := res.Schema["namespace"]
	assert.True(t, ok, "namespace should exist in schema")
	assert.True(t, namespace.Optional)
	assert.True(t, namespace.ForceNew)
	assert.Equal(t, "default", namespace.Default)

	// Mutable fields are Optional and not ForceNew.
	for _, field := range []string{"enable", "qualifier", "trigger_desc", "description", "custom_argument"} {
		f, ok := res.Schema[field]
		assert.True(t, ok, "%s should exist in schema", field)
		assert.True(t, f.Optional, "%s should be optional", field)
		assert.False(t, f.ForceNew, "%s should not be ForceNew", field)
	}

	// Computed fields.
	for _, field := range []string{"available_status", "add_time", "mod_time"} {
		f, ok := res.Schema[field]
		assert.True(t, ok, "%s should exist in schema", field)
		assert.True(t, f.Computed, "%s should be computed", field)
	}

	// Importer.
	assert.NotNil(t, res.Importer)
}
