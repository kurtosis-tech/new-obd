# New Online Boutique Demo

This demo is based on the [Google cloud microservices demo](https://github.com/GoogleCloudPlatform/microservices-demo) 

It's a smaller version with just a few GO micro-services using REST APIs

## Kardinal steps

1. Starts a local K8s cluster like Minikube

```bash
minikube start --driver=docker --cpus=10 --memory 8192 --disk-size 32g
minikube addons enable ingress
```

2. Both prod.app.localhost and dev.app.localhost defined in the host file
```bash
# Add these entries in the '/private/etc/hosts' file
127.0.0.1 prod.app.localhost
127.0.0.1 dev.app.localhost
```

3. Install Istio resources in the local cluster

```bash
istioctl install --set profile=default -y
```

4. Deploy Kardinal Manager in the local cluster

```bash
kardinal manager deploy local-minikube
```

5. Deploy the online boutique app with Kardinal

```bash
kardinal deploy --k8s-manifest ./release/obd-kardinal.yaml
```

6. Start the tunnel to access the services (you may have to provide you password for the underlying sudo access)
```bash
minikube tunnel
```

7. Open the [production page](http://prod.app.localhost/) in the browser to see the production online boutique


## Development Guide

This doc explains how to build and run the OnlineBoutique source code locally using the `tilt` command-line tool.

### Prerequisites

- [Docker for Desktop](https://www.docker.com/products/docker-desktop)
- kubectl (can be installed via `gcloud components install kubectl`)
- [tilt **0.22.2+**](https://docs.tilt.dev/install.html) (latest version recommended)
- [Minikube](https://minikube.sigs.k8s.io/docs/start/) (optional - see Local Cluster)

### Local Cluster

1. Launch a local Kubernetes cluster with one of the following tools:

    - To launch **Minikube** (tested with Ubuntu Linux). Please, ensure that the
      local Kubernetes cluster has at least:
        - 4 CPUs
        - 4.0 GiB memory
        - 32 GB disk space

      ```shell
      minikube start --cpus=4 --memory 4096 --disk-size 32g
      ```

    - To launch **Docker for Desktop** (tested with Mac/Windows). Go to Preferences:
        - choose “Enable Kubernetes”,
        - set CPUs to at least 3, and Memory to at least 6.0 GiB
        - on the "Disk" tab, set at least 32 GB disk space

2. Run `kubectl get nodes` to verify you're connected to the respective control plane.

3. Two options:
   1. Run `sudo tilt up`.
         To deploy the app using the `./release/obd-kardinal.yaml` file, with Kardinal annotations, in the cluster. Take into account that it will use the container images defined in the YAML, it will try to pull them from the cloud. The sudo privileges are necessary in order to port-forward the port "80"
   2. Run `sudo tilt up -- --build frontend --build productcatalogservice`.
         To deploy the app using the `./release/obd-kardinal.yaml` file and also create a new 'dev' flow with dev images version for the services specified with the `build` flag (valid values: 'frontend', 'cartservice', 'productcatalogservice', convine these as you want).
         Edit the source code and check the changes in the dev URL, `Tilt` will trigger the hot-reload for it

## Cleanup

If you've deployed the application with `tilt up` command, you can run
`tilt down --delete-namespaces` to clean up the deployed resources, the `--delete-namespaces` flas is important because otherwise it won't delete the namespace.
