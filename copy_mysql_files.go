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
	copyMyCNFCommand := ("sshpass -f \"/home/" + uAccount + "/j\" scp /home/hch/mysql/my.cnf " + uAccount + "@" + mysqlS + ":/tmp")
	//fmt.Printf("copyMyCNFCommand: %s", copyMyCNFCommand)
	_, err := exec.Command("bash", "-c", copyMyCNFCommand).Output()

	if err != nil {
		fmt.Println("\n ❌  Couldn't copy my.cnf into MySQL Server /tmp")
		fmt.Println("\n ❌  Please check: ")
		fmt.Println("\n ❌  1) Is sshpass installed? yum install -y sshpass")
		fmt.Println("\n ❌  2) Have you ssh to the server before?")
		panic(err)
	}
}
func copyRMPs(uAccount, mysqlS string) {
	copyRMPsCommand := ("sshpass -f \"/home/" + uAccount + "/j\" scp /home/hch/mysql/*.rpm " + uAccount + "@" + mysqlS + ":/tmp")
	_, err := exec.Command("bash", "-c", copyRMPsCommand).Output()

	if err != nil {
		fmt.Println("\n ❌  Couldn't copy RPM files into MySQL Server /tmp")
		panic(err)
	}
}

func copyClustercheck(uAccount, mysqlS string) {
	copyClustercheckCommand := ("sshpass -f \"/home/" + uAccount + "/j\" scp /home/hch/mysql/clustercheck* " + uAccount + "@" + mysqlS + ":/tmp")
	//	fmt.Println(copyClustercheckCommand)
	_, err := exec.Command("bash", "-c", copyClustercheckCommand).Output()

	if err != nil {
		fmt.Println("\n ❌  Couldn't copy clustercheck files into MySQL Server /tmp")
		panic(err)
	}
}
func copyNrpe(uAccount, mysqlS string) {
	copyNrpeCommand := ("sshpass -f \"/home/" + uAccount + "/j\" scp /home/hch/nrpe.service " + uAccount + "@" + mysqlS + ":/tmp")
	//	fmt.Println(copyNrpeCommand)
	_, err := exec.Command("bash", "-c", copyNrpeCommand).Output()

	if err != nil {
		fmt.Println("\n ❌  Couldn't copy nrpe.service into MySQL Server /tmp")
		panic(err)
	}
}

func createPWDFile(acc string) {
	var u string
	f, err1 := os.Create("/home/" + acc + "/j")
	if err1 != nil {
		fmt.Printf("\n ❌  Could not create \"/home/%s/j\" file! \n", acc)
		panic(err1)
	}
	fmt.Printf("\nEnter your password for ssh: ")
	fmt.Scanln(&u)
	f.WriteString(u)
	err2 := f.Close()
	if err2 != nil {
		panic(err2)
	}
	// check if file exists
	_, err3 := os.Stat("/tmp/j")
	if os.IsNotExist(err3) {
		fmt.Printf("\n ❌  \"/home/%s/j\"  file does not exist! ", acc)
		panic(err3)
	}
}

func main() {
	fmt.Println("\nThis script copies all the files located in /home/hch/mysql directory into MySQL server /tmp directory.")
	// destinationServer - takes the User Account and the MySQL Server as an input
	userAccount, mysqlServer := destinationServer()
	// createPWDFile - creates /tmp/j file with the user account password  (required for SSH). This file will be deleted at the end of the script
	createPWDFile(userAccount)
	// all the following copy the required files from the jumphost to the MySQL server
	copyMyCNF(userAccount, mysqlServer)
	copyRMPs(userAccount, mysqlServer)
	copyClustercheck(userAccount, mysqlServer)
	copyNrpe(userAccount, mysqlServer)

	// deleting the /tmp/j containing the password file
	e := os.Remove("/home/" + userAccount + "/j")
	if e != nil {
		fmt.Println()
		log.Fatal(e)
	}
}
