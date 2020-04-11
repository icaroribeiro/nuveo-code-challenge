# Nuveo Code Challenge

## 1 - Introduction

The purpose of this file is to present information about the work developed to solve the code challenge prepared by the company **Nuveo** that can be founded in the following link: 

*Website*: https://github.com/nuveo/backend-test/blob/master/README.md

In order to summarise, the project comprehends the implementation of a back-end application for the management of workflows. It is composed by a **REST** API developed using **Go** programming language, in addition to a **Postgres** database and a **RabbitMQ** message broker.

Throughout this documentation, a few aspects will be highlighted, such as, the configuration of environment variables of te **Postgres** database and the **RabbitMQ** message broker, and the procedures adopted to run the project with **Docker** containers.

Finally, the last section named **Project Dynamics** illustrates a brief report of how the solution works in practice.

## 2 - API Documentation

The documentation of the API implemented in **Go** programming language was developed following the **OpenAPI 3.0** specification.

Inside the directory **api-docs** there is a script named **swagger-json-to-html.py**. When running it using the **openapi.json** file, it is generated a HTML page named **index.html** within the current directory that illustrates details about the API *endpoints*.

## 3 - Project Organization

The developed solution is organized according to the structure of directories and files summarized below:

### 3.1 - Back-end

The following directories contain the **REST** API implementation using **Go** programming language.

**back-end/handlers**: it contains the handling of API requests, as well as the elaboration of API responses.

**back-end/handlers_test**: it contains the tests for handling requests to the API *endpoint* using the **Go** language test package.

**back-end/middlewares**: it contains intermediate validations of parameters transmitted through API requests.

**back-end/models**: it contains the definition of the data entities used by both the API and the database.

**back-end/postgresdb**: it contains the implementation directed to the database configuration along with some **CRUD** operations (*create*, *read* and *update*).

**back-end/postgresdb_test**: it contains the tests of the implementation of some **CRUD** operations using the **Go** language test package.

**back-end/rabbitmq**: it contains the implementation aimed to the message broker configuration together with the publishing and consumption operations.

**back-end/rabbitmq_test**: it contains the tests of the implementation of the publishing and consumption operations using the **Go** language test package.

**back-end/router**: it contains a router that exposes the routes associated with the API *endpoints*.

**back-end/router/routes**: it contains the routes associated with the API *endpoints*.

**back-end/server**: it contain an abstraction of the server that allows to "attach" some resources in order to make them available during the API requests. Here, it's used to store other structure that holds attributes to manage the data, the message queues and the storage directory where the csv files will be saved.

**back-end/services**: it contains the implementation of the generation of csv files using **yukithm/json2csv** package which is a package for converting json to csv where the csv header is generated from each json key path in the dot-bracket style. (Please, refer to https://github.com/yukithm/json2csv).

**back-end/utils**: it contains supporting functions, such as, to elaborate the API responses in JSON-like format and to generate random data used during the tests.

**back-end/.env**: it contains the environment variables for the configuration of the **development** environment.

**back-end/.test.env**: it contains the environment variables for the configuration of the **test** environment.

The **back-end/.env** file contains the environment variables referring to the connection to **Postgres** database, the exposure of the access address for HTTP communication as well as the directory where the csv files will be saved, as indicated below:

```
DB_USERNAME=user
DB_PASSWORD=password
DB_HOST=db
DB_PORT=5432
DB_NAME=db
```

```
HTTP_SERVER_HOST=0.0.0.0
HTTP_SERVER_PORT=8080
```

In order to not compromise the integrity of the database used by the project in terms of data generated from the execution of the test cases, two Postgres databases will be used.

In this sense, to facilitate future explanations regarding the details of the databases, consider that the database used for the storage of data in a "normal" actions is the **development** database and the one used for the storage of data resulting from the test cases is the **test** database named **db** and **test-db** by the **DB_NAME** environment variables defined in the **back-end/.env** and **back-end/.test.env** files, respectively. This way, it is necessary to pay special attention to the database environment variables defined in these two previous files.

### 3.2 - Postgres

The **postgresdb/scripts/1-create_tables.sql** file contains instructions for creating the **workflows** table, as detailed below:

#### 3.2.1 - Table

**Workflows**

In the **workflows** table each record contains the data of a workflows.

This way, the **id** field refers to the unique identifier of the workflow and the **status**, **data** and **steps** fields refer to its status, input and a list of names of all its steps, respectively.

| Fields           | Data type |
|:-----------------|:----------|
| id               | UUID      |
| status           | String    |
| data             | JSONB     |
| steps            | Array     |

The list of names of all its steps is configured as follows:

```
*application/json*

"steps": [
        <The name of a step>
    ]
```

#### 3.2.2 - Configurations of Docker database containers

To execute the solution through **Docker** containers, it is necessary to relate the environment variables of the **postgresdb/.env** and **postgresdb/.test.env** files with the corresponding environment variables directed to the development and test databases defined in the **back-end** application settings.

To do so, the environment variables of the **postgresdb/.env** and **postgresdb/.test.env** files must be associated with the environment variables of the **back-end/.env** and **back-end/.test.env** files, respectively.

Additionally, it is necessary to indicate that the environment variable **DB_HOST** of the **back-end/.env** and **back-end/.test.env** files must be related to the database **services** defined in the **docker-compose.yml** file.

The **docker-compose.yml** file contains the database services:

```
services:
  ...

  db:
    container_name: db
    build:
      context: ./postgresdb
      dockerfile: Dockerfile
    env_file:
      - ./postgresdb/.env
    ...

  test-db:
    container_name: test-db
    build:
      context: ./postgresdb
      dockerfile: Dockerfile
    env_file:
      - ./postgresdb/.test.env
    ...
```

**Development**

The **postgresdb/.env** file contains the database environment variables:

```
POSTGRES_USER=user
POSTGRES_PASSWORD=password
POSTGRES_DB=db
```

The **back-end/.env** file contains the database environment variables:

```
DB_USERNAME=user
DB_PASSWORD=password
DB_HOST=db
DB_PORT=5432
DB_NAME=db
```

**Test**

The **postgresdb/.test.env** file contains the database environment variables:

```
POSTGRES_USER=user
POSTGRES_PASSWORD=password
POSTGRES_DB=test-db
```

The **back-end/.test.env** file contains the database environment variables:

```
DB_USERNAME=user
DB_PASSWORD=password
DB_HOST=test-db
DB_PORT=5432
DB_NAME=test-db
```

**Important note**

After the project has been successfully executed, it is possible to check the data of the development and test databases resulting from the operations carried out at a command prompt with access to instructions directed to Docker:

```
$ docker exec -it <The id of the container of the corresponding database> /bin/bash
```

To do this, we need to inform the username and password that were previously defined by the database environment variables prior to accessing data.

```
$ psql -U <Username> <Database name>
```

In the case of the envinronment variables are kept as they were delivered, if the **id** of the container corresponds to the service named **db**, the data are obtained from the development database:

```
$ psql -U user db
```

On the other hand, if the **id** of the container corresponds to the service named **test-db**, the data are obtained from the test database:

```
$ psql -U user test-db
```

### 3.3 - RabbitMQ

Similarly to the configurations of Docker database containers, when executing the solution through **Docker** containers, it is necessary to relate the environment variables of the **rabbitmq/.env** and **rabbitmq/.test.env** files with the corresponding environment variables directed to the development and test message brokers defined in the **back-end** application settings.

To do so, the environment variables of the **rabbitmq/.env** and **rabbitmq/.test.env** files must be associated with the environment variables of the **back-end/.env** and **back-end/.test.env** files, respectively.

Additionally, it is necessary to indicate that the environment variable **MB_HOST** of the **back-end/.env** and **back-end/.test.env** files must be related to the database **services** defined in the **docker-compose.yml** file.

The **docker-compose.yml** file contains the message broker services:

```
services:
  ...

  mq:
    container_name: mq
    build:
      context: ./rabbitmq
      dockerfile: Dockerfile
    env_file:
      - ./rabbitmq/.env
    ...

  test-mq:
    container_name: test-mq
    build:
      context: ./rabbitmq
      dockerfile: Dockerfile
    env_file:
      - ./rabbitmq/.test.env
    ...
```

**Development**

File **rabbitmq/.env**:

```
RABBITMQ_DEFAULT_USER=user
RABBITMQ_DEFAULT_PASS=password
```

File **back-end/.env**:

```
MB_USERNAME=user
MB_PASSWORD=password
MB_HOST=mq
MB_PORT=5672
MB_NAME=mq
```

**Test**

File **rabbitmq/.test.env**:

```
RABBITMQ_DEFAULT_USER=user
RABBITMQ_DEFAULT_PASS=password
```

File **back-end/.test.env**:

```
MB_USERNAME=user
MB_PASSWORD=password
MB_HOST=test-mq
MB_PORT=5672
MB_NAME=test-mq
```

## 4 - How to execute the project?

It is possible to execute the project in two ways through commands from the terminal:

### 4.1 Run locally (without using Docker)

As previously described, the project was developed using the **Go** programming language and **Postgres** database.

So, for more information on how to install the software, please check the official websites:

i) https://golang.org

ii) https://www.postgresql.org

iii) https://www.rabbitmq.com

For the RabbitMQ, there is a management plugin to check the message queues through UI that should also be installed:

iv) https://www.rabbitmq.com/management.html  

After installing the softwares, create a directory structure as **github.com/icaroribeiro** and inside the **icaroribeiro** directory, download the project from Github:

```
$ git clone git@github.com:icaroribeiro/nuveo-code-challenge.git
```

Then, navigate to the back-end directory and configure the database and message-broker environment variables of the .env and .test.env files accordingly, as previously explained.

Finally, execute the project by running the command:

```
$ go run main.go
```

If there are no errors, the API endpoits will be accessed according to the host configured locally. For example:

```
http://localhost:8080
```

or even, 

```
http://127.0.0.1:8080
```

In continuity, the message queues will be also checked based on the host configured locally. For example:

```
http://localhost:15672
```

or even,

```
http://127.0.0.1:15672
```


### 4.2 Run with Docker containers

**It's recommended to run the project with Docker containers since in the case of the environment variables of all the .env and .test.env files from all directories are kept as they were delivered I strongly believe that it will not be necessary any change before executing the project.**.

At a command prompt with access to instructions directed to **Docker**, navigate to the project's root directory where the **docker-compose.yml** file is and run the command:

```
$ docker-compose up -d
```

If there are no errors, the API *endpoints* will be accessed using the address composed by the *host* and the HTTP server port **8080**. For example:

```
http://{host}:8080
```

The *host* corresponds to the value returned when executing a command at a command prompt with access to instructions directed to Docker:

```
$ docker-machine ip
```

In continuity, suppose the *host* is: 192.168.99.100. As a result, the API requests can be performed through a front-end client or test tool like Postman using the address as:

```
http://192.168.99.100:8080
```

Lastly, the message queues will be also checked based on the host configured by the Docker. However, since it was proposed two message brokers containers, they will be accessed using two different ports. For example:

**Development**

```
http://192.168.99.100:15673
```

**Test**

```
http://192.168.99.100:15674
```

(P.S.: For some unknown reason, I could not open the Rabbitmq management UI using the previous addresses when accessing from Microsoft Edge browser, however it works as expected using the Google Chrome browser).

In addition, it is also worth emphasizing that the entire configuration related to **Docker** was evaluated in this documentation based on the **DockerToolbox** tool for Windows.

## 5 - How to use the API *endpoints*?

The API request are performed through the HTTP server port **8080** and the API responses can be viewed by means of a **front-end** client or test tool, for example **Postman**.

In what follows, there is a guide that includes API requests for creating, obtaining and updating data from the database.

(P.S., before checking the following examples, consider that no data is recorded prior to this explanation).

### Status

Request:

```
Method: HTTP GET
```

```
URL: http://{host}:8080/status
```

Response:

```
Code: 200 OK - In the case of the service has started up correctly and is ready to accept requests.
```

### Management of Workflows

#### Creation of Workflows

Request:

```
Method: HTTP POST
```

```
URL: http://{host}:8080/workflow
```

```
*application/json*

Body: {
    "data": {
        "array": [
        	true,
        	1.1,
            1,
            {
                "key": "value"
            },
            "string"
        ],
        "boolean": true,
        "float": 1.1,
        "integer": 1,
        "object": {
            "key": "value"
        },
        "string": "string"
    },
    "steps": [
        "Step1"
    ]
}
```

Response:

```
Code: 200 OK - In the case of the workflow is successfully created.
```

```
*application/json*

Body: {
    "id": "f7e29b15-da15-4197-8d77-baee49826b3a",
    "status": "inserted",
    "data": {
        "array": [
            true,
            1.1,
            1,
            {
                "key": "value"
            },
            "string"
        ],
        "boolean": true,
        "float": 1.1,
        "integer": 1,
        "object": {
            "key": "value"
        },
        "string": "string"
    },
    "steps": [
        "Step1"
    ]
}
```

**Important note**:

Whenever the API request to create a workflow is performed, the related workflow is recorded in the database with the **status** field defined automatically to **inserted**. Therefore, the only parameters required in the request body are the workflow **data** and **steps**.

#### Listing of Workflows

Request:

```
Method: HTTP GET
```

```
URL: http://{host}:8080/workflows
```

Response:

```
Code: 200 OK - In the case of the list of all workflows is successfully obtained.
```

```
*application/json*

Body: [
    {
        "id": "f7e29b15-da15-4197-8d77-baee49826b3a",
        "status": "inserted",
        "data": {
            "array": [
                true,
                1.1,
                1,
                {
                    "key": "value"
                },
                "string"
            ],
            "boolean": true,
            "float": 1.1,
            "integer": 1,
            "object": {
                "key": "value"
            },
            "string": "string"
        },
        "steps": [
            "Step1"
        ]
    }
]
```

#### Consumption of Workflows

Request:

```
Method: HTTP GET
```

```
URL: http://{host}:8080/workflows/consume
```

Response:

```
Code: 200 OK - In the case of the workflow from the head of the queue is successfully consumed.
```

```
*application/json*

Body: {
    "id": "f7e29b15-da15-4197-8d77-baee49826b3a",
    "status": "inserted",
    "data": {
        "array": [
            true,
            1.1,
            1,
            {
                "key": "value"
            },
            "string"
        ],
        "boolean": true,
        "float": 1.1,
        "integer": 1,
        "object": {
            "key": "value"
        },
        "string": "string"
    },
    "steps": [
        "Step1"
    ]
}
```

#### Updating of a Workflow by its id

Request:

```
Method: HTTP PUT
```

```
URL: http://{host}:8080/workflows/f7e29b15-da15-4197-8d77-baee49826b3a
```

```
*application/json*

Body: {
	"status": "consumed"
}
```

Response:

```
Code: 200 OK - In the case of the workflow is successfully updated.
```

```
*application/json*

Body: {
    "id": "f7e29b15-da15-4197-8d77-baee49826b3a",
    "status": "consumed",
    "data": {
        "array": [
            true,
            1.1,
            1,
            {
                "key": "value"
            },
            "string"
        ],
        "boolean": true,
        "float": 1.1,
        "integer": 1,
        "object": {
            "key": "value"
        },
        "string": "string"
    },
    "steps": [
        "Step1"
    ]
}
```

**Important note**:

Whenever the API request to updated a workflow is performed, the **status** field of the related workflow is changed from **inserted** to **consumed** in the database. Therefore, the only parameter required in the request body is the workflow **status** and it must be configured to **consumed**.

## 6 - CSV file generation with workflow data

As previously commented, all csv files generated during a "normal" execution of the project will be saved in the directory informed in the environment variable named **STORAGE_DIR** of the **back-end/.env** file. The name of each csv file is composed by the id of the related workflow and the .csv extension. (By default, they wll be placed in a directory named **storage** inside the **back-end** directory).

## 7 - Tests

In order to test the solution three **test sets** were developed.

(P.S., only the first two sets involve creating and editing records from the test database).

### 7.1 Database

The tests are related to some **CRUD** operations (*create*, *read* and *update*) in the test database.

To execute them, navigate to the **back-end/postgresdb_test** directory.

So, if you prefer to evaluate all tests at once, run the command:

```
$ go test -v
```

However, it is also possible to run each test separately using the commands:

**Tests of the some CRUD operations directed to Workflows**

```
$ go test -v -run=TestCreateWorkflow
```

```
$ go test -v -run=TestGetAllWorkflows
```

```
$ go test -v -run=TestGetWorkflow
```

```
$ go test -v -run=TestUpdateWorkflow
```

### 7.2 Handlers

These tests are related to the API requests.

In this regard, navigate to the **back-end/handlers_test** directory.

So, if you prefer to evaluate all tests at once, run the command:

```
$ go test -v
```

Nevertheless, it is also possible to run each test separately using the commands:

**Tests of the API requests directed to Workflows**

```
$ go test -v -run=TestGetAllWorkflows
```

```
$ go test -v -run=TestConsumeWorkflow
```

```
$ go test -v -run=TestCreateWorkflow
```

```
$ go test -v -run=TestUpdateWorkflow
```

**Important note**:

The test of the API request directed to the consumption of a workflow comprehends to consume it from queue in addition to generate a csv file with its related data. In this case, the csv files will be stored in the directory informed in the environment variable named **STORAGE_DIR** of the **back-end/.test.env** file. (By default the files will be placed in the **storage** directory inside the **back-end/handlers_test** directory).

### 7.3 Message Broker

These tests are related to the publishing and .

To achieve this, navigate to the **back-end/rabbitmq_test** directory.

So, if you prefer to evaluate all tests at once, run the command:

```
$ go test -v
```

Even so, it is also possible to run each test separately using the commands:

```
$ go test -v -run=TestPublish
```

```
$ go test -v -run=TestConsume
```

## 7 - Project Dynamics

In what follows, there is a brief account of how the solution works in practice meeting the requirements specified in the comments of the code challenge.

(P.S., consider that no data is recorded prior to this explanation).

#### Creation of Workflows

Request:

```
Method: HTTP POST
```

```
URL: http://{host}:8080/workflow
```

```
*application/json*

Body: {
    "data": {
        "array": [
        	true,
        	1.1,
            1,
            {
                "key": "value"
            },
            "string"
        ],
        "boolean": true,
        "float": 1.1,
        "integer": 1,
        "object": {
            "key": "value"
        },
        "string": "string"
    },
    "steps": [
        "Step1"
    ]
}
```

Response:

```
Code: 200 OK - In the case of the workflow is successfully created.
```

```
*application/json*

Body: {
    "id": "f7e29b15-da15-4197-8d77-baee49826b3a",
    "status": "inserted",
    "data": {
        "array": [
            true,
            1.1,
            1,
            {
                "key": "value"
            },
            "string"
        ],
        "boolean": true,
        "float": 1.1,
        "integer": 1,
        "object": {
            "key": "value"
        },
        "string": "string"
    },
    "steps": [
        "Step1"
    ]
}
```

Database:

```
postgres=# \c db

db=# select * from workflows;
                  id                  |  status  |                                                                        data                                                                        |  steps
--------------------------------------+----------+----------------------------------------------------------------------------------------------------------------------------------------------------+---------
 f7e29b15-da15-4197-8d77-baee49826b3a | inserted | {"array": [true, 1.1, 1, {"key": "value"}, "string"], "float": 1.1, "object": {"key": "value"}, "string": "string", "boolean": true, "integer": 1} | {Step1}
```

Message Broker:

Whenever the API request to create a workflow is performed, a workflow is created in the database and added to the queue. With this in mind, to check the workflow in the queue, it is necessary to access the Rabbitmq management UI using a Web browser at the URL **http://{hostname}:15672** (This is the URL when running the project locally).

#### Listing of Workflows

Request:

```
Method: HTTP GET
```

```
URL: http://{host}:8080/workflows
```

Response:

```
Code: 200 OK - In the case of the list of all workflows is successfully obtained.
```

```
*application/json*

Body: [
    {
        "id": "f7e29b15-da15-4197-8d77-baee49826b3a",
        "status": "inserted",
        "data": {
            "array": [
                true,
                1.1,
                1,
                {
                    "key": "value"
                },
                "string"
            ],
            "boolean": true,
            "float": 1.1,
            "integer": 1,
            "object": {
                "key": "value"
            },
            "string": "string"
        },
        "steps": [
            "Step1"
        ]
    }
]
```

#### Consumption of Workflows

Request:

```
Method: HTTP GET
```

```
URL: http://{host}:8080/workflows/consume
```

Response:

```
Code: 200 OK - In the case of the workflow from the head of the queue is successfully consumed.
```

```
*application/json*

Body: {
    "id": "f7e29b15-da15-4197-8d77-baee49826b3a",
    "status": "inserted",
    "data": {
        "array": [
            true,
            1.1,
            1,
            {
                "key": "value"
            },
            "string"
        ],
        "boolean": true,
        "float": 1.1,
        "integer": 1,
        "object": {
            "key": "value"
        },
        "string": "string"
    },
    "steps": [
        "Step1"
    ]
}
```

As part of the process of the consumption of a workflow, a csv file is created with the workflow data and named **f7e29b15-da15-4197-8d77-baee49826b3a.csv**:

```
> cat f7e29b15-da15-4197-8d77-baee49826b3a.csv 
boolean,float,integer,string,array[0],array[1],array[2],array[4],object.key,array[3].key
true,1.1,1,string,true,1.1,1,string,value,value
```

#### Updating of a Workflow by its id

Request:

```
Method: HTTP PUT
```

```
URL: http://{host}:8080/workflows/f7e29b15-da15-4197-8d77-baee49826b3a
```

```
*application/json*

Body: {
	"status": "consumed"
}
```

Response:

```
Code: 200 OK - In the case of the workflow is successfully updated.
```

```
*application/json*

Body: {
    "id": "f7e29b15-da15-4197-8d77-baee49826b3a",
    "status": "consumed",
    "data": {
        "array": [
            true,
            1.1,
            1,
            {
                "key": "value"
            },
            "string"
        ],
        "boolean": true,
        "float": 1.1,
        "integer": 1,
        "object": {
            "key": "value"
        },
        "string": "string"
    },
    "steps": [
        "Step1"
    ]
}
```

Database:

```
postgres=# \c db

db=# select * from workflows;
                  id                  |  status  |                                                                        data                                                                        |  steps
--------------------------------------+----------+----------------------------------------------------------------------------------------------------------------------------------------------------+---------
 f7e29b15-da15-4197-8d77-baee49826b3a | consumed | {"array": [true, 1.1, 1, {"key": "value"}, "string"], "float": 1.1, "object": {"key": "value"}, "string": "string", "boolean": true, "integer": 1} | {Step1}
```