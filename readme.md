## Polybius Bifid Cipher 

#### Comes with a built-in non-pseudorandom polybius generator. The polybius generator only uses a latin alphabet where "J" is exchanged by "I" for all occurrences. Non-pseudo-randomness comes from generating 26 numbers in the range of 0-25. Using the random integer sequence, each integer used to index the alphabet and get the letter to build up the matrix. The generator service used is [https://www.random.org/](https://www.random.org/). 

#### The project comes with 10 randomly generated polybius squares in the generated_polybius_squares directory

### Encrypt/Decrypt text:
#### Run the executable with a path to your polybius file:
`./build/bifid.sh -pb=path/to/a/polybius/square.json`

#### An encrypt operation is used by prefixing the input with a "+"
#### A decrypt operation is used by prefixing the input with a "-"

```
| ./build/bifid.sh -pb=./generated_polybius_squares/polybius_1603628969.json 
| ------------
| input  > +HELLO
| 🔐 Encrypting message...
| output > XHZSO
| ------------ 
```

#### To Decrypt a message simply change the input prefix to "-"
```
| ------------
| input  > -XHZSO
| 🔓 Decrypting message...
| output > HELLO
```

#### Generate a polybius matrix:
#### Note: you need an API key for [https://www.random.org/](https://www.random.org/) in order to generate the random integer sequence
```
./build/bifid.sh -gen -apiKey=<your_api_key> 
```
#### Output
```
------------
🆗 Created new random polybius at path: ./generated_polybius_squares/polybius_1603628969.json
Re-run the program with the new polybius file
------------
```
