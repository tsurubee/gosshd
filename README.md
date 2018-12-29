# gosshd
Simple SSH server implemented by Golang.  

## Usage
### Installation
```
$ go get github.com/tsurubee/gosshd
```

### Launch the SSH server
For host authentication of the SSH server, create a public key / private key pair and put it in the root directory of the repository.  
```
ssh-keygen -t rsa -N '' -f ./id_rsa
```
Register the generated public key (id_rsa.pub) in the known_hosts of the local PC like below.  
```
[localhost]:2222 ssh-rsa AAAAB3・・・・
```
The SSH server starts up on Docker with the command below.  
```
$ docker-compose up
Starting gosshd_gosshd_1 ... done
Attaching to gosshd_gosshd_1
gosshd_1  | ==> Installing Dependencies
gosshd_1  | go get -u github.com/golang/dep/...
gosshd_1  | dep ensure
gosshd_1  | go run main.go
gosshd_1  | 2018/07/28 12:44:31 Listening on 2222...
```
Connecting SSH from the local PC to port 2222 leads to the inside of the container.  
(Since user authentication function is turned off, any user name is OK.)  
```
$ ssh tsurubee@localhost -p 2222
root@9cd2bdaf33c0:/go/src/gosshd#
```
Commands can be executed like ordinary connections to the server via SSH.  
```
root@9cd2bdaf33c0:/go/src/gosshd# pwd
/go/src/gosshd
root@9cd2bdaf33c0:/go/src/gosshd# ls
Gopkg.lock  Makefile   docker-compose.yml  id_rsa.pub  vendor
Gopkg.toml  README.md  id_rsa          main.go
```

## Blog
- [Golangで軽量なSSHサーバを実装する](https://blog.tsurubee.tech/entry/2018/07/28/225520)