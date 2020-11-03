# web-bomberman

In diesem Projekt soll ein Spiel im Bombermanstil entstehen, welches mit bis zu x-Personen gleichzeitig im Browser gespielt werden kann.

Server: Go 
Clientkommunikation: WebSockets

Für die erfolgreiche Ausführung ist es notwendig, im root-Verzeichnis des Projekts eine go-Datei mit dem package "main" zu erstellen, in der folgende Konstanten vorhanden sein müssen:
+ "DB_URL": Adresse mit Port des Datenbankservers
+ "DB_NAME": Name der Datenbank
+ "DB_USERNAME": Username des Users des Datenbankservers
+ "DB_PASSWORD": Passwort des Users des Datenbankservers

### Das SQL-Skript, um die Datenbank aufzubauen, ist in "databse.sql" zu finden.

Nach dem Start des Servers ist der Webclient auf Port http://localhost:2100/ zu erreichen.



