package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	_ "time"
	"unicode"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/changepwd/{user}/{oldpwd}/{newpwd}", ChpwdHandler)
	router.HandleFunc("/checkpwd/{user}", checkpwd)
	router.HandleFunc("/checkpwdexp/{user}", checkpwdage)

	log.Fatal(http.ListenAndServe(":9877", router))
}

func ChpwdHandler(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	user := vars["user"]
	oldpwd := vars["oldpwd"]
	newpwd := vars["newpwd"]
	log.Println(user, oldpwd, newpwd)
	er := validPassword(newpwd)
	if er != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "P"+er.Error())
		return
	}
	obuf, e := CallPS(user, oldpwd, newpwd, w)
	if e != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, e.Error())
		return
	} else {
		outMsg := obuf.String()
		w.WriteHeader(200)
		fmt.Fprint(w, outMsg)
		CallPS(user, oldpwd, newpwd, w)

	}
	//As it allows two times...it is a work around

}
func CallPS(user string, oldpwd string, newpwd string, w http.ResponseWriter) (bytes.Buffer, error) {
	var obuf bytes.Buffer
	var oerr bytes.Buffer
	//Below path is in Khind AD
	//cwd, _ := os.Getwd()
	_, err := os.ReadFile("Powershell/change_pwd.ps1")
	if err != nil {
		fmt.Println("Error reading script1.ps1:", err)

	}

	cmd := exec.Command("Powershell", "-file", "Powershell/change_pwd.ps1", user, oldpwd, newpwd)
	cmd.Stdout = &obuf
	cmd.Stderr = &oerr
	err = cmd.Run()
	defer cmd.Process.Kill()
	if err != nil {
		cmd.Process.Kill()
		return bytes.Buffer{}, err
	} else {
		cmd.Process.Kill()
		fmt.Println("User:", user)
		fmt.Println("Old_Pswd:", oldpwd)
		fmt.Println("New_Pswd:", newpwd)
		return obuf, nil

	}

}
func validPassword(s string) error {
next:
	for name, classes := range map[string][]*unicode.RangeTable{
		"Upper case":        {unicode.Upper, unicode.Title},
		"Lower case":        {unicode.Lower},
		"Numeric":           {unicode.Number, unicode.Digit},
		"Special character": {unicode.Space, unicode.Symbol, unicode.Punct},
	} {
		for _, r := range s {
			if unicode.IsOneOf(classes, r) {
				continue next
			}
		}
		return fmt.Errorf("assword must contain atleast one  %s ", name)
	}
	return nil
}
func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	// (w).Header().Set("Access-Control-Allow-Origin", "https://")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}
func checkpwd(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	user := vars["user"]
	var obuf bytes.Buffer
	var oerr bytes.Buffer
	_, err := os.ReadFile("Powershell/checkpwd.ps1")
	if err != nil {
		fmt.Println("Error reading script1.ps1:", err)

	}

	cmd := exec.Command("Powershell", "-file", "Powershell/checkpwd.ps1", user)
	cmd.Stdout = &obuf
	cmd.Stderr = &oerr
	err = cmd.Run()
	defer cmd.Process.Kill()
	if err != nil {
		cmd.Process.Kill()
		fmt.Fprint(w, err)
	} else {
		cmd.Process.Kill()
		outMsg := obuf.String()
		w.WriteHeader(200)
		outMsg = strings.TrimSpace(outMsg)
		fmt.Println(outMsg)
		fmt.Fprint(w, outMsg)
	}

}
func checkpwdage(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	vars := mux.Vars(r)
	user := vars["user"]

	var obuf bytes.Buffer
	var oerr bytes.Buffer

	//cwd, _ := os.Getwd()
	_, err := os.ReadFile("Powershell/pwdexpdate.ps1")
	if err != nil {
		fmt.Println("Error reading script1.ps1:", err)

	}
	cmd := exec.Command("Powershell", "-file", "Powershell/pwdexpdate.ps1", user)
	cmd.Stdout = &obuf
	cmd.Stderr = &oerr
	err = cmd.Run()
	defer cmd.Process.Kill()
	if err != nil {
		cmd.Process.Kill()
		fmt.Println(err)
		fmt.Fprint(w, err)
	} else {
		cmd.Process.Kill()
		outMsg := obuf.String()
		w.WriteHeader(200)
		outMsg = strings.TrimSpace(outMsg)
		expdate := strings.Split(outMsg, "----------")
		//fmt.Println(expdate)
		warn := checkprompt()
		if len(expdate) == 2 {
			if len(strings.TrimSpace(expdate[1])) > 0 {
				// expdate, err := time.Parse("1/2/2006 3:04:05 PM", strings.TrimSpace(expdate[1]))
				// if err != nil {
				// 	fmt.Println(err)
				// }
				//warning := expdate.Sub(time.Now())

				//days := warning.Hours() / 24

				fmt.Println("Expiry=", expdate)

				fmt.Fprintf(w, strings.TrimSpace(expdate[1])+"|"+strings.TrimSpace(warn))
				//fmt.Println(user + "=" + strings.TrimSpace(fmt.Sprintf("%d", expdate))+"|"+strings.TrimSpace(fmt.Sprintf("%d", warning)))
				fmt.Println("Warning=" + strings.TrimSpace(warn))
				fmt.Println("-----------------------------------")
			} else {
				fmt.Fprintf(w, "0"+"|"+strings.TrimSpace(warn))
				fmt.Println(user + "=" + "0")
				fmt.Println("Warning=" + strings.TrimSpace(warn))
				fmt.Println("-----------------------------------")
			}

		} else {
			fmt.Fprintf(w, "0"+"|"+strings.TrimSpace(warn))
			fmt.Println(user + "=" + "0")
			fmt.Println("Warning=" + strings.TrimSpace(warn))
			fmt.Println("-----------------------------------")

		}

	}

}
func checkprompt() string {
	var obuf bytes.Buffer
	var oerr bytes.Buffer
	var outMsg string

	//cwd, _ := os.Getwd()
	_, err := os.ReadFile("Powershell/Pswd_interactivelogon_prompt.ps1")
	if err != nil {
		fmt.Println("Error reading script1.ps1:", err)

	}

	cmd := exec.Command("Powershell", "-file", "Powershell/Pswd_interactivelogon_prompt.ps1")

	cmd.Stdout = &obuf
	cmd.Stderr = &oerr
	err = cmd.Run()
	defer cmd.Process.Kill()
	if err != nil {
		cmd.Process.Kill()
		fmt.Println(err)
	} else {
		cmd.Process.Kill()
		outMsg = obuf.String()
	}
	return strings.TrimSpace(outMsg)
}
func createTempFile(content []byte) string {
	tempFile, err := ioutil.TempFile("", "script_*.ps1")
	if err != nil {
		fmt.Println("Error creating temp file:", err)
		return ""
	}
	defer tempFile.Close()

	tempFile.Write(content)
	return tempFile.Name()
}
