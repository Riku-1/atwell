### Mock
Use mockery.  
Mockery generates automatically mock types.
```shell
docker run -v "$PWD":/src -w /src vektra/mockery --all --keeptree --case underscore
```