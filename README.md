## utelnet - get information about your Ubiquiti EdgeMAX devices

utelnet will check to see what information it can gather from your Ubiquiti EdgeMAX devices. Currently this has been tested on Ubiquiti EdgeRouters only.

### How to use

Clone the repo:

```
git clone https://github.com/sheran/utelnet.git
```

Build the tool

```
go build
```

Run the tool by pointing it to your EdgeRouter Web GUI:

```
$ utelnet https://10.0.0.1
host_url: https://10.0.0.1
model: ER-X
lib_date: 2023-06-15 16:29:23 +0800 +08
hostname: Host is "EdgeRouter-X-5-Port"
websocket: true
$
```

You can also run without building:

```
$ go run main.go https://10.0.0.1
host_url: https://10.0.0.1
model: ER-X
lib_date: 2023-06-15 16:29:23 +0800 +08
hostname: Host is "EdgeRouter-X-5-Port"
websocket: true
$
```
