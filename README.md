# Hailstorm

# Installation
For now starting Hailstorm is as easy as:
```
make run
```

# Creating an application
Post to `/app` to create an Hailstorm application
```json
{
	"name": "My first app
}
```

Should give you a response like this:
```json
{
  	"id": "0744b007-681d-4473-8db5-1db9631de4df",
  	"name": "My first app",
  	"createAt": "2023-12-23T09:42:49.913529415+01:00"
}
```

# Deploying functions 
To deploy a function you post a binary WASM blob to `/app/{id}/deploy`. That will give you a response like this:
```
{
  "id": "a40fafe4-8b67-4c4a-992c-a34590c227e5",
  "appId": "0744b007-681d-4473-8db5-1db9631de4df",
  "createdAt": "2023-12-23T09:45:21.526261201+01:00"
}
```

Your function should be reachable on `/proxy/{deploy_id}`

