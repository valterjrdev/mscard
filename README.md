# ms card

Card transaction authorizer

**This project is for studies not used in production**

## Usage

Install Dependencies

* [docker engine](https://docs.docker.com/engine/install/ubuntu/)
* [docker-compose](https://docs.docker.com/compose/install/) 

Generate `.env` file
```bash
cp .env.example .env
```
Run using docker-compose
```bash
make run
```

Generate seeds
```bash
make seeds
```

Access api doc:
```
http://127.0.0.1:8000/docs/index.html
```

## Tools
Install development tools
```bash
make install
```

Generate mock files
```bash
make mock
```

Unit tests
```bash
make test
```

Generate coverage
```bash
make coverage
```

Lint
```bash
make lint
```


## Help

Input currency format

Example
* 1000 -> R$10,00
* 10  -> R$0,10

If your system works with floating point:
```
value/100
```