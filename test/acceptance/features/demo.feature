Feature: Service Catalog Life-cycle

    Background:
        # Multi-cluster environment setup
        Given Primaza Cluster "main" is running
        And Worker Cluster "worker" for ClusterEnvironment "worker" is running
        And On Primaza Cluster "main", Worker "worker"'s ClusterContext secret "primaza-kw" for ClusterEnvironment "worker" is published
        And On Worker Cluster "worker", application namespace "applications" for ClusterEnvironment "worker" exists
        And Fail

    @wip
    Scenario: On Cluster Environment creation, Primaza Application Agent is successfully deployed to applications namespace
        # Cluster Environment setup
        When On Primaza Cluster "main", Resource is created
        """
        apiVersion: primaza.io/v1alpha1
        kind: ClusterEnvironment
        metadata:
            name: worker
            namespace: primaza-system
        spec:
            environmentName: dev
            clusterContextSecret: primaza-kw
            applicationNamespaces:
            - "applications"
        """
        And On Primaza Cluster "main", ClusterEnvironment "worker" state will eventually move to "Online"
        And On Primaza Cluster "main", ClusterEnvironment "worker" status condition with Type "Online" has Status "True"
        And On Primaza Cluster "main", ClusterEnvironment "worker" status condition with Type "ApplicationNamespacePermissionsRequired" has Status "False"
        And On Primaza Cluster "main", ClusterEnvironment "worker" status condition with Type "ServiceNamespacePermissionsRequired" has Status "False"
        # Agents propagation
        Then On Worker Cluster "worker", Primaza Application Agent exists into namespace "applications"
        # Service catalog propagation
        # And On Primaza Cluster "main", ServiceCatalog "dev" exists
        And On Worker Cluster "worker", ServiceCatalog "dev" exists in "applications"
        # Feeding the Service Catalog
        When On Primaza Cluster "main", Resource is created
        """
        apiVersion: primaza.io/v1alpha1
        kind: RegisteredService
        metadata:
          name: primaza-rsdb
          namespace: primaza-system
        spec:
          serviceClassIdentity:
            - name: type
              value: psqlserver
            - name: provider
              value: aws
          serviceEndpointDefinition:
            - name: host
              value: mydavphost.io
            - name: port
              value: "5432"
            - name: user
              value: davp
            - name: password
              value: quedicelagente
            - name: database
              value: davpdata
          sla: L3
        """
        And On Primaza Cluster "main", RegisteredService "primaza-rsdb" state will eventually move to "Available"
        Then On Primaza Cluster "main", ServiceCatalog "primaza-service-catalog" will contain RegisteredService "primaza-rsdb"
        # Consume the service catalog
        When On Primaza Cluster "main", Resource is created
        """
        apiVersion: primaza.io/v1alpha1
        kind: ServiceClaim
        metadata:
          name: sc-test
          namespace: primaza-system
        spec:
          serviceClassIdentity:
          - name: type
            value: psqlserver
          - name: provider
            value: aws
          serviceEndpointDefinitionKeys:
          - host
          - port
          - user
          - password
          - database
          environmentTag: stage
          application:
            kind: Deployment
            apiVersion: apps/v1
            selector:
              matchLabels:
                app: myapp
        """
        Then On Primaza Cluster "main", RegisteredService "primaza-rsdb" state will eventually move to "Claimed"
        And On Primaza Cluster "main", ServiceCatalog "primaza-service-catalog" will not contain RegisteredService "primaza-rsdb"
        # Release the Registered Service
        When On Primaza Cluster "main", Resource is deleted
        """
        apiVersion: primaza.io/v1alpha1
        kind: ServiceClaim
        metadata:
          name: sc-test
          namespace: primaza-system
        """
        Then On Primaza Cluster "main", RegisteredService "primaza-rsdb" state will eventually move to "Available"
        And  On Primaza Cluster "main", ServiceCatalog "primaza-service-catalog" will contain RegisteredService "primaza-rsdb"
        # Delete the Registered Service
        When On Primaza Cluster "main", RegisteredService "primaza-rsdb" is deleted
        Then On Primaza Cluster "main", ServiceCatalog "primaza-service-catalog" will not contain RegisteredService "primaza-rsdb"
        # Delete the Cluster Environment
        When On Primaza Cluster "main", Resource is deleted
        """
        apiVersion: primaza.io/v1alpha1
        kind: ClusterEnvironment
        metadata:
            name: worker
            namespace: primaza-system
        """
        Then On Primaza Cluster "main", ServiceCatalog "dev" does not exist
        Then On Worker Cluster "worker", Service Catalog "dev" does not exist in "applications"
