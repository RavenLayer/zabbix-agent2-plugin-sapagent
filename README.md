# Zabbix Agent 2 plugin for SAP® Agent
This plugin provides a native Zabbix solution to monitor SAP® systems through the SAP Host Agent.
It communicates with the SAP Host Agent via SOAP (HTTP/HTTPS or Unix sockets) to collect various performance metrics and host information.
The plugin supports multiple SAP instances simultaneously, either remote or local to the Zabbix Agent.

## Requirements
* Zabbix Agent 2 version 6.0.0 or newer
* SAP Host Agent (standard SOAP interface 7.20+)
* If building from source: Go version v1.22 or newer

## Plugin setup
`Plugins.SAPAGENT.System.Path` variable needs to be set in the Zabbix agent 2 configuration file with the path to the SAP Agent plugin executable. By default the variable is set in the **plugin** configuration file `sapagent.conf` and then included in the **agent** configuration file `zabbix_agent2.conf`.
You should add the following option to the **plugin** configuration file:

*For example*:

    Plugins.SAPAGENT.System.Path=/path/to/executable/sapagent

Then, the configuration file needs to be included in the main Zabbix agent 2 configuration file via the
`Include` directive. You should add the following option to the **agent** configuration file:

*For example*:

    Include=/path/to/config/sapagent.conf

## Options
SAP Agent plugin can be executed on its own with these parameters:
* `-h`, `--help` displays help message;
* `-V`, `--version` displays the plugin version and license information.

## Installation
The plugin requires credentials to authenticate with the SAP Host Agent. Usually, the `sapadm` user is used.
For HTTPS connections, ensure the appropriate TLS certificates are configured if validation is required.

## Configuration
To configure plugins, use Zabbix agent configuration file.

`Plugins.SAPAGENT.Timeout` — the maximum amount of time to wait for a server to respond to SOAP requests.  
*Default value:* equals the global Timeout configuration parameter defined in Zabbix agent 2 configuration file.  
*Limits:* 1-30

`Plugins.SAPAGENT.CacheTTL` — the maximum time in seconds for keeping cached metric values in memory.  
*Default value:* 60 sec.

`Plugins.SAPAGENT.Sessions.<session_name>.Uri` — Uri to connect.  
*Default value:* http://localhost:1128  
Accepted schemes: http, https, unix (e.g., unix:///tmp/.sapstream51128)

`Plugins.SAPAGENT.Sessions.<session_name>.Username` — Username for SAP Host Agent authentication.  
*Default value:* sapadm

`Plugins.SAPAGENT.Sessions.<session_name>.Password` — Password for SAP Host Agent authentication.  
*Default value:* empty

`Plugins.SAPAGENT.Sessions.<session_name>.TLSRootCAFile` — full pathname of a file containing the top-level CA(s) certificates for SAP Agent.  
*Default value:* empty

`Plugins.SAPAGENT.Sessions.<session_name>.TLSServerName` — TLS ServerName to verify the hostname.  
*Default value:* empty

`Plugins.SAPAGENT.Sessions.<session_name>.TLSInsecureSkipVerify` — Controls whether a client verifies the server's certificate chain and host name.  
*Default value:* false  
*Accepted values:* false, true

### Configuring connections
A connection can be configured using either keys' parameters or named sessions.

If you leave any connection parameter empty, a default hard-coded value of `http://localhost:1128` with user `sapadm` will be used.

#### Using keys' parameters
The common parameters (later used as `<common_params>`) for all keys are: `[ConnString][,User][,Password]` where *ConnString* can be either a URI or session name.   
*ConnString* will be treated as a URI if no session with the given name are found.  
If you use *ConnString* as a session name, just skip the rest of the connection parameters.  

*Examples*:

* Using a standard connection URI: `sap.agent.ping[http://192.168.1.10:1128]`
* Using a named session (configured in `sapagent.conf`): `sap.agent.mem[Prod]`
* Using a Unix socket: `sap.agent.version[unix:///tmp/.sapstream51128]`

#### Using named sessions
Named sessions allow you to define specific parameters for each SAP instance.

For example, if you have two SAP instances: "Prod" and "Test",
you should add the following options to the plugin configuration file:

    Plugins.SAPAGENT.Sessions.Prod.Uri=https://sap-prod:1129
    Plugins.SAPAGENT.Sessions.Prod.Username=sapadm
    Plugins.SAPAGENT.Sessions.Prod.Password=prod_password
    Plugins.SAPAGENT.Sessions.Prod.TLSInsecureSkipVerify=true

    Plugins.SAPAGENT.Sessions.Test.Uri=http://sap-test:1128
    Plugins.SAPAGENT.Sessions.Test.Username=sapadm
    Plugins.SAPAGENT.Sessions.Test.Password=test_password

Then, you will be able to use these names as the first parameter (ConnString) in keys instead of URIs.

*For example*:

    sap.agent.ping[Prod]
    sap.agent.cpu.util[Test]

## Supported keys

`sap.agent.ping[<common_params>]` — tests if a connection to SAP Host Agent is alive.
*Returns:*
- `1` if a connection is alive.
- `0` if a connection is broken.

`sap.agent.version[<common_params>]` — returns SAP Host Agent version information.

`sap.agent.cpu.discovery[<common_params>]` — returns JSON discovery data of logical CPUs.

`sap.agent.cpu.util[<common_params>[,cpu],<type>]` — returns CPU utilization percentage.  
*Parameters:*
- `cpu`: CPU number or `all` (default).
- `type`: utilization type (`user`, `system`, `idle`, `iowait`, `steal` ). `steal` time is only available for aggregate (`all`).

`sap.agent.cpu.load[<common_params>[,mode]]` — returns CPU load average.  
*Parameters:*
- `mode`: load average mode (`avg1` (default), `avg5`, `avg15`).

`sap.agent.cpu.num[<common_params>]` — returns total number of logical CPUs.

`sap.agent.cpu.intr[<common_params>]` — returns interrupt rate per second.

`sap.agent.cpu.switches[<common_params>]` — returns context switch rate per second.

`sap.agent.mem[<common_params>,<type>,<mode>]` — returns host memory information.  
*Parameters:*
- `type`: memory type (`phy`, `virt`).
- `mode`: memory mode (`available`, `free`, `size`). `virt` memory type only supports `size` mode.

`sap.agent.swap.size[<common_params>,<mode>]` — returns swap size information.  
*Parameters:*
- `mode`: swap mode (`free`, `used`, `total`).

`sap.agent.swap.in[<common_params>,<mode>]` — returns swap-in rate information.  
*Parameters:*
- `mode`: swap mode (`bytes`, `pages`).

`sap.agent.swap.out[<common_params>,<mode>]` — returns swap-out rate information.  
*Parameters:*
- `mode`: swap mode (`bytes`, `pages`).

`sap.agent.net.discovery[<common_params>]` — returns JSON discovery data of network interfaces.

`sap.agent.net.in[<common_params>,<iface>,<mode>]` — returns network incoming traffic information.  
*Parameters:*
- `iface`: network interface name.
- `mode`: traffic mode (`bytes`, `packets`, `errors`).

`sap.agent.net.out[<common_params>,<iface>,<mode>]` — returns network outgoing traffic information.  
*Parameters:*
- `iface`: network interface name.
- `mode`: traffic mode (`bytes`, `packets`, `errors`).

`sap.agent.fs.discovery[<common_params>]` — returns JSON discovery data of mounted filesystems.

`sap.agent.fs.size[<common_params>,<fs>,<mode>]` — returns host filesystem information.  
*Parameters:*
- `fs`: filesystem mount point.
- `mode`: size mode (`total`, `free`).

`sap.agent.dev.discovery[<common_params>]` — returns JSON discovery data of block devices.

`sap.agent.dev[<common_params>,<dev>,<mode>]` — returns host block devices information.  
*Parameters:*  
- `dev`: block device name.
- `mode`: device mode (`util`, `queue`, `svc`, `wait`, `io`, `tput`).

`sap.agent.os.name[<common_params>]` — returns host name information.

`sap.agent.os.version[<common_params>]` — returns host version information.

`sap.agent.db.discovery[<common_params>]` — returns JSON discovery data of databases.

`sap.agent.proc.discovery[<common_params>]` — returns JSON discovery data of instance processes.

`sap.agent.proc.get[<common_params>,<instance>,<process>]` — returns JSON data of a specific instance process.  
*Parameters:*  
- `instance`: instance name.
- `process`: process name.

## Troubleshooting
The plugin uses Zabbix agent's logs. You can increase the debugging level of the Zabbix Agent if you need more details about what is happening.

If you encounter the error *"Cannot fetch data: Post "https://sapagent:1129": tls: failed to verify certificate: x509: certificate is not valid for any names, but wanted to match sapagent."*, you may need to set the `TLSInsecureSkipVerify` option to **true** or use a valid certificate signed by a trusted certificate authority.

## Feedback
If you encounter any issues or have suggestions for improvements, feel free [contact me](https://ravenlayer.com/contact/).

## License
This project is licensed under the Apache License 2.0 - see the LICENSE file for details.

## Disclaimer
**RavenLayer** is an independent software provider and is not affiliated with, sponsored by, or endorsed by **Zabbix LLC** or **SAP SE**. Zabbix is a trademark of Zabbix LLC. SAP and other SAP products and services mentioned herein as well as their respective logos are trademarks or registered trademarks of SAP SE in Germany and other countries. The software is provided "as is", without warranty of any kind.
