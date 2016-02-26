#
# Official docker file for building Gossha server
#

FROM fedora:23

# Upgrade dependencies and clean dnf cache to reduce image size
RUN dnf upgrade -y && dnf clean all

# Listen on 22 port
ENV GOSSHA_PORT=22
EXPOSE 22

# Set one of possible database storages. We can use only one method

# Set database connection string for MySQL database
#ENV GOSSHA_DRIVER=mysql 
#ENV GOSSHA_CONNECTIONSTRING=user:password@localhost/dbname?charset=utf8&parseTime=True&loc=Local

# Set database connection string for PostgreSQL database
#ENV GOSSHA_DRIVER=postgres 
#ENV GOSSHA_CONNECTIONSTRING='user=gorm dbname=gorm sslmode=disable'
#ENV GOSSHA_CONNECTIONSTRING='postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full'

# Set database connection string for SQLite database
ENV GOSSHA_DRIVER=sqlite3
#ENV GOSSHA_CONNECTIONSTRING=':memory:'
ENV GOSSHA_CONNECTIONSTRING=/root/.gossha/gossha.db


# Create home directory
ENV GOSSHA_HOMEDIR=/root/.gossha

# Inject code of your application
ADD build/gossha /usr/bin/gossha

# Create root user
RUN /usr/bin/gossha passwd root root

# Create first ordinary user
RUN /usr/bin/gossha passwd user1 user1

# Create second ordinary user
RUN /usr/bin/gossha passwd user2 user2

# Inject SSHD server keys
ADD build/id_rsa /root/.ssh/
ADD build/id_rsa.pub /root/.ssh/

# Run the image process. Point second argument to your entry point of application
CMD ["/usr/bin/gossha"]