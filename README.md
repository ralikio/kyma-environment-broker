[![REUSE status](https://api.reuse.software/badge/github.com/kyma-project/kyma-environment-broker)](https://api.reuse.software/info/github.com/kyma-project/kyma-environment-broker)
# Kyma Environment Broker

## Overview

Kyma Environment Broker (KEB) is a component that allows you to provision [SAP BTP, Kyma runtime](https://kyma-project.io/#/?id=kyma-and-sap-btp-kyma-runtime) on clusters provided by third-party providers. In the process, KEB first uses Provisioner to create a cluster. Then, it uses Reconciler and Lifecycle Manager to install Kyma runtime on the cluster.

## Configuration

KEB binary allows you to override some configuration parameters. You can specify the following environment variables:

| Name | Description | Default value |
|-----|---------|:--------:|
| **APP_PORT** | Specifies the port on which the HTTP server listens. | `8080` |
| **APP_STATUS_PORT** | | |

| **APP_KYMA_VERSION** | Specifies the default Kyma version. | None |

| **APP_VERSION_CONFIG_NAMESPACE** | Defines the namespace with the ConfigMap that contains Kyma versions for global accounts configuration. | None |
| **APP_VERSION_CONFIG_NAME** | Defines the name of the ConfigMap that contains Kyma versions for global accounts configuration. | None |

| **APP_TRIAL_REGION_MAPPING_FILE_PATH** | Defines a path to the file which contains a mapping between the platform region and the Trial plan region. | None |

| **APP_PROFILER_MEMORY** | Enables memory profiling every sampling period with the default location `/tmp/profiler`, backed by a persistent volume. | `false` |

| **APP_BROKER_EXPOSE_SCHEMA_WITH_REGION_REQUIRED** | Extends Open Service Broker API schema with additional `region` parameter for specifying what region Kyma is provisioned in. | `false` |
| **APP_BROKER_REGION_PARAMETER_IS_REQUIRED** | TODO: what is the difference between this and APP_BROKER_EXPOSE_SCHEMA_WITH_REGION_REQUIRED parameter? | |
| **APP_DISABLE_PROCESS_OPERATIONS_IN_PROGRESS** | DisableProcessOperationsInProgress allows to disable processing operations which are in progress on starting application. Set to true if you are running in a separate testing deployment but with the production DB. | `false` |

| **APP_BROKER_ONLY_SINGLE_TRIAL_PER_GA** | Allows to specify if more than one trial instance can be provisioner in given SAP BTP global account | `true` |
| **APP_BROKER_URL** | URL used for concatenation of KEB URL returned in response body returned from KEB API. | `kyma-env-broker.localhost` |

| **APP_BROKER_INCLUDE_ADDITIONAL_PARAMS_IN_SCHEMA** | If set to true Open Service Broker API schema is rendered with addtional `oidc` and `administrators` parameters. | |
| **APP_BROKER_SHOW_TRIAL_EXPIRATION_INFO** | Conditionally enrich trial instances returned from KEB API with expiration data based on accounts' ids defined in `APP_BROKER_SUBACCOUNTS_IDS_TO_SHOW_TRIAL_EXPIRATION_INFO` parameter. | `false` |
| **APP_BROKER_SUBACCOUNTS_IDS_TO_SHOW_TRIAL_EXPIRATION_INFO** | List of BTP Subaccounts that to show additional data for (if the are using trial plan) | |
| **APP_BROKER_TRIAL_DOCS_URL** | URL to trial docs returned from KEB API for expired clusters. | |


| **APP_EU_ACCESS_WHITELISTED_GLOBAL_ACCOUNTS_FILE_PATH** | Path to a file with whitelisted accounts that can provision runtime in eu restricted regions. | |
| **APP_EU_ACCESS_REJECTION_MESSAGE** | Message returned from KEB API in case account is not whitelisted but attempt to create an instance in restricted region has been made for it. | |

| **APP_FREEMIUM_PROVIDERS** | | |
| **APP_CATALOG_FILE_PATH** | | |

| **APP_DEFAULT_REQUEST_REGION** | Values used when region is not given in the request context. | |
| **APP_UPDATE_PROCESSING_ENABLED** | | |
| **APP_DOMAIN_NAME** | Used for rendering domain inside of Swagger API specification. | |
| **APP_SKR_OIDC_DEFAULT_VALUES_YAML_FILE_PATH** | Path to file with default oidc configuration. | |
| **APP_SKR_DNS_PROVIDERS_VALUES_YAML_FILE_PATH** | Default DNS provider configuration file path.  | |
| **APP_ORCHESTRATION_CONFIG_NAMESPACE** | Namespace of where maintenance configuration rules are defined in. | |
| **APP_ORCHESTRATION_CONFIG_NAME** | Name of the config map that maintenance configuration rules are defined in. | |
| **APP_KYMA_DASHBOARD_CONFIG_LANDSCAPE_URL** | Busola URL returned in response body to KEB API request. | |

### Infrastructure Manager

| **APP_INFRASTRUCTURE_MANAGER_INTEGRATION_DISABLED** | A feature flag for enabling integration with Infrastructure Manager. | |

### Gardener

| **APP_GARDENER_PROJECT** | Defines the project in which the cluster is created. | `kyma-dev` |
| **APP_GARDENER_SHOOT_DOMAIN** | Defines the domain for clusters created in Gardener. | `shoot.canary.k8s-hana.ondemand.com` |
| **APP_GARDENER_KUBECONFIG_PATH** | Defines the path to the kubeconfig file for Gardener. | `/gardener/kubeconfig/kubeconfig` |

### Lifecycle Manager

| **APP_LIFECYCLE_MANAGER_INTEGRATION_DISABLED** | A feature flag for enabling integration with Lifecycle Manager. | |

### Provisioner


| **APP_PROVISIONER_DEFAULT_GARDENER_SHOOT_PURPOSE** | Specifies the purpose of the created cluster. The possible values are: `development`, `evaluation`, `production`, `testing`. | `development` |
| **APP_PROVISIONER_URL** | Specifies a URL to the Runtime Provisioner's API. | None |
| **APP_PROVISIONER_MACHINE_IMAGE** | Defines the Gardener machine image used in a provisioned node. | None |
| **APP_PROVISIONER__MACHINE_IMAGE_VERSION** | Defines the Gardener image version used in a provisioned cluster. | None |

| **APP_PROVISIONER_SAP_CONVERGED_CLOUD_FLOATING_POOL_NAME** | Name of a floating pool name to use for Converged Cloud plans. | |
| **APP_PROVISIONER_DEFAULT_TRIAL_PROVIDER** | Default hyperscaler for trial plans. | |
| **APP_PROVISIONER_KUBERNETES_VERSION** | | |
| **APP_PROVISIONER_TRIAL_NODES_NUMBER** | Number of nodes created for trial clusters. | |
| **APP_PROVISIONER_AUTO_UPDATE_KUBERNETES_VERSION** | Configuration send to Gardener for enabling K8S version update. | |
| **APP_PROVISIONER_AUTO_UPDATE_MACHINE_IMAGE_VERSION** | Configuration send to Gardener for enabling maching image update. | |
| **APP_PROVISIONER_MULTI_ZONE_CLUSTER** | ?If set to true provisioned cluster for preconfigured plans are setup as multizone clusters. ?| |
| **APP_PROVISIONER_CONTROL_PLANE_FAILURE_TOLERANCE** | | |

### Database

| **APP_DATABASE_USER** | Defines the database username. | `postgres` |
| **APP_DATABASE_PASSWORD** | Defines the database user password. | `password` |
| **APP_DATABASE_HOST** | Defines the database host. | `localhost` |
| **APP_DATABASE_PORT** | Defines the database port. | `5432` |
| **APP_DATABASE_NAME** | Defines the database name. | `broker` |
| **APP_DATABASE_SSLMODE** | Specifies the SSL Mode for PostgreSQL. See [all the possible values](https://www.postgresql.org/docs/9.1/libpq-ssl.html).  | `disable`|
| **APP_DATABASE_SSLROOTCERT** | Specifies the location of CA cert of PostgreSQL. (Optional)  | None |

| **APP_DATABASE_SECRET_KEY** | Cipher used for encryption/decryption of sensitive data. | |

### EDP



| **APP_EDP_AUTH_URL** | Authorization URL for EDP integration. | |
| **APP_EDP_ADMIN_URL** | An URL for API EDP integration. | |
| **APP_EDP_NAMESPACE** | EDP namespace used for request sent from KEB. | |
| **APP_EDP_ENVIRONMENT** | | |
| **APP_EDP_REQUIRED** | A flag that makes all EDP errors non-fatal allowing for further operation processing. | |
| **APP_EDP_SECRET** | API secret to EDP system. | |
| **APP_EDP_DISABLED** | A feature flag for EDP integration. | |

### IAS

| **APP_IAS_URL** | | |
| **APP_IAS_USER_ID** | | |
| **APP_IAS_USER_SECRET** | | |
| **APP_IAS_IDENTITY_PROVIDER** | | |
| **APP_IAS_TLS_RENEGOTIATION_ENABLE** | | |
| **APP_IAS_TLS_SKIP_CERT_VERIFICATION** | | |
| **APP_IAS_DISABLED** | | |


### AVS

| **APP_AVS_REGION_TAG_CLASS_ID** | Specifies the **TagClassId** of the tag that contains Gardener cluster's region. | None |
| **APP_AVS_OAUTH_TOKEN_ENDPOINT** | | |
| **APP_AVS_OAUTH_USERNAME** | | |
| **APP_AVS_OAUTH_PASSWORD** | | |
| **APP_AVS_API_ENDPOINT** | | |
| **APP_AVS_OAUTH_CLIENT_ID** | | |
| **APP_AVS_API_KEY** | | |
| **APP_AVS_INTERNAL_TESTER_ACCESS_ID** | | |
| **APP_AVS_EXTERNAL_TESTER_ACCESS_ID** | | |
| **APP_AVS_INTERNAL_TESTER_SERVICE** | | |
| **APP_AVS_EXTERNAL_TESTER_SERVICE** | | |
| **APP_AVS_GROUP_ID** | | |
| **APP_AVS_PARENT_ID** | | |
| **APP_AVS_TRIAL_API_KEY** | | |
| **APP_AVS_TRIAL_INTERNAL_TESTER_ACCESS_ID** | | |
| **APP_AVS_TRIAL_GROUP_ID** | | |
| **APP_AVS_TRIAL_PARENT_ID** | | |
| **APP_AVS_INSTANCE_ID_TAG_CLASS_ID** | | |
| **APP_AVS_GLOBAL_ACCOUNT_ID_TAG_CLASS_ID** | | |
| **APP_AVS_SUB_ACCOUNT_ID_TAG_CLASS_ID** | | |
| **APP_AVS_LANDSCAPE_TAG_CLASS_ID** | | |
| **APP_AVS_EXTERNAL_TESTER_DISABLED** | | |
| **APP_AVS_PROVIDER_TAG_CLASS_ID** | | |
| **APP_AVS_SHOOT_NAME_TAG_CLASS_ID** | | |
| **APP_AVS_MAINTENANCE_MODE_DURING_UPGRADE_DISABLED** | | |
| **APP_AVS_MAINTENANCE_MODE_DURING_UPGRADE_ALWAYS_DISABLED_GLOBAL_ACCOUNTS_FILE_PATH** | | |

### Director


| **APP_DIRECTOR_URL** | Specifies the Director's URL. | `http://compass-director.compass-system.svc.cluster.local:3000/graphql` |
| **APP_DIRECTOR_OAUTH_TOKEN_URL** | Specifies the URL for OAuth authentication. | None |
| **APP_DIRECTOR_OAUTH_CLIENT_ID** | Specifies the client ID for OAuth authentication. | None |
| **APP_DIRECTOR_OAUTH_SCOPE** | Specifies the scopes for OAuth authentication. | `runtime:read runtime:write` |

| **APP_DIRECTOR_DEFAULT_TENANT** | | |
| **APP_DIRECTOR_OAUTH_CLIENT_SECRET** | | |

### Feature Flags

| **APP_ENABLE_ON_DEMAND_VERSION** | If set to `true`, a user can specify a Kyma version in a provisioning request. | `false` |
| **APP_BROKER_ENABLE_PLANS** | EnablePlans defines the plans that should be available for provisioning. | `azure,gcp,azure_lite,trial` |
| **APP_BROKER_ENABLE_KUBECONFIG_URL_LABEL** | A feature flag for attaching url for Kubeconfig download in the response of KEB API calls. | `false` |

### Timeouts

| **APP_OPERATION_TIMEOUT** | Maximum amount of time that given operation can take after which it fails. | |
| **APP_PROVISIONER_PROVISIONING_TIMEOUT** | | |
| **APP_PROVISIONER_DEPROVISIONING_TIMEOUT** | | |

### Reconciler

| **APP_RECONCILER_URL** | URL to Reconciler used for forwarding Kyma installation to. | |
| **APP_RECONCILER_INTEGRATION_DISABLED** | A feature flag for enabling integration with Reconciler. | |
| **APP_RECONCILER_PROVISIONING_TIMEOUT** | | |


## Read More

To learn more about how to use KEB, read the documentation in the [`user`](./docs/user/) directory.
For more technical details on KEB, go to the [`contributor`](./docs/contributor/) directory.

## Contributing
<!--- mandatory section - do not change this! --->

See the [Contributing](CONTRIBUTING.md) guidelines.

## Code of Conduct
<!--- mandatory section - do not change this! --->

See the [Code of Conduct](CODE_OF_CONDUCT.md) document.

## Licensing
<!--- mandatory section - do not change this! --->

See the [license](./LICENSE) file.
