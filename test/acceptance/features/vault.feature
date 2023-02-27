Feature: vault
    @wip
    Scenario: vault is installed
        Given Primaza Cluster "primaza-main" is running
        And   On Primaza Cluster "primaza-main", vault is installed
        When  On Primaza Cluster "primaza-main", Resource is created
        """
        apiVersion: primaza.io/v1alpha1
        kind: ClusterEnvironment
        metadata:
            name: primaza-worker
            namespace: primaza-system
        spec:
            environmentName: dev
            clusterContextSecret: primaza-kw
            secretManager:
                vault:
                    address: http://vault:8200/
                    auth:
                        useServiceAccount: true
        """
