------------
-- Goqite - https://github.com/maragudk/goqite/blob/main/schema.sql
------------
create trigger if not exists goqite_updated_timestamp after update on goqite begin
  update goqite set updated = strftime('%Y-%m-%dT%H:%M:%fZ') where id = old.id;
end;
