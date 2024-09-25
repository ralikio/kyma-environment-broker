# Kyma Environment Broker Endpoints

Kyma Environment Broker (KEB) implements the [Open Service Broker API (OSB API)](https://github.com/openservicebrokerapi/servicebroker/blob/v2.14/spec.md). All the OSB API endpoints are served with the following prefixes:

| Prefix | Description |
|---|---|
| `/oauth` | Defines a prefix for the endpoint secured with the OAuth2 authorization. The value for the SAP BTP region is specified under the **broker.defaultRequestRegion** parameter in the [`values.yaml`](https://github.com/kyma-project/kyma-environment-broker/blob/main/resources/keb/values.yaml) file. |
| `/oauth/{region}` | Defines a prefix for the endpoint secured with the OAuth2 authorization. The SAP BTP region value is specified in the request. |

> [!NOTE] 
> When the `{region}` value is one of EU Access BTP regions, the EU Access restrictions apply. For more information, see [EU Access](../contributor/03-20-eu-access.md).

Besides OSB API endpoints, KEB exposes the REST `/info/runtimes` endpoint that provides information about all created Runtimes, both succeeded and failed. This endpoint is secured with the OAuth2 authorization.

For more details on KEB APIs, see [this file](../../files/swagger/index.html).

## Kyma Binding Creation

One of the Broker API endpoints related to bindings allow for generation of credentials for accessing given service. The endpoints in question include all subpaths of `v2/service_instances/<service_id>/service_bindings`. In case of Kyma Environment Broker, the generated credentials are in the form of a Kubeconfig for managed SKR cluster. To generated a kubeconfig for a given SKR instance sent the following request to the Broker API:

```
PUT http://localhost:8080/oauth/v2/service_instances/{{instance_id}}/service_bindings/{{binding_id}}?accepts_incomplete=true&accepts_incomplete=true&service_id={{service_id}}&plan_id={{plan_id}}
Content-Type: application/json
X-Broker-API-Version: 2.14

{
  "service_id": "{{service_id}}",
  "plan_id": "{{plan_id}}"
  "parameters": {
  }
}
```

As a result of the above call the Broker will return a Kubeconfig file in the response body. The Kubeconfig file will contain the necessary information to access the managed SKR cluster. By default, Kyma Environment Broker uses [`shoots/adminkubeconfig`](https://github.com/gardener/gardener/blob/master/docs/usage/shoot_access.md#shootsadminkubeconfig-subresource) subresources to generate a Kubeconfig that uses certificates to authenticate its user. To customize format of returned kubeconfig the following parameters (used in the `parameters` section of the request body) can be used:

| Name | Default | Description |
|---|---|---|
| `token_request` | `false` | If set to `true` the Broker will return a kubeconfig with JWT token used as user authentication mechanism. The token itself ise generated using K8S's TokenRequest that is attached to a service account, clusterroel and clusterrolebinding all named `kyma-binding-{{binding_id}}`. Such approach allows to easily modify permissions granted to the Kubeconfig. |