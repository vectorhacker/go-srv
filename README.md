SRV
===

This is a simple naming resolver that helps resolve SRV DNS records. This is
useful for when your running service discovery that's exposed through DNS.

Example
-------

Using standard consul SRV DNS

```golang
package main

// exposes the 'srv' package namespace
import "github.com/vectorhacker/go-srv"

func main() {
  changes := srv.Resolver.Resolve("hello.service.consul")

  // Blocks until you get the next set of changes or an error
  addrs, err := changes.Next()

  if err != nil {
    panic(err)
  }

  con, err := net.Dial("tcp", addrs[0].Addr)
  if err != nil {
    // handle error
  }
  fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
  status, err := bufio.NewReader(conn).ReadString('\n')

  // ...

  // Closes the watcher
  changes.Close()
}
```