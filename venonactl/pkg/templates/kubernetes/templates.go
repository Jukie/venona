// Code generated by go generate; DO NOT EDIT.
// using data from templates/kubernetes
package kubernetes

func TemplatesMap() map[string]string {
	templatesMap := make(map[string]string)

	templatesMap["cluster-role-binding.dind-volume-provisioner.vp.yaml"] = `---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: volume-provisioner-{{ .AppName }}-{{ .Namespace }}
  labels:
    app: dind-volume-provisioner-{{ .AppName }}
subjects:
  - kind: ServiceAccount
    name: volume-provisioner-{{ .AppName }}
    namespace: {{ .Namespace }}
roleRef:
  kind: ClusterRole
  name: volume-provisioner-{{ .AppName }}-{{ .Namespace }}
  apiGroup: rbac.authorization.k8s.io
`

	templatesMap["cluster-role-binding.engine.yaml"] = `kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: {{ .AppName }}-{{ .Namespace }}-engine
subjects:
- kind: ServiceAccount
  name: engine
  namespace: {{ .Namespace }}
roleRef:
  kind: ClusterRole
  name: cluster-admin
  apiGroup: rbac.authorization.k8s.io`

	templatesMap["cluster-role-binding.venona.yaml"] = `kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: {{ .AppName }}-{{ .Namespace }}
subjects:
- kind: ServiceAccount
  name: {{ .AppName }}
  namespace: {{ .Namespace }}
roleRef:
  kind: ClusterRole
  name: system:discovery
  apiGroup: rbac.authorization.k8s.io`

	templatesMap["cluster-role.dind-volume-provisioner.vp.yaml"] = `kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: volume-provisioner-{{ .AppName }}-{{ .Namespace }}
  labels:
    app: dind-volume-provisioner
rules:
  - apiGroups: [""]
    resources: ["persistentvolumes"]
    verbs: ["get", "list", "watch", "create", "delete", "patch"]
  - apiGroups: [""]
    resources: ["persistentvolumeclaims"]
    verbs: ["get", "list", "watch", "update"]
  - apiGroups: ["storage.k8s.io"]
    resources: ["storageclasses"]
    verbs: ["get", "list", "watch"]
  - apiGroups: [""]
    resources: ["events"]
    verbs: ["list", "watch", "create", "update", "patch"]
  - apiGroups: [""]
    resources: ["secrets"]
    verbs: ["get", "list"]
  - apiGroups: [""]
    resources: ["nodes"]
    verbs: ["get", "list", "watch"]
  - apiGroups: [""]
    resources: ["pods"]
    verbs: ["get", "list", "watch", "create", "delete", "patch"]
  - apiGroups: [""]
    resources: ["endpoints"]
    verbs: ["get", "list", "watch", "create", "update", "delete"]
`

	templatesMap["cluster-role.engine.yaml"] = `kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: {{ .AppName }}-{{ .Namespace }}-engine
rules:
  - apiGroups: [""]
    resources: ["pods", "persistentvolumeclaims"]
    verbs: ["get", "create", "delete", "list"]
`

	templatesMap["codefresh-certs-server-secret.re.yaml"] = `apiVersion: v1
type: Opaque
kind: Secret
metadata:
  labels:
    app: venona
  name: codefresh-certs-server
  namespace: {{ .Namespace }}
data:
  server-cert.pem: {{ .ServerCert.Cert }}
  server-key.pem: {{ .ServerCert.Key }}
  ca.pem: {{ .ServerCert.Ca }}

`

	templatesMap["cron-job.dind-volume-cleanup.vp.yaml"] = `{{- if not (eq .Storage.Backend "local") }}
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: dind-volume-cleanup-{{ .AppName }}
  namespace: {{ .Namespace }}
  labels:
    app: dind-volume-cleanup
spec:
  schedule: "0,10,20,30,40,50 * * * *"
  concurrencyPolicy: Forbid
  {{- if eq .Storage.Backend "local" }}
  suspend: true
  {{- end }}
  jobTemplate:
    spec:
      template:
        spec:
          serviceAccountName: volume-provisioner-{{ .AppName }}
          restartPolicy: Never
          containers:
            - name: dind-volume-cleanup
              image: {{ if ne .DockerRegistry ""}} {{- .DockerRegistry }}/codefresh/dind-volume-cleanup {{- else }}codefresh/dind-volume-cleanup {{- end}}
              env:
              - name: PROVISIONED_BY
                value: codefresh.io/dind-volume-provisioner-{{ .AppName }}-{{ .Namespace }}
{{- end }}`

	templatesMap["daemonset.dind-lv-monitor.vp.yaml"] = `{{- if eq .Storage.Backend "local" -}}
{{- $localVolumeParentDir := ( .Storage.LocalVolumeParentDir | default "/var/lib/codefresh/dind-volumes" ) }}
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: dind-lv-monitor-{{ .AppName }}
  namespace: {{ .Namespace }}
  labels:
    app: dind-lv-monitor
spec:
  selector:
    matchLabels:
      app: dind-lv-monitor
  template:
    metadata:
      labels:
        app: dind-lv-monitor
      annotations:
        prometheus_port: "9100"
        prometheus_scrape: "true"
    spec:
      serviceAccountName: volume-provisioner-{{ .AppName }}
      # Debug:
      # hostNetwork: true
      # nodeSelector:
      #   kubernetes.io/role: "node"
      tolerations:
        - key: 'codefresh/dind'
          operator: 'Exists'
          effect: 'NoSchedule'
      {{ if ne .Tolerations "" }}
        {{ .Tolerations | indent 8 }}
      {{ end }}
      containers:
        - image: {{ if ne .DockerRegistry ""}} {{- .DockerRegistry }}/codefresh/dind-volume-utils:v5 {{- else }}codefresh/dind-volume-utils:v5{{- end}}
          name: lv-cleaner
          imagePullPolicy: Always
          command:
          - /bin/local-volumes-agent
          env:
            - name: NODE_NAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
            - name: VOLUME_PARENT_DIR
              value: {{ $localVolumeParentDir }}
#              Debug:
#            - name: DRY_RUN
#              value: "1"
#            - name: DEBUG
#              value: "1"
#            - name: SLEEP_INTERVAL
#              value: "3"
#            - name: LOG_DF_EVERY
#              value: "60"
#            - name: KB_USAGE_THRESHOLD
#              value: "20"

          volumeMounts:
          - mountPath: {{ $localVolumeParentDir }}
            readOnly: false
            name: dind-volume-dir
      volumes:
      - name: dind-volume-dir
        hostPath:
          path: {{ $localVolumeParentDir }}
{{- end -}}`

	templatesMap["deployment.dind-volume-provisioner.vp.yaml"] = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: dind-volume-provisioner-{{ .AppName }}
  namespace: {{ .Namespace }}
  labels:
    app: dind-volume-provisioner
spec:
  selector:
    matchLabels:
      app: dind-volume-provisioner
  replicas: 1
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: dind-volume-provisioner
    spec:
      {{ if .Storage.VolumeProvisioner.NodeSelector }}
      nodeSelector: 
{{ .Storage.VolumeProvisioner.NodeSelector | nodeSelectorParamToYaml | indent 8 | unescape}}
      {{ end }}
      serviceAccount: volume-provisioner-{{ .AppName }}
      tolerations:
      - effect: NoSchedule
        key: node-role.kubernetes.io/master
        operator: "Exists"
      {{ if ne .Tolerations "" }}
        {{ .Tolerations | indent 6 }}
      {{ end }}
      containers:
      - name: dind-volume-provisioner
        image: {{ if ne .DockerRegistry ""}} {{- .DockerRegistry }}/{{ .Storage.VolumeProvisioner.Image }} {{- else }} {{- .Storage.VolumeProvisioner.Image }} {{- end}}
        imagePullPolicy: Always
        resources:
          requests:
            cpu: "200m"
            memory: "200Mi"
          limits:
            cpu: "1000m"
            memory: "6000Mi"
        command:
          - /usr/local/bin/dind-volume-provisioner
          - -v=4
          - --resync-period=50s
        env:
        - name: PROVISIONER_NAME
          value: codefresh.io/dind-volume-provisioner-{{ .AppName }}-{{ .Namespace }}
        {{- if ne .DockerRegistry "" }}
        - name: DOCKER_REGISTRY
          value: {{ .DockerRegistry }}
        {{- end }}
        {{- if .Storage.AwsAccessKeyId }}
        - name: AWS_ACCESS_KEY_ID
          valueFrom:
            secretKeyRef:
              name: dind-volume-provisioner-{{ .AppName }}
              key: aws_access_key_id
        {{- end }}
        {{- if .Storage.AwsSecretAccessKey }}
        - name: AWS_SECRET_ACCESS_KEY
          valueFrom:
            secretKeyRef:
              name: dind-volume-provisioner-{{ .AppName }}
              key: aws_secret_access_key
        {{- end }}
      {{- if .Storage.GoogleServiceAccount }}
        - name: GOOGLE_APPLICATION_CREDENTIALS
          value: /etc/dind-volume-provisioner/credentials/google-service-account.json
        volumeMounts:
        - name: credentials
          readOnly: true
          mountPath: "/etc/dind-volume-provisioner/credentials"
      volumes:
      - name: credentials
        secret:
          secretName: dind-volume-provisioner-{{ .AppName }}
      {{- end }}
`

	templatesMap["deployment.monitor.yaml"] = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Monitor.AppName }}
  namespace: {{ .Namespace }}
  labels:
    app: {{ .AppName }}
    version: {{ .Version }}
spec:
  replicas: 1
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 50%
      maxSurge: 50%
  selector:
    matchLabels:
      app: {{ .Monitor.AppName }}
  template:
    metadata:
      labels:
        app: {{ .Monitor.AppName }}
        version: {{ .Version }}
    spec:
      {{- if .Monitor.RbacEnabled}}
      serviceAccountName: {{ .Monitor.AppName }}
      {{- end }}
      containers:
      - name: {{ .Monitor.AppName }}
        image: {{ if ne .DockerRegistry ""}} {{- .DockerRegistry }}/{{ .Monitor.Image.Name }}:{{ .Monitor.Image.Tag }} {{- else }} {{- .Monitor.Image.Name }}:{{ .Monitor.Image.Tag }} {{- end}}
        imagePullPolicy: Always
        env:
          - name: SERVICE_NAME
            value: {{ .Monitor.AppName }}
          {{- if .Monitor.UseNamespaceWithRole }}
          - name: ROLE_BINDING
            value: "true"
          {{- end }}
          - name: PORT
            value: "9020"
          - name: API_TOKEN
            value: {{ .Token }}
          - name: CLUSTER_ID
            value: {{ .ClusterId }}
          - name: API_URL
            value: {{ .CodefreshHost }}/api/k8s-monitor/events
          - name: ACCOUNT_ID
            value: user
          - name: HELM3
            value: "{{ .Monitor.Helm3 }}"
          - name: NAMESPACE
            value: "{{ .Namespace }}"
          - name: NODE_OPTIONS
            value: "--max_old_space_size=4096"
        ports:
        - containerPort: 9020
          protocol: TCP
        readinessProbe:
          httpGet:
            path: /api/ping
            port: 9020
          periodSeconds: 5
          timeoutSeconds: 5
          successThreshold: 1
          failureThreshold: 5
`

	templatesMap["deployment.venona.yaml"] = `apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: {{ .AppName }}
    version: {{ .Version }} 
  name: {{ .AppName }}
  namespace: {{ .Namespace }}
spec:
  selector:
    matchLabels:
      app: {{ .AppName }}
      version: {{ .Version }}
  replicas: 1
  revisionHistoryLimit: 5
  strategy:
    rollingUpdate:
      maxSurge: 50%
      maxUnavailable: 50%
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: {{ .AppName }}
        version: {{ .Version }}
    spec:
      volumes:
        - name: runnerconf
          secret:
            secretName: runnerconf
      {{ if ne .NodeSelector "" }}
      nodeSelector:
{{ .NodeSelector | nodeSelectorParamToYaml | indent 8 | unescape }}
      {{ end }}
      {{ if ne .Tolerations "" }}
      tolerations:
        {{ .Tolerations | indent 8 }}
      {{ end }}
      containers:
      - env:
        {{- if $.AdditionalEnvVars }}
        {{- range $key, $value := $.AdditionalEnvVars }}
        - name: {{ $key }}
          value: "{{ $value}}"
        {{- end}}
        {{- end}}
        - name: SELF_DEPLOYMENT_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: CODEFRESH_TOKEN
          valueFrom:
            secretKeyRef:
              name: {{ .AppName }}
              key: codefresh.token
        - name: CODEFRESH_HOST
          value: {{ .CodefreshHost }}
        - name: AGENT_MODE
          value: {{ .Mode }}
        - name: AGENT_NAME
          value: {{ .AppName }}
        - name: "AGENT_ID"
          value: {{ .AgentId }}
        - name: VENONA_CONFIG_DIR
          value: "/etc/secrets"
        {{- if ne .DockerRegistry "" }}
        - name: DOCKER_REGISTRY
          value: {{ .DockerRegistry }}
        {{- end }}
        image: {{ if ne .DockerRegistry ""}} {{- .DockerRegistry }}/{{ .Image.Name }} {{- else }} {{- .Image.Name }}{{- end}}:{{ .Image.Tag | default "latest"}}
        volumeMounts:
        - name: runnerconf
          mountPath: "/etc/secrets"
          readOnly: true
        imagePullPolicy: Always
        name: {{ .AppName }}
      restartPolicy: Always
`

	templatesMap["dind-daemon-conf.re.yaml"] = `---
apiVersion: v1
kind: ConfigMap
metadata:
  name: codefresh-dind-config
  namespace: {{ .Namespace }}
data:
  daemon.json: |
    {
      "hosts": [ "unix:///var/run/docker.sock",
                 "tcp://0.0.0.0:1300"],
      "storage-driver": "overlay2",
      "tlsverify": true,  
      "tls": true,
      "tlscacert": "/etc/ssl/cf-client/ca.pem",
      "tlscert": "/etc/ssl/cf/server-cert.pem",
      "tlskey": "/etc/ssl/cf/server-key.pem",
      "insecure-registries" : ["192.168.99.100:5000"],
      "metrics-addr" : "0.0.0.0:9323",
      "experimental" : true
    }
`

	templatesMap["dind-headless-service.re.yaml"] = `---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: dind
  name: dind
  namespace: {{ .Namespace }}
spec:
  ports:
  - name: "dind-port"
    port: 1300
    protocol: TCP

  # This is a headless service, Kubernetes won't assign a VIP for it.
  # *.dind.default.svc.cluster.local
  clusterIP: None
  selector:
    app: dind

`

	templatesMap["role-binding.re.yaml"] = `kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: {{ .AppName }}
  namespace: {{ .Namespace }}
subjects:
- kind: ServiceAccount
  name: {{ .AppName }}
  namespace: {{ .Namespace }}
roleRef:
  kind: Role
  name: {{ .AppName }}
  apiGroup: rbac.authorization.k8s.io`

	templatesMap["role.monitor.yaml"] = `{{- if .Monitor.RbacEnabled }}
{{- if .Monitor.UseNamespaceWithRole }}
kind: Role
{{- else }}
kind: ClusterRole
{{- end }}
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: {{ .Monitor.AppName }}-cluster-reader
  namespace: {{ .Namespace }}
  labels:
    app: {{ .Monitor.AppName }}
    version: {{ .Version }}
rules:
- apiGroups:
  - ""
  resources: ["*"]
  verbs:
  - get
  - list
  - watch
  - create
  - delete
- apiGroups:
    - ""
  resources: ["pods"]
  verbs:
    - get
    - list
    - watch
    - create
    - deletecollection
- apiGroups:
  - extensions
  resources: ["*"]
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - apps
  resources: ["*"]
  verbs:
  - get
  - list
  - watch
{{- end }}
`

	templatesMap["role.re.yaml"] = `kind: Role
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: {{ .AppName }}
  namespace: {{ .Namespace }}
rules:
- apiGroups: [""]
  resources: ["pods", "persistentvolumeclaims"]
  verbs: ["get", "create", "delete"]
`

	templatesMap["rolebinding.monitor.yaml"] = `{{- if .Monitor.RbacEnabled }}
{{- if .Monitor.UseNamespaceWithRole }}
kind: RoleBinding
{{- else }}
kind: ClusterRoleBinding
{{- end }}
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: {{ .Monitor.AppName }}-cluster-reader
  namespace: {{ .Namespace }}
  labels:
    app: {{ .Monitor.AppName }}
    version: {{ .Version }}
subjects:
- kind: ServiceAccount
  namespace: {{ .Namespace }}
  name: {{ .Monitor.AppName }}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  {{- if .Monitor.UseNamespaceWithRole }}
  kind: Role
  {{- else }}
  kind: ClusterRole
  {{- end }}
  name: {{ .Monitor.AppName }}-cluster-reader
{{- end }}
`

	templatesMap["rollback-role-binding.monitor.yaml"] = `{{- if .Monitor.RbacEnabled }}
{{- if .Monitor.UseNamespaceWithRole }}
kind: RoleBinding
{{- else }}
kind: ClusterRoleBinding
{{- end }}
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: {{ .Monitor.AppName }}-rollback
  namespace: {{ .Namespace }}
  labels:
    app: {{ .Monitor.AppName }}
    version: {{ .Version }}
subjects:
  - kind: ServiceAccount
    namespace: {{ .Namespace }}
    name: {{ .Monitor.AppName }}-rollback
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
  {{- end }}
`

	templatesMap["rollback-serviceaccount.monitor.yaml"] = `{{- if and .Monitor.RbacEnabled (not .Monitor.UseNamespaceWithRole) }}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ .Monitor.AppName }}-rollback
  namespace: {{ .Namespace }}
  labels:
    app: {{ .Monitor.AppName }}
    version: {{ .Version }}
{{- end }}
`

	templatesMap["secret.dind-volume-provisioner.vp.yaml"] = `apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: dind-volume-provisioner-{{ .AppName }}
  namespace: {{ .Namespace }}
  labels:
    app: dind-volume-provisioner
data:
{{- if .Storage.GoogleServiceAccount }}
  google-service-account.json: {{ .Storage.GoogleServiceAccount | b64enc }}
{{- end }}
{{- if .Storage.AwsAccessKeyId }}
  aws_access_key_id: {{ .Storage.AwsAccessKeyId | b64enc }}
{{- end }}
{{- if .Storage.AwsSecretAccessKey }}
  aws_secret_access_key: {{ .Storage.AwsSecretAccessKey | b64enc }}
{{- end }}`

	templatesMap["secret.runtime-attach.yaml"] = `apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: {{ .AppName }}conf
  namespace: {{ .Namespace }}
data:
{{ range $key, $value := .runnerConf }}
  {{ $key }}: {{ $value }}
{{ end }}`

	templatesMap["secret.venona.yaml"] = `apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: {{ .AppName }}
  namespace: {{ .Namespace }}
data:
  codefresh.token: {{ .AgentToken | b64enc }}`

	templatesMap["service-account.dind-volume-provisioner.vp.yaml"] = `---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: volume-provisioner-{{ .AppName }}
  namespace: {{ .Namespace }}
  labels:
    app: dind-volume-provisioner
`

	templatesMap["service-account.engine.yaml"] = `apiVersion: v1
kind: ServiceAccount
metadata:
  name: engine
  namespace: {{ .Namespace }}`

	templatesMap["service-account.monitor.yaml"] = `{{- if .Monitor.RbacEnabled }}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ .Monitor.AppName }}
  namespace: {{ .Namespace }}
  labels:
    app: {{ .Monitor.AppName }}
    version: {{ .Version }}
{{- end }}
`

	templatesMap["service-account.re.yaml"] = `apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ .AppName }}
  namespace: {{ .Namespace }}`

	templatesMap["service.monitor.yaml"] = `apiVersion: v1
kind: Service
metadata:
  name: {{ .Monitor.AppName }}
  namespace: {{ .Namespace }}
  labels:
    app: {{ .Monitor.AppName }}
    version: {{ .Version }}
spec:
  type: ClusterIP
  ports:
  - name: "http"
    port: 80
    protocol: TCP
    targetPort: 9020
  selector:
    app: {{ .Monitor.AppName }}
`

	templatesMap["storageclass.dind-volume-provisioner.vp.yaml"] = `---
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: {{ .Storage.StorageClassName }}
  labels:
    app: dind-volume-provisioner
provisioner: codefresh.io/dind-volume-provisioner-{{ .AppName }}-{{ .Namespace }}
parameters:
{{- if eq .Storage.Backend "local" }}
  volumeBackend: local
  volumeParentDir: {{ .Storage.LocalVolumeParentDir | default "/var/lib/codefresh/dind-volumes" }}
{{- else if eq .Storage.Backend "gcedisk" }}
  volumeBackend: {{ .Storage.Backend }}
  #  pd-ssd or pd-standard
  type: {{ .Storage.VolumeType | default "pd-ssd" }}
  # Valid zone in GCP
  zone: {{ .Storage.AvailabilityZone }}
  # ext4 or xfs (default to ext4 because xfs is not installed on GKE by default )
  fsType: {{ .Storage.FsType | default "ext4" }}
{{- else if or (eq .Storage.Backend "ebs") (eq .Storage.Backend "ebs-csi")}}
  # ebs or ebs-csi
  volumeBackend: {{ .Storage.Backend }}
  #  gp2 or io1
  VolumeType: {{ .Storage.VolumeType | default "gp2" }}
  # Valid zone in aws (us-east-1c, ...)
  AvailabilityZone: {{ .Storage.AvailabilityZone }}
  # ext4 or xfs (default to ext4 )
  fsType: {{ .Storage.FsType | default "ext4" }}
  # "true" or "false" (default - "false")
  encrypted: "{{ .Storage.Encrypted | default "false" }}"
  {{ if .Storage.KmsKeyId }}
  # KMS Key ID
  kmsKeyId: {{ .Storage.KmsKeyId }}
  {{- end }}
{{- end }}`

	templatesMap["venonaconf.secret.venona.yaml"] = `apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: {{ .AppName }}conf
  namespace: {{ .Namespace }}
data:
{{ range $key, $value := .runnerConf }}
  {{ $key }}: {{ $value }}
{{ end }}`

	return templatesMap
}
