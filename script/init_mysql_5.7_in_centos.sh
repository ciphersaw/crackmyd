#!/bin/bash

# Define color variables
GREEN='\033[0;32m'
RED='\033[0;31m'
DEFAULT='\033[0m' # Reset to default color

# Check command execution result
check_command_result() {
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}[$(date "+%Y-%m-%d %H:%M:%S")][INFO] $@ successfully!${DEFAULT}"
    else
        echo -e "${RED}[$(date "+%Y-%m-%d %H:%M:%S")][ERROR] $@ unsuccessfully and exit.${DEFAULT}"
        exit 1
    fi
}

# Print information formatted
print_info() {
    echo -e "${GREEN}[$(date "+%Y-%m-%d %H:%M:%S")][INFO] $@${DEFAULT}"
}

# Print error formatted
print_error() {
    echo -e "${RED}[$(date "+%Y-%m-%d %H:%M:%S")][ERROR] $@${DEFAULT}"
}

# 1. Download and install MySQL official Yum Repository
print_info "Downloading MySQL Yum Repository..."
wget -c http://repo.mysql.com/mysql57-community-release-el7-10.noarch.rpm
check_command_result "Download MySQL Yum Repository"

# 2. Install Yum Repository
print_info "Installing MySQL Yum Repository..."
yum -y install mysql57-community-release-el7-10.noarch.rpm
check_command_result "Install MySQL Yum Repository"

# 3. Import RPM GPG Key
print_info "Importing MySQL RPM GPG Key..."
rpm --import https://repo.mysql.com/RPM-GPG-KEY-mysql-2022
check_command_result "Install MySQL RPM GPG Key"

# 4. Install MySQL server
print_info "Installing MySQL server..."
yum -y install mysql-community-server
check_command_result "Install MySQL server"

# 5. Start MySQL service
print_info "Starting MySQL service..."
systemctl start mysqld.service
check_command_result "Start MySQL service"

# 6. Check MySQL service status
print_info "Checking MySQL service status..."
systemctl status mysqld.service
check_command_result "Check MySQL service status"

# 7. Retrieve MySQL initial password
print_info "Retrieving MySQL initial password..."
initial_password=$(grep 'temporary password' /var/log/mysqld.log | awk '{print $NF}')
if [ -z "$initial_password" ]; then
    print_error "Initial password not found, please check MySQL log file: /var/log/mysqld.log"
    exit 1
else
    print_info "MySQL initial password is:" $initial_password
fi

# 8. Log in to MySQL and change password
# Prompt user to enter a new password
read -p "Please enter a new MySQL root password: " mysql_password
# Log in with the initial password and change the password
print_info "Changing MySQL root password..."
mysql -uroot -p${initial_password} --connect-expired-password -e "ALTER USER 'root'@'localhost' IDENTIFIED BY '${mysql_password}';"
check_command_result "Change MySQL root password"

# 9. Enable MySQL to start on boot
print_info "Enabling MySQL to start on boot..."
systemctl enable mysqld.service
check_command_result "Enable MySQL to start on boot"

# 10. Modify password policy and length
print_info "Modifying MySQL password policy and length to allow simple passwords..."
mysql -uroot -p${mysql_password} -e "SET GLOBAL validate_password_policy=LOW;"
check_command_result "Set MySQL password policy to LOW"
mysql -uroot -p${mysql_password} -e "SET GLOBAL validate_password_length=4;"
check_command_result "Set MySQL minimum password length to 4"

# 11. Create users
print_info "Creating MySQL users..."
mysql -uroot -p${mysql_password} -e "CREATE USER 'centos'@'%' IDENTIFIED BY 'qwerty';"
check_command_result "Create user 'centos'@'%'"
mysql -uroot -p${mysql_password} -e "CREATE USER 'centosnew'@'%' IDENTIFIED BY 'q1w2e3r4';"
check_command_result "Create user 'centosnew'@'%'"
mysql -uroot -p${mysql_password} -e "CREATE USER 'app'@'127.0.0.1' IDENTIFIED BY 'app123';"
check_command_result "Create user 'app'@'127.0.0.1'"
mysql -uroot -p${mysql_password} -e "CREATE USER 'appnew'@'127.0.0.1' IDENTIFIED BY 'appnew@gmail.com';"
check_command_result "Create user 'appnew'@'127.0.0.1'"
mysql -uroot -p${mysql_password} -e "CREATE USER 'crackmyd'@'%' IDENTIFIED BY 'crackmyd';"
check_command_result "Create user 'crackmyd'@'%'"

# 12. Check user list
print_info "Checking MySQL user list..."
mysql -uroot -p${mysql_password} -e "SELECT host,user,authentication_string FROM mysql.user;"
check_command_result "Check MySQL user list"

# 13. Completion message
print_info "MySQL 5.7 installation completed and set to start on boot."
print_info "New password has been set to:" $mysql_password
