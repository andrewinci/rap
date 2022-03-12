<h1 align="center">RAP</h2>
<h2 align="center">Random avro producer</h2>

RAP is an Avro generator that allows to fully control the data generation via configurations and with a built in Kafka producer.
## Configuration
Use a `.yaml` file to configure the avro generation. The format is the following:
```yaml
kafka:
  clusterEndpoint: exampleEndpoint #the kafka endpoint

producers:
  - name: producer1 # an identifier for this producer
    numberOfMessages: 2000  # number of messages to generate from this producer
    avro:
      schema: {} # the avro schema in json format
      generationRules: # set of rules to configure the generation of specific fields
        .Name: nameGen 
        .SubRecord.Email: emailGen
      generators: # map of generators that can be used in the rules 
        nameGen: "{string}[Mary|James|Patricia|Robert]{1}"
        emailGen: "{string}[a-z]{10}[@]{1}[a-z]{10}[.org|.com]{1}"
        ageGen: "{int}[2|3|4|5|6]{1}[0-9]{1}"
  - name: producer2
...
```

## Field generation syntax
To customize the generation of the fields it is possible to provide a pattern.
The generic structure of a data gen pattern is:
```
{type}[content-restrictions]{count}[content-restrictions2]{count2}....
```
Where `type` can be one of the following avro type:  
`boolean` `int` `long` `float` `double` `bytes` `string`

`content-restriction` can be:
- an interval: `a-z`, `A-Z`, `a-Z`, `0-9`
- a constant value: `testvalue`
- a combination of intervals and constants: `a-z | 0-9 | test`

The field `count` tells the generator how many times the generation should be performed accordingly to the `content-restriction`. The result of each generation is concatenated.

### Examples

**Generate a constant value**

`{string}[test]{1}`: will always generate the constant value test

**Generate alphanumeric string**

`{string}[a-Z|0-9]{10}`: will generate a random alphanumeric string of length 10

**Generate a random email**

`{string}[a-z]{10}[@]{1}[a-z]{10}[.org|.com]{1}`

**Generate a random number**
`{int}[0-9]{5}`: will generate a random number of 5 digits.

## Development

Use the `make` to perform common operations on this project.

Run tests with `make test`

## TODO:
- [ ] support union in field gen
- [ ] support logical types in field gen
- [ ] support split yaml file