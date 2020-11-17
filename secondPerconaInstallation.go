/*
Description: The script has to be executed on MySQL server for installing Percona
Date: 11/16/2020

9.updatesAndGrants()
	grant replication client,replication slave on *.* to 'replication'@'%' identified by 'makincopies';
	grant select on *.* to 'clustercheckuser'@'%.ooma.internal' identified by 'clustercheckpassword!';
	grant select on *.* to 'clustercheckuser'@'localhost' identified by 'clustercheckpassword!';
	grant usage,replication client on *.* to 'monitor'@'%.ooma.internal' identified by 'candiisdandi';
	grant usage,replication client on *.* to 'monitor'@'localhost' identified by 'candiisdandi';
	flush privileges;

10.
	cp /tmp/nrpe.service /etc/systemd/system
	systemctl stop nrpe
	systemctl daemon-reload
	systemctl disable nrpe
	systemctl enable nrpe
	systemctl restart nrpe

11. Add Percona monitoring hooks
	cd /tmp
	rpm -ivh Percona-Server-server-57-5.7.21-21.3.el7.x86_64.rpm

12. Run the following commands to create these functions:
	mysql -u root -p -e "CREATE FUNCTION fnv1a_64 RETURNS INTEGER SONAME 'libfnv1a_udf.so'"
	mysql -u root -p -e "CREATE FUNCTION fnv_64 RETURNS INTEGER SONAME 'libfnv_udf.so'"
	mysql -u root -p -e "CREATE FUNCTION murmur_hash RETURNS INTEGER SONAME 'libmurmur_udf.so'"
	(See http://www.percona.com/doc/percona-server/5.7/management/udf_percona_toolkit.html for more details)

13. OS level optimizations.  Disable swap:  sysctl -w vm.swappiness=1

14. Add monitoring hooks  -
	cp /tmp/clustercheck /usr/bin
	cp /tmp/clustercheck.socket /etc/systemd/system
	cp /tmp/clustercheck.service /etc/systemd/system
	NO: 	cp /home/hch/mysql/clustercheck@.service /etc/systemd/system
	NO:  	cd /etc/systemd/system/sockets.target.wants; ln -s /etc/systemd/system/clustercheck.socket
	systemctl daemon-reload
	systemctl enable clustercheck.socket
	systemctl enable clustercheck@.service
	systemctl start clustercheck.socket
	systemctl start clustercheck@.service
*/

package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	_ "github.com/go-sql-driver/mysql"
)

/*
	updatesAndGrants() :
	grant replication client,replication slave on *.* to 'replication'@'%' identified by 'makincopies';
	grant select on *.* to 'clustercheckuser'@'%.ooma.internal' identified by 'clustercheckpassword!';
	grant select on *.* to 'clustercheckuser'@'localhost' identified by 'clustercheckpassword!';
	grant usage,replication client on *.* to 'monitor'@'%.ooma.internal' identified by 'candiisdandi';
	grant usage,replication client on *.* to 'monitor'@'localhost' identified by 'candiisdandi';
	flush privileges;
*/
func updatesAndGrants() {
	db, err2 := sql.Open("mysql", "root:wejk23#1@tcp(127.0.0.1:3306)/mysql")
	// if there is an error opening the connection, handle it
	if err2 != nil {
		fmt.Println("\n❌  Couldn't connect to MySQL.")
		panic(err2.Error())
	}
	defer db.Close()

	// perform GRANT 1
	_, err6 := db.Exec("GRANT replication client,replication slave on *.* to 'replication'@'%' identified by 'makincopies';")
	if err6 != nil {
		fmt.Println("\n❌  Couldn't execute first GRANT;")
		panic(err6.Error())
	}

	// perform GRANT 2
	_, err7 := db.Exec("grant select on *.* to 'clustercheckuser'@'%.ooma.internal' identified by 'clustercheckpassword!';")
	if err7 != nil {
		fmt.Println("\n❌  Couldn't execute second GRANT;")
		panic(err7.Error())
	}

	// perform GRANT 3
	_, err8 := db.Exec("grant select on *.* to 'clustercheckuser'@'localhost' identified by 'clustercheckpassword!';")
	if err8 != nil {
		fmt.Println("\n❌  Couldn't execute third GRANT;")
		panic(err8.Error())
	}

	// perform GRANT 4
	_, err9 := db.Exec("grant usage,replication client on *.* to 'monitor'@'%.ooma.internal' identified by 'candiisdandi';")
	if err9 != nil {
		fmt.Println("\n❌  Couldn't execute fourth GRANT;")
		panic(err9.Error())
	}

	// perform GRANT 5
	_, err10 := db.Exec("grant usage,replication client on *.* to 'monitor'@'localhost' identified by 'candiisdandi';")
	if err10 != nil {
		fmt.Println("\n❌  Couldn't execute fifth GRANT;")
		panic(err10.Error())
	}

	// perform FLUSH PRIVILEGES
	_, err11 := db.Exec("FLUSH PRIVILEGES;")
	if err11 != nil {
		fmt.Println("\n❌  Couldn't execute the last FLUSH PRIVILEGES;")
		panic(err11.Error())
	}

	fmt.Println("\n✅  All UPDATE and GRANT completed OK!")
}

// copyNRPE(): cp /tmp/nrpe.service /etc/systemd/system
func copyNRPE() {
	source, err1 := os.Open("/tmp/nrpe.service")
	if err1 != nil {
		fmt.Println("\n❌  Copying /tmp/nrpe.service to /etc/systemd/ -> Error: couldn't open /tmp/nrpe.service")
		fmt.Println(err1)
		panic(err1)
	}
	defer source.Close()

	destination, err2 := os.Create("/etc/systemd/system/nrpe.service")
	if err2 != nil {
		fmt.Println("\n❌  Copying /tmp/nrpe.service to /etc/systemd/system -> Error: couldn't open  /etc/systemd/system")
		fmt.Println(err2)
		panic(err2)
	}
	defer destination.Close()
	_, err3 := io.Copy(destination, source)
	if err3 != nil {
		fmt.Println("\n❌  Error copying /tmp/nrpe.service to /etc/systemd/system")
		fmt.Println(err3)
		panic(err3)
	}

}

/*
restartEnableNRPE():
	systemctl stop nrpe
	systemctl daemon-reload
	systemctl disable nrpe
	systemctl enable nrpe
	systemctl restart nrpe
*/
func restartEnableNRPE() {
	if err1 := exec.Command("systemctl", "stop", "nrpe").Run(); err1 != nil {
		fmt.Println("\n❌  Couldn't stop NRPE.")
		fmt.Println(err1)
		panic(err1)
	}
	if err2 := exec.Command("systemctl", "daemon-reload").Run(); err2 != nil {
		fmt.Println("\n❌  Couldn't run daemon-reload.")
		fmt.Println(err2)
		panic(err2)
	}
	if err3 := exec.Command("systemctl", "disable", "nrpe").Run(); err3 != nil {
		fmt.Println("\n❌  Couldn't disable NRPE.")
		fmt.Println(err3)
		panic(err3)
	}
	if err4 := exec.Command("systemctl", "enable", "nrpe").Run(); err4 != nil {
		fmt.Println("\n❌  Couldn't enable NRPE.")
		fmt.Println(err4)
		panic(err4)
	}
	if err5 := exec.Command("systemctl", "restart", "nrpe").Run(); err5 != nil {
		fmt.Println("\n❌  Couldn't restart NRPE.")
		fmt.Println(err5)
		panic(err5)
	}
	fmt.Println("\n✅  NRPE enabled OK!")
}

func addPerconaMonitoring() {
	installRPMCommand := ("rpm -ivh --nodeps /tmp/Percona-Server-server-57-5.7.21-21.3.el7.x86_64.rpm")
	_, err := exec.Command("bash", "-c", installRPMCommand).Output()
	if err != nil {
		fmt.Println("\nIt looks like Percona-Server-server-57-5.7.21-21.3.el7.x86_64.rpm is already installed")
	}
}

// See http://www.percona.com/doc/percona-server/5.7/management/udf_percona_toolkit.html for more details
func functionsCreation() {
	db, err1 := sql.Open("mysql", "root:wejk23#1@tcp(127.0.0.1:3306)/mysql")
	if err1 != nil {
		fmt.Println("\n❌  Couldn't connect to MySQL for Functions creation.")
		panic(err1.Error())
	}
	defer db.Close()

	create1, err2 := db.Query("CREATE FUNCTION fnv1a_64 RETURNS INTEGER SONAME 'libfnv1a_udf.so';")
	if err2 != nil {
		fmt.Println("\n❌  Couldn't create the first function.")
		panic(err2.Error())
	}
	defer create1.Close()

	create2, err3 := db.Query("CREATE FUNCTION fnv_64 RETURNS INTEGER SONAME 'libfnv_udf.so'")
	if err3 != nil {
		fmt.Println("\n❌  Couldn't create the second function.")
		panic(err3.Error())
	}
	defer create2.Close()

	create3, err4 := db.Query("CREATE FUNCTION murmur_hash RETURNS INTEGER SONAME 'libmurmur_udf.so'")
	if err4 != nil {
		fmt.Println("\n❌  Couldn't create the third function.")
		panic(err4.Error())
	}
	defer create3.Close()
	fmt.Println("\n✅  Percona Functions creation OK!")
}

func disableSwap() {
	if err := exec.Command("sysctl", "-w", "vm.swappiness=1").Run(); err != nil {
		fmt.Println("\n❌  Couldn't disable SWAP.")
		fmt.Println(err)
		panic(err)
	}
}

// clustercheckCopyFiles(): copying files from /tmp into /user/bin/  and  /etc/systemd/system/
func clustercheckCopyFiles() {

	// renaming /tmp/clustercheck@.service to /tmp/clustercheck.service
	OriginalPath := "/tmp/clustercheck@.service"
	NewPath := "/tmp/clustercheck.service"
	e := os.Rename(OriginalPath, NewPath)
	if e != nil {
		fmt.Println("\n❌  Couldn't rename /tmp/clustercheck@.service to /tmp/clustercheck.service")
		log.Fatal(e)
	}

	//Changing Accept=true to false on  /home/hch/mysql/clustercheck.socket
	input, err := ioutil.ReadFile("/tmp/clustercheck.socket")
	if err != nil {
		fmt.Println("\n❌  Error reading /tmp/clustercheck.socket")
		fmt.Println(err)
		os.Exit(1)
	}

	output := bytes.Replace(input, []byte("true"), []byte("false"), -1)
	if err = ioutil.WriteFile("/tmp/clustercheck.socket", output, 0666); err != nil {
		fmt.Println("\n❌  \"Accept=true\" value couldn't be changed in file /tmp/clustercheck.socket")
		fmt.Println(err)
		os.Exit(1)
	}

	// Copying /tmp/clustercheck to /etc/clustercheck
	source1, err1 := os.Open("/tmp/clustercheck")
	if err1 != nil {
		fmt.Println("\n❌  Copying /tmp/clustercheck into /usr/bin/-> Error: couldn't open /tmp/clustercheck")
		fmt.Println(err1)
		panic(err1)
	}
	defer source1.Close()

	destination1, err2 := os.Create("/usr/bin/clustercheck")
	if err2 != nil {
		fmt.Println("\n❌  Copying /tmp/clustercheck into /usr/bin/-> Error: couldn't create /etc/clustercheck")
		fmt.Println(err2)
		panic(err2)
	}
	defer destination1.Close()
	_, err3 := io.Copy(destination1, source1)
	if err3 != nil {
		fmt.Println("\n❌  Copying /tmp/clustercheck to /etc/clustercheck")
		fmt.Println(err3)
		panic(err3)
	}

	//Copying /tmp/clustercheck.socket /etc/systemd/system
	source2, err4 := os.Open("/tmp/clustercheck.socket")
	if err4 != nil {
		fmt.Println("\n❌  Copying /tmp/clustercheck.socket into /etc/systemd/system-> Error: couldn't open /tmp/clustercheck.socket")
		fmt.Println(err4)
		panic(err4)
	}
	defer source2.Close()

	destination2, err5 := os.Create("/etc/systemd/system/clustercheck.socket")
	if err5 != nil {
		fmt.Println("\n❌  Copying /tmp/clustercheck.socket into /etc/systemd/system")
		fmt.Println(err5)
		panic(err5)
	}
	defer destination2.Close()
	_, err6 := io.Copy(destination2, source2)
	if err6 != nil {
		fmt.Println("\n❌  Copying /tmp/clustercheck to /etc/clustercheck")
		fmt.Println(err6)
		panic(err6)
	}
	//Copying /tmp/clustercheck.service /etc/systemd/system
	source3, err7 := os.Open("/tmp/clustercheck.service")
	if err7 != nil {
		fmt.Println("\n❌  Copying /tmp/clustercheck.service into /etc/systemd/system-> Error: couldn't open /tmp/clustercheck.service")
		fmt.Println(err7)
		panic(err7)
	}
	defer source3.Close()

	destination3, err8 := os.Create("/etc/systemd/system/clustercheck.service")
	if err8 != nil {
		fmt.Println("\n❌  Copying /tmp/clustercheck.socket into /etc/systemd/system")
		fmt.Println(err8)
		panic(err8)
	}
	defer destination3.Close()
	_, err9 := io.Copy(destination3, source3)
	if err9 != nil {
		fmt.Println("\n❌  Copying /tmp/clustercheck.service /etc/systemd/system")
		fmt.Println(err9)
		panic(err9)
	}

}

// clustercheckServiceEnable(): enablig clustercheck service
func clustercheckServiceEnable() {
	// executing systemctl daemon-reload
	if err := exec.Command("systemctl", "daemon-reload").Run(); err != nil {
		fmt.Printf("\n❌  Couldn't run \"systemctl daemon-reload\"")
	}
	if err1 := exec.Command("systemctl", "enable", "clustercheck.socket").Run(); err1 != nil {
		fmt.Println("\n❌  Couldn't enable clustercheck.socket.")
		fmt.Println(err1)
		panic(err1)
	}
	if err2 := exec.Command("systemctl", "enable", "clustercheck.service").Run(); err2 != nil {
		fmt.Println("\n❌  Couldn't enable clustercheck.service.")
		fmt.Println(err2)
		panic(err2)
	}
	if err3 := exec.Command("systemctl", "start", "clustercheck.socket").Run(); err3 != nil {
		fmt.Println("\n❌  Couldn't start clustercheck.socket.")
		fmt.Println(err3)
		panic(err3)
	}
	if err4 := exec.Command("systemctl", "start", "clustercheck.service").Run(); err4 != nil {
		fmt.Println("\n❌  Couldn't start clustercheck.service.")
		fmt.Println(err4)
		panic(err4)
	}
	fmt.Println("\n✅  Clustercheck service is running OK.")
}

func main() {
	updatesAndGrants()
	copyNRPE()
	restartEnableNRPE()
	addPerconaMonitoring()
	functionsCreation()
	disableSwap()
	clustercheckCopyFiles()
	clustercheckServiceEnable()
}
