# Implementation Tasks: Add Lighthouse Blueprints DataSource

**Change ID**: `add-lighthouse-blueprints-datasource`  
**Total Tasks**: 62  
**Estimated Time**: 3.5 hours

---

## Phase 1: Service Layer Implementation (10 tasks, ~30 min)

### Task 1.1: Add DescribeLighthouseBlueprintsByFilter method to service
**File**: `tencentcloud/services/lighthouse/service_tencentcloud_lighthouse.go`  
**Action**: Add new method after existing blueprint-related methods (around line 1100)

```go
func (me *LightHouseService) DescribeLighthouseBlueprintsByFilter(ctx context.Context, param map[string]interface{}) (blueprints []*lighthouse.Blueprint, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = lighthouse.NewDescribeBlueprintsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	// TODO: Implement parameter mapping
	return
}
```

**Verify**: Method signature exists and compiles.

---

### Task 1.2: Implement BlueprintIds parameter
**File**: `tencentcloud/services/lighthouse/service_tencentcloud_lighthouse.go`  
**Action**: Map `BlueprintIds` from param map

```go
if v, ok := param["BlueprintIds"]; ok {
	blueprintIds := v.([]*string)
	request.BlueprintIds = blueprintIds
}
```

**Verify**: BlueprintIds are passed to API request.

---

### Task 1.3: Implement Filters parameter
**File**: `tencentcloud/services/lighthouse/service_tencentcloud_lighthouse.go`  
**Action**: Map `Filters` from param map

```go
if v, ok := param["Filters"]; ok {
	filters := v.([]*lighthouse.Filter)
	request.Filters = filters
}
```

**Verify**: Filters are passed to API request.

---

### Task 1.4: Implement pagination loop
**File**: `tencentcloud/services/lighthouse/service_tencentcloud_lighthouse.go`  
**Action**: Add pagination logic (DO NOT expose offset/limit to users)

```go
var offset int64 = 0
var pageSize int64 = 100
blueprints = make([]*lighthouse.Blueprint, 0)

for {
	request.Offset = &offset
	request.Limit = &pageSize
	ratelimit.Check(request.GetAction())
	
	response, err := me.client.UseLighthouseClient().DescribeBlueprints(request)
	if err != nil {
		errRet = err
		return
	}
	
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	
	if response == nil || response.Response == nil || len(response.Response.BlueprintSet) < 1 {
		break
	}
	
	blueprints = append(blueprints, response.Response.BlueprintSet...)
	
	if len(response.Response.BlueprintSet) < int(pageSize) {
		break
	}
	
	offset += pageSize
}
```

**Verify**: Pagination fetches all results, offset/limit hidden from users.

---

### Task 1.5: Add logging
**File**: `tencentcloud/services/lighthouse/service_tencentcloud_lighthouse.go`  
**Action**: Ensure proper logging in service method

**Verify**: Logs include request/response bodies and errors.

---

### Task 1.6-1.10: Code quality checks
**Actions**:
- Format code: `gofmt -w service_tencentcloud_lighthouse.go`
- Check for compilation errors
- Verify method name follows convention
- Check nil safety for all pointer fields
- Ensure consistency with other service methods

**Verify**: No errors, follows project patterns.

---

## Phase 2: DataSource Schema Definition (20 tasks, ~45 min)

### Task 2.1: Create datasource file
**File**: `tencentcloud/services/lighthouse/data_source_tc_lighthouse_blueprints.go`  
**Action**: Create new file with basic structure

```go
package lighthouse

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudLighthouseBlueprints() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceTencentCloudLighthouseBlueprintsRead,
		Schema: map[string]*schema.Schema{},
	}
}

func dataSourceTencentCloudLighthouseBlueprintsRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}
```

**Verify**: File compiles with basic structure.

---

### Task 2.2: Add blueprint_ids argument
**File**: `data_source_tc_lighthouse_blueprints.go`  
**Action**: Add schema field

```go
"blueprint_ids": {
	Optional: true,
	Type:     schema.TypeSet,
	Elem: &schema.Schema{
		Type: schema.TypeString,
	},
	Description: "Blueprint ID list.",
},
```

**Verify**: Field defined correctly.

---

### Task 2.3: Add filters argument
**File**: `data_source_tc_lighthouse_blueprints.go`  
**Action**: Add filters schema (similar to bundle datasource)

```go
"filters": {
	Optional: true,
	Type:     schema.TypeList,
	Description: "Filter list.\n" +
		"- `blueprint-id`: Filter by blueprint ID.\n" +
		"- `blueprint-type`: Filter by blueprint type. Values: `APP_OS`, `PURE_OS`, `DOCKER`, `PRIVATE`, `SHARED`.\n" +
		"- `platform-type`: Filter by platform type. Values: `LINUX_UNIX`, `WINDOWS`.\n" +
		"- `blueprint-name`: Filter by blueprint name.\n" +
		"- `blueprint-state`: Filter by blueprint state.\n" +
		"- `scene-id`: Filter by scene ID.\n" +
		"NOTE: The upper limit of Filters per request is 10. The upper limit of Filter.Values is 100. Parameter does not support specifying both BlueprintIds and Filters.",
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Field to be filtered.",
			},
			"values": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "Filter value of field.",
			},
		},
	},
},
```

**Verify**: Filters schema matches API documentation.

---

### Task 2.4: Add blueprint_set computed output
**File**: `data_source_tc_lighthouse_blueprints.go`  
**Action**: Add blueprint_set schema with all Blueprint fields

```go
"blueprint_set": {
	Computed:    true,
	Type:        schema.TypeList,
	Description: "List of blueprint details.",
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			// TODO: Add all Blueprint fields
		},
	},
},
```

**Verify**: Output structure defined.

---

### Task 2.5-2.22: Add all Blueprint fields to blueprint_set
**File**: `data_source_tc_lighthouse_blueprints.go`  
**Action**: Add each field individually (18 fields total)

Fields to add:
1. `blueprint_id` (string, required field)
2. `display_title` (string)
3. `display_version` (string)
4. `description` (string, nullable)
5. `os_name` (string)
6. `platform` (string)
7. `platform_type` (string)
8. `blueprint_type` (string)
9. `image_url` (string)
10. `required_system_disk_size` (int)
11. `blueprint_state` (string)
12. `created_time` (string, nullable)
13. `blueprint_name` (string)
14. `support_automation_tools` (bool)
15. `required_memory_size` (int)
16. `image_id` (string, nullable)
17. `community_url` (string)
18. `guide_url` (string)
19. `scene_id_set` ([]string, nullable)
20. `docker_version` (string, nullable)

**Example**:
```go
"blueprint_id": {
	Type:        schema.TypeString,
	Computed:    true,
	Description: "Blueprint ID, which is the unique identifier of Blueprint.",
},
```

**Verify**: All 18+ fields defined with correct types and descriptions.

---

### Task 2.23: Add result_output_file argument
**File**: `data_source_tc_lighthouse_blueprints.go`  
**Action**: Add standard output file field

```go
"result_output_file": {
	Type:        schema.TypeString,
	Optional:    true,
	Description: "Used to save results.",
},
```

**Verify**: Standard field added.

---

## Phase 3: DataSource Read Function (15 tasks, ~45 min)

### Task 3.1: Implement read function basics
**File**: `data_source_tc_lighthouse_blueprints.go`  
**Action**: Add defer statements and context

```go
func dataSourceTencentCloudLighthouseBlueprintsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_lighthouse_blueprints.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	
	// TODO: Build paramMap, call service, map response
	
	return nil
}
```

**Verify**: Function structure correct.

---

### Task 3.2: Parse blueprint_ids input
**File**: `data_source_tc_lighthouse_blueprints.go`  
**Action**: Convert schema input to API format

```go
if v, ok := d.GetOk("blueprint_ids"); ok {
	blueprintIdsSet := v.(*schema.Set).List()
	paramMap["BlueprintIds"] = helper.InterfacesStringsPoint(blueprintIdsSet)
}
```

**Verify**: IDs parsed correctly.

---

### Task 3.3: Parse filters input
**File**: `data_source_tc_lighthouse_blueprints.go`  
**Action**: Convert filters to API format

```go
if v, ok := d.GetOk("filters"); ok {
	filtersList := v.([]interface{})
	filters := make([]*lighthouse.Filter, 0, len(filtersList))
	for _, item := range filtersList {
		filterMap := item.(map[string]interface{})
		filter := lighthouse.Filter{}
		if v, ok := filterMap["name"]; ok {
			filter.Name = helper.String(v.(string))
		}
		if v, ok := filterMap["values"]; ok {
			valuesSet := v.(*schema.Set).List()
			filter.Values = helper.InterfacesStringsPoint(valuesSet)
		}
		filters = append(filters, &filter)
	}
	paramMap["Filters"] = filters
}
```

**Verify**: Filters parsed correctly.

---

### Task 3.4: Call service layer
**File**: `data_source_tc_lighthouse_blueprints.go`  
**Action**: Invoke DescribeLighthouseBlueprintsByFilter with retry

```go
service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

var blueprintSet []*lighthouse.Blueprint

err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
	result, e := service.DescribeLighthouseBlueprintsByFilter(ctx, paramMap)
	if e != nil {
		return tccommon.RetryError(e)
	}
	blueprintSet = result
	return nil
})

if err != nil {
	return err
}
```

**Verify**: Service called with retry logic.

---

### Task 3.5: Map response to schema (initialize)
**File**: `data_source_tc_lighthouse_blueprints.go`  
**Action**: Create result lists

```go
ids := make([]string, 0, len(blueprintSet))
tmpList := make([]map[string]interface{}, 0, len(blueprintSet))
```

**Verify**: Lists initialized.

---

### Task 3.6-3.14: Map each Blueprint field
**File**: `data_source_tc_lighthouse_blueprints.go`  
**Action**: Loop through blueprintSet and map all fields

```go
if blueprintSet != nil {
	for _, blueprint := range blueprintSet {
		blueprintMap := map[string]interface{}{}

		if blueprint.BlueprintId != nil {
			blueprintMap["blueprint_id"] = blueprint.BlueprintId
			ids = append(ids, *blueprint.BlueprintId)
		}

		if blueprint.DisplayTitle != nil {
			blueprintMap["display_title"] = blueprint.DisplayTitle
		}

		// ... map all other fields with nil checks ...

		if blueprint.SceneIdSet != nil {
			blueprintMap["scene_id_set"] = blueprint.SceneIdSet
		}

		if blueprint.DockerVersion != nil {
			blueprintMap["docker_version"] = blueprint.DockerVersion
		}

		tmpList = append(tmpList, blueprintMap)
	}

	_ = d.Set("blueprint_set", tmpList)
}
```

**Verify**: All 18+ fields mapped with nil checks.

---

### Task 3.15: Set resource ID
**File**: `data_source_tc_lighthouse_blueprints.go`  
**Action**: Generate ID from blueprint IDs

```go
d.SetId(helper.DataResourceIdsHash(ids))
```

**Verify**: ID set correctly.

---

### Task 3.16: Handle result_output_file
**File**: `data_source_tc_lighthouse_blueprints.go`  
**Action**: Write results to file if specified

```go
output, ok := d.GetOk("result_output_file")
if ok && output.(string) != "" {
	if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
		return e
	}
}
return nil
```

**Verify**: Output file written when specified.

---

## Phase 4: Testing (9 tasks, ~1 hour)

### Task 4.1: Create test file
**File**: `tencentcloud/services/lighthouse/data_source_tc_lighthouse_blueprints_test.go`  
**Action**: Create test file with basic structure

```go
package lighthouse_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseBlueprintsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseBlueprintsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_lighthouse_blueprints.blueprints"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_lighthouse_blueprints.blueprints", "blueprint_set.#"),
				),
			},
		},
	})
}

const testAccLighthouseBlueprintsDataSource = `
data "tencentcloud_lighthouse_blueprints" "blueprints" {
}
`
```

**Verify**: Test compiles and structure is correct.

---

### Task 4.2: Add filter test
**File**: `data_source_tc_lighthouse_blueprints_test.go`  
**Action**: Add test case with filters

```go
func TestAccTencentCloudLighthouseBlueprintsDataSource_filter(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseBlueprintsDataSourceFilter,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_lighthouse_blueprints.linux"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_lighthouse_blueprints.linux", "blueprint_set.#"),
				),
			},
		},
	})
}

const testAccLighthouseBlueprintsDataSourceFilter = `
data "tencentcloud_lighthouse_blueprints" "linux" {
  filters {
    name   = "platform-type"
    values = ["LINUX_UNIX"]
  }
}
`
```

**Verify**: Filter test defined.

---

### Task 4.3: Add multiple filters test
**File**: `data_source_tc_lighthouse_blueprints_test.go`  
**Action**: Test with multiple filters

```go
const testAccLighthouseBlueprintsDataSourceMultiFilter = `
data "tencentcloud_lighthouse_blueprints" "app_linux" {
  filters {
    name   = "blueprint-type"
    values = ["APP_OS"]
  }
  
  filters {
    name   = "platform-type"
    values = ["LINUX_UNIX"]
  }
}
`
```

**Verify**: Multiple filter test added.

---

### Task 4.4: Compile tests
**Action**: Run `go test -c ./tencentcloud/services/lighthouse/... -o /dev/null`

**Verify**: Tests compile without errors.

---

### Task 4.5: Run unit tests (if applicable)
**Action**: Run tests without TF_ACC

**Verify**: No test failures.

---

### Task 4.6-4.9: Acceptance test checks
**Actions**:
- Verify test names follow convention: `TestAccTencentCloud*`
- Check all assertions are meaningful
- Ensure tests are parallel (`t.Parallel()`)
- Add test for checking specific fields exist

**Verify**: All test best practices followed.

---

## Phase 5: Documentation (8 tasks, ~30 min)

### Task 5.1: Create source documentation
**File**: `tencentcloud/services/lighthouse/data_source_tc_lighthouse_blueprints.md`  
**Action**: Create markdown documentation

```markdown
Provides a list of Lighthouse blueprints (images).

Use this data source to query available blueprints for Lighthouse instances.

Example Usage

Query all blueprints:

\`\`\`hcl
data "tencentcloud_lighthouse_blueprints" "all" {
}

output "blueprints" {
  value = data.tencentcloud_lighthouse_blueprints.all.blueprint_set
}
\`\`\`

Filter by platform type:

\`\`\`hcl
data "tencentcloud_lighthouse_blueprints" "linux" {
  filters {
    name   = "platform-type"
    values = ["LINUX_UNIX"]
  }
}
\`\`\`

Filter by blueprint type:

\`\`\`hcl
data "tencentcloud_lighthouse_blueprints" "app_os" {
  filters {
    name   = "blueprint-type"
    values = ["APP_OS"]
  }
}
\`\`\`

Query specific blueprints by ID:

\`\`\`hcl
data "tencentcloud_lighthouse_blueprints" "specific" {
  blueprint_ids = ["lhbp-xxx", "lhbp-yyy"]
}
\`\`\`

## Argument Reference

The following arguments are supported:

* `blueprint_ids` - (Optional, Set: [`String`]) Blueprint ID list.
* `filters` - (Optional, List) Filter list. Cannot be used with `blueprint_ids`. Each filter supports:
  - `name` - (Required, String) Filter name. Valid values: `blueprint-id`, `blueprint-type`, `platform-type`, `blueprint-name`, `blueprint-state`, `scene-id`.
  - `values` - (Required, Set) Filter values.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `blueprint_set` - List of blueprint details. Each element contains:
  - `blueprint_id` - Blueprint ID.
  - `display_title` - Blueprint display title.
  - `display_version` - Blueprint display version.
  - `description` - Blueprint description.
  - `os_name` - Operating system name.
  - `platform` - Operating system platform.
  - `platform_type` - Platform type (LINUX_UNIX or WINDOWS).
  - `blueprint_type` - Blueprint type (APP_OS, PURE_OS, DOCKER, PRIVATE, or SHARED).
  - `image_url` - Blueprint image URL.
  - `required_system_disk_size` - Required system disk size in GB.
  - `blueprint_state` - Blueprint state.
  - `created_time` - Creation time (ISO 8601 format).
  - `blueprint_name` - Blueprint name.
  - `support_automation_tools` - Whether the blueprint supports automation tools.
  - `required_memory_size` - Required memory size in GB.
  - `image_id` - CVM image ID (if shared from CVM).
  - `community_url` - Community URL.
  - `guide_url` - Guide documentation URL.
  - `scene_id_set` - List of associated scene IDs.
  - `docker_version` - Docker version (for Docker blueprints).
```

**Verify**: Documentation complete with all fields and examples.

---

### Task 5.2: Generate website documentation
**Action**: Run `make doc` to generate HTML documentation

**Verify**: `website/docs/d/lighthouse_blueprints.html.markdown` generated.

---

### Task 5.3: Review generated documentation
**Action**: Check generated HTML documentation is correct

**Verify**: All fields, examples present and formatted correctly.

---

### Task 5.4-5.8: Documentation quality checks
**Actions**:
- Check all examples are valid HCL
- Verify filter names match API exactly
- Ensure descriptions are clear and helpful
- Check links/references are correct
- Spell check documentation

**Verify**: Documentation is high quality and accurate.

---

## Phase 6: Code Quality & Integration (5 tasks, ~15 min)

### Task 6.1: Format all code
**Action**: Run `gofmt -w tencentcloud/services/lighthouse/`

**Verify**: All files formatted consistently.

---

### Task 6.2: Run linter
**Action**: Run `golangci-lint run tencentcloud/services/lighthouse/data_source_tc_lighthouse_blueprints.go`

**Verify**: No linter errors or warnings.

---

### Task 6.3: Compile provider
**Action**: Run `go build -o /tmp/terraform-provider-tencentcloud .`

**Verify**: Provider compiles successfully.

---

### Task 6.4: Register datasource
**File**: `tencentcloud/provider.go` (or relevant registration file)  
**Action**: Add datasource to provider's datasource map

**Note**: Need to check where datasources are registered in this project.

**Verify**: Datasource registered and accessible.

---

### Task 6.5: Run all lighthouse tests
**Action**: Run `go test ./tencentcloud/services/lighthouse/... -v`

**Verify**: All tests pass (including new tests).

---

## Phase 7: Final Validation (5 tasks, ~15 min)

### Task 7.1: Manual acceptance test
**Action**: Set TF_ACC=1 and run: `go test -v -run TestAccTencentCloudLighthouseBlueprintsDataSource`

**Verify**: Acceptance tests pass with real API.

---

### Task 7.2: Test with Terraform CLI
**Action**: Create test configuration and run `terraform plan`

```hcl
terraform {
  required_providers {
    tencentcloud = {
      source = "tencentcloudstack/tencentcloud"
    }
  }
}

provider "tencentcloud" {}

data "tencentcloud_lighthouse_blueprints" "test" {
  filters {
    name   = "platform-type"
    values = ["LINUX_UNIX"]
  }
}

output "blueprints" {
  value = data.tencentcloud_lighthouse_blueprints.test.blueprint_set
}
```

**Verify**: Terraform can read the datasource and returns data.

---

### Task 7.3: Verify pagination
**Action**: Test with scenarios that return >100 results (if possible)

**Verify**: All results returned, pagination transparent.

---

### Task 7.4: Verify filter combinations
**Action**: Test multiple filter combinations

**Verify**: Filters work correctly.

---

### Task 7.5: Final checklist
**Actions**:
- [ ] All code formatted and passes linting
- [ ] All tests pass (unit and acceptance)
- [ ] Documentation complete and generated
- [ ] No compilation errors or warnings
- [ ] Datasource registered in provider
- [ ] Manual testing successful
- [ ] Code follows project conventions
- [ ] All Blueprint fields mapped correctly
- [ ] Pagination works (offset/limit hidden)
- [ ] Filters work correctly
- [ ] No nil pointer dereferences

**Verify**: All checkboxes checked.

---

## Validation Checklist

### Code Quality
- [ ] All files formatted with gofmt
- [ ] No golangci-lint errors or warnings
- [ ] Provider compiles successfully
- [ ] All imports are used
- [ ] No commented-out code
- [ ] Consistent naming conventions

### Functionality
- [ ] Datasource queries all blueprints
- [ ] Filters work (type, platform, state, name, scene)
- [ ] Blueprint IDs query works
- [ ] Cannot use both IDs and filters (API constraint respected)
- [ ] Pagination is automatic and hidden
- [ ] All 18+ Blueprint fields mapped
- [ ] Nil checks for all pointer fields
- [ ] result_output_file works

### Testing
- [ ] Unit tests pass
- [ ] Acceptance tests pass
- [ ] Tests follow naming convention
- [ ] Tests are parallel
- [ ] Test coverage is adequate

### Documentation
- [ ] Source documentation complete (.md)
- [ ] Website documentation generated (.html.markdown)
- [ ] All examples are valid
- [ ] All fields documented
- [ ] Filter descriptions match API

### Integration
- [ ] Datasource registered in provider
- [ ] No conflicts with existing datasources
- [ ] Follows project patterns
- [ ] Service layer follows conventions
- [ ] Error handling is consistent

---

## Notes

1. **Offset/Limit**: Hidden from users - handled automatically in service layer
2. **Filters**: Cannot use with blueprint_ids (API limitation) - document clearly
3. **Nil Safety**: All pointer fields must have nil checks
4. **Pagination**: Loop until all results fetched (API returns max 100 per call)
5. **Consistency**: Follow patterns from existing datasources (bundle, zone, region)

---

**Total Estimated Time**: 3.5 hours (broken down by phase above)
