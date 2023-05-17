# Endpoints
If you use the vscode editor, you can refer to [vscode.http](/vscode.http) to easily test the rest-api.

If you dont use the vscode editor, you can use the following endpoints to test:



#
###  Home or Welcome Page
GET `/zample`

```bash
curl --request GET \
  --url http://localhost:4000/zample/ 
```

#
###  Example (method not allowed)
POST `/zample`


```bash
curl --request POST \
  --url http://localhost:4000/zample \
  --header 'content-type: application/json' 
```


#
###  Example (url not found, on router level not app level)
GET `/zample/{param}`


Request:
- `param`*: `any`. demo to show it doesnt exists
```bash
curl --request GET \
  --url http://localhost:4000/zample/{does-not-exist} \
  --header 'content-type: application/json' 
```


#
### Validation endpoint: here validation done via validation middleware
POST `/zample/validation-via-middleware`


Request:
- `name`: ``string`` name
- `phone`: ``string`` phone number in e164 format
- `email`: ``string`` email
```bash
curl --request POST \
  --url http://localhost:4000/zample/validation-via-middleware \
  --header 'content-type: application/json'  \
  --data '{"name": "abcde fghi","phone": "+230","email": "example@domain.com"}'
```


#
###  Validation endpoint: here validation done in handler
POST `/zample/validation-in-handler`

Request:
- `name`: ``string`` name
- `phone`: ``string`` phone number in e164 format
- `age`: ``int`` age 
- `email`: ``string`` email
```bash
curl --request POST \
  --url http://localhost:4000/zample/validation-in-handler \
  --header 'content-type: application/json'  \
  --data '{"name": "abcde fghi","phone": "+230","Age": 2,"email": "example@domain.com"}'
```


#
### List using db helpers
GET `/zample/list-via-db-helpers`

Request:
```bash
curl --request GET \
  --url http://localhost:4000/zample/list-via-db-helpers \
  --header 'content-type: application/json' 
```


#
### List using orm in datastore
GET `/zample/list-via-repo`

Request:
```bash
curl --request GET \
  --url http://localhost:4000/zample/list-via-repo \
  --header 'content-type: application/json' 
```


#
### Example to show error-500
GET `/zample/error-500` 

Request:
- `body`: any.  can be anything, just trying to demo example of wrong body parsing
```bash
curl --request GET \
  --url http://localhost:4000/zample/error-500 \
  --header 'content-type: application/json'  \
  --data 'kj body send to log'
```


#
### Example to show error returned
GET `/zample/error-returned`

Request:
```bash
curl --request GET \
  --url http://localhost:4000/zample/error-returned \
  --header 'content-type: application/json' 
```



 
 