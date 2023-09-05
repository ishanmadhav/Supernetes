# Supernetes
A container orchestration system based on Kubernetes, made all the way from scratch. Supernetes tries to mimic/closely follow the actual Kubernetes architecture and components.

## Components of Supernetes-Container Orchestration System are:
1. Superlet: Core container creation and running functionality is provided by this component. Closely resembles Kubelet component in Kubernetes
2. SuperController: Tracks our cluster state (state of pods) and compares it to the desired/supplied state of the deployments, it controls the actual state of our cluster by sending requests to Superlet to create or delete pods/containers. Resembles Controller Manager component of Kubernetes.
3. SuperCache: Cache for storing key value pairs containing state of our cluster.
4. SuperAPIServer: API Server that provides a centralized API for orchestration. Similar to KubeAPIServer
5. SuperCTL: CLI for efficient cluster/node management. Resembles KubeCTL

Below, we'll discuss how each of our components function. 

## Superlet:
The core of our container orchestration system. This component is responsible for creation, updation and deletion of our pods based on the requests it receives. In its current implementation, Superlet uses Docker API to deal with containers. Our initial implementation was based on an even lower level library called containerd (which itself is based on runc). But, since containerd was never made to be very developer friendly, and due to a few more issues such as data cleanup, and other limitations, we switched to Docker API. Considering, Kubernetes itself was Docker based in its early iterations, we felt it was an okay decision to move forward with it. 

## SuperController:
This component closely resembles the functioning of Controller Manager in Kubernetes. It compares the desried state of our cluster with the current state of the cluster. And if both of these states do not match, SuperController sends requests to the Superlet to create or delete certain pods. All this is done via te use of 

## SuperCache:

## SuperAPIServer:

## SuperCTL:


