Feature: Publish Application Agent to worker cluster

    @wip
    Scenario: On Cluster Environment creation, Primaza Application Agent is successfully deployed to application namespace

        Given Primaza Cluster "primaza-main" is running
        And   Worker Cluster "primaza-worker" for "primaza-main" is running
        And   Clusters "primaza-main" and "primaza-worker" can communicate
        And   On Primaza Cluster "primaza-main", Worker "primaza-worker"'s ClusterContext secret "primaza-kw" is published
        And   On Worker Cluster "primaza-worker", applications namespace "applications" exists
        When On Primaza Cluster "primaza-main", Resource is created
        """
        apiVersion: primaza.io/v1alpha1
        kind: ClusterEnvironment
        metadata:
            name: primaza-worker
            namespace: primaza-system
        spec:
            environmentName: dev
            clusterContextSecret: primaza-kw
            applicationNamespaces:
            - applications
        """
        Then On Worker Cluster "primaza-worker", Primaza Application Agent is deployed into "applications" namespace

    Scenario: On Cluster Environment update, Primaza Application Agent is successfully removed from application namespace

        Given Primaza Cluster "primaza-main" is running
        And   Worker Cluster "primaza-worker" for "primaza-main" is running
        And   Clusters "primaza-main" and "primaza-worker" can communicate
        And   On Primaza Cluster "primaza-main", Worker "primaza-worker"'s ClusterContext secret "primaza-kw" is published
        And   On Worker Cluster "primaza-worker", applications namespace "applications" exists
        And   On Primaza Cluster "primaza-main", Resource is created
        """
        apiVersion: primaza.io/v1alpha1
        kind: ClusterEnvironment
        metadata:
            name: primaza-worker
            namespace: primaza-system
        spec:
            environmentName: dev
            clusterContextSecret: primaza-kw
            applicationNamespaces:
            - applications
        """
        And On Worker Cluster "primaza-worker", Primaza Application Agent is deployed into "applications" namespace
        When On Primaza Cluster "primaza-main", Resource is updated
        """
        apiVersion: primaza.io/v1alpha1
        kind: ClusterEnvironment
        metadata:
            name: primaza-worker
            namespace: primaza-system
        spec:
            environmentName: dev
            clusterContextSecret: primaza-kw
            applicationNamespaces: []
        """
        Then On Worker Cluster "primaza-worker", Primaza Application Agent is removed from namespace "applications"

    Scenario: On Cluster Environment update, Primaza Application Agent is successfully published into application namespace

        Given Primaza Cluster "primaza-main" is running
        And   Worker Cluster "primaza-worker" for "primaza-main" is running
        And   Clusters "primaza-main" and "primaza-worker" can communicate
        And   On Primaza Cluster "primaza-main", Worker "primaza-worker"'s ClusterContext secret "primaza-kw" is published
        And   On Worker Cluster "primaza-worker", applications namespace "applications" exists
        And   On Primaza Cluster "primaza-main", Resource is created
        """
        apiVersion: primaza.io/v1alpha1
        kind: ClusterEnvironment
        metadata:
            name: primaza-worker
            namespace: primaza-system
        spec:
            environmentName: dev
            clusterContextSecret: primaza-kw
            applicationNamespaces: []
        """
        When On Primaza Cluster "primaza-main", Resource is updated
        """
        apiVersion: primaza.io/v1alpha1
        kind: ClusterEnvironment
        metadata:
            name: primaza-worker
            namespace: primaza-system
        spec:
            environmentName: dev
            clusterContextSecret: primaza-kw
            applicationNamespaces:
            - applications
        """
        Then On Worker Cluster "primaza-worker", Primaza Application Agent is deployed into "applications" namespace

    Scenario: On Cluster Environment deletion, Primaza Application Agent is successfully published into application namespace

        Given Primaza Cluster "primaza-main" is running
        And   Worker Cluster "primaza-worker" for "primaza-main" is running
        And   Clusters "primaza-main" and "primaza-worker" can communicate
        And   On Primaza Cluster "primaza-main", Worker "primaza-worker"'s ClusterContext secret "primaza-kw" is published
        And   On Worker Cluster "primaza-worker", namespace "applications" exists
        And   On Primaza Cluster "primaza-main", Resource is created
        """
        apiVersion: primaza.io/v1alpha1
        kind: ClusterEnvironment
        metadata:
            name: primaza-worker
            namespace: primaza-system
        spec:
            environmentName: dev
            clusterContextSecret: primaza-kw
            applicationNamespaces:
            - applications
        """
        When On Primaza Cluster "primaza-main", Resource is deleted
        """
        apiVersion: primaza.io/v1alpha1
        kind: ClusterEnvironment
        metadata:
            name: primaza-worker
            namespace: primaza-system
        spec:
            environmentName: dev
            clusterContextSecret: primaza-kw
            applicationNamespaces:
            - applications
        """
        Then On Worker Cluster "primaza-worker", Primaza Application Agent is removed from namespace "applications"
