1.  Single integer input
```rust:main.rs
use proconio::input;

fn main() {
    input! {
        n: i32,
    }
    println!("You entered: {}", n);
}
```


Input Exapmple:
```
42
```
Output Exapmple:
```
You entered: 42
```


2.  Multiple space-separated integer inputs

```rust:main.rs
use proconio::input;

fn main() {
    input! {
        a: i32,
        b: i32,
        c: i32,
    }
    println!("You entered: {} {} {}", a, b, c);
}
```

Input Exapmple:
```
10 20 30
```

Output Exapmple:
```
You entered: 10 20 30
```

3.  String input
```rust:main.rs
use proconio::input;

fn main() {
    input! {
        s: String,
    }
    println!("You entered: {}", s);
}
```

Input Exapmple:
```
hello
```

Output Exapmple:
```
You entered: hello
```

4.  Multiple space-separated string inputs

```rust:main.rs
use proconio::input;

fn main() {
    input! {
        s1: String,
        s2: String,
        s3: String,
    }
    println!("You entered: {}, {}, {}", s1, s2, s3);
}
```

入力例
```
hello world program
```

Output Exapmple:
```
You entered: hello, world, program
```


5.  Single-dimensional array (vector) input
```rust:main.rs
use proconio::input;

fn main() {
    input! {
        n: usize,
        vec: [i32; n],
    }
    println!("Vector: {:?}", vec);
}
```

Input Exapmple:
```
5
1 2 3 4 5
```

Output Exapmple:
```
Vector: [1, 2, 3, 4, 5]
```

6.  Two-dimensional array (vector of vectors) input
```rust:main.rs
use proconio::input;

fn main() {
    input! {
        h: usize,
        w: usize,
        matrix: [[i32; w]; h],
    }
    for row in matrix {
        println!("{:?}", row);
    }
}
```

Input Exapmple:
```
3 3
1 2 3
4 5 6
7 8 9

```

Output Exapmple:
```
[1, 2, 3]
[4, 5, 6]
[7, 8, 9]
```

7.  List of integers in a single line input

```rust:main.rs
use proconio::input;

fn main() {
    input! {
        n: usize,
        list: [i32; n],
    }
    println!("List: {:?}", list);
}
```


Input Exapmple:
```
4
10 20 30 40
```

Output Exapmple:
```
List: [10, 20, 30, 40]
```

8.  List of integers over multiple lines input
```rust:main.rs
use proconio::input;

fn main() {
    input! {
        n: usize,
        m: usize,
        lists: [[i32; m]; n],
    }
    for list in lists {
        println!("{:?}", list);
    }
}
```


Input Exapmple:
```
2 3
1 2 3
4 5 6
```

Output Exapmple:
```
[1, 2, 3]
[4, 5, 6]
```


9.  Tuple input
```rust:main.rs
use proconio::input;

fn main() {
    input! {
        pair: (i32, i32),
    }
    println!("Pair: {:?}", pair);
}
```


Input Exapmple:
```
100 200
```

Output Exapmple:
```
Pair: (100, 200)
```


10. Multiple types of inputs

```rust:main.rs
use proconio::input;

fn main() {
    input! {
        a: i32,
        b: f64,
        c: char,
    }
    println!("Integer: {}, Float: {}, Char: {}", a, b, c);
}
```

Input Exapmple:
```
42 3.14 A
```

Output Exapmple:
```
Integer: 42, Float: 3.14, Char: A
```


11. Input multiple lines of text
```rust:main.rs
use proconio::input;

fn main() {
    input! {
        n: usize,
        lines: [String; n],
    }
    for line in lines {
        println!("{}", line);
    }
}
```


Input Exapmple:
```
3
hello
world
rust
```


Output Exapmple:
```
hello
world
rust
```
