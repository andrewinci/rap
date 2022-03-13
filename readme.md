<h1 align="center">RAP</h2>
<h2 align="center">Random avro producer</h2>

RAP is an Avro generator that allows to fully control the data generation via configurations and with a built in Kafka producer.
## Configuration
Use a `.yaml` file to configure the avro generation. The format is the following:
```yaml
kafka:
  clusterEndpoint: exampleEndpoint # the kafka endpoint (can also be passed with an env variable like $KAFKA_ENDPOINT)
  security: sasl #one of: none, sasl, mtls is expected here
  # the schema registry is optional and only required if schemaName is used in a producer
  schemaRegistry:
    # the schema registry endpoint (can also be passed with an env variable like $SR_ENDPOINT)
    endpoint: schemaRegistryEndpoint 
    username: $SR_USER
    password: $SR_PASSWORD
    
  sasl:
    # retrieve username and password from an env variable
    username: $KAFKA_USERNAME 
    password: $KAFKA_PASSWORD
  # mtls: todo

producers:
  - name: producer1 # an identifier for this producer
    numberOfMessages: 2000  # number of messages to generate from this producer
    topic: mytest-topic
    avro:
      schema: 
        id:     # the id registered in the schema registry
        raw: {} # the avro schema in json format
      generationRules: # set of rules to configure the generation of specific fields
        key: keyGen # special generation rule used to generate the record key
        .Name: nameGen 
        .SubRecord.Email: emailGen
        .SubRecord.Age: ageGen
      generators: # map of generators that can be used in the rules 
        nameGen: "{string}[Mary|James|Patricia|Robert]{1}"
        emailGen: "{string}[a-z]{10}[@]{1}[a-z]{10}[.org|.com]{1}"
        ageGen: "{int}[2|3|4|5|6]{1}[0-9]{1}"
        keyGen: "{string}[a-Z | 0-9]{10}"
  - name: producer2
    numberOfMessages: 20000
    topic: mytest-topic
    avro:
        schemaName: test-value # retrieve the schema from schema registry
        generationRules: # set of rules to configure the generation of specific fields
            key: keyGen # special generation rule used to generate the record key
        generators: # map of generators that can be used in the rules 
            keyGen: "{string}[0-9]{10}"
...
```
**NOTE:** Only the Kafka configurations can be passed via env variables
 
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
- a function: `uuid()`, `timestamp_ms()`
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

**Generate a random v4 uuid**  
`{string}[uuid()]{1}`

## Development

Run tests with `go test ./...`

## TODO:
- [ ] support union in field gen
- [ ] support logical types in field gen
- [ ] support split yaml file