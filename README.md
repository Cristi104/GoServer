# Go Server

## Description

Go server is a simple instant messaging web app similar to WhatsApp and other messaging apps. 

This app has not been created for actual usage, I only created it to improve my development skills and capabilities (mostly to try and understand fullstack web development) and to learn the Go programing language

The application's stack consists of:
- a MySql database used to store account and message data
- the backend written in Go
- the frontend created with plain JavaScript HTML and CSS

**This app is currently under active development. Features, functionality, and design are subject to change.**

## Instalation

### Prerequesites

Before installing the actual app you first need to have a MySql database

#### Using Docker

Repalce <CONTAINER_NAME> and <DATABASE_PASSWORD> with your desired name and password

```sh
sudo docker run --name <CONTAINER_NAME> -e MYSQL_ROOT_PASSWORD=<DATABASE_PASSWORD> -p 3306:3306 -d mysql
```

### Running the app

Clone or download the repository and cd into its root directory

```sh
git clone https://github.com/Cristi104/GoServer.git
cd Goserver
```

Before running the app you first need to create the **DB.json** inside the **config/** directory, this json file will hold the information requiered to connect to the MySql database.

```sh
mkdir config
cd config
echo "{\"User\":\"root\",\"Passwd\":\"<DATABASE_PASSWORD>\",\"Net\":\"tcp\",\"Addr\":\"127.0.0.1:3306\",\"DBName\":\"mysql\"}" > DB.json
cd ..
```

The DB.json file:

```json
{
    "User":"root",
    "Passwd":"<DATABASE_PASSWORD>",
    "Net":"tcp",
    "Addr":"127.0.0.1:3306",
    "DBName":"mysql"
}
```

To run the app just run **go run .** in the root directory of the app this will start the app which will search for the database connection information inside the **config/DB.json** 

```sh
go run .
```

If everything works as intended the app will start and output "Connected to database" after witch it will start listening for connections on the **8080** port and can be navigated by going to 127.0.0.1:8080 in your web browser