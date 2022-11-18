# Topic implementation using SIP conference
- sample data can be seen in [`data.json`](/initializers/data.json)
- the API can be seen in [`api.yml`](/api.yml)

## Convention based on [RFC4575](https://datatracker.ietf.org/doc/rfc4575/)
- Each conference has a Conference Information (`ConfInfo`)
- The `ConfInfo` has a list of participants (many `User` in a slice called `Users`)
- Each `ConfInfo` is used as a `Topic` entity in a pub/sub system.

## Usage

To run test, docker and docker-compose are required

- pull project
- to start: `docker-compose up --build --force-recreate`
- try CRUD commands via curl
- to stop: `docker-compose down --remove-orphans --volumes --rmi 'local'`


## Example curl commands
- Get all ConfInfos: `curl --location --request GET 'http://localhost:8080/api/v1/confInfos'`
- Add a User to a Topic(id=3): 
    ```curl --location --request PATCH 'http://localhost:8080/api/v1/topicMode/confInfos/3' \
    --header 'Content-Type: application/json' \
    --data-raw '{
    "EntityURL": "newUser@url.com",
    "Role": "Subscriber"
    }'```
- Delete a User (id=2) from a Topic(id=2):
    ```
    curl --location --request DELETE 'http://localhost:8080/api/v1/topicMode/confInfos/subject2/users/2' \
    --data-raw ''
    ```