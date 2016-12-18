Simple library for solve extra verbose error handling in golang.

It allow to create wrapper object to handle error, then call operation through it. The wrapper object before every call
will check about errors on previous step. If true - it do nothing (skip really work) and return error.

You can check error only on last step of work and doesn't worry about mix ok and error calls - all calls after first error will be skip.

It is unstable now and API may change often without notice. If you use the library in production - vendor it.