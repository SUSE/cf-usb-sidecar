{{/* This files contains templates for the variable name transition */}}

{{/* Set a default value in a dict */}}
{{- define "_setdefault" -}}{{/* (list <dict> <name> <default>) */}}
    {{- with $v := . -}}
        {{- if hasKey (index $v 0) (index $v 1) | not -}}
            {{- with set (index $v 0) (index $v 1) (index $v 2) -}}
                {{/* We don't actually need the result of the "set", but need to discard it */}}
            {{- end -}}
        {{- end -}}
    {{- end -}}
{{- end -}}

{{/* Get a variable, with an error message if it's missing. Arguments are passed as a list. */}}
{{- define "getvar" }}{{/* (dict "ctx" . "names" (list "A" "B") "quote" true "b64" false) */}}
    {{- with $v := . -}}
        {{- template "_setdefault" (list $v "quote" true) -}}
        {{- template "_setdefault" (list $v "b64" false) -}}
        {{- range $name := index $v "names" -}}
            {{- if hasKey (index $v "ctx" "Values" "env") $name -}}
                {{- if typeIs "<nil>" (index $v "ctx" "Values" "env" $name) | not -}}
                    {{- template "_setdefault" (list $v "result" (index $v "ctx" "Values" "env" $name)) -}}
                {{- end -}}
            {{- end -}}
        {{- end -}}
        {{- if hasKey $v "default" -}}
            {{- template "_setdefault" (list $v "result" (index $v "default")) -}}
        {{- end -}}
        {{- if (index $v "quote") -}}
            {{- if (index $v "b64") -}}
                {{- required (printf "env.%s configuration missing" (index $v "names" 0) ) (index $v "result") | b64enc | quote -}}
            {{- else -}}
                {{- required (printf "env.%s configuration missing" (index $v "names" 0) ) (index $v "result") | quote -}}
            {{- end -}}
        {{- else -}}
            {{- if (index $v "b64") -}}
                {{- required (printf "env.%s configuration missing" (index $v "names" 0) ) (index $v "result") | b64enc -}}
            {{- else -}}
                {{- required (printf "env.%s configuration missing" (index $v "names" 0) ) (index $v "result") -}}
            {{- end -}}
        {{- end -}}
    {{- end -}}
{{- end -}}
