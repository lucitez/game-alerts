This program checks a pickup soccer league website to check if games have been scheduled, then sends emails to subscribers.

## Migrations
```shell
$ supabase migration new <migration_name>
```

```shell
$ supabase login
```

```shell
$ supabase db push
```

## Development

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

## Helpful Docs

- [Supabase](https://supabase.com/docs/guides/database/overview)
