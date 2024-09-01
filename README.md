# TODOS

- Write tests

# Development

Spin up the supabase db container

```shell
$ supabase start
```

Migrate and seed local db
```shell
$ supabase db reset
```

Connect to local db
```shell
$ supabase status
$ psql DB_URL
```

Run the main function:
`go run .`

# Helpful Docs

- [Supabase](https://supabase.com/docs/guides/database/overview)