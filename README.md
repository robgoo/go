# webdavhttps.go

webdav via https in go:

create the user to start webdavproces as non-root user:
The -r option creates a system user, and -s /bin/false sets the user’s shell to /bin/false, preventing direct login.
sudo useradd -r -s /bin/false webdavuser

Adjust the permissions (755 or otherwise) as necessary, depending on whether the server needs write access or not.

sudo chown webdavuser:webdavuser /path/to/davdir

sudo chmod -R 755 /path/to/davdir

Initialize a new Go module:
go mod init securewebdav

Fetch the required dependencies:
go get -u golang.org/x/net/webdav
Build WebDAV server:
go build securewebdav.go

Limit Network Exposure: Bind the service to a non-privileged port (above 1024) to avoid the need for root privileges to bind to ports 80 or 443

To access WebDAV server, use a WebDAV client or a compatible browser, navigate to https://localhost:8443/ . Since using a self-signed certificate, 
browser warn you about the site's security. proceed with the warning or add an exception for development environment.

sudo ufw allow from 123.45.67.89 to any port 8443
sudo ufw limit 8443

openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
When prompted for the Common Name enter the IP address or domain name of WebDAV server, 192.168.178.249.

Run WebDAV server with the required username and password for Basic Authentication, and specify the directory want to serve:
not as root!
sudo -u webdavuser ./securewebdav -addr=":8443" -cert="cert.pem" -key="key.pem" -user="pi" -pass="qpzm0123" -davdir="/tmp/webdavtest"

sudo mount -t davfs -o trust_server_cert https://192.168.178.249:8443/ /mnt/wd
Please enter the username to authenticate with server
https://192.168.178.249:8443/ or hit enter for none.
  Username: pi
Please enter the password to authenticate user pi with server
https://192.168.178.249:8443/ or hit enter for none.
  Password:  

improvement?
-data in rest encryption
-mkcert

*/
