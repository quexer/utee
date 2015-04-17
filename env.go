package utee
import (
    "os"
    "log"
)

func Env(name string, required bool, showLog ...bool) string {
    s := os.Getenv(name)
    if required && s == "" {
        panic("env variable required, " + name)
    }
    if len(showLog) == 0 || showLog[0] {
        log.Printf("[env][%s] %s\n", name, s)
    }
    return s
}