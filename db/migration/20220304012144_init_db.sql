-- +goose Up
-- +goose StatementBegin
CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "password" varchar NOT NULL,
  "created_at" datetime NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE "scrape" (
  "id" bigserial PRIMARY KEY,
  "user_id" int NOT NULL,
  "url" varchar NOT NULL,
  "scrapped" varchar NOT NULL,
  "created_at" datetime NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_scrape_url 
ON "scrape"("url");
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_scrape_date 
ON "scrape"("created_at");
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE "scrape" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE scrape;
-- +goose StatementEnd
