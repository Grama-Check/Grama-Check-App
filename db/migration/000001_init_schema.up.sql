CREATE TABLE "users" (
  "nic" varchar PRIMARY KEY NOT NULL,
  "name" varchar NOT NULL,
  "address" varchar NOT NULL,
  "email" varchar NOT NULL,
  "idcheck" boolean NOT NULL,
  "addresscheck" boolean NOT NULL,
  "policecheck" boolean NOT NULL,
  "failed" boolean NOT NULL
);

CREATE INDEX ON "users" ("name");
