from steps.clusterprovisioner import ClusterProvisioner
from steps.util import get_api_client_from_kubeconfig


class Cluster(object):
    """
    Base class for managing a kubernetes cluster provisioned through a ClusterProvisioner
    """
    cluster_name: str = None
    cluster_provisioner: ClusterProvisioner = None

    def __init__(self, cluster_provisioner: ClusterProvisioner, cluster_name: str):
        self.cluster_provisioner = cluster_provisioner
        self.cluster_name = cluster_name

    def start(self):
        """
        Starts the cluster via the cluster provisioner
        """
        output, ec = self.cluster_provisioner.start()
        assert ec == 0, f'Worker Cluster "{self.cluster_name}" failed to start: {output}'
        print(f'Worker "{self.cluster_name}" started')

    def get_api_client(self):
        """
        Build and returns a client for the kubernetes API server of the cluster
        using the administrator user
        """
        kubeconfig = self.cluster_provisioner.kubeconfig()
        api_client = get_api_client_from_kubeconfig(kubeconfig)
        return api_client

    def delete(self):
        """
        Deletes the cluster via the cluster provisioner
        """
        self.cluster_provisioner.delete()

    def get_admin_kubeconfig(self, internal=False):
        """
        Returns the cluster admin kubeconfig
        """
        return self.cluster_provisioner.kubeconfig(internal)

    def install_vault(self):
        kubeconfig = self.cluster_provisioner.kubeconfig()
        with tempfile.NamedTemporaryFile(prefix=f"kubeconfig-{self.cluster_name}-") as t:
            t.write(kubeconfig.encode("utf-8"))
            t.flush()

            cmd = """
helm repo add hashicorp https://helm.releases.hashicorp.com && \
    helm repo update && \
    helm install vault hashicorp/vault --set "injector.enabled=false" --set "server.dev.enabled=true" && \
    until [ "$(kubectl get pods vault-0 --output=jsonpath='{.status.phase}')" = "Running" ]; do echo "waiting for pod vault-0 to have status 'Running'"; sleep 5; done && \
    kubectl exec vault-0 -- vault secrets enable -path=internal kv-v2 && \
    kubectl apply -f test/acceptance/resources/vault_ingress.yaml
"""
            out, err = Command() \
                .setenv("HOME", os.getenv("HOME")) \
                .setenv("USER", os.getenv("USER")) \
                .setenv("KUBECONFIG", t.name) \
                .setenv("GOCACHE", os.getenv("GOCACHE", "/tmp/gocache")) \
                .setenv("GOPATH", os.getenv("GOPATH", "/tmp/go")) \
                .run(cmd)

    def install_nginx_ingress(self):
        """
        install NGINX ingress controller in the cluster
        """
