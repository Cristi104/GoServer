# Go Server

## Description

Go server is a simple instant messaging web app similar to WhatsApp and other messaging apps. 

This app has not been created for actual usage, I only created it to improve my development skills and capabilities (mostly to try and understand fullstack web development) and to learn the Go programing language

The application's stack consists of:
- a PostreSQL database used to store account and message data
- the backend written in Go
- the frontend created with plain React.js and tailwindcss

**This app is under active development. Features, functionality, and design are subject to change.**

## Demo

Curently the app is hosted as a demo at [goserverproject.duckdns.org](https://goserverproject.duckdns.org) (it is self hosted with a self certified ssl so your browser will warn you about the site not beeing secure).
You can login with one of the test users:

- testuser1@email.com Password1!
- testuser2@email.com Password2!
- testuser3@email.com Password3!


## Instalation

### Using Docker Compose

Clone or download the repository and cd into its root directory.

```sh
git clone https://github.com/Cristi104/GoServer.git
cd Goserver
```

Run docker compose to build and start the app.

```sh
docker compose up -d
```

Once the containers are created you will need to create the database objects.

```sh
# Default password: root_pass
./scripts/run_migrations.sh
```

Once this is done the app will be open on **localhost:8080**.
