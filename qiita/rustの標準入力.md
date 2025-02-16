1.  ​単一の整数入力
```rust:main.rs
use proconio::input;

fn main() {
    input! {
        n: i32,
    }
    println!("You entered: {}", n);
}
```


入力例:
```
42
```
出力例:
```
You entered: 42
```


2.  ​複数のスペース区切り整数入力
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

入力例:
```
10 20 30
```

出力例:
```
You entered: 10 20 30
```

3.  ​文字列入力
```rust:main.rs
use proconio::input;

fn main() {
    input! {
        s: String,
    }
    println!("You entered: {}", s);
}
```

入力例:
```
hello
```

出力例:
```
You entered: hello
```

4.  ​複数のスペース区切り文字列入力
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

出力例:
```
You entered: hello, world, program
```


5.  ​1次元配列（ベクタ）入力
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

入力例:
```
5
1 2 3 4 5
```

出力例:
```
Vector: [1, 2, 3, 4, 5]
```

6.  ​2次元配列（ベクタのベクタ）入力
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

入力例:
```
3 3
1 2 3
4 5 6
7 8 9

```

出力例:
```
[1, 2, 3]
[4, 5, 6]
[7, 8, 9]
```

7.  ​単一行に含まれる整数のリスト入力
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


入力例:
```
4
10 20 30 40
```

出力例:
```
List: [10, 20, 30, 40]
```

8.  ​複数行の整数のリスト入力
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


入力例:
```
2 3
1 2 3
4 5 6
```

出力例:
```
[1, 2, 3]
[4, 5, 6]
```


9.  ​タプル入力
```rust:main.rs
use proconio::input;

fn main() {
    input! {
        pair: (i32, i32),
    }
    println!("Pair: {:?}", pair);
}
```


入力例:
```
100 200
```

出力例:
```
Pair: (100, 200)
```


10. ​複数の異なる型の入力
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

入力例:
```
42 3.14 A
```

出力例:
```
Integer: 42, Float: 3.14, Char: A
```


11. ​複数行テキスト入力
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


入力例:
```
3
hello
world
rust
```


出力例:
```
hello
world
rust
```