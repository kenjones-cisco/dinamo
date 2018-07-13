## dinamo generate

Generate files

### Synopsis

Generate files using go templates and multiple data sources.
Use any combination of the datasources JSON/YAML data files, environment variables, and key-value pairs arguments.


```
dinamo generate [flags]
```

### Examples

```

# create output.txt from config.tmpl using the key-value pairs
dinamo gen -t config.tmpl -f output.txt key1=value1 key2=value2

# create output.txt from config.tmpl using the JSON data file source.json
dinamo gen -t config.tmpl -f output.txt -d source.json

# create output.txt from config.tmpl using the YAML data file source.yaml
dinamo gen -t config.tmpl -f output.txt -d source.yaml

# create output.txt from config.tmpl using environment variables
dinamo gen -t config.tmpl -f output.txt -e

# create output.txt from config.tmpl using the key-value pairs and the YAML data file source.yml
dinamo gen -t config.tmpl -f output.txt -d source.yml key1=value1 key2=value2

# create output.txt from config.tmpl using the key-value pairs and the JSON data file source.json
dinamo gen -t config.tmpl -f output.txt -d source.json key1=value1 key2=value2

# create output.txt from config.tmpl using the key-value pairs and environment variables
dinamo gen -t config.tmpl -f output.txt -e key1=value1 key2=value2

# create output.txt from config.tmpl using the key-value pairs, environment variables, and the YAML data file source.yml
dinamo gen -t config.tmpl -f output.txt -e -d source.yml key1=value1 key2=value2

```

### Options

```
  -d, --data string       Path to data file of type ("json", "yaml", "yml")
  -e, --env               Use environment variables for placeholders
  -f, --file string       Path to generated file
  -h, --help              help for generate
  -t, --template string   Template file path
```

### Options inherited from parent commands

```
  -D, --debug              Enable debug mode
  -l, --log-level string   Set the logging level ("debug", "info", "warn", "error", "fatal") (default "info")
```

### SEE ALSO

* [dinamo](dinamo.md)	 - Dynamic Generator

