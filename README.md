# ctf-challenge-operator

ctf-challenge-operator is a Kubernetes operator designed to automate the deployment and management of Capture The Flag
(CTF) challenges within a Kubernetes cluster. It streamlines the process of running CTF events by handling challenge
lifecycle, configuration, and scaling, making it easier for organizers to manage and participants to engage with
challenges.

**NOTE: This project is currently in early development and is not yet in a state to be actually used in a CTF event
or even in a proof-of-concept situation.**

## Description

This project provides three Kubernetes custom resource definitions to help with running a CTF event:

- `ChallengeDescription`: This resource describes a single CTF challenge with the Kubernetes resources required to
  provision a single instance of that challenge. No workload is actually provisioned for this resource.
- `ChallengeInstance`: This resource is a specific instance of a ChallengeDescription with specific workload being
  provisioned. Each instance is provisioned into its own namespace and has an automated lifetime for cleaning up
  all workload provisioned for that instance.
- `APIKey`: This resource is an API key for accessing APIs. Automated lifetime management removes the APIKey once the
  lifetime is over.

## Getting Started

To deploy this operator into your Kubernetes cluster:

```shell
kubectl apply -k https://github.com/backbone81/ctf-challenge-operator/manifests?ref=main
```

**NOTE: As there is not yet a real release of this operator, the docker image referenced in that manifest does not
exist yet.**

### ChallengeDescription CR

The `ChallengeDescription` custom resource defines the template for a single CTF challenge. It specifies the metadata,
configuration, and Kubernetes resources required to provision an instance of the challenge, but does not itself create
any workloads. Instead, it acts as a blueprint that can be instantiated by `ChallengeInstance` resources. This resource
enables consistent, repeatable deployment of challenge instances.

For details about available fields, see [`api/v1alpha1/challenge_description.go`](api/v1alpha1/challenge_description.go).
For a concrete example, see [`examples/challenge-description-sample.yaml`](examples/challenge-description-sample.yaml).

### ChallengeInstance CR

The `ChallengeInstance` custom resource represents a specific, provisioned instance of a CTF challenge based on a
`ChallengeDescription`. When a `ChallengeInstance` is created, the operator provisions all necessary Kubernetes
resources in a dedicated namespace, ensuring isolation and automated lifecycle management for each challenge instance.
This resource allows organizers to spin up, manage, and clean up individual challenge environments for participants.

For details about available fields, see [`api/v1alpha1/challenge_instance.go`](api/v1alpha1/challenge_instance.go).
For a concrete example, see [`examples/challenge-instance-sample.yaml`](examples/challenge-instance-sample.yaml).

### APIKey CR

The `APIKey` custom resource manages API keys used for accessing various APIs within the CTF environment. Each `APIKey`
resource represents a single key, including metadata and configuration for its usage and automated lifecycle management.
The operator ensures that API keys are created, distributed, and deleted according to their defined lifetimes, helping
maintain security and automate access control for challenge instances or other components.

For details about available fields, see [`api/v1alpha1/api_key.go`](api/v1alpha1/api_key.go).
For a concrete example, see [`examples/api-key-sample.yaml`](examples/api-key-sample.yaml).

### Operator Command Line Parameters

The operator provides the following command line parameters:

```text
This operator manages CTF challenge instances.

Usage:
  ctf-challenge-operator [flags]

Flags:
      --enable-developer-mode              This option makes the log output friendlier to humans.
      --health-probe-bind-address string   The address the probe endpoint binds to. (default "0")
  -h, --help                               help for ctf-challenge-operator
      --kubernetes-client-burst int        The number of burst queries the Kubernetes client is allowed to send against the Kubernetes API. (default 10)
      --kubernetes-client-qps float32      The number of queries per second the Kubernetes client is allowed to send against the Kubernetes API. (default 5)
      --leader-election-enabled            Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.
      --leader-election-id string          The ID to use for leader election. (default "ctf-challenge-operator")
      --leader-election-namespace string   The namespace in which leader election should happen. (default "ctf-challenge-operator")
      --log-level int                      How verbose the logs are. Level 0 will show info, warning and error. Level 1 and up will show increasing details.
      --metrics-bind-address string        The address the metrics endpoint binds to. Use :8443 for HTTPS or :8080 for HTTP, or leave as 0 to disable the metrics service. (default "0")
```

## Development

This project intends to be run on cloud provider infrastructure. As cloud providers provide new Kubernetes version only
after some time, this project aims at the oldest supported Kubernetes version. See
[Supported versions](https://kubernetes.io/releases/version-skew-policy/#supported-versions) of the official Kubernetes
documentation.

This project uses tools like `controller-gen` which was sensitive to Go version updates in the past. To reduce the
likelyhood of Go versions breaking the toolchain, this project aims at the oldest supported Go version. See
[Release Policy](https://go.dev/doc/devel/release#policy) of the official Go documentation.

You need to have the following tools available

- Go
- Docker
- make
- Linux

**NOTE: Windows is currently not supported. MacOS might work but is untested.**

For setting up your local Kubernetes development cluster:

```shell
make init-local
```

This will create a Kubernetes cluster with kind and install the CRDs.

To run the code:

```shell
make run
```

To run the tests:

```shell
make test
```

If you changed the data types of the custom resources, you can install the updated version with:

```shell
make install
```

To clean up everything (including the kind cluster you created with `make init-local`):

```shell
make clean
```

### Third Party Tools

All third party tools required for development are provided through shims in the `bin` directory of this project. Those
shims are shell scripts which download the required tool on-demand into the `tmp` directory and forward any arguments
to the real executable. If you want to interact with those tools outside of makefile targets, add the `bin` directory to
your `PATH` environment variable like this:

```shell
export PATH=${PWD}/bin:${PATH}
```
