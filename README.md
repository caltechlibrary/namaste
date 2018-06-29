
# namaste [![DOI](https://data.caltech.edu/badge/79394591.svg)](https://data.caltech.edu/badge/latestdoi/79394591)


This is a Go language port of Namaste (Name as Text) library inspired by the Python
implementation of [namaste](https://github.com/mjgiarlo/namaste) which in turn was 
ported from [Perl implementation](https://metacpan.org/pod/File::Namaste).


## Go Package

The Go package seeks to mimic the Python implementation as much as possible while
expanding on eventual storage options.


## Command Line implementation

The [command line implementation](docs/)'s purpose is to exercise the Go package
and test practical usage and integration into other Caltech Library projects.
It is intended to be largely command line compatible with the Python 
implementation. 


## Extensions

There are a few things that extend beyond the existing Python namaste 
implementation which are desirable to test. 

+ S3 storage support
+ Google Cloud Storage support
+ Output namaste metadata in JSON format

