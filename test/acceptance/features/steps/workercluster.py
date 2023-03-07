import os
import polling2
import tempfile
from behave import given, step
from typing import Dict, List
from kubernetes import client
from kubernetes.client.rest import ApiException
from steps.command import Command
from steps.cluster import Cluster
from steps.primazacluster import PrimazaCluster


class WorkerCluster(Cluster):
    """
    Base class for instances of Worker clusters.
    Implements functionalities for configuration of kubernetes clusters that will act as Primaza workers,
    like CertificateSigningRequest approval and Primaza CRD installation.
    """

    def __init__(self, cluster_provisioner: str, cluster_name: str):
        super().__init__(cluster_provisioner, cluster_name)

    def start(self):
        super().start()

    def install_primaza_crd(self):
        """
        Installs Primaza's CRD into the cluster
        """
        self.__install_crd(component="primaza")

    def install_agentapp_crd(self):
        """
        Installs Application Agent's CRD into the cluster
        """
        self.__install_crd(component="agentapp")

    def install_agentsvc_crd(self):
        """
        Installs Service Agent's CRD into the cluster
        """
        self.__install_crd(component="agentsvc")

    def __install_crd(self, component: str):
        """
        Installs component's CRD into the cluster
        """
        kubeconfig = self.cluster_provisioner.kubeconfig()
        with tempfile.NamedTemporaryFile(prefix=f"kubeconfig-{self.cluster_name}-") as t:
            t.write(kubeconfig.encode("utf-8"))
            t.flush()

            out, err = Command() \
                .setenv("HOME", os.getenv("HOME")) \
                .setenv("USER", os.getenv("USER")) \
                .setenv("KUBECONFIG", t.name) \
                .setenv("GOCACHE", os.getenv("GOCACHE", "/tmp/gocache")) \
                .setenv("GOPATH", os.getenv("GOPATH", "/tmp/go")) \
                .run(f"make {component} install")

            print(out)
            assert err == 0, f"error installing {component}'s manifests"

    def create_primaza_user(self, csr_pem: bytes, timeout: int = 60):
        """
        Creates a CertificateSigningRequest for user primaza, approves it,
        and creates the needed roles and role bindings.
        """
        super().create_csr_user("primaza", csr_pem, timeout)

    def create_application_namespace(self, namespace: str):
        api_client = self.get_api_client()
        corev1 = client.CoreV1Api(api_client)

        self.install_agentapp_crd()

        # create application namespace
        ns = client.V1Namespace(metadata=client.V1ObjectMeta(name=namespace))
        corev1.create_namespace(ns)

        self.__create_application_agent_identity(namespace)
        self.__allow_primaza_access_to_namespace(namespace)

    def create_service_namespace(self, namespace: str):
        api_client = self.get_api_client()
        corev1 = client.CoreV1Api(api_client)

        self.install_agentsvc_crd()

        # create service namespace
        ns = client.V1Namespace(metadata=client.V1ObjectMeta(name=namespace))
        corev1.create_namespace(ns)

        self.__create_service_agent_identity(namespace)
        self.__allow_primaza_access_to_namespace(namespace)

    def __allow_primaza_access_to_namespace(self, namespace: str):
        api_client = self.get_api_client()
        rbacv1 = client.RbacAuthorizationV1Api(api_client)

        r = client.V1Role(
            metadata=client.V1ObjectMeta(name="primaza-role", namespace=namespace),
            rules=[
                client.V1PolicyRule(
                    api_groups=[""],
                    resources=["services"],
                    verbs=["get", "list", "watch", "create", "update", "patch", "delete"]),
                client.V1PolicyRule(
                    api_groups=["apps"],
                    resources=["deployments"],
                    verbs=["create"]),
                client.V1PolicyRule(
                    api_groups=["apps"],
                    resources=["deployments"],
                    verbs=["delete"],
                    resource_names=["primaza-controller-agentapp", "primaza-controller-agentsvc"]),
                client.V1PolicyRule(
                    api_groups=["primaza.io"],
                    resources=["servicebindings", "serviceclasses"],
                    verbs=["get", "list", "watch", "create", "update", "patch", "delete"])
            ])
        rbacv1.create_namespaced_role(namespace, r)

        # bind role to service account
        rb = client.V1RoleBinding(
            metadata=client.V1ObjectMeta(name="primaza-rolebinding", namespace=namespace),
            role_ref=client.V1RoleRef(
                api_group="rbac.authorization.k8s.io",
                kind="Role",
                name="primaza-role"),
            subjects=[
                client.V1Subject(
                    api_group="",
                    kind="User",
                    name="primaza")
            ])
        rbacv1.create_namespaced_role_binding(namespace=namespace, body=rb)

    def __create_application_agent_identity(self, namespace: str):
        rules = [
            client.V1PolicyRule(
                api_groups=["apps"],
                resources=["deployments"],
                verbs=["get", "list", "watch", "update", "patch"]),
            client.V1PolicyRule(
                api_groups=["primaza.io"],
                resources=["serviceclaims"],
                verbs=["get", "list", "watch"]),
            client.V1PolicyRule(
                api_groups=["primaza.io"],
                resources=["servicebindings", "servicebinding/status"],
                verbs=["get", "list", "watch", "update", "patch"])
        ]
        self.__create_agent_identity(namespace, "agentapp", rules)

    def __create_service_agent_identity(self, namespace: str):
        rules = [
            client.V1PolicyRule(
                api_groups=["apps"],
                resources=["deployments"],
                verbs=["get", "list"]),
            client.V1PolicyRule(
                api_groups=["primaza.io"],
                resources=["serviceclasses"],
                verbs=["get", "list", "watch"])
        ]
        self.__create_agent_identity(namespace, "agentsvc", rules)

    def __create_agent_identity(self, namespace: str, agent_name: str, rules: List[client.V1PolicyRule]):
        api_client = self.get_api_client()
        corev1 = client.CoreV1Api(api_client)
        rbacv1 = client.RbacAuthorizationV1Api(api_client)

        # create service account
        sa = client.V1ServiceAccount(metadata=client.V1ObjectMeta(name=f"primaza-controller-{agent_name}"))
        corev1.create_namespaced_service_account(namespace=namespace, body=sa)

        cr = client.V1ClusterRole(
            metadata=client.V1ObjectMeta(name=f"primaza-{agent_name}-{namespace}-role"),
            rules=[
                client.V1PolicyRule(
                    api_groups=["authentication.k8s.io"],
                    resources=["tokenreviews", "subjectaccessreviews"],
                    verbs=["create"])
            ])
        rbacv1.create_cluster_role(cr)

        # bind cluster role to service account
        crb = client.V1RoleBinding(
            metadata=client.V1ObjectMeta(name=f"primaza-{agent_name}-{namespace}-rolebinding", namespace=namespace),
            role_ref=client.V1RoleRef(
                api_group="rbac.authorization.k8s.io",
                kind="ClusterRole",
                name=f"primaza-{agent_name}-{namespace}-role"),
            subjects=[
                client.V1Subject(
                    kind="ServiceAccount",
                    name=f"primaza-controller-{agent_name}",
                    namespace=namespace)
            ])
        rbacv1.create_cluster_role_binding(body=crb)

        # create role
        r = client.V1Role(
            metadata=client.V1ObjectMeta(name=f"primaza-{agent_name}-role", namespace=namespace),
            rules=[
                client.V1PolicyRule(
                    api_groups=["coordination.k8s.io"],
                    resources=["leases"],
                    verbs=["get", "list", "watch", "create", "update", "patch", "delete"]),
                client.V1PolicyRule(
                    api_groups=[""],
                    resources=["configmaps"],
                    verbs=["get", "list", "watch", "create", "update", "patch", "delete"]),
                client.V1PolicyRule(
                    api_groups=[""],
                    resources=["events"],
                    verbs=["create", "patch"]),
                client.V1PolicyRule(
                    api_groups=[""],
                    resources=["secrets"],
                    verbs=["get", "list"])
            ] + rules)
        rbacv1.create_namespaced_role(namespace, r)

        # bind role to service account
        rb = client.V1RoleBinding(
            metadata=client.V1ObjectMeta(name=f"{agent_name}-rolebinding", namespace=namespace),
            role_ref=client.V1RoleRef(
                api_group="rbac.authorization.k8s.io",
                kind="Role",
                name=f"primaza-{agent_name}-role"),
            subjects=[
                client.V1Subject(
                    kind="ServiceAccount",
                    name=f"primaza-controller-{agent_name}",
                    namespace=namespace)
            ])
        rbacv1.create_namespaced_role_binding(namespace=namespace, body=rb)

    def create_primaza_user(self, csr_pem: bytes, timeout: int = 60):
        super().create_csr_user("primaza", csr_pem, timeout)
        self.__configure_primaza_user_permissions()

    def __configure_primaza_user_permissions(self):
        api_client = self.get_api_client()
        rbac = client.RbacAuthorizationV1Api(api_client)
        try:
            rbac.read_cluster_role_binding(name="primaza-primaza")
            return
        except ApiException as e:
            if e.reason != "Not Found":
                raise e

        role = client.V1ClusterRole(
            metadata=client.V1ObjectMeta(name="primaza"),
            rules=[
                client.V1PolicyRule(
                    api_groups=[""],
                    resources=["pods"],
                    verbs=["list", "get", "create"]),
                client.V1PolicyRule(
                    api_groups=[""],
                    resources=["secrets"],
                    verbs=["create"]),
                client.V1PolicyRule(
                    api_groups=["primaza.io/v1alpha1"],
                    resources=["servicebindings"],
                    verbs=["create"]),
            ])
        rbac.create_cluster_role(role)

        role_binding = client.V1ClusterRoleBinding(
            metadata=client.V1ObjectMeta(name="primaza-primaza"),
            role_ref=client.V1RoleRef(api_group="rbac.authorization.k8s.io", kind="ClusterRole", name="primaza"),
            subjects=[
                client.V1Subject(api_group="rbac.authorization.k8s.io", kind="User", name="primaza"),
            ])
        rbac.create_cluster_role_binding(role_binding)

    def get_primaza_user_kubeconfig(self, certificate_key: str) -> Dict:
        """
        Generates the kubeconfig for the CertificateSignignRequest "primaza".
        The key used when creating the CSR is also needed.
        """
        return super().get_csr_kubeconfig(certificate_key, "primaza")

    def get_primaza_user_kubeconfig_yaml(self, certificate_key: str) -> str:
        """
        Generates the kubeconfig for the CertificateSignignRequest "primaza".
        The key used when creating the CSR is also needed.
        Returns the YAML string
        """
        return self.get_user_kubeconfig_yaml(certificate_key, "primaza")

    def is_app_agent_deployed(self, namespace: str) -> bool:
        api_client = self.get_api_client()
        appsv1 = client.AppsV1Api(api_client)

        appsv1.read_namespaced_deployment(name="primaza-controller-agentapp", namespace=namespace)
        return True

    def is_svc_agent_deployed(self, namespace: str) -> bool:
        api_client = self.get_api_client()
        appsv1 = client.AppsV1Api(api_client)

        appsv1.read_namespaced_deployment(name="primaza-controller-agentsvc", namespace=namespace)
        return True

    def deploy_agentapp(self, namespace: str):
        """
        Deploys Application Agent into a cluster's namespace
        """

    def deploy_agentsvc(self, namespace: str):
        """
        Deploys the Service Agent into a cluster's namespace
        """


# Behave steps
@given('Worker Cluster "{cluster_name}" for "{primaza_cluster_name}" is running')
@given('Worker Cluster "{cluster_name}" for "{primaza_cluster_name}" is running with kubernetes version "{version}"')
def ensure_worker_cluster_running(context, cluster_name: str, primaza_cluster_name: str, version: str = None):
    worker_cluster = context.cluster_provider.create_worker_cluster(cluster_name, version)
    worker_cluster.start()

    primaza_cluster = context.cluster_provider.get_primaza_cluster(primaza_cluster_name)
    p_csr_pem = primaza_cluster.create_certificate_signing_request_pem()
    worker_cluster.create_primaza_user(p_csr_pem)


@given(u'On Worker Cluster "{cluster_name}", application namespace "{namespace}" exists')
@given(u'On Worker Cluster "{cluster_name}", application namespace "{namespace}" exists linked to "{primaza_cluster_name}"')
def ensure_application_namespace_exists(context, cluster_name: str, namespace: str, primaza_cluster_name: str = "primaza-main"):
    worker = context.cluster_provider.get_worker_cluster(cluster_name)  # type: WorkerCluster
    worker.create_application_namespace(namespace)

    primaza = context.cluster_provider.get_primaza_cluster(primaza_cluster_name)  # type: PrimazaCluster
    # w_csr_pem = worker.create_certificate_signing_request_pem("agentapp")
    # primaza.create_csr_user(w_csr_pem)

    # cc_kubeconfig = primaza.get_user_kubeconfig_yaml(worker.certificate_private_key, "agentapp")
    cc_kubeconfig = primaza.get_admin_kubeconfig(internal=True)
    worker.create_clustercontext_secret("primaza-kubeconfig", cc_kubeconfig, namespace)


@step(u'On Worker Cluster "{cluster_name}", Primaza Application Agent is deployed into namespace "{namespace}"')
def application_agent_is_deployed(context, cluster_name: str, namespace: str):
    worker = context.cluster_provider.get_worker_cluster(cluster_name)  # type: WorkerCluster
    worker.deploy_agentapp(namespace)
    polling2.poll(
        target=lambda: worker.is_app_agent_deployed(namespace),
        step=1,
        timeout=30)


@step(u'On Worker Cluster "{cluster_name}", Primaza Application Agent is not deployed into namespace "{namespace}"')
def application_agent_is_not_deployed(context, cluster_name: str, namespace: str):
    worker = context.cluster_provider.get_worker_cluster(cluster_name)  # type: WorkerCluster

    def is_not_found():
        try:
            worker.is_app_agent_deployed(namespace)
        except ApiException as e:
            return e.reason == "Not Found"
        return False

    polling2.poll(
        target=lambda: is_not_found(),
        step=1,
        timeout=30)


@given(u'On Worker Cluster "{cluster_name}", service namespace "{namespace}" exists')
@given(u'On Worker Cluster "{cluster_name}", service namespace "{namespace}" exists linked to "{primaza_cluster_name}"')
def ensure_service_namespace_exists(context, cluster_name: str, namespace: str, primaza_cluster_name: str = "primaza-main"):
    worker = context.cluster_provider.get_worker_cluster(cluster_name)  # type: WorkerCluster
    worker.create_service_namespace(namespace)

    primaza = context.cluster_provider.get_primaza_cluster(primaza_cluster_name)  # type: PrimazaCluster
    # w_csr_pem = worker.create_certificate_signing_request_pem("agentsvc")
    # primaza.create_csr_user(w_csr_pem)

    # cc_kubeconfig = primaza.get_user_kubeconfig_yaml(worker.certificate_private_key, "agentsvc")
    cc_kubeconfig = primaza.get_admin_kubeconfig(internal=True)
    worker.create_clustercontext_secret("primaza-kubeconfig", cc_kubeconfig, namespace)


@step(u'On Worker Cluster "{cluster_name}", Primaza Service Agent is deployed into namespace "{namespace}"')
def service_agent_is_deployed(context, cluster_name: str, namespace: str):
    worker = context.cluster_provider.get_worker_cluster(cluster_name)  # type: WorkerCluster
    worker.deploy_agentsvc(namespace)
    polling2.poll(
        target=lambda: worker.is_svc_agent_deployed(namespace),
        step=1,
        timeout=30)


@step(u'On Worker Cluster "{cluster_name}", Primaza Service Agent is not deployed into namespace "{namespace}"')
def service_agent_is_not_deployed(context, cluster_name: str, namespace: str):
    worker = context.cluster_provider.get_worker_cluster(cluster_name)  # type: WorkerCluster

    def is_not_found():
        try:
            worker.is_svc_agent_deployed(namespace)
        except ApiException as e:
            return e.reason == "Not Found"
        return False

    polling2.poll(
        target=lambda: is_not_found(),
        step=1,
        timeout=30)


@given('Worker Cluster "{cluster_name}" is running')
@given('Worker Cluster "{cluster_name}" is running with kubernetes version "{version}"')
def ensure_worker_cluster_is_running(context, cluster_name: str, version: str = None):
    worker_cluster = context.cluster_provider.create_worker_cluster(cluster_name, version)
    worker_cluster.start()


@step(u'On Worker Cluster "{cluster_name}", the secret "{secret_name}" in namespace "{namespace}" has the key "{key}" with value "{value}"')
def ensure_secret_key_has_the_right_value(context, cluster_name: str, secret_name: str, namespace: str, key: str, value: str):
    primaza_cluster = context.cluster_provider.get_worker_cluster(cluster_name)
    polling2.poll(
        target=lambda: primaza_cluster.read_secret_resource_data(namespace, secret_name, key) == bytes(value, 'utf-8'),
        step=1,
        timeout=30)
