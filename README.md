# Charm
My very own tree-walker interpreter built with Go.

### Motivation
For a long time, I have wondered how code is transformed into an executable program. Specifically, I was curious about how an interpreter/compiler gives meaning to code, 
executes it, and reports any error or feedback to the user. I undertook this project to demystify these processes. By building my own interpreter, I have applied many core 
concepts I have learned during my degree. I have explored the complexities in lexical analysis, grammar, error handling/reporting, custom language features, and so much more.
I came away with a much deeper understanding of the tradeoffs between performance, memory, code complexity, and usability.

### Installation
1. Clone the repository and navigate to the project directory.
 ```bash
 git https://github.com/Nafis-Anjoom/Charm.git
 cd charm
   ``` 
2. Install [Go][go]. If already installed, make sure it's up-to-date.
3. To start the interpreter:
  ```bash
  # to run the interpeter directly from source code and code in the repl
  go run main.go
  
  # to execute a charm source file
  go run main.go sourcefile.ch
  
  # to build and run via an executable binary
  go build
  ./charm
  ./charm sourcefile.ch
  ```


### Syntax
```python
# Variables and Data Types
name = "Alice";       # String
age = 18;             # Integer
height = 5.6;         # Float
isStudent = true;     # Boolean

print("Name:", name);
print("Age:", age);
print("Height:", height);
print("Is Student:", isStudent);

# Conditional Statements
if (age >= 18) {
    print(name + " is an adult.");
} else {
    print(name + " is not an adult.");
}

# Functions
func greet(personName) {
    return "hello, " + personName + "!";
}

greet = func(personName) {
  return "hello, " + personName + "!";
};

# Closure
func generatePrinter(personAge) {
  age = 99;

  return func() {
    print("this person is", age, "years old.");
  };

}

# Calling the function
greeting = greet(name);
print(greeting);

agePrinter = generatePrinter(99);
agePrinter();

# Loops (while loop)
print("Counting down from 5:");
count = 5;
while (count > 0) {
    print(count);
    count = count - 1;
}

# Lists (similar to Python lists)
fruits = ["apple", "banana", "cherry"];
print("Fruit list:", fruits);

# Adding an element to the list
push(fruits, "orange");
print("After adding orange:", fruits);

# HashMap
map = {"hello": "world", 1: greet, true: age};
map["hello"];
delete(map, 1);
mapKeys = keys(map);
```

### References
1. [Writing an Interpreter in GO][ball] by Thorsten Ball
2. [Crafting Interpreters][nystorm] by Robert Nystorm

[ball]: https://interpreterbook.com/
[nystorm]: https://craftinginterpreters.com/
[go]: https://go.dev/doc/install
