# Tiny Proxy

**Tiny Proxy** is a highly customizable and dynamic TCP layer 4 proxy server written in Golang. It supports mapping multiple ports to different forwarding addresses and allows domain-to-IP mappings through a simple HTTP API.

## Features

- **Dynamic Proxy Mapping**: Forward incoming TCP traffic on specific ports to other addresses and ports.
- **Custom Domain Mapping**: Use the API to map custom domain names to specific IP addresses.
- **"Me" Replacement**: Use `me` as a placeholder to forward traffic to the client IP address making the request.
- **Concurrency Handling**: Efficient data transfer between source and destination connections with proper connection management.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [API Documentation](#api-documentation)
  - [Domain Endpoints](#domain-endpoints)
    - [GET /domains](#get-domains)
    - [POST /domains](#post-domains)
    - [DELETE /domains](#delete-domains)
  - [Mapping Endpoints](#mapping-endpoints)
    - [GET /mappings](#get-mappings)
    - [POST /mappings](#post-mappings)
    - [PUT/PATCH /mappings](#putpatch-mappings)
    - [DELETE /mappings](#delete-mappings)
- [Examples](#examples)

## Installation

Make sure you have [Go installed](https://golang.org/doc/install). Then, follow these steps to install and run the Tiny Proxy application.

1. Clone the repository:

```bash
git clone https://github.com/nebulamods/tiny-proxy.git
cd tiny-proxy
```

2. Build the application:

Linux/MacOS:

```bash
cd src
../build/build.sh
```

Windows:

```powershell
cd src
..\build\build.cmd
```

3. Run the application:

Linux:

```bash
../bin/linux/tiny-proxy <api_address:port>
```

MacOS:

```bash
../bin/darwin/tiny-proxy <api_address:port>
```

Windows:

```powershell
..\bin\windows\tiny-proxy.exe <api_address:port>
```

- Replace `<api_address:port>` with the address and port where you want the HTTP API server to run (e.g., `:6969`).

## Usage

Once the application is running, you can interact with it via HTTP requests to the API to manage domain mappings and proxy mappings.

## API Documentation

### Domain Endpoints

#### **GET /domains**

Retrieve all existing domain mappings.

**Request**:  
`GET /domains`

**Response**:

```json
{
  "example.local": "192.168.1.100",
  "another.local": "192.168.1.101"
}
```

#### **POST /domains**

Add or update a domain-to-IP mapping.

**Request**:  
`POST /domains`

**Body**:

```json
{
  "domain": "example.local",
  "ip": "192.168.1.100"
}
```

- **domain**: The custom domain name to map (string).
- **ip**: The IP address to associate with the domain. Use `"me"` to map to the client's IP address.

**Response**:  
`204 No Content` on success.

#### **DELETE /domains**

Remove a domain-to-IP mapping.

**Request**:  
`DELETE /domains`

**Body**:

```json
{
  "domain": "example.local"
}
```

- **domain**: The custom domain name to remove (string).

**Response**:  
`204 No Content` on success.

### Mapping Endpoints

#### **GET /mappings**

Retrieve all existing TCP port mappings.

**Request**:  
`GET /mappings`

**Response**:

```json
{
  ":8080": {
    "listen_addr": ":8080",
    "forward_addr": "192.168.1.100:9090"
  },
  ":9090": {
    "listen_addr": ":9090",
    "forward_addr": "192.168.1.101:7070"
  }
}
```

#### **POST /mappings**

Add a new TCP port mapping.

**Request**:  
`POST /mappings`

**Body**:

```json
{
  "listen_addr": ":8080",
  "forward_addr": "192.168.1.100:9090"
}
```

- **listen_addr**: The local address and port on which the proxy will listen (string).
- **forward_addr**: The address and port to which incoming traffic will be forwarded. Use `"me:<port>"` to replace `"me"` with the client's IP address.

**Response**:  
`204 No Content` on success.

#### **PUT/PATCH /mappings**

Update the forwarding address of an existing TCP port mapping.

**Request**:  
`PUT /mappings` or `PATCH /mappings`

**Body**:

```json
{
  "listen_addr": ":8080",
  "forward_addr": "192.168.1.100:9091"
}
```

- **listen_addr**: The local address and port to update (string).
- **forward_addr**: The new forwarding address and port. Use `"me:<port>"` to replace `"me"` with the client's IP address.

**Response**:  
`204 No Content` on success.

#### **DELETE /mappings**

Remove an existing TCP port mapping.

**Request**:  
`DELETE /mappings`

**Body**:

```json
{
  "listen_addr": ":8080"
}
```

- **listen_addr**: The local address and port to remove (string).

**Response**:  
`204 No Content` on success.

## Examples

### Add a Domain Mapping

```bash
curl -X POST -H "Content-Type: application/json" -d '{"domain": "example.local", "ip": "me"}' http://localhost:6969/domains
```

Maps the domain `example.local` to the IP address of the client making the request.

### Add a Proxy Mapping

```bash
curl -X POST -H "Content-Type: application/json" -d '{"listen_addr": ":8080", "forward_addr": "192.168.1.100:9090"}' http://localhost:6969/mappings
```

Forwards traffic received on port `8080` to `192.168.1.100:9090`.

### Update a Proxy Mapping

```bash
curl -X PUT -H "Content-Type: application/json" -d '{"listen_addr": ":8080", "forward_addr": "me:9090"}' http://localhost:6969/mappings
```

Updates the forwarding address for port `8080` to the IP of the client making the request on port `9090`.

### Delete a Domain Mapping

```bash
curl -X DELETE -H "Content-Type: application/json" -d '{"domain": "example.local"}' http://localhost:6969/domains
```

Deletes the domain mapping for `example.local`.

### Delete a Proxy Mapping

```bash
curl -X DELETE -H "Content-Type: application/json" -d '{"listen_addr": ":8080"}' http://localhost:6969/mappings
```

Removes the proxy mapping for port `8080`.

---

Enjoy using Tiny Proxy! If you have any questions or issues, feel free to contribute or raise an issue in the repository.
