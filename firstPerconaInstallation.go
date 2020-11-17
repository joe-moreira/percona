/*
Description: The script has to be executed on MySQL server for installing Percona
Date: 11/16/2020

Procedure:
1. from login1-eqix-sv5.ooma.internal,  cd /home/hch/mysql,
	scp my.cnf  <user_account>@<destination_server>:/tmp
	scp *.rpm  <user_account>@<destination_server>:/tmp
	scp clustercheck*  <user_account>@<destination_server>:/tmp
	scp  /home/hch/nrpe.service <user_account>@<destination_server>:/tmp

1. On destination_server install RPM:
	As root: cd /tmp
	rpm -ivh --nodeps Percona-Server-shared-compat-57-5.7.21-21.3.el7.x86_64.rpm
	rpm -ivh Percona-Server-shared-57-5.7.21-21.3.el7.x86_64.rpm
	rpm -ivh Percona-Server-client-57-5.7.21-21.3.el7.x86_64.rpm

2. safe to remove mariadb now
	rpm -e mariadb-libs-5.5.60-1.el7_5.x86_64

3. Install server
	rpm -ivh Percona-Server-server-57-5.7.21-21.3.el7.x86_64.rpm

4. create directories
	cd /db
	mkdir /db/mysql
	mkdir /db/mysql/binlogs
	mkdir /db/mysql/etc
	mkdir /db/mysql/logs
	mkdir /db/mysql/tmp
	chown -R mysql.mysql mysql

5. copy my.cnf template for each DC into /etc
	mv /etc/my.cnf /etc/my.cnf.orig  -->> Percona doesn't come with a my.cnf file. It has to be created.
	cp /tmp/my.cnf /etc

	In my.cnf add "serverid" to something uniq using IP address and convert to integer.
	Use the IP address without periods as an unique integer number? e.g. 172.22.13.24 --> 172221324

6. Turn off auto-restarts: vi /etc/systemd/system/mysql.service
	#Restart=on-failure
	#RestartPreventExitStatus=1
	Execute: systemctl daemon-reload

7. Startup MySQL on this server
	systemctl start mysql   (wait 8 seconds)

8. root pwd is defined on "sudo grep "temporary password" /db/mysql/logs/error.log "
	Set password if expired immediately and can't extract from logs file. In /etc/my.cnf , under [mysqld] add "skip-grant-tables"
	restart "systemctl restart  mysql"
	Login: mysql -uroot -p(PWD_from_log_file)

9.
	SET GLOBAL validate_password_policy = LOW;
	UPDATE mysql.user SET authentication_string = PASSWORD('wejk23#1'), password_expired = 'N' WHERE User = 'root' AND Host = 'localhost'; FLUSH PRIVILEGES;
	grant replication client,replication slave on *.* to 'replication'@'%' identified by 'makincopies';
*/

package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func installRPM1() {
	installRPMCommand1 := ("rpm -ivh --nodeps /tmp/Percona-Server-shared-compat-57-5.7.21-21.3.el7.x86_64.rpm")
	_, err1 := exec.Command("bash", "-c", installRPMCommand1).Output()
	if err1 != nil {
		fmt.Println("\nIt looks like Percona-Server-shared-compat-57-5.7.21-21.3.el7.x86_64.rpm is already installed")
	}
	installRPMCommand2 := ("rpm -ivh /tmp/Percona-Server-shared-57-5.7.21-21.3.el7.x86_64.rpm")
	_, err2 := exec.Command("bash", "-c", installRPMCommand2).Output()
	if err2 != nil {
		fmt.Println("\nIt looks like Percona-Server-shared-57-5.7.21-21.3.el7.x86_64.rpm is already installed")
	}
	installRPMCommand3 := ("rpm -ivh /tmp/Percona-Server-client-57-5.7.21-21.3.el7.x86_64.rpm")
	_, err3 := exec.Command("bash", "-c", installRPMCommand3).Output()
	if err3 != nil {
		fmt.Println("\nIt looks like Percona-Server-client-57-5.7.21-21.3.el7.x86_64.rpm is already installed")
	}
}
func removingMariaDB() {
	removeMariaDBCommand := ("rpm -e mariadb-libs-5.5.60-1.el7_5.x86_64")
	exec.Command("bash", "-c", removeMariaDBCommand).Output()
}
func installRPM2() {
	installRPMCommand := ("rpm -ivh --nodeps /tmp/Percona-Server-server-57-5.7.21-21.3.el7.x86_64.rpm")
	_, err := exec.Command("bash", "-c", installRPMCommand).Output()
	if err != nil {
		fmt.Println("\nIt looks like Percona-Server-server-57-5.7.21-21.3.el7.x86_64.rpm is already installed")
	}
	fmt.Printf("\n✅  All the RPMs were installed fine.")
}
func dataDirCreation() {
	if _, err1 := os.Stat("/db/mysql"); os.IsNotExist(err1) {
		os.MkdirAll("/db/mysql", 0777)
		err1 = os.Chmod("/db/mysql", 0777)
		if err1 != nil {
			fmt.Println("\n❌  Couldn't create /db/mysql")
			panic(err1)
		}
	}
	if _, err2 := os.Stat("/db/mysql/binlogs"); os.IsNotExist(err2) {
		os.MkdirAll("/db/mysql/binlogs", 0777)
		err2 = os.Chmod("/db/mysql/binlogs", 0777)
		if err2 != nil {
			fmt.Println("\n❌  Couldn't create /db/mysql/binlogs")
			panic(err2)
		}
	}
	if _, err3 := os.Stat("/db/mysql/etc"); os.IsNotExist(err3) {
		os.MkdirAll("/db/mysql/etc", 0777)
		err3 = os.Chmod("/db/mysql/etc", 0777)
		if err3 != nil {
			fmt.Println("\n❌  Couldn't create /db/mysql/etc")
			panic(err3)
		}
	}
	if _, err4 := os.Stat("/db/mysql/logs"); os.IsNotExist(err4) {
		os.MkdirAll("/db/mysql/logs", 0777)
		err4 = os.Chmod("/db/mysql/logs", 0777)
		if err4 != nil {
			fmt.Println("\n❌  Couldn't create /db/mysql/logs")
			panic(err4)
		}
	}
	if _, err5 := os.Stat("/db/mysql/tmp"); os.IsNotExist(err5) {
		os.MkdirAll("/db/mysql/tmp", 0777)
		err5 = os.Chmod("/db/mysql/tmp", 0777)
		if err5 != nil {
			fmt.Println("\n❌  Couldn't create /db/mysql/tmp")
			panic(err5)
		}
	}
	chownCommand := ("chown -R mysql.mysql /db/mysql")
	_, err6 := exec.Command("bash", "-c", chownCommand).Output()
	if err6 != nil {
		fmt.Println("\n❌  Couldn't chown -R mysql.mysql /db/mysql")
	}
	fmt.Printf("\n✅  Directories created.")
}
func createMyCNFFile() {
	//renaming /etc/my.cnf
	/*	oldName := "/etc/my.cnf"
		newName := "/etc/my.cnf.orig"
		err1 := os.Rename(oldName, newName)
		if err1 != nil {
			fmt.Println("\n❌  Couldn't rename /etc/my.cnf")
			panic(err1)
		} */
	// copying /tmp/my.cnf to /etc/my.cnf
	source, err2 := os.Open("/tmp/my.cnf")
	if err2 != nil {
		fmt.Println("\n❌  Copying /tmp/my.cnf to /etc/my.cnf -> Error: couldn't open /tmp/my.cnf")
		panic(err2)
	}
	defer source.Close()

	destination, err3 := os.Create("/etc/my.cnf")
	if err3 != nil {
		fmt.Println("\n❌  Copying /tmp/my.cnf to /etc/my.cnf -> Error: couldn't open /etc/my.cnf")
		panic(err3)
	}
	defer destination.Close()
	_, err4 := io.Copy(destination, source)
	if err4 != nil {
		fmt.Println("\n❌  Couldn't copy /tmp/my.cnf to /etc/my.cnf")
		panic(err4)
	}
}

// removeIPPeriods - converts the IP address to an unique integer and writes the number in "serverid" variable in my.cnf file
func removeIPPeriods() int {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("\n❌  Error getting server IP: " + err.Error() + "\n")
		os.Exit(1)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				//os.Stdout.WriteString(ipnet.IP.String() + "\n")
				IPnoPeriods := strings.Replace(ipnet.IP.String(), ".", "", -1)
				IPnoPeriodsInt, _ := strconv.Atoi(IPnoPeriods)
				return IPnoPeriodsInt
			}
		}
	}
	return 0
}

// changingMyCNFFile - in my.cnf file changes "server_id" to the IP address converted to an unique integer. Then executing systemctl daemon-reload
func changingMyCNFFile(s int) {
	IPString := strconv.Itoa(s)

	input, err := ioutil.ReadFile("/etc/my.cnf")
	if err != nil {
		fmt.Println("\n❌  Error reading /etc/my.cnf")
		fmt.Println(err)
		os.Exit(1)
	}

	output := bytes.Replace(input, []byte("CHANGEME"), []byte(IPString), -1)
	if err = ioutil.WriteFile("/etc/my.cnf", output, 0666); err != nil {
		fmt.Println("\n❌  \"server_id\" value could not be changed in file /etc/my.cnf")
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("\n✅  The \"server_id\" value has been updated in my.cnf file.")
}

func turnOFFautorestarts() {
	input, err := ioutil.ReadFile("/etc/systemd/system/mysql.service")
	if err != nil {
		fmt.Println("\n❌  Error reading /etc/systemd/system/mysql.service")
		fmt.Println(err)
		os.Exit(1)
	}

	output1 := bytes.Replace(input, []byte("Restart"), []byte("#Restart"), -1)
	if err = ioutil.WriteFile("/etc/systemd/system/mysql.service", output1, 0666); err != nil {
		fmt.Println("\n❌  \"Restart\" could not be commented out in file /etc/systemd/system/mysql.service")
		fmt.Println(err)
		os.Exit(1)
	}
	output2 := bytes.Replace(input, []byte("RestartPreventExitStatus"), []byte("#RestartPreventExitStatus"), -1)
	if err = ioutil.WriteFile("/etc/systemd/system/mysql.service", output2, 0666); err != nil {
		fmt.Println("\n❌  \"RestartPreventExitStatus\" could not be commented out in file /etc/systemd/system/mysql.service")
		fmt.Println(err)
		os.Exit(1)
	}

	// executing systemctl daemon-reload
	if err := exec.Command("systemctl", "daemon-reload").Run(); err != nil {
		fmt.Printf("\n✅  Auto restart has been turned off (/etc/systemd/system/mysql.service).")
	}
}

func startMySQL() {
	startMySQL, err := exec.Command("systemctl", "start", "mysql").Output()
	_ = startMySQL
	if err == nil {
		fmt.Println("\n✅  MySQL is starting...")
		time.Sleep(7 * time.Second)
		restartMySQL, _ := exec.Command("systemctl", "restart", "mysql").Output()
		_ = restartMySQL
		time.Sleep(5 * time.Second)
		fmt.Println("\n✅  MySQL has been started.")
		/*	db, err2 := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/mysql")
			// if there is an error opening the connection, handle it
			if err2 != nil {
				fmt.Println("\n❌  Could not login into MySQL.")
				panic(err2.Error())
			}
			// defer the close till after the main function has finished
			defer db.Close()*/
	} else {
		fmt.Println("\n❌  \"systemctl start mysql\" didn't work.")
		fmt.Println(err)
		os.Exit(1)
	}

}

func appendSkipGrantTables() {
	f, err := os.OpenFile("/etc/my.cnf", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()

	if _, err = f.WriteString("skip-grant-tables"); err != nil {
		fmt.Println("\n❌  Could not append skip-grant-tables in etc/my.cnf")
		panic(err)
	}
}

/*
updatesAndGrants() :
	UPDATE mysql.user SET authentication_string = PASSWORD('wejk23#1'), password_expired = 'N' WHERE User = 'root' AND Host = 'localhost';
	FLUSH PRIVILEGES;
	SET GLOBAL validate_password_policy = LOW;
	grant replication client,replication slave on *.* to 'replication'@'%' identified by 'makincopies';
*/
func updatesAndGrants() {
	// get the temporary pwd defined on /db/mysql/logs/error.log
	grepCommand := ("grep \"temporary password\" /db/mysql/logs/error.log")
	s, err1 := exec.Command("bash", "-c", grepCommand).Output()
	if err1 != nil {
		fmt.Println("\n❌  Couldn't grep the PWD from /db/mysql/logs/error.log")
	}
	result := strings.Split(string(s), " ")
	pwd := strings.TrimSpace(string(result[len(result)-1]))
	fmt.Printf("\n  Temporary password taken from /db/mysql/logs/error.log is: %s", pwd)

	// Connecting to MySQL by unsing the temporary password from /db/mysql/logs/error.log
	db, err2 := sql.Open("mysql", fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/mysql", pwd))
	// if there is an error opening the connection, handle it
	if err2 != nil {
		fmt.Println("\n❌  Couldn't connect to MySQL.")
		panic(err2.Error())
	}
	defer db.Close()

	fmt.Println("\n✅ MySQL is now open for connections. Execute the following below. Then exit and execute secondPerconaInstall file. ")
	fmt.Println("\n✅ mysql -uroot -p   <with the pwd provided above>")
	fmt.Println("\n✅ SET GLOBAL validate_password_policy = LOW;")
	fmt.Println("\n✅ UPDATE mysql.user SET authentication_string = PASSWORD('wejk23#1'), password_expired = 'N' WHERE User = 'root' AND Host = 'localhost';")
	fmt.Println("\n✅ FLUSH PRIVILEGES;")
}

func main() {
	installRPM1()
	removingMariaDB()
	installRPM2()
	dataDirCreation()
	createMyCNFFile()
	serverid := removeIPPeriods()
	changingMyCNFFile(serverid)
	turnOFFautorestarts()
	startMySQL()
	updatesAndGrants()
}
