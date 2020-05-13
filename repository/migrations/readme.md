Database setup 
-
- Ensure that you have a db named _innersource_ first.
- Run the following command in cmd to get the postgres db setup.  
_https://github.com/golang-migrate/migrate/releases_

``migrate.exe -database postgres://postgres:{postgres_user_pwd}@localhost:5432/innersource?sslmode=disable -path .  up
``

