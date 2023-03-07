@wip
Feature: Setup Multicluster environment

    Scenario: On Cluster Environment creation, Primaza Application Agent is successfully deployed to applications namespace

        Given Primaza Cluster "primaza-main" is running
        And Worker Cluster "primaza-worker" for "primaza-main" is running
        And Clusters "primaza-main" and "primaza-worker" can communicate
        And On Primaza Cluster "primaza-main", Worker "primaza-worker"'s ClusterContext secret "primaza-kw" is published
        And On Worker Cluster "primaza-worker", application namespace "applications" exists
        And On Worker Cluster "primaza-worker", service namespace "services" exists
        And On Worker Cluster "primaza-worker", Resource is created
        """
        apiVersion: apps/v1
        kind: Deployment
        metadata:
          name: demo-postgresql
          namespace: services
          labels:
            app: demo-postgresql
        spec:
          replicas: 1
          selector:
            matchLabels:
              app: demo-postgresql
          template:
            metadata:
              labels:
                app: demo-postgresql
            spec:
              containers:
                - name: postgresql
                  image: postgres:13
                  imagePullPolicy: IfNotPresent
                  env:
                    - name: POSTGRES_DB
                      value: demo-psql-db
                    - name: POSTGRES_PASSWORD
                      value: Sup3rSecret!
                    - name: POSTGRES_USER
                      value: demo-psql-admin
                    - name: PGDATA
                      value: /tmp/data
                  ports:
                    - name: postgresql
                      containerPort: 5432
        """
        And On Worker Cluster "primaza-worker", Resource is created
        """
        apiVersion: apps/v1
        kind: Deployment
        metadata:
          labels:
            app: demo-app
          name: demo-app
          namespace: applications
        spec:
          replicas: 1
          selector:
            matchLabels:
              app: demo-app
          template:
            metadata:
              creationTimestamp: null
              labels:
                app: demo-app
            spec:
              containers:
              - image: bash:latest
                name: bash
                command:
                - sleep
                - infinite
        """
        And Fail
        And On Primaza Cluster "primaza-main", Resource is created
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
        And On Worker Cluster "primaza-worker", Primaza Application Agent is deployed into namespace "applications"
        And On Worker Cluster "primaza-worker", Primaza Service Agent is deployed into namespace "services"
        And On Primaza Cluster "primaza-main", Resource is created
        """
        apiVersion: primaza.io/v1alpha1
        kind: ServiceClass
        metadata:
            name: demo-psql-serviceclass
            namespace: primaza-system
        spec:
            constraints:
              environments:
              - demo
            resource:
                apiVersion: apps/v1
                kind: Deployment
                serviceEndpointDefinitionMapping:
                - name: port
                  jsonPath: .spec.template.spec.containers[0].ports[0].containerPort
                - name: user
                  jsonPath: .spec.template.spec.containers[0].env[?(@.name == "POSTGRES_USER")].value
                - name: password
                  jsonPath: .spec.template.spec.containers[0].env[?(@.name == "POSTGRES_PASSWORD")].value
                - name: database
                  jsonPath: .spec.template.spec.containers[0].env[?(@.name == "POSTGRES_DB")].value
            serviceClassIdentity:
            - name: type
              value: psqlserver
            - name: provider
              value: hosted
        """
        When On Primaza Cluster "primaza-main", Resource is created
        """
        apiVersion: primaza.io/v1alpha1
        kind: ServiceClaim
        metadata:
          name: demo-app-psql
          namespace: primaza-system
        spec:
          serviceClassIdentity:
          - name: type
            value: psqlserver
          - name: provider
            value: hosted
          serviceEndpointDefinitionKeys:
          - port
          - user
          - password
          - database
          environmentTag: demo
          application:
            kind: Deployment
            apiVersion: apps/v1
            selector:
              matchLabels:
                app: demo-app
        """
        Then  On Worker Cluster "primaza-worker", ServiceBinding "demo-app-psql" on namespace "applications" state will eventually move to "Ready"
        And   On Worker Cluster "primaza-worker", in pods with label "app=demo-app" running in namespace "applications" file "/bindings/demo-app-psql/password" has content "Sup3rSecret!"
