# Network Isolation

## Default Tenant Network Policies

Tenants will, by default, have L2 connectivitly between all namespaces owned by the tenant. Namespaces not owned by a tenant, will not have connectivity.

### Enable cross tenant connectivity

Supposed `Tenant` `A` wanted access to a database owned by `Tenant` `B`. In order to properly enable this, `Tenant` `A` would need to create a new `NetworkPolicy` in the pod's namespace allowing `egress` traffic to the pod. Addionally, `Tenant` `B` would need to create a `NetworkPolicy` to allow `ingress` trafic from `Tenant` `A`s application pod.

## Cluster Setup

Because pods in `Tenant` `Namespaces` will now be controlled by a `NetworkPolicy`, specific permissions need to be granted to those pods to access system level resources. The `Tenant` `Namespaces` will contain a `Network` policy that will allow traffic to/from `Namespaces` with the label `multitenant.k8s.io=public`.

By labeling `Namespaces` with this label, cluster admins will allow common resources to be accessed by all Tenants. For instance, to allow `Tenants` to continue to use `KubeDNS` in Minikube, the `kube-system` namespace needs to be appropriately labeled. Additional `NetworkPolicy`s have been applied to `kube-system` in the [default Minikube setup](minikube/)
