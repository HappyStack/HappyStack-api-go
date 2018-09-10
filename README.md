# HappyStack-api-go
HappyStack api written in Go



## Api

### List items
GET `users/:userId/items`  

#### Response
```json
[{
		"id": 1,
		"name": "Magnesium",
		"dosage": "some",
		"takenToday": false,
		"servingSize": 1,
		"servingType": "pill",
		"timing": "0001-01-01T00:00:00Z"
	},
	{
		"id": 7,
		"name": "Vitamine k",
		"dosage": "",
		"takenToday": false,
		"servingSize": 1,
		"servingType": "pill",
		"timing": "0001-01-01T00:00:00Z"
	}
]
```

### Add item
POST `users/:userId/items`  
```json
{
  "name": "Vitamin D",
  "dosage": "2000UI",
  "takenToday": true,
  "servingSize": 3,
  "servingType": "pill",
  "timing": "0001-01-01T00:00:00Z"
}
```
*"servingType" is "pill", "drop" or "scoop"*

## Run
```
go run main.go item.go handlers.go routes.go logger.go router.go repo.go auth.go
```

## build for linux
```
GOOS=linux GOARCH=amd64 go build
```


## Connect to the remote Database with Postico (Postgres SQL Client)
In [Postico](https://eggerapps.at/postico/) add a `New 
Favorite`(connexion)

**Nickname**: `HappyStackServer` for example  
**Host**: `[leave blank]` (localhost by default)  
**User**: `happystack`  
**Password**: `[Use database password]`  
**Database**: `happystack`  
Click `Option` the choose `Connect via SSH`  
**SSH Host**: `104.248.56.250`
**User**: `happystack`  
**Private key** `[Select your private key that is authorized on the server]`


## Deploy Server
- In the terminal, connect in ssh to server `ssh happystack@104.248.56.250`
- Look for existing tmux sessions with `tmux ls`

 *Runing on `tmux` so that the server still runs once you shut down the ssh session.*

- If none is listed, open one session with `tmux`
- Otherwise link to existing session with `tmux attach -t [session name or id]`
- Navigate to the source code `cd go/src/HappyStack-api-go`
- Pull the repo `git pull`
- Run the binary program with `go run *.go` (TODO use install binay later)
- Detach tmux sesison by typing `ctrl+b` then `d`
- You can safely close the ssh session and the process will keep running in background \o/