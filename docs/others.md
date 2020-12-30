### Mock
Use mockery.  
Mockery generates automatically mock types.
```shell
docker run -v "$PWD":/src -w /src vektra/mockery --all --keeptree --case underscore
```

### Migrations
Use goose.  
By following command, you can migrate all sql to database.
```shell
cd migration
goose mysql "root:example@/sample_project?parseTime=true" up
```

Please see [documentation of goose](https://github.com/pressly/goose) if you want to do other operations.
