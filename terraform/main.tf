resource "docker_image" "redis" {
    name = "redis:7.4.1"
}

resource "docker_container" "redis" {
    image = docker_image.redis.image_id
    name  = "redis"

    restart = "always"

    # ports {
    #     internal = 6379
    #     external = 6379
    # }

    # networks_advanced {
    #     name = docker_network.private_network.name
    # }
}

resource "docker_image" "mysql" {
    name = "mysql:8"
}

resource "docker_container" "mysql" {
    image = docker_image.mysql.image_id
    name = "mysql"

    env = [
        "MYSQL_ROOT_PASSWORD=root",
        "MYSQL_DATABASE=test",
    ]

    restart = "always"

    # ports {
    #     internal = 3306
    #     external = 3306
    # }

    # networks_advanced {
    #     name = docker_network.private_network.name
    # }
}

resource "docker_image" "api_server" {
    name = "ghcr.io/toysin/api_server:v0.0.0"
}

# locals {
#     # mysql_host = docker_container.mysql.name
#     # redis_addr = join(":", [docker_container.redis.name, docker_container.redis.ports[0].external])
#     mysql_host = docker_container.mysql.network_data[0].ip_address
#     redis_addr = join(":", [docker_container.redis.network_data[0].ip_address, docker_container.redis.ports[0].external])
# }

resource "docker_container" "api_server" {
    image = docker_image.api_server.image_id
    name  = "api_server"

    ports {
        internal = 8080
        external = 8080
    }

    env = [
        "DATABASE_TYPE=mysql",
        "DATABASE_HOST=localhost",
        "DATABASE_PORT=3306",
        "DATABASE_USER=root",
        "DATABASE_PASSWORD=root",
        "DATABASE_NAME=test",

        "REDIS_ADDR=localhost:6379",
    ]

    restart = "always"

    network_mode = "host"

    # networks_advanced {
    #     name = docker_network.private_network.name
    # }

    depends_on = [
        docker_container.redis,
        docker_container.mysql,
    ]
}

resource "docker_image" "process_server" {
    name = "ghcr.io/toysin/process_server:v0.0.0"
}
