## CSV Parser written in Go

```bash
go build main.go
./main -fp test.csv
```

Output
```bash
2024/04/24 22:43:17 COLUMNS [{id int} {name string} {age int} {address string}]
2024/04/24 22:43:17 ROWS [[1 bob 19 London, United Kingdom] [2 tom 18 Paris, France] [3 john 22 London, United Kingdom]]
2024/04/24 22:43:17 NAMES [bob tom john]
2024/04/24 22:43:17 ADDRESSES [London, United Kingdom Paris, France London, United Kingdom]
```
