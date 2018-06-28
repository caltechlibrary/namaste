
# Namaste

Namaste (i.e. "NAMe AS TExt") is a file naming convention to support primitive directory-level metadata tags 
exposed directly via filenames. As such, Namaste tags greet visitors who request a directory listing 
(e.g., Linux 'ls') with a glimpse of what the directory holds. This approach can be readily implemented
as a library or as a command line tool. Prior implementations exist in Perl and Python. This implementation
was written in [Go](https://golang.org). It includes support for namaste hosted on S3.

## namaste command

The command `namaste` can be used to set, check or retrieve the "name as text" associated with
a directory's type, who, what, when and where metadata. If you use the `-verbose` option
the values set or retrieved with be display to standard out.

### operations

+ [get](get.html) - retrieves _namaste_ data
+ [gettypes](gettypes.html) - retreives _namaste_ type information
+ [type](type.html) - sets type information for a directory
+ [what](what.html) - sets the content description for a directory
+ [when](when.html) - sets an associated date string with a directory
+ [where](where.html) - sets a location string with a directory
+ [who](who.html) - sets a person's name associated with a directory

### options

+ -d, -directory sets the directory to operate on (default is current directory)
+ -verbose - display verbose output


### Reference

+ [Namaste](https://confluence.ucop.edu/display/Curation/Namaste)
+ [Perl Implementation](https://metacpan.org/pod/File::Namaste)
+ [Python Implementation](http://github.com/mjgiarlo/namaste)


