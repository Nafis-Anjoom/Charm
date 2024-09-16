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
