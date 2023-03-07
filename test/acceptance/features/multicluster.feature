@wip
Feature: Setup Multicluster environment

    Scenario: On Cluster Environment creation, Primaza Application Agent is successfully deployed to applications namespace

        Given Primaza Cluster "primaza-main" is running
        And Worker Cluster "primaza-worker" for "primaza-main" is running
        And Clusters "primaza-main" and "primaza-worker" can communicate
        And On Primaza Cluster "primaza-main", Worker "primaza-worker"'s ClusterContext secret "primaza-kw" is published
        And On Worker Cluster "primaza-worker", application namespace "applications" exists
        And On Worker Cluster "primaza-worker", service namespace "services" exists
        When On Primaza Cluster "primaza-main", Resource is created
        """
        apiVersion: primaza.io/v1alpha1
        kind: ClusterEnvironment
        metadata:
            name: primaza-worker
            namespace: primaza-system
        spec:
            environmentName: demo
            clusterContextSecret: primaza-kw
            applicationNamespaces:
            - "applications"
            serviceNamespaces:
            - "services"
        """
        Then On Worker Cluster "primaza-worker", Primaza Application Agent is deployed into namespace "applications"
        Then On Worker Cluster "primaza-worker", Primaza Service Agent is deployed into namespace "services"
        And fail
