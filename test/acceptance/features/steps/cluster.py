from steps.clusterprovisioner import ClusterProvisioner
from steps.util import get_api_client_from_kubeconfig
from typing import Dict
import base64
import time
import yaml
from datetime import datetime, timezone, timedelta
from kubernetes import client
from kubernetes.client.rest import ApiException


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

    def get_csr_kubeconfig(self, certificate_key: str, csr: str) -> Dict:
        """
        Generates the kubeconfig for the CertificateSignignRequest `csr`.
        The key used when creating the CSR is also needed.
        """

        # retrieve primaza's certificate
        api_client = self.get_api_client()
        certs = client.CertificatesV1Api(api_client)
        v1csr = certs.read_certificate_signing_request(name=csr)
        certificate = v1csr.status.certificate

        # building kubeconfig
        kubeconfig = self.cluster_provisioner.kubeconfig(internal=True)
        kcd = yaml.safe_load(kubeconfig)
        kcd["contexts"][0]["context"]["user"] = csr
        kcd["users"][0]["name"] = csr
        kcd["users"][0]["user"]["client-key-data"] = base64.b64encode(certificate_key.encode("utf-8")).decode("utf-8")
        kcd["users"][0]["user"]["client-certificate-data"] = certificate  # yet in base64 encoding

        return kcd

    def create_csr_user(self, csr_name: str, csr_pem: bytes, timeout: int = 60):
        api_client = self.get_api_client()
        certs = client.CertificatesV1Api(api_client)

        # Check if CertificateSigningRequest has yet been created and approved
        try:
            s = certs.read_certificate_signing_request_status(name=csr_name)
            if s == "Approved":
                print(f"cluster '{self.cluster_name}' already has an approved CertificateSigningRequest '{csr_name}'")
                return
        except ApiException as e:
            if e.reason != "Not Found":
                raise e

        # Create CertificateSigningRequest
        v1csr = client.V1CertificateSigningRequest(
            metadata=client.V1ObjectMeta(name=csr_name),
            spec=client.V1CertificateSigningRequestSpec(
                signer_name="kubernetes.io/kube-apiserver-client",
                request=base64.b64encode(csr_pem).decode("utf-8"),
                expiration_seconds=86400,
                usages=["client auth"]))
        certs.create_certificate_signing_request(v1csr)

        # Approve CertificateSigningRequest
        v1csr = certs.read_certificate_signing_request(name=csr_name)
        approval_condition = client.V1CertificateSigningRequestCondition(
            last_update_time=datetime.now(timezone.utc).astimezone(),
            message='This certificate was approved by Acceptance tests',
            reason='Acceptance tests',
            type='Approved',
            status='True')
        v1csr.status.conditions = [approval_condition]
        certs.replace_certificate_signing_request_approval(name=csr_name, body=v1csr)

        # Wait for certificate emission
        tend = datetime.now() + timedelta(seconds=timeout)
        while datetime.now() < tend:
            v1csr = certs.read_certificate_signing_request(name=csr_name)
            status = v1csr.status
            if hasattr(status, 'certificate') and status.certificate is not None:
                print(f"CertificateSignignRequest '{csr_name}' certificate is ready")
                return
            print(f"CertificateSignignRequest '{csr_name}' certificate is not ready")
            time.sleep(5)
        assert False, f"Timed-out waiting CertificateSignignRequest '{csr_name}' certificate to become ready"
