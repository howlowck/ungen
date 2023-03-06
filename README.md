# Ungen

Flipping a code generators on its head.

## Design Principles

* Prototype First (no custom generator code)
* Works in any programming language
* Simple and human-readable DSL

## How to Use the CLI

Easiest way is to use the docker image, "howlowck/ungen". You can also download the executible for your OS.

There are a few options associated with the tool, you can run ungen to see the text, but here is the output below:

```
-i string
    InputDirectory (Required)
-o string
    OutputDirectory (Required)
-keep
    Keep the UNGEN line
-zip
    Zip the output directory into a file
-var value
    Set Variables (ex. -var foo=bar -var baz=qux)
```

You can set `-var` option multiple times (just like Terraform var flags)

## Writing UNGEN commands

Ungen is meant to be human readable and simple to use, but it is a programming language.

To create a Ungen command, you just have to tag your line comments with "UNGEN:" (just like you would with `TODO:`) in any code you are writing. Like so:

```js
// UNGEN: replace "World" with kebabCase(var.appName)
app.get('/', (req, res) => res.send('Hello World!'));
```

The command above will replace all occurance of "World" with the kebab-case form of the value of the `appName` variable.

## Language Features:

* `if / then / else` - conditionals. `else` is optional
* `copy or cut` - "clipboard" operations
* `"hello"` - string literals are in double quotes. This would be a value of a string `Hello`.
* `var.<variable-name>` - variable values are prefix with `var.`
* `camelCase(<value>)` - string functions are called like so, and will return a string value.

## Operations (more to come)

### Copy or Cut into the Clipboard
Specify the next x number of lines (`next 2 lines`) or line number (`ln 2`) or line number range (`ln 3-6`), then specify the name of the clipboard value (`cb.myText`)

Example: `// UNGEN: copy next 1 line to cb.myText` it will store the next line into the `myText` cb value.

### String Replace
Looks at the next line, and replace any string value to another value.

Example: `// UNGEN: replace "World" with var.appName`

### String Delete
Delete the n number of lines.

Example: `// UNGEN: delete 2 lines`

### File Delete
Delete the current file.

Example: `// UNGEN: delete file`

### Directory Delete
Delete the current directory

Example: `// UNGEN: delete folder`

## String Functions
* `kebabCase`: "Hello World" -> `hello-world`
* `snakeCase`: "Hello World" -> `hello_world`
* `camelCase`: "Hello World" -> `helloWorld`
* `upperCase`: "Hello World" -> `HELLO WORLD`
* `lowerCase`: "Hello World" -> `hello world`
* `substitute`: subtitute("Hello World", "lo", "ping") -> `Helping World`
* `concat`: concat("Hello ", "World", "!") -> `Hello World!`

## Language Feature: Injection

In some "languages", line comments are not allowed (like JSON) or not common (like Markdown), UnGen allows you to "inject" Ungen commands into those files before the parser processes the files. All you need to do is create a `.ungen` file in the same directory, and write your injection commands like so:

`// UNGEN: inject file:<filename> on ln <line-number> '<ungen-command'`

An example can be found in `examples/simple-nodejs/.ungen`

## Developer Notes
* Run Test: `go test -timeout 30s -run ^TestLexer$ github.com/howlowck/ungen -v`
