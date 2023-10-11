# gotmpl

A CLI tool for templating files using input data from a json
file. Gotmpl features the [Sprig functions](http://github.com/Masterminds/sprig)
for richer templating. You can find the function documentation at
<http://masterminds.github.io/sprig/>.

```consoletext
$ gotmpl -h
usage: gotmpl [option] [values_json] file_to_template
If values_json is omitted, values are read from stdin (as json)
Sprig functions are supported, see https://masterminds.github.io/sprig/
gotmpl also supports the following aditional function(s):
shellescape string (returns string)
Flags:
  -o string
        output file, default is stdout

$ echo '{"hello":"world"}' > values.json

$ echo 'Hello {{ .hello }}' > file_to_template

$ gotmpl values.json file_to_template 
Hello world

$ gotmpl -o output values.json file_to_template
```

You can also use environment variables as template values via the
Sprig `env` function and mix with values from `values.json`.

```consoletext
$ export MY_ENV_VAR=world

$ echo 'Hello {{ env "MY_ENV_VAR" }}' > file_to_template

$ echo '{}' | gotmpl file_to_template
Hello world
```

The `shellescape` function can be used to escape arguments for shell scripts, e.g:

```consoletext
$ echo '{"name":"Eugene Belford;RM -RF /"}' > values.json

$ echo 'Hello {{ .name }}, type the following in a terminal: `echo {{ .name | shellescape }}`' > file_to_template

$ gotmpl values.json file_to_template
Hello username;RM -RF, type the following in a terminal: `echo 'username;RM -RF'`
```

## Building and installing

Please use Go 1.21.2 or later...

```consoletext
go install github.com/sa6mwa/gotmpl@latest

# or use the Makefile...

make install
```
