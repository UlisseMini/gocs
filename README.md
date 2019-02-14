# gocs - <b>go</b> <b>c</b>ommit <b>s</b>udoku
gocs is a powerful project generation tool written in go,
template files can be put in ~/.gocs/templates and used without changing any code.
![gif](https://github.com/UlisseMini/gocs/raw/master/pictures/example.gif)

## Installation
gocs is fully self contained, simply download a binary from the releases page and you're good to go!

## Usage
if template is not supplied the default for golang is used.
the default can be changed in ~/.gocs/templates
```bash
gocs <projectname> [template]
```

when creating templates you can hook into the following variables inside files and filenames.
```
{{.Project}} | name of the project being created.
{{.Year}}    | the current year, useful for licenses
{{.Github}}  | Github username, stored in ~/.gocs/config.yaml
{{.Author}}  | Full name of the author, stored in ~/.gocs/config.yaml
```

## Examples

```bash
# look inside the ~/.gocs/templates/go directory and use templates from there to create
# files.
gocs foobar go
```

```bash
# use the default template to create a project called foobar
gocs foobar
```

example template file for another language (python) using the MIT license
to use this you could create a directory `~/.gocs/templates/python` then generate
a project using it with `gocs myproject python` where python says "use the python template"
```python
# Copyright {{.Year}} {{.Author}}
# Permission is hereby granted, free of charge, to any person obtaining
# a copy of this software and associated documentation files... (license continues)

def main():
	# TODO Write code
	pass

if __name__ == "__main__":
	main()
```

## Development
```bash
go get -u github.com/UlisseMini/gocs
cd $GOPATH/src/github.com/UlisseMini/gocs
GO111MODULE=on go mod download
```

## Contributing

1. Fork it (<https://github.com/ulissemini/gocs/fork>)
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request

## Contributors

- [Ulisse Mini](https://github.com/UlisseMini) - creator and maintainer
