# json-parser

A command-line tool to parse and validate JSON from a file.

## Install

### From source

```bash
git clone https://github.com/your-username/json-parser.git
cd json-parser
go build -o json-parser

# Move to PATH
sudo mv json-parser /usr/local/bin/
```

## Usage

`json-parser` takes a file path as an argument and checks if the file contains valid JSON and parses it into go data structures.

```bash
json-parser my_file.json
```
