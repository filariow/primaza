import os
import tempfile
from steps.clusterprovisioner import ClusterProvisioner
from steps.util import get_api_client_from_kubeconfig
from steps.command import Command
from kubernetes import client


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

    def get_admin_kubeconfig(self):
        """
        Returns the cluster admin kubeconfig
        """
        return self.cluster_provisioner.kubeconfig()

    def __deploy_agentapp(self, kubeconfig_path: str, img: str, namespace: str):
        out, err = self.__build_install_agentapp_base_cmd(kubeconfig_path, img).setenv("NAMESPACE", namespace).run("make agentapp deploy")
        print(out)
        assert err == 0, f"error deploying Agent app's controller into cluster {self.cluster_name}"

    def __deploy_agentsvc(self, kubeconfig_path: str, img: str, namespace: str):
        out, err = self.__build_install_agentapp_base_cmd(kubeconfig_path, img).setenv("NAMESPACE", namespace).run("make agentsvc deploy")
        print(out)
        assert err == 0, f"error deploying Agent app's controller into cluster {self.cluster_name}"

    def __install_crd_and_build_app_image(self, kubeconfig_path: str, img: str):
        out, err = self.__build_install_agentapp_base_cmd(kubeconfig_path, img).run("make agentapp install docker-build")
        print(out)
        assert err == 0, "error installing manifests and building agent app  controller"

    def __install_crd_and_build_svc_image(self, kubeconfig_path: str, img: str):
        out, err = self.__build_install_agentsvc_base_cmd(kubeconfig_path, img).run("make agentsvc install docker-build")
        print(out)
        assert err == 0, "error installing manifests and building agent app  controller"

    def __build_install_agentapp_base_cmd(self, kubeconfig_path: str, img: str) -> Command:
        return Command() \
            .setenv("KUBECONFIG", kubeconfig_path) \
            .setenv("GOCACHE", os.getenv("GOCACHE", "/tmp/gocache")) \
            .setenv("GOPATH", os.getenv("GOPATH", "/tmp/go")) \
            .setenv("IMG", img)

    def __build_install_agentsvc_base_cmd(self, kubeconfig_path: str, img: str) -> Command:
        return Command() \
            .setenv("KUBECONFIG", kubeconfig_path) \
            .setenv("GOCACHE", os.getenv("GOCACHE", "/tmp/gocache")) \
            .setenv("GOPATH", os.getenv("GOPATH", "/tmp/go")) \
            .setenv("IMG", img)

    def deploy_agentapp(self, namespace: str):
        """
        Deploys Application Agent into a cluster's namespace
        """
        kubeconfig = self.get_admin_kubeconfig()
        img = "agentapp:latest"
        with tempfile.NamedTemporaryFile(prefix=f"kubeconfig-{self.cluster_name}-") as t:
            t.write(kubeconfig.encode("utf-8"))
            self.__install_crd_and_build_app_image(t.name, img)
            self.load_docker_image(img)
            self.__deploy_agentapp(t.name, img, namespace)

    def deploy_agentsvc(self, namespace: str):
        """
        Deploys the Service Agent into a cluster's namespace
        """
        kubeconfig = self.get_admin_kubeconfig()
        img = "agentsvc:latest"
        with tempfile.NamedTemporaryFile(prefix=f"kubeconfig-{self.cluster_name}-") as t:
            t.write(kubeconfig.encode("utf-8"))
            self.__install_crd_and_build_svc_image(t.name, img)
            self.__load_docker_image(img)
            self.__deploy_agentsvc(t.name, img, namespace)

    def is_app_agent_deployed(self, namespace: str) -> bool:
        api_client = self.get_api_client()
        appsv1 = client.AppsV1Api(api_client)

        appsv1.read_namespaced_deployment(name="controller-agentapp", namespace=namespace)
        return True

    def is_svc_agent_deployed(self, namespace: str) -> bool:
        api_client = self.get_api_client()
        appsv1 = client.AppsV1Api(api_client)

        appsv1.read_namespaced_deployment(name="controller-agentsvc", namespace=namespace)
        return True

