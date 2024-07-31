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

This doc explains how to build and run the OnlineBoutique source code locally using the `skaffold` command-line tool.

### Prerequisites

- [Docker for Desktop](https://www.docker.com/products/docker-desktop).
- kubectl (can be installed via `gcloud components install kubectl`)
- [skaffold **1.27+**](https://skaffold.dev/docs/install/) (latest version recommended), a tool that builds and deploys Docker images in bulk.
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

3. Run `skaffold run` (first time will be slow, it can take ~20 minutes).
   This will build and deploy the application. If you need to rebuild the images
   automatically as you refactor the code, run `skaffold dev` command.

4. Run `kubectl get pods` to verify the Pods are ready and running.

5. Access the web frontend through your browser
    - **Minikube** requires you to run a command to access the frontend service:

    ```shell
    minikube service frontend-external
    ```

    - **Docker For Desktop** should automatically provide the frontend at http://localhost:80

## Cleanup

If you've deployed the application with `skaffold run` command, you can run
`skaffold delete` to clean up the deployed resources.
