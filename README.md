# Supernetes
A container orchestration system based on Kubernetes, made all the way from scratch. Supernetes tries to mimic/closely follow the actual Kubernetes architecture and components.

# Components of Supernetes-Container Orchestration System are:
1. Superlet: Core container creation and running functionality is provided by this component. Closely resembles Kubelet component in Kubernetes
2. SuperController: Tracks our cluster state (state of pods) and compares it to the desired/supplied state of the deployments, it controls the actual state of our cluster by sending requests to Superlet to create or delete pods/containers. Resembles Controller Manager component of Kubernetes.
3. SuperCache: Cache for storing key value pairs containing state of our cluster.
4. SuperAPIServer: API Server that provides a centralized API for orchestration. Similar to KubeAPIServer
5. SuperCTL: CLI for efficient cluster/node management. Resembles KubeCTL
   
