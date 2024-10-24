-- Create "app_models" table
CREATE TABLE "app_models" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create "athletes" table
CREATE TABLE "athletes" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "first_name" text NULL,
  "last_name" text NULL,
  "birthday" date NULL,
  "country" text NULL,
  "gender" text NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "pokemons" table
CREATE TABLE "pokemons" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "pokemon_name" text NULL,
  "age" bigint NULL,
  "email" character varying(100) NOT NULL,
  "attack" text NOT NULL,
  "defense" text NULL,
  PRIMARY KEY ("id")
);
-- Create "disciplines" table
CREATE TABLE "disciplines" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "name" text NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "athletes_disciplines" table
CREATE TABLE "athletes_disciplines" (
  "discipline_id" bigint NOT NULL,
  "athlete_id" bigint NOT NULL,
  PRIMARY KEY ("discipline_id", "athlete_id"),
  CONSTRAINT "fk_athletes_disciplines_athlete" FOREIGN KEY ("athlete_id") REFERENCES "athletes" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_athletes_disciplines_discipline" FOREIGN KEY ("discipline_id") REFERENCES "disciplines" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "roles" table
CREATE TABLE "roles" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "name" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_roles_name" UNIQUE ("name")
);
-- Create "users" table
CREATE TABLE "users" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "email" text NOT NULL,
  "password" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_users_email" UNIQUE ("email")
);
-- Create "user_roles" table
CREATE TABLE "user_roles" (
  "user_id" bigint NOT NULL,
  "role_id" bigint NOT NULL,
  PRIMARY KEY ("user_id", "role_id"),
  CONSTRAINT "fk_user_roles_role" FOREIGN KEY ("role_id") REFERENCES "roles" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_user_roles_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
