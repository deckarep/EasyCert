package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
)

var (
	flagCertificateAuthorityName = flag.String("cn", "", "")
	flagHostName                 = flag.String("h", "", "")
)

var usage = `Usage: EasyCert [options...]

Options:
  -cn Certificate Authority Name (can be any name, but should reflect your company name.)
  -h  Hostname of TLS server to install the private cert/key
`

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage)
	}

	flag.Parse()

	certName := *flagCertificateAuthorityName
	hostName := *flagHostName

	if certName == "" || hostName == "" {
		usageAndExit("You must supply both a -cn (certificate name) and -h (host name) parameter")
	}

	createPrivateCA(certName)
	createServerCertKey(hostName)

	log.Println("*** Operation Completed Succesfully ***")
	log.Println("Private root certificate created: ", "myCA.cer")
	log.Println("Web server certificate created: ", "mycert1.cer")
	log.Println("Web server key created: ", "mycer1.key")
}

func createPrivateCA(certificateAuthorityName string) {
	_, err := callCommand("openssl", "genrsa", "-out", "myCA.key", "2048")
	if err != nil {
		log.Fatal("Could not create private Certificate Authority key")
	}

	_, err = callCommand("openssl", "req", "-x509", "-new", "-key", "myCA.key", "-out", "myCA.cer", "-days", "730", "-subj", "/CN=\""+certificateAuthorityName+"\"")
	if err != nil {
		log.Fatal("Could not create private Certificate Authority certificate")
	}
}

func createServerCertKey(host string) {
	_, err := callCommand("openssl", "genrsa", "-out", "mycert1.key", "2048")
	if err != nil {
		log.Fatal("Could not create private server key")
	}

	_, err = callCommand("openssl", "req", "-new", "-out", "mycert1.req", "-key", "mycert1.key", "-subj", "/CN="+host)
	if err != nil {
		log.Fatal("Could not create private server certificate signing request")
	}

	_, err = callCommand("openssl", "x509", "-req", "-in", "mycert1.req", "-out", "mycert1.cer", "-CAkey", "myCA.key", "-CA", "myCA.cer", "-days", "365", "-CAcreateserial", "-CAserial", "serial")
	if err != nil {
		log.Fatal("Could not create private server certificate")
	}

}

func callCommand(command string, arg ...string) (string, error) {
	out, err := exec.Command(command, arg...).Output()

	if err != nil {
		log.Println("callCommand failed!")
		log.Println("")
		log.Println(string(debug.Stack()))
		return "", err
	}
	return string(out), nil
}

func usageAndExit(message string) {
	if message != "" {
		fmt.Fprintf(os.Stderr, message)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}
