# Sales Tax Problem
The idea of this project creates example application solving "Sales Tax Problem."

## Problem

Basic sales tax is applicable at a rate of 10% on all goods, except books, food, and medical products that are exempt. Import duty is an additional sales tax applicable to all imported goods at a rate of 5%, with no exemptions.

When I purchase items, I receive a receipt which lists the name of all the items and their price (including tax), finishing with the total cost of the items, and the total amounts of sales taxes paid. The rounding rules for sales tax are that for a tax rate of n%, a shelf price of p contains (np/100 rounded up to the nearest 0.05) amount of sales tax.


## Limitations

* Data consistency between modules is limitted. There is no controll between modules about usage. You can delete data easily.

* Configuration management in plans and because of that  currently some configuration is hardcoded.

* Rest API is primitive so you can get errors that you shouldn't see. Error management is missing in API side

## Prerequests

This project need following dependencies :

* go installation
* `dep` dependency management
* `build-essentials` for able to run __MakeFile__

## How to Use
  Some commands listed below but you can also check MakeFile for more info.

* __Building Binaries:__

You can run flowing commands for creating binaries. Binaries will be created under `$(Project_Dir)/release` folder.

```bash
make
```
or

```bash
 make build
```

* __Running App:__

You can starting application without building. Application Server will be start at `8080`

```bash
make run
```

* __Running Tests:__

For running test

```bash
make test
```

* __Clean:__

Cleaning development store and release folder.

```bash
make clean
```

* __Load Demo:__

Triggers demo scenario in server. Nothing will return here. After you run this endpoint you can check demo data in server.

__NOTE__ Load demo imports spesific data. You should run only once. You can see it's output in application logs. After loading data you can use rest api.

```bash
curl -X POST http://localhost:8080/demo
```

## Docs

You can find more info in [docs](https://documenter.getpostman.com/view/5717174/RzZ1rNiG) about rest api