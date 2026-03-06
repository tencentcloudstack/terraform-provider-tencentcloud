# Implementation Tasks

## 1. Schema Definition
- [x] 1.1 Add `band_width` field to resource schema
  - Type: `schema.TypeInt`
  - Optional: true
  - Description: "Public network bandwidth in Mbps. Only takes effect when enable_public_access is true."
- [x] 1.2 Add `enable_public_access` field to resource schema
  - Type: `schema.TypeBool`
  - Optional: true
  - Description: "Whether to enable public network access. Default is false."

## 2. Create Operation
- [x] 2.1 Add `band_width` field handling in Create function
  - Check if field exists with `d.GetOkExists("band_width")`
  - Set `request.Bandwidth` to `helper.IntUint64(v.(int))`
- [x] 2.2 Add `enable_public_access` field handling in Create function
  - Check if field exists with `d.GetOkExists("enable_public_access")`
  - Set `request.EnablePublicAccess` to `helper.Bool(v.(bool))`

## 3. Read Operation
- [x] 3.1 Read `band_width` from API response
  - Extract from `rabbitmqVipInstance.ClusterSpecInfo.PublicNetworkTps`
  - Set to state with `d.Set("band_width", value)`
- [x] 3.2 Read `enable_public_access` from API response
  - Extract from `rabbitmqVipInstance.ClusterNetInfo.PublicDataStreamStatus`
  - Convert string to bool: "ON" → true, "OFF" → false
  - Set to state with `d.Set("enable_public_access", value)`

## 4. Update Operation
- [x] 4.1 Add `band_width` to immutableArgs list
  - Ensures field cannot be modified after creation
- [x] 4.2 Add `enable_public_access` to immutableArgs list
  - Ensures field cannot be modified after creation

## 5. Documentation
- [x] 5.1 Update resource documentation file
  - Add `band_width` field description
  - Add `enable_public_access` field description
  - Add usage example showing public access configuration
  - Document field immutability constraint

## 6. Testing and Validation
- [x] 6.1 Run `go fmt` on modified files
- [ ] 6.2 Run `make lint` to check code quality
- [ ] 6.3 Manual testing: Create instance with public access enabled
- [ ] 6.4 Manual testing: Create instance with public access disabled (default)
- [ ] 6.5 Manual testing: Verify fields are read correctly from API
- [ ] 6.6 Manual testing: Verify update attempt fails with immutable error

## 7. Changelog
- [x] 7.1 Create changelog entry in `.changelog/` directory
  - Format: `resource/tencentcloud_tdmq_rabbitmq_vip_instance: add band_width and enable_public_access fields for public network access configuration`
