
# List all keys from default bucket

slide ls

# List all keys from `todo` bucket

slide ls @todo

# List all keys from `todo` bucket with template

slide ls @todo --template '{{ .Key }} -> {{ .Value }}'

# Export all keys from `env` bucket as environment variables

$(slide ls @env --template 'export {{ .Key }}="{{ .Value }}"')
