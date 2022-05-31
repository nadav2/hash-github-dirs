```
First build the docker image with the following command:
    docker-compose build

then run the docker image with the following command:
    docker-compose up
```

```
The initialize serviec is running on port 8081.
And the api service is running on port 8080.
```

```
    To use the service make a post request to the initialize service.
    Examople of the json body:
    {
        "ref": "google/go-github"
    }
```

```
    You finished set up the service.
    Now you can start sending api requests to the service.
    To do so, you need to make a post requests to the api services. 
``` 


*For calling the getFileContent api, you need to make a post request to the api service.*


```
    For example in python:
        data = { "fileName": "go.sum" }
        file_content = requests.post(f"{main_api_url}/get_file_content", json=data).json()  
        print(file_conten["fileContent"])
    Or in curl:
        curl -X POST -H "Content-Type: application/json" -d '{"fileName": "go.sum"}' http://localhost:8080/get_file_content
``` 


*For calling the hashFiles api is mush the same as the getFileContent api.*
    
    
```
    For example in curl:
        curl -X POST -H "Content-Type: application/json" -d '{"files": ["go.sum", "README.md"]}' http://localhost:8080/hash_files
    
```

