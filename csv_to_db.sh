#!/bin/bash

# this script is used to fill disciplines, athletes and athletes_disciplines tables in the database using csv files that were scrapped from the web
# example call: CSV_PATH=/home/user/scrapped_csvs ./csv_to_db.sh
disciplines_csv="$CSV_PATH/disciplines.csv"
athletes_csv="$CSV_PATH/athletes.csv"
athletes_disciplines_csv="$CSV_PATH/athletes_disciplines.csv"

psql -h $DB_HOST -U $DB_USER -d $DB_NAME <<EOF
BEGIN;
\copy disciplines(id, name, type) FROM '$disciplines_csv' DELIMITER ',' CSV HEADER;
\copy athletes(id, first_name, last_name, birthday, country, gender) FROM '$athletes_csv' DELIMITER ',' CSV HEADER;
\copy athletes_disciplines(discipline_id, athlete_id) FROM '$athletes_disciplines_csv' DELIMITER ',' CSV HEADER;
UPDATE disciplines SET created_at = NOW(), updated_at = NOW();
UPDATE athletes SET created_at = NOW(), updated_at = NOW();
COMMIT;
EOF