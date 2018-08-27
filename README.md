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
