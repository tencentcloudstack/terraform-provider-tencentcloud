# Migration Examples

This document shows examples of migrating from the old pagination-based approach to the new automatic retrieval approach.

## Before (Old Approach - Manual Pagination)

### Example 1: Basic Query with Pagination

```hcl
# Old way - users had to manually handle pagination
data "tencentcloud_ckafka_instances" "example" {
  offset = 0
  limit  = 50  # User must choose appropriate limit
}

# Problem: If there are > 50 instances, user needs to manually paginate
# by calling this data source multiple times with different offsets
```

### Example 2: Filtered Query with Pagination

```hcl
# Old way - filtering with manual pagination
data "tencentcloud_ckafka_instances" "vpc_instances" {
  filters {
    name   = "VpcId"
    values = ["vpc-abc123"]
  }
  
  offset = 0
  limit  = 100
}

# Problem: User must be aware of pagination to get all results
```

## After (New Approach - Automatic Retrieval)

### Example 1: Basic Query (Automatic Pagination)

```hcl
# New way - all instances retrieved automatically
data "tencentcloud_ckafka_instances" "example" {
  # No offset or limit needed!
  # The data source automatically retrieves all instances
}

# Benefit: No risk of missing data due to pagination limits
```

### Example 2: Filtered Query (Automatic Pagination)

```hcl
# New way - filtering with automatic retrieval
data "tencentcloud_ckafka_instances" "vpc_instances" {
  filters {
    name   = "VpcId"
    values = ["vpc-abc123"]
  }
  
  # offset and limit are deprecated
  # All matching instances are retrieved automatically
}

# Benefit: Simpler configuration, guaranteed complete results
```

### Example 3: Using Multiple Filters

```hcl
data "tencentcloud_ckafka_instances" "running_instances" {
  status = [1]  # Running instances only
  
  filters {
    name   = "InstanceType"
    values = ["profession"]
  }
  
  filters {
    name   = "VpcId"
    values = ["vpc-abc123", "vpc-def456"]
  }
  
  search_word = "prod"
  
  # No pagination parameters needed
  # All matching instances retrieved automatically
}
```

## Deprecation Warnings

If you still use `offset` or `limit` parameters, you'll see warnings like:

```
Warning: Argument is deprecated

  with data.tencentcloud_ckafka_instances.example,
  on main.tf line 5, in data "tencentcloud_ckafka_instances" "example":
   5:   offset = 0

This parameter is deprecated and will be removed in a future version. The
data source now automatically retrieves all instances.
```

## Migration Checklist

- [ ] Remove `offset` parameter from all `tencentcloud_ckafka_instances` data sources
- [ ] Remove `limit` parameter from all `tencentcloud_ckafka_instances` data sources
- [ ] Test that all expected instances are retrieved
- [ ] Update any scripts or automation that relied on manual pagination

## Benefits of Migration

1. **Simpler Configuration**: No need to manage pagination parameters
2. **Complete Data**: Guaranteed to retrieve all instances, no risk of missing data
3. **Better Performance**: Optimized page size (100 vs default 10)
4. **Reliability**: Built-in retry logic handles transient API failures
5. **Consistency**: Aligns with other Terraform provider data sources

## Backward Compatibility

✅ **Your existing configurations will continue to work** during the deprecation period. However, you'll see warnings encouraging you to remove the deprecated parameters.

The `offset` and `limit` parameters are now ignored internally - the data source always retrieves all instances regardless of these values.
