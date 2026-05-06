# glirc

glirc is a CLI-based IRC client written in Go.

---

# What the Project Is

- Learning project for Go and IRC

---

# How It Works

- `type Message struct` for incoming and outgoing messages
- parser/serializer for messages
- `type Client struct` with handler maps for client and server message processing, goroutines and channels for I/O

---

# Engineering Decisions

## `type Client struct`

Took all code out of the main method and put it in `client.go` for separation of concerns and easier debugging.

## `*_handlers.go`

Parse input/output and call the appropriate method based on message command. Much easier to create, update and delete
commands and their matching handlers.

# Tech Stack

- Languages used: Go
- Frameworks/libraries: tview (TUI)

---

# What I Learned

- IRC is really cool!
- Go is very interesting, glad to have got the experience with it.

---

# How to Run the Project

Note: you will need a `config.json` file containing your nickname, username, server and port. Place this in the root
directory after cloning the repository.

```
{
    "nick": "yournick",
    "user": "Your Name",
    "server": "irc.libera.chat",
    "port": 6697
}
```

```
git clone https://github.com/liampallett/glirc.git 
cd glirc
go build
./glirc
```

---

# Project Structure

```
├── LICENSE.md
├── README.md
├── client.go
├── client_handlers.go
├── config.go
├── go.mod
├── main.go
├── message.go
├── parser_test.go
└── server_handlers.go
```

# Future Improvements

- TUI integration with an existing library.
- More supported commands.

---
