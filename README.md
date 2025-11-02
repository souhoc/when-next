# When next?
That's the question we ask ourself with my friends at the end of each sessions.
Usually we create a discord poll with the following days by hand. But hey, why
taking 2 minutes to do it if I can take 20 minutes to make a CLI that do it for
me?

## Installation
```fish
go install github.com/souhoc/when-next@latest
```

## How to run
```fish
when-next -webhook https://discord.com/api/webhooks/****
```
Or with a config:
```ini
# config.ini
-webhook=https://discord.com/api/webhooks/****
-layout=Mon _2 Jan
```
Run it with the `config` flag.
```fish
when-next -config config.ini
```
