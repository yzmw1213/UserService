[
  {
    "essential": true,
    "memory": 256,
    "name": "${CONTAINER_NAME}",
    "cpu": 256,
    "image": "${REPOSITORY_URL}",
    "portMappings": [
        {
            "containerPort": 8888,
            "hostPort": 8888
        }
    ],
    "environment": [
      {
        "name": "DB_ADRESS",
        "value": "${DB_ADRESS}"
      },
      {
        "name": "DB_NAME",
        "value": "${DB_NAME}"
      },
      {
        "name": "DB_PASSWORD",
        "value": "${DB_PASSWORD}"
      },
      {
        "name": "DB_USER",
        "value": "${DB_USER}"
      }
    ]
  }
]
