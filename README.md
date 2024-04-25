## CSV Parser written in Go

```bash
go build main.go
./main -fp test.csv
```

Output
```bash
2024/04/24 21:48:53 COLUMNS [{id int} {name string} {age int}]
2024/04/24 21:48:53 ROWS [[1 bob 19] [2 tom 18] [3 john 22]]
2024/04/24 21:48:53 NAMES [bob tom john]
```
