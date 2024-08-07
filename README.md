# Abbreviation webapp

A webapp to search abbreviations & acronyms in the public sector.

## Stack

- Go, Chi, Templ
- Htmx
- Tailwind, DaisyUI
- SQLite for simple storage/retrieval
- Jaro-Winkler for fuzzy string matching
- (Dev) Air, Prettier, Makefile (includes Init() if building on Windows)

## Alternative implementation

- Tried to rewrite this some time back in Sveltekit (https://github.com/zachczx/svelte-letters) after getting frustrated with Go. Dropped it when Go wasn't so bad!

## Note to self:

- For doing local devt, hot reload via Air, Tailwind generation - on localhost:7331

  ```sh
  make dev
  ```

- For local docker containerization

  ```sh
  docker compose -f ./docker-compose-dev.yaml build && docker compose -f ./docker-compose-dev.yaml up
  ```

- For building via Coolify (if Traefik or Caddy config disappear, generate domain at first instance of adding a resource and edit it from there)
  ```sh
  $ docker compose build
  $ docker compose up
  ```
