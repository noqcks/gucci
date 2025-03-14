{{- $l := list $.testTwo $.testTwo -}}
{{- range $i := $l }}
  {{- range $j := $i }}
    V ({{ kindOf $j }}): {{ $j | mustToJson | nindent 2 }}
  {{- end }}
{{- end }}

{{- /* Test case for numeric and boolean keys */ -}}
Numeric Keys (toJson): {{ $.testThree.numericKeys | toJson }}
Boolean Keys (mustToJson): {{ $.testThree.booleanKeys | mustToJson }}

{{- /* Test case for deeply nested structures */ -}}
Deeply Nested (toJson): {{ $.testThree.deeplyNested | toJson }}

{{- /* Test case for mixed array types */ -}}
Mixed Array (mustToJson): {{ $.testThree.mixedArray | mustToJson }}

{{- /* Test the entire structure */ -}}
Complete testThree (toJson): {{ $.testThree | toJson }}