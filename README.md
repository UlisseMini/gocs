# goc
goc is a powerful project generation tool written in go,
template files can be put in ~/.goc/templates and used without changing any code.

## Installation
Download a binary from the releases page

## Usage
if template is not supplied the default for golang is used.
the default can be changed in ~/.goc/templates
```bash
goc <projectname> [template]
```

## Examples
```bash
# look inside the ~/.goc/templates directory and use templates from there to create
# files.
goc foobar go
```

```bash
# use the default template to create a project called foobar
goc foobar
```

## Development
```bash
go get -u github.com/UlisseMini/goc
cd $GOPATH/src/github.com/UlisseMini/goc
GO111MODULE=on go mod download
```

## Contributing

1. Fork it (<https://github.com/ulissemini/goc/fork>)
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request

## Contributors

- [Ulisse Mini](https://github.com/UlisseMini) - creator and maintainer
