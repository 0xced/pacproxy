pacproxy
========

[![Build Status](https://travis-ci.org/williambailey/pacproxy.svg)](https://travis-ci.org/williambailey/pacproxy)

A no-frills local HTTP proxy server powered by a [proxy auto-config (PAC) file](https://web.archive.org/web/20070602031929/http://wp.netscape.com/eng/mozilla/2.0/relnotes/demo/proxy-live.html). Especially handy when you are working in an environment with many different proxy servers and your applications don't support proxy auto-configuration.

```
$ ./pacproxy -h
pacproxy v2.0.6

A no-frills local HTTP proxy server powered by a proxy auto-config (PAC) file
https://github.com/williambailey/pacproxy

Usage:
  -c string
        PAC file name, url or javascript to use (required)
  -l string
        Interface and port to listen on (default "127.0.0.1:8080")
  -s string
        Scheme to use for the URL passed to FindProxyForURL
  -v    send verbose output to STDERR
```

```bash
# shell 1
pacproxy -l 127.0.0.1:8080 -s http -c 'function FindProxyForURL(url, host){ console.log("hello pac world!"); return "PROXY random.example.com:8080"; }'
# shell 2
pacproxy -l 127.0.0.1:8443 -s https -c 'function FindProxyForURL(url, host){ console.log("hello pac world!"); return "PROXY random.example.com:8080"; }'
# shell 3
export http_proxy="127.0.0.1:8080"
export https_proxy="127.0.0.1:8443"
curl -I "http://www.example.com"
curl -I "https://www.example.com"
```

## Using pacproxy as a Windows Service

> [!IMPORTANT]  
> The commands in this section should be executed with PowerShell.

1. Configure `http_proxy` and `https_proxy` in the **Environment Variables** control panel (also optionally `no_proxy`)

<img src=".assets/env_http_proxy.png" alt="Edit User Variable for http_proxy window" style="zoom:75%;" />

<img src=".assets/env_https_proxy.png" alt="Edit User Variable for https_proxy window" style="zoom:75%;" />

<img src=".assets/env_no_proxy.png" alt="Edit User Variable for no_proxy window" style="zoom:75%;" />

Here are the values used in this example:

```
http_proxy = http://127.0.0.1:24944
https_proxy = http://127.0.0.1:24945
no_proxy = localhost,127.0.0.1
```

2. Install the `pacproxy.exe` binary

At that point the `https_proxy` environment variable must set manually with the correct proxy. In order to find out which proxy should be used to download the `pacproxy` dependencies, you can use the `pactester.exe` tool that comes with the [pacparser](https://github.com/manugarg/pacparser/releases/) library.

```powershell
$pacUrl = Get-ItemPropertyValue "HKCU:\Software\Microsoft\Windows\CurrentVersion\Internet Settings" "AutoConfigURL"
& curl --silent $pacUrl | pactester.exe -p - -u https://github.com
```

It should print the proxy that you will have to use to successfully download the dependencies.

```
PROXY proxy.example.com
```

Use that proxy to set the `https_proxy` environment variable.

```powershell
${Env:https_proxy} = "proxy.example.com"
```

Finally, install the `pacproxy.exe` binary.
```powershell
go install
```

3. Ensure that the `pacproxy.exe` binary has been successfully installed by checking the contents of the `$GOPATH\bin` directory:

   ```powerhell
   $appDir = "$(go env GOPATH)\bin"
   dir $appDir
   ```

4. In an administrator PowerShell, create two Windows Services with [NSSM](https://nssm.cc/)

   For http:
   
   ```powershell
   nssm install PacProxyHttpSvc "$appDir\pacproxy.exe"
   nssm set PacProxyHttpSvc AppParameters "-l 127.0.0.1:24944 -s http -c $pacUrl -v"
   nssm set PacProxyHttpSvc AppDirectory "$appDir"
   nssm set PacProxyHttpSvc AppStderr "$appDir\pacproxy_http.log"
   nssm set PacProxyHttpSvc DisplayName "PAC Proxy for HTTP requests"
   nssm set PacProxyHttpSvc Description "HTTP proxy server running on 127.0.0.1:24944 for applications that don't support proxy auto-configuration"
   nssm start PacProxyHttpSvc
   ```
   
   For https:
   
   ```powershell
   nssm install PacProxyHttpsSvc "$appDir\pacproxy.exe"
   nssm set PacProxyHttpsSvc AppParameters "-l 127.0.0.1:24945 -s https -c $pacUrl -v"
   nssm set PacProxyHttpsSvc AppDirectory "$appDir"
   nssm set PacProxyHttpsSvc AppStderr "$appDir\pacproxy_https.log"
   nssm set PacProxyHttpsSvc DisplayName "PAC Proxy for HTTPS requests"
   nssm set PacProxyHttpsSvc Description "HTTPS proxy server running on 127.0.0.1:24945 for applications that don't support proxy auto-configuration"
   nssm start PacProxyHttpsSvc
   ```

## Using pacproxy as a macOS Launch Agent

Requirements

```sh
brew install go
brew install pacparser
```

1. Configure the `http_proxy` and `https_proxy`  environment variables (also optionally `no_proxy`), for example in the `~/.zprofile` file.

```sh
export http_proxy=http://127.0.0.1:24944
export https_proxy=http://127.0.0.1:24945
export no_proxy=localhost,127.0.0.1
```

2. Install the `pacproxy` binary

At that point the `https_proxy` environment variable must set manually with the correct proxy. In order to find out which proxy should be used to download the `pacproxy` dependencies, you can use the `pactester` tool that comes with the [pacparser](https://github.com/manugarg/pacparser/releases/) library.

```sh
unset http_proxy
unset https_proxy
unset no_proxy
pacUrl=$(scutil --proxy | grep "ProxyAutoConfigURLString : " | awk -F "ProxyAutoConfigURLString : " '{print $2}')
curl --silent $pacUrl | pactester -p - -u https://github.com
```

It should print the proxy that you will have to use to successfully download the dependencies.

```
PROXY proxy.example.com:80
```

Use that proxy to set the `https_proxy` environment variable.

```sh
export https_proxy="proxy.example.com:80"
```

Finally, install the `pacproxy` binary.
```sh
go install
```

3. Ensure that the `pacproxy` binary has been successfully installed by checking the contents of the `$GOPATH/bin` directory:

```sh
appDir="$(go env GOPATH)/bin"
ls -al $appDir
```

4. Create then load the Launch Agents

For http:

```sh
cat <<EOF > ${HOME}/Library/LaunchAgents/PacProxyHttpSvc.plist
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>EnvironmentVariables</key>
	<dict>
		<key>PATH</key>
		<string>${appDir}</string>
	</dict>
	<key>KeepAlive</key>
	<dict>
		<key>SuccessfulExit</key>
		<true/>
	</dict>
	<key>Label</key>
	<string>PacProxyHttpSvc</string>
	<key>ProgramArguments</key>
	<array>
		<string>./pacproxy</string>
		<string>-l</string>
		<string>127.0.0.1:24944</string>
		<string>-s</string>
		<string>http</string>
		<string>-c</string>
		<string>${pacUrl}</string>
		<string>-v</string>
	</array>
	<key>RunAtLoad</key>
	<true/>
	<key>StandardErrorPath</key>
	<string>pacproxy_http.log</string>
	<key>WorkingDirectory</key>
	<string>${appDir}</string>
</dict>
</plist>
EOF
```

```sh
launchctl load -w ${HOME}/Library/LaunchAgents/PacProxyHttpSvc.plist
```

For https:

```sh
cat <<EOF > ${HOME}/Library/LaunchAgents/PacProxyHttpsSvc.plist
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>EnvironmentVariables</key>
	<dict>
		<key>PATH</key>
		<string>${appDir}</string>
	</dict>
	<key>KeepAlive</key>
	<dict>
		<key>SuccessfulExit</key>
		<true/>
	</dict>
	<key>Label</key>
	<string>PacProxyHttpsSvc</string>
	<key>ProgramArguments</key>
	<array>
		<string>./pacproxy</string>
		<string>-l</string>
		<string>127.0.0.1:24945</string>
		<string>-s</string>
		<string>https</string>
		<string>-c</string>
		<string>${pacUrl}</string>
		<string>-v</string>
	</array>
	<key>RunAtLoad</key>
	<true/>
	<key>StandardErrorPath</key>
	<string>pacproxy_https.log</string>
	<key>WorkingDirectory</key>
	<string>${appDir}</string>
</dict>
</plist>
EOF
```

```sh
launchctl load -w ${HOME}/Library/LaunchAgents/PacProxyHttpsSvc.plist
```

> [!NOTE]  
> The Launch Agents plists were generated with [Lingon](https://www.peterborgapps.com/lingon/).

## License

> Copyright 2020 William Bailey
>
> Licensed under the Apache License, Version 2.0 (the "License");
> you may not use this file except in compliance with the License.
> You may obtain a copy of the License at
>
>     http://www.apache.org/licenses/LICENSE-2.0
>
> Unless required by applicable law or agreed to in writing, software
> distributed under the License is distributed on an "AS IS" BASIS,
> WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
> See the License for the specific language governing permissions and
> limitations under the License.
