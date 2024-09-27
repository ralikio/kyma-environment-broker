# SAP BTP, Kyma Runtime Operations

Kyma Environment Broker (KEB) allows you to configure operations you can run on SAP BTP, Kyma runtime. Each operation is processed by several steps arranged in stages and ordered in a queue. As every step can be re-launched multiple times, you should determine a behavior for each step in case of a processing failure. It can:

- Return an error, which interrupts the entire process, or skip step execution.
- Repeat the entire operation after a specified period.

> [!NOTE] 
> It's important to set lower timeouts for the Kyma installation in Runtime Provisioner.

## Stages

A stage is a grouping unit for steps. An operation can consist of multiple stages, and a stage can consist of multiple steps. Once all the steps in a stage are successfully executed, the stage is marked as finished and never repeated, even if the next stage fails. If a step fails at a given stage, the whole stage is repeated from the beginning.

## Provisioning

The provisioning process is executed when the instance is created, or an unsuspension is triggered.
Each provisioning step is responsible for a separate part of preparing Kyma runtime. For example, in a step you can provide tokens, credentials, or URLs to integrate SAP BTP, Kyma runtime with external systems.
You can find all the provisioning steps in the [provisioning](../../cmd/broker/provisioning.go) file.

> [!NOTE] 
> The timeout for processing this operation is set to `24h`.

## Deprovision

Each deprovisioning step is responsible for a separate part of cleaning Kyma runtime dependencies. To properly deprovision all the dependencies, you need the data used during the Kyma runtime provisioning. The first step finds the previous operation and copies the data.

None of the deprovisioning steps should block the entire deprovisioning operation. Use the `RetryOperationWithoutFail` function from the `DeprovisionOperationManager` struct to skip a step in case of a retry timeout. Set a 5-minute, at the most, timeout for retries in a step.
Only one step may fail the operation, namely `Check_Runtime_Removal`. It fails the operation in case of a timeout while checking for the Provisioner to remove the shoot. 
Once the step is successfully executed, it isn't retried (every deprovisioning step is defined in a separate stage). If a step has been skipped due to a retry timeout or error, the [Cron Job](../contributor/06-50-deprovision-retrigger-cronjob.md) tries to deprovision all remaining Kyma runtime dependencies again at a scheduled time.
You can find all the deprovisioning steps in the [deprovisioning](../../cmd/broker/deprovisioning.go) file.

> [!NOTE] 
> The timeout for processing this operation is set to `24h`.

## Update

The update process is triggered by an [OSB API update operation](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#updating-a-service-instance) request.
You can find all the updating steps in the [update](../../cmd/broker/update.go) file.

## Upgrade Cluster

The upgrade cluster process is triggered by upgrade cluster orchestration.
You can find all the upgrading cluster steps in the [upgrade_cluster](../../cmd/broker/upgrade_cluster.go) file.

## Provide Additional Steps

You can configure SAP BTP, Kyma runtime operations by providing additional steps. Every operation (see the implementation of `internal.Operation` structure in [model.go](../../internal/model.go)) is based on the same Operation structure. The following examples present how to extend the KEB process based on provisioning operation. Extensions for other processes follow the same steps but require their specific structures.

<div tabs name="runtime-provisioning-deprovisioning" group="runtime-provisioning-deprovisioning">
  <details>
  <summary label="provisioning">
  Provisioning
  </summary>

1. Create a new file in [this directory](../../internal/process/provisioning).

2. Implement the following interface in your provisioning or deprovisioning step:

    ```go
    type Step interface {
        Name() string
        Run(operation internal.Operation, logger logrus.FieldLogger) (internal.Operation, time.Duration, error)
    }
    ```

   - `Name()` method returns the name of the step that is used in logs.
   - `Run()` method implements the functionality of the step. The method receives operations as an argument to which it can add appropriate overrides or save other used variables. You must always return the modified operation from the method.

    ```go
    operation.InputCreator.AppendOverrides(COMPONENT_NAME, []*gqlschema.ConfigEntryInput{
        {
            Key:   "path.to.key",
            Value: SOME_VALUE,
        },
        {
            Key:    "path.to.secret",
            Value:  SOME_VALUE,
            Secret: ptr.Bool(true),
        },
    })
    ```

    If your functionality requires saving data in the storage, you can do it by adding fields to the generic `internal.Operation`, a specific implementation of that structure, or the InstanceDetails, all of which are defined in [model.go](../../internal/model.go). The difference is that for a specific operation implementation, new fields are only visible for that specific type and InstanceDetails is copied during operation initialization across all operations that concert a given runtime. The example below shows how to extend operation with additional fields:
    ```go
    type Operation struct {

        // These fields are serialized to JSON and stored in the storage
        RuntimeVersion RuntimeVersionData `json:"runtime_version"`

        // These fields are not stored in the storage
        InputCreator ProvisionerInputCreator `json:"-"`
    }
    ```

    See the example of the step implementation:

    ```go
    package provisioning

    import (
        "encoding/json"
        "net/http"
        "time"

        "github.com/kyma-incubator/compass/components/provisioner/pkg/gqlschema"
        "github.com/kyma-incubator/compass/components/kyma-environment-broker/internal"
        "github.com/kyma-incubator/compass/components/kyma-environment-broker/internal/storage"

        "github.com/sirupsen/logrus"
    )

    type HelloWorldStep struct {
        operationStorage storage.Operations
        client           *http.Client
    }

    type ExternalBodyResponse struct {
        data  string
        token string
    }

    func NewHelloWorldStep(operationStorage storage.Operations, client *http.Client) *HelloWorldStep {
        return &HelloWorldStep{
            operationStorage: operationStorage,
            client:           client,
        }
    }

    func (s *HelloWorldStep) Name() string {
        return "Hello_World"
    }

    // Your step can be repeated in case any other step fails, even if your step has already done its job
    func (s *HelloWorldStep) Run(operation internal.Operation, log *logrus.Entry) (internal.Operation, time.Duration, error) {
        log.Info("Start step")

        // Check whether your step should be run or if its job has been done in the previous iteration
        // All non-save operation data are empty (e.g. InputCreator overrides)

        // Add your logic here

        // Add a call to an external service (optional)
        response, err := s.client.Get("http://example.com")
        if err != nil {
            // Error during a call to an external service may be temporary so you should return time.Duration
            // All steps will be repeated in X seconds/minutes
            return operation, 1 * time.Second, nil
        }
        defer response.Body.Close()

        body := ExternalBodyResponse{}
        err = json.NewDecoder(response.Body).Decode(&body)
        if err != nil {
            log.Errorf("error: %s", err)
            // Handle a process failure by returning an error or time.Duration
        }

        // If a call or any other action is time-consuming, you can save the result in the operation
        // If you need an extra field in the Operation structure, add it first
        // In the following step, you can check beforehand if a given value already exists in the operation
        operation.HelloWorlds = body.data
        updatedOperation, err := s.operationStorage.UpdateOperation(operation)
        if err != nil {
            log.Errorf("error: %s", err)
            // Handle a process failure by returning an error or time.Duration
        }

        // If your step finishes with data which should be added to override used during the Runtime provisioning,
        // add an extra value to operation.InputCreator, then return the updated version of the Application
        updatedOperation.InputCreator.AppendOverrides("component-name", []*gqlschema.ConfigEntryInput{
            {
                Key:   "some.key",
                Value: body.token,
            },
        })

        // Return the updated version of the Application
        return *updatedOperation, 0, nil
    }
    ```

3. Add the step to the [`/cmd/broker/provisioning.go`](../../cmd/broker/provisioning.go) file:

    ```go
    provisioningSteps := []struct {
   		stage   string
   		step     provisioning.Step
   	}{
   		{
   			stage: "create_runtime",
   			step:   provisioning.NewHelloWorldStep(db.Operations(), &http.Client{}),
   		},
    }
    ```

   Once all the steps in the stage have run successfully, the stage is  not retried even if the application is restarted.
  </details>
