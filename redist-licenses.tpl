# Dependencies Licenses

Slide uses the following open source software and associated licenses.

| Name                   | Version                | License                |
|------------------------|------------------------|------------------------|
{{ range . -}}
| {{ .Name }} | {{ .Version }} | [{{ .LicenseName }}]({{ .LicenseURL }}) |
{{ end -}}

## Licenses 
{{ range . }}
### {{ .Name }}

* Name: {{ .Name }}
* Version: {{ .Version }}
* License: [{{ .LicenseName }}]({{ .LicenseURL }})

```
{{ .LicenseText }}
```
{{ end }}