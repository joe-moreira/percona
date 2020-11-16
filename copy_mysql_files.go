package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func destinationServer() (u string, s string) {
	fmt.Printf("\nEnter your user account: ")
	fmt.Scanln(&u)
	fmt.Printf("\nEnter MySQL server name: ")
	fmt.Scanln(&s)
	return u, s
}

func copyMyCNF(uAccount, mysqlS string) {
	copyMyCNFCommand := ("sshpass -f \"/tmp/j\" scp /home/hch/mysql/my.cnf " + uAccount + "@" + mysqlS + ":/tmp")
	//	fmt.Println(copyMyCNFCommand)
	_, err := exec.Command("bash", "-c", copyMyCNFCommand).Output()

	if err != nil {
		fmt.Println("\n ❌  Couldn't copy my.cnf into MySQL Server /tmp")
	}
}
func copyRMPs(uAccount, mysqlS string) {
	copyRMPsCommand := ("sshpass -f \"/tmp/j\" scp /home/hch/mysql/*.rpm " + uAccount + "@" + mysqlS + ":/tmp")
	//	fmt.Println(copyRMPsCommand)
	_, err := exec.Command("bash", "-c", copyRMPsCommand).Output()

	if err != nil {
		fmt.Println("\n ❌  Couldn't copy RPM files into MySQL Server /tmp")
	}
}

func copyClustercheck(uAccount, mysqlS string) {
	copyClustercheckCommand := ("sshpass -f \"/tmp/j\" scp /home/hch/mysql/clustercheck* " + uAccount + "@" + mysqlS + ":/tmp")
	//	fmt.Println(copyClustercheckCommand)
	_, err := exec.Command("bash", "-c", copyClustercheckCommand).Output()

	if err != nil {
		fmt.Println("\n ❌  Couldn't copy clustercheck files into MySQL Server /tmp")
	}
}
func copyNrpe(uAccount, mysqlS string) {
	copyNrpeCommand := ("sshpass -f \"/tmp/j\" scp /home/hch/nrpe.service " + uAccount + "@" + mysqlS + ":/tmp")
	//	fmt.Println(copyNrpeCommand)
	_, err := exec.Command("bash", "-c", copyNrpeCommand).Output()

	if err != nil {
		fmt.Println("\n ❌  Couldn't copy nrpe.service into MySQL Server /tmp")
	}
}

func createPWDFile() {
	var u string
	f, err := os.Create("/tmp/j")
	if err != nil {
		fmt.Println("\n ❌  Could not create \"/tmp/j\" file! ")
		return
	}
	fmt.Printf("\nEnter your password for ssh: ")
	fmt.Scanln(&u)
	f.WriteString(u)
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	// check if file exists
	_, err = os.Stat("/tmp/j")
	if os.IsNotExist(err) {
		fmt.Println("\n ❌  \"/tmp/j\" file does not exist! ")
		return
	}
}

func main() {
	fmt.Println("\nThis script copies all the files located in /home/hch/mysql directory into MySQL server /tmp directory.")
	// destinationServer - takes the User Account and the MySQL Server as an input
	userAccount, mysqlServer := destinationServer()
	// createPWDFile - creates /tmp/j file with the user account password  (required for SSH). This file will be deleted at the end of the script
	createPWDFile()
	// all the following copy the required files from the jumphost to the MySQL server
	copyMyCNF(userAccount, mysqlServer)
	copyRMPs(userAccount, mysqlServer)
	copyClustercheck(userAccount, mysqlServer)
	copyNrpe(userAccount, mysqlServer)

	// deleting the /tmp/j containing the password file
	e := os.Remove("/tmp/j")
	if e != nil {
		log.Fatal(e)
	}
}