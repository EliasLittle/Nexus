# For Ubuntu
#
# install golang
sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt update
sudo apt install golang-go

# install protoc
apt install -y protobuf-compiler

# install protoc golang plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
export PATH="$PATH:$(go env GOPATH)/bin"


make all
make install-server

mkdir /usr/local/etc/nexus/
mkdir /usr/local/etc/nexus/logs
touch /usr/local/etc/nexus/index.json
