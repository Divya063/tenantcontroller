# Networking Tests

This test suite is designed to validate the `NetworkPolicy` objects defined [here](docs/networkpolicy.md) are consistent. The test framework performs the following operations:

# Step 00

Create two tenants, `foo` and `bar`. Creates `TenantNamespaces` `foo` and `baz` for Tenant foo, and TenantNamespace `bar` for Tenant bar.

# Step 01

Deploys an [Request Logging API](https://github.com/runyontr/request-log) in namespace `foo`.

# Step 02

Creates 3 jobs to query the API. One in namespace `foo`, one in namespace `baz` and one in namespace `bar`. The jobs in namespaces `foo` and `baz` should complete succesfully, as their run in the same `Tenant` as the API, but the job in namespace `baz` should fail since the API is hosted in a different Namespace.
