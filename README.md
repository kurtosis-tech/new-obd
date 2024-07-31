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