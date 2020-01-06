# netlog
logging package

## How to use?

1. Determine the log output method.
* output to any log file
    * ex. `/var/log/example.log`
    ```
    var logOutput string
    logOutput = "file:///var/log/example.log"
    ```
* output via syslog
    * ex. facility: sys, tag: example
    ```
    var logOutput string
    logOutput = "net:///?facility=sys&tag=example"
    ```
2. Initialize using log output destination.
    ```
    import (
    	"log"
    )
    
    func main() {
    	if logOutput != "" {
    		if err := netlog.SetOutputURL(logOutput); err != nil {
    			log.Fatal("log output destination:", err)
    		}
    	}
    }
    ```
3. Output logs.
    ```
    import (
    	"github.com/lufia/netlog"
    )
    
    func example() {
        // netlog.SetOutputURL(logOutput, debugMode)
        netlog.Debug("output if the debug flag is enabled at initialization")

        netlog.Info("this is info message")
        netlog.Warning("this is warning message")
        netlog.Err("this is error message")
        netlog.Crit("this is critical error message. process AbEnd here")
    }
    ```

## Log level
| func | Level (linux/unix) | Msg ID (windows) |
| --- | --- | --- |
| netlog.Crit() | LOG_CRIT | 4001 |
| netlog.Err() | LOG_ERR | 3001 |
| netlog.Warning() | LOG_WARNING | 2001 |
| netlog.Info() | LOG_INFO |1001 |
| netlog.Debug() | LOG_DEBUG | 1001 |
