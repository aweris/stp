# Sales Tax Problem [![GitHub license](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=flat-square)](https://github.com/aweris/std/blob/master/LICENSE)
The idea of this project to create an example application solving "Sales Tax Problem."

## Problem

Basic sales tax is applicable at a rate of 10% on all goods, except books, food, and medical products that are exempt. Import duty is an additional sales tax applicable to all imported goods at a rate of 5%, with no exemptions.

When I purchase items, I receive a receipt which lists the name of all the items and their price (including tax), finishing with the total cost of the items, and the total amounts of sales taxes paid. The rounding rules for sales tax are that for a tax rate of n%, a shelf price of p contains (np/100 rounded up to the nearest 0.05) amount of sales tax.


## Limitations

* Data consistency between modules is limited. There is no control between modules about usage. You can delete data easily.

* Application configuration is mostly hardcoded but adding configuration management is in the TODO list.

* Rest API is primitive so you can get errors that you shouldn't see. Error management is missing in API side

## Prerequisites
 For building or developing this application you need to install and configure following `go`, `dep` and `build-essentials`

## Installing / Getting started

A quick introduction of the minimal setup you need to get  STP up & running.

for OSX :

```bash
> make & bash -c release/stp-darwin-amd64
```

for Linux :

```bash
> make & bash -c release/stp-linux-amd64
```

When you execute this commands, you'll create binaries under `$(PROJECT_DIR)`/release folder and execute binary.

## Developing

#### Setting up:

```bash
git clone https://github.com/aweris/stp.git
cd stp/
make get
```

#### Running App :

```bash
make run
```

#### Running Tests :

```bash
make test
```

#### Building Binaries :

##### Building All :

```bash
make
```

or

```
make build
```

the main difference between these 2 commands is `make` is running `clean`,`get,` `test` before starting the build

##### Building only osx :

```bash
make build-darwin
```

##### Building only Linux :

```bash
make build-linux
```

#### Clean resources :

Cleaning development store and release folder.

```bash
make clean
```

#### Load Demo:

Triggers a demo scenario in the server. After you run this endpoint, you can check demo data in the server.

__NOTE__ Load demo imports specific data. You should run only once. You can see it's output in application logs. After loading data, you can use Rest API. You need to run `make clean` for reset development db.

```bash
curl -X POST http://localhost:8080/demo
```

## Docs

You can find more info in [docs](https://documenter.getpostman.com/view/5717174/RzZ1rNiG) about rest API