# Fix MongoDB SRV Connection Domain Suffix Truncation

## What

In `resourceTencentCloudMongodbInstanceSrvConnectionRead`, truncate the suffix from the `domain` value returned by the API before writing it to state.

The user inputs a plain custom prefix such as `123asasdfasdf`, but the API returns the fully-qualified domain `123asasdfasdf.gz.tencentmdb.com`. Storing the full FQDN in state causes a perpetual diff because Terraform compares it against the user-supplied short prefix.

## Why

- State drift: every plan after the first apply would show a diff on `domain`, degrading UX.
- The user-facing contract for `domain` is the short custom prefix, not the system-appended FQDN suffix.

## Root Cause

`DescribeSRVConnectionDomain` returns the full FQDN (e.g. `<prefix>.gz.tencentmdb.com`). The Read handler writes this verbatim to state, but the user only ever writes the prefix.

## Fix Strategy

After receiving the API response, strip everything from the first `.` onward, keeping only the prefix segment. Store only that prefix in state.
