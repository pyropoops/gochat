# gochat

## Description
A PoC for a future chat application.

gochat compiles both the backend and frontend into a single native executable, so there is no extra installation required.

## Usage [Console Commands]
### Registering Users
```register <username> <password>```
### Deregistering Users
```unregister <username>```
### Kicking users
```kick <username>```
### Viewing Active Connections
```sessions```
### Broadcasting messages from console
```broadcast <message>```

## Building

### Prerequisites 
- Golang 
- npm (NodeJS)
- Windows/macOS (Linux untested.)

### Compiling
```
npm i
./build.ps1 <addressable IP>
```

### Running
```
.\server.exe
```