package wireless_sensor_tags

import (
    "bufio"
    "fmt"
    "github.com/JulianSauer/WaterCan/wireless_sensor_tags/api/client"
    "golang.org/x/crypto/ssh/terminal"
    "os"
    "strings"
    "syscall"
)

var username *string
var password *string

func Login() error {
    if isSignedIn, e := client.IsSignedIn();
        e != nil || isSignedIn {
        return e
    }

    username, password, e := getUsernameAndPassword()
    if e != nil {
        return e
    }

    return client.SignIn(username, password)
}

func getUsernameAndPassword() (string, string, error) {
    if username != nil && password != nil {
        return *username, *password, nil
    } else {
        u, p, e := loginPromt()
        username = &u
        password = &p
        return u, p, e
    }
}

func loginPromt() (string, string, error) {
    reader := bufio.NewReader(os.Stdin)

    fmt.Println("Please enter your Wireless Sensor Tags credentials")
    fmt.Print("Username: ")
    username, e := reader.ReadString('\n')
    if e != nil {
        return "", "", e
    }

    fmt.Print("Password: ")
    password, e := terminal.ReadPassword(int(syscall.Stdin))
    if e != nil {
        return "", "", e
    }
    fmt.Println()

    return strings.TrimSpace(username), strings.TrimSpace(string(password)), nil
}
