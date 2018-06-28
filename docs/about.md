
# What is Namaste?

The name "Namaste" in this context is derived from the phrase
"name as text". It is a way to encode a directory with basic
metadata (e.g. directory type, who created it, what is it, 
what it is, when was it created, where it was created).

You can see "namaste" metadata by looking at the directory
contents without any special software. Namaste fields start with
zero (type), one (who), two (what), three (when) or four (where).
This is followed by an equal sign then the value of the metadata
field.

## Example

```
    0=bagit_0.1
    1=Twain,M.
    2=Hamlet
    3=2005
    4=Seattle
```

