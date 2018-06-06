{{/* This files contains templates for the variable name transition */}}

{{/* Get a variable, with an error message if it's missing. Arguments are passed as a list. */}}
{{- define "getvar"}}{{/* (dict "ctx" . "names" (list "A" "B") "quote" true "b64" false) */}}
    {{- with $v := . -}}
        {{- if hasKey $v "quote" | not -}}
            {{- with set $v "quote" true -}}
                {{/* We don't actually need the result of the "set", but need to discard it */}}
            {{- end -}}
        {{- end -}}
        {{- if hasKey $v "b64" | not -}}
            {{- with set $v "b64" false -}}
                {{/* We don't actually need the result of the "set", but need to discard it */}}
            {{- end -}}
        {{- end -}}
        {{- range $name := index $v "names" | reverse -}}
            {{- if hasKey (index $v "ctx" "Values" "env") $name -}}
                {{- with set $v "result" (index $v "ctx" "Values" "env" $name) -}}
                    {{/* We don't actually need the result of the "set", but need to discard it */}}
                {{- end -}}
            {{- end -}}
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
