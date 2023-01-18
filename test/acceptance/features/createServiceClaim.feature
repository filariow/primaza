Feature: Create a service claim with label selector

    Scenario: Create a service claim
        Given Primaza Cluster "primaza-main" is running
        When On Primaza Cluster "primaza-main", Resource is created
        """
        apiVersion: primaza.io/v1alpha1
        kind: RegisteredService
        metadata:
          name: primaza-rsdb
          namespace: primaza-system
        spec:
          constraints:
            environments:
            - stage
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
        And On Primaza Cluster "primaza-main", Primaza ClusterContext secret "primaza-km" is published
        And On Primaza Cluster "primaza-main", Resource is created
        """
        apiVersion: primaza.io/v1alpha1
        kind: ClusterEnvironment
        metadata:
            name: primaza-main
            namespace: primaza-system
        spec:
            environmentName: stage
            clusterContextSecret: primaza-km
        """
        And On Primaza Cluster "primaza-main", RegisteredService "primaza-rsdb" state will eventually move to "Available"
        And On Primaza Cluster "primaza-main", Resource is created
        """
        apiVersion: primaza.io/v1alpha1
        kind: ServiceClaim
        metadata:
          name: sc-test
          namespace: primaza-system
        spec:
          name: some-name
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
                a: b
                c: d
        """
        Then On Primaza Cluster "primaza-main", the status of ServiceClaim is "Resolved"
