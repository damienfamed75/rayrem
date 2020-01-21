# RayRem

---

## Dependencies

- [Go 1.13.3](#Golang)
- [pkger](#pkger)

## Installation

### pkger

```sh
go get github.com/markbates/pkger/cmd/pkger
```

### Golang

Follow Go's official [install instructions.](https://golang.org/doc/install)

### Clone repository

```sh
git clone https://github.com/hecate-tech/endorem.git
```

```sh
cd endorem/
```

### Package assets

```sh
pkger
```

### Build the binary

```sh
go build .
```

### Encrypt settings

```sh
go run ./cmd/encrypter
```

### Package the game

The structure of the game should be like this:

```txt
 ┌─ config/
 │     └ game.config
 └─ endorem.exe
```
