
# USAGE

	namaste [OPTIONS] [ACTION] [ACTION PARAMETERS...]

## SYNOPSIS

The name "namaste" in this context is derived from the phrase
"Name as text". It is a way to encode a directory with basic
metadata (e.g. directory type, who created it,
what is it, when was it created).

You can see "namaste" metadata by looking at the directory
contents without any special software. Namaste fields start with
zero to 4 followed by an equal sign. I.e.  0 - type, 1 - who,
2 - what, 3 - when, 4 - where.

```
    0=bagit_0.1
    1=Twain,M.
    2=Hamlet
    3=2005
    4=Seattle
```

## OPTIONS

Options are shared between all actions and must precede the action on the command line.

```
    -V, -verbose              verbose output
    -d, -directory            directory
    -generate-markdown-docs   output documentation in Markdown
    -h, -help                 display help
    -json                     output in JSON format
    -l, -license              display license
    -v, -version              display version
```


## ACTIONS

```
    get        returns all the namaste metadata of a directory if known
    gettypes   returns the types of a directory if known
    type       returns the type of a directory if known
    what       returns the what value of a directory if known
    when       returns the when value of a directory if known
    where      returns the where value of a directory if known
    who        returns the who value of a directory if known
```


Related: [get](get.html), [gettypes](gettypes.html), [type](type.html), [what](what.html), [when](when.html), [where](where.html), [who](who.html)

namaste v0.0.1
