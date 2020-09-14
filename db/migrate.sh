#!/usr/bin/env sh

touch -a ./db/.applied_migrations

for filename in ./db/migration/*.sql; do
  if ! grep -Fxq "$filename" ./db/.applied_migrations
  then
    sqlite3 db/leader_board.db "$(<$filename)"
    echo $filename >> ./db/.applied_migrations
  fi
done
