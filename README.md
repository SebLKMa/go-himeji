# go-himeji
My interpreter and my compiler for my programming language ***himeji***.  

![himeji](./images/himeji.jpg) 

## Compiler
Build the compiler:  
```sh
cd cmd/compiler
make build
```

Compile a source file to a bytecodes file:  
```sh
cd cmd/himeji

./compiler codes.txt
Source codes:
21 + 21

i:0, width:2, offset:1
instruction: [0 0 1]
instruction: [1]
instruction: [1]
175 bytes written to codes.bin
```

## Runtime
Build the runtime virtual machine:  
```sh
cd cmd/runtime
make build
```

Run the bytecodes file:  
```sh
cd cmd/himeji

./runtime codes.bin
Read 175 bytes from codes.bin
Deserialized bytecode: &{Instructions:0000 OpConstant 0
0003 OpConstant 1
0006 OpAdd
 Constants:[0x400000e2c0 0x400000e2c8]}
Result: 42
```

## References

Inspired by the books:  
https://interpreterbook.com/  
https://compilerbook.com/ 

https://interpreterbook.com/waiig_code_1.3.zip  
https://compilerbook.com/waiig_code_1.2.zip  

## Reference Implementation

https://github.com/SaladinoBelisario/Compiler_Go  

