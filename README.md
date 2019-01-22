# AnyPort [![GoDoc](https://godoc.org/github.com/thekondor/anyport?status.svg)](http://godoc.org/github.com/thekondor/anyport/)

`anyport` is a `Go` package to provide with a tool to bind any available or random port to listen for incoming TCP connections. The use case does make sense when a started (micro)service should not be pinned to a fixed TCP port but could be found through Service Discovery (like Consul) endpoints.

# Usage

## Random TCP Port (any range)

  Plain TCP:
  ```go
    anyPort, err := anyport.ListenInsecure("localhost")
    if nil != err {
        panic(err)
    }
    defer anyPort.Listener.Close()
    log.Printf("Incoming port: %d", anyPort.PortNumber)
  ```

  Over TLS:
  ```go
    tlsConfig := tls.Config{...}
    anyPort, err := anyport.ListenSecure("localhost", &tlsConfig)
    if nil != err {
        panic(err)
    }
    defer anyPort.Listener.Close()
    log.Printf("Incoming TLS port: %d", anyPort.PortNumber)
  ```

## Random TCP port (in a fixed range)

  Plain TCP:
  ```go
    anyPort, err := anyport.ListenInsecure("localhost:3000-5000")
    if nil != err {
        panic(err)
    }
    defer anyPort.Listener.Close()
    log.Printf("Incoming port: %d", anyPort.PortNumber)
  ```

  Over TLS:
  ```go
    tlsConfig := tls.Config{...}
    anyPort, err := anyport.ListenSecure("localhost:3000-5000", &tlsConfig)
    if nil != err {
        panic(err)
    }
    defer anyPort.Listener.Close()
    log.Printf("Incoming TLS port: %d", anyPort.PortNumber)
  ```
  
# License


The library is released under the MIT license. See LICENSE file.

