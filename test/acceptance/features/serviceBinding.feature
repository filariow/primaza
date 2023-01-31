Feature: Bind application to  the secret pushed

    Scenario: Service binding projection works

        Given   Primaza Cluster "primaza-main" is running
        And     On Primaza Cluster "primaza-main", application namespace "applications" exists
        Given   On Primaza Cluster "primaza-main" application "newapp" is running in namespace "applications"
        And     On Primaza Cluster "primaza-main", Resource is created
        """
        apiVersion: v1
        kind: Secret
        metadata:
            name: demo
            namespace: applications
        stringData:
            username: AzureDiamond
        """
        When    On Primaza Cluster "primaza-main", Resource is created
        """
        apiVersion: primaza.io/v1alpha1
        kind: ServiceBinding
        metadata:
            name: newapp-binding
            namespace: applications
        spec:
            serviceEndpointDefinitionSecret: demo
            application:
                name: newapp
                apiVersion: apps/v1
                kind: Deployment
        """
        Then  On Primaza Cluster "primaza-main", ServiceBinding "newapp-binding" on namespace "applications" state will eventually move to "Ready"


    Scenario: Agents bind application on Worker Cluster

        Given Worker Cluster "primaza-worker" is running
        And   On Worker Cluster "primaza-worker", application namespace "applications" exists
        And   On Worker Cluster "primaza-worker", Primaza Application Agent is deployed into namespace "applications"
        And   On Worker Cluster "primaza-worker" application "app" is running in namespace "applications"
        And   On Worker Cluster "primaza-worker", Resource is created
        """
        apiVersion: v1
        kind: Secret
        metadata:
            name: demo
            namespace: applications
        stringData:
            username: AzureDiamond
            password: pass
        """
        When On Worker Cluster "primaza-worker", Resource is created
        """
        apiVersion: primaza.io/v1alpha1
        kind: ServiceBinding
        metadata:
            name: app-binding
            namespace: applications
        spec:
            serviceEndpointDefinitionSecret: demo
            application:
                name: app
                apiVersion: apps/v1
                kind: Deployment
        """
        Then On Worker Cluster "primaza-worker", ServiceBinding "app-binding" on namespace "applications" state will eventually move to "Ready"

    Scenario: Service binding projection works for label selector

        Given   Primaza Cluster "primaza-main" is running
        And     On Primaza Cluster "primaza-main", application namespace "applications" exists
        Given   On Primaza Cluster "primaza-main" application "applicationone" is running in namespace "applications"
        And     On Primaza Cluster "primaza-main" application "applicationtwo" is running in namespace "applications"
        And     On Primaza Cluster "primaza-main", Resource is created
        """
        apiVersion: v1
        kind: Secret
        metadata:
            name: demo
            namespace: applications
        stringData:
            username: AzureDiamond
        """
        When    On Primaza Cluster "primaza-main", Resource is created
        """
        apiVersion: primaza.io/v1alpha1
        kind: ServiceBinding
        metadata:
            name: application-binding
            namespace: applications
        spec:
            serviceEndpointDefinitionSecret: demo
            application:
                apiVersion: apps/v1
                kind: Deployment
                selector:
                    matchLabels:
                        app: myapp
        """
        Then  On Primaza Cluster "primaza-main", ServiceBinding "application-binding" on namespace "applications" state will eventually move to "Ready"