{{- define "requip.name" -}}
{{- default .Release.Name .Values.nameOverride | trunc 40 | trimSuffix "-" -}}
{{- end -}}

{{- define "requip.labels" -}}
app.kubernetes.io/part-of: requip
app.kubernetes.io/managed-by: {{ .Release.Service }}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version }}
{{- end -}}

{{- define "requip.backend.fullname" -}}{{ include "requip.name" . }}-backend{{- end -}}
{{- define "requip.frontend.fullname" -}}{{ include "requip.name" . }}-frontend{{- end -}}
{{- define "requip.mysql.fullname" -}}{{ include "requip.name" . }}-mysql{{- end -}}

{{- define "requip.backend.selectorLabels" -}}
app.kubernetes.io/name: {{ include "requip.backend.fullname" . }}
app.kubernetes.io/component: api
{{- end -}}
{{- define "requip.frontend.selectorLabels" -}}
app.kubernetes.io/name: {{ include "requip.frontend.fullname" . }}
app.kubernetes.io/component: web
{{- end -}}
{{- define "requip.mysql.selectorLabels" -}}
app.kubernetes.io/name: {{ include "requip.mysql.fullname" . }}
app.kubernetes.io/component: database
{{- end -}}

{{- define "requip.backend.image" -}}
{{ .Values.backend.image.repository }}:{{ .Values.backend.image.tag | default .Chart.AppVersion }}
{{- end -}}
{{- define "requip.frontend.image" -}}
{{ .Values.frontend.image.repository }}:{{ .Values.frontend.image.tag | default .Chart.AppVersion }}
{{- end -}}

{{/* Restricted Pod Security Standard settings shared by every container. */}}
{{- define "requip.containerSecurityContext" -}}
allowPrivilegeEscalation: false
readOnlyRootFilesystem: true
capabilities:
  drop: ["ALL"]
{{- end -}}
