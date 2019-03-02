# h2memcache
In-memory cache server/client with HTTP/2.0 REST API written in Go

## Server

### Environmental Variables
- **CACHE_SIZE** Maximum cache size in MB
- **HOSTNAME** Hostname for Let's Encrypt certs
- **PORT** Port (Default: 8080)
- **TLS_CERT_FILE** TLS Certificate file (only if you need no Let's Encrypt)
- **TLS_KEY_FILE** TLS Key file for (only if you need no Let's Encrypt)
- **TLS_CERT_DIR** Let's Encrypt cache directory
- **API_KEY** API Key
- **GOGC** Garbage collection target percentage (Recommend value: 20)

## Client

### Usage

```
import (
	"github.com/akosmarton/h2memcache"
)

func main() {
    url := "http://localhost:8080"
    apikey := "veryverysecret"
	cache := h2memcache.NewCache(&http.Client{}, url, apikey)

    if err := cache.Set("key", []byte("value"), 0); err != nil {
        log.Fatal(err)
    }

    if value, err := cache.Get("key"); err == nil {
        log.Println(string(value))
    } else {
        log.Fatal(err)
    }

    if err := cache.Delete("key"); err != nil {
        log.Fatal(err)
    }
}
```
