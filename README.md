# GO-CPSV

## Document

please refer the link below:
- [https://github.com/ianchen0119/GO-CPSV/wiki](https://github.com/ianchen0119/GO-CPSV/wiki)

## Usage

### Run Example

```sh
$ git clone https://github.com/ianchen0119/GO-CPSV.git
$ cd GO-CPSV
$ go get github.com/ianchen0119/GO-CPSV/cpsv
$ go run main.go
```

### Go module

```sh
$ cd YOUR_PROJECT
$ go get github.com/ianchen0119/GO-CPSV/cpsv
```

in your program:
```go=
import (
  "fmt"
  "github.com/ianchen0119/GO-CPSV/cpsv"
  // ...
)

// ...
```

### Benchmark

- [test script](https://github.com/ianchen0119/GO-CPSV/tree/master/compose/test)


`sync.Map`:
```
SC-1  | sync.Map W: 12 R: 0 Times: 10000
```

single worker:
```
SC-1  | CPSV-W Start: 1670419919729 End: 1670419919730 Spent: 1 Times: 100
SC-1  | CPSV-R Start: 1670419922732 End: 1670419922780 Spent: 48 Times: 100
SC-1  | CPSV-W Start: 1670419922780 End: 1670419923372 Spent: 592 Times: 1000
SC-1  | CPSV-R Start: 1670419926374 End: 1670419926820 Spent: 446 Times: 1000
SC-1  | CPSV-W Start: 1670419926820 End: 1670419933789 Spent: 6969 Times: 10000
SC-1  | CPSV-R Start: 1670419936791 End: 1670419941260 Spent: 4469 Times: 10000
```

worker pool (3 worker):
```
SC-1  | CPSV-W Start: 1670420319096 End: 1670420319096 Spent: 0 Times: 100
SC-1  | CPSV-R Start: 1670420322096 End: 1670420322138 Spent: 42 Times: 100
SC-1  | CPSV-W Start: 1670420322138 End: 1670420322466 Spent: 328 Times: 1000
SC-1  | CPSV-R Start: 1670420325468 End: 1670420325892 Spent: 424 Times: 1000
SC-1  | CPSV-W Start: 1670420325892 End: 1670420329943 Spent: 4051 Times: 10000
SC-1  | CPSV-R Start: 1670420332945 End: 1670420337278 Spent: 4333 Times: 10000
```

worker pool (5 worker):
```
SC-1  | CPSV-W Start: 1670420456171 End: 1670420456171 Spent: 0 Times: 100
SC-1  | CPSV-R Start: 1670420459174 End: 1670420459215 Spent: 41 Times: 100
SC-1  | CPSV-W Start: 1670420459215 End: 1670420459489 Spent: 274 Times: 1000
SC-1  | CPSV-R Start: 1670420462491 End: 1670420462896 Spent: 405 Times: 1000
SC-1  | CPSV-W Start: 1670420462896 End: 1670420466656 Spent: 3760 Times: 10000
SC-1  | CPSV-R Start: 1670420469657 End: 1670420473802 Spent: 4145 Times: 10000
```

worker pool (10 worker):
```
SC-1  | CPSV-W Start: 1670420571694 End: 1670420571694 Spent: 0 Times: 100
SC-1  | CPSV-R Start: 1670420574697 End: 1670420574740 Spent: 43 Times: 100
SC-1  | CPSV-W Start: 1670420574740 End: 1670420574979 Spent: 239 Times: 1000
SC-1  | CPSV-R Start: 1670420577981 End: 1670420578369 Spent: 388 Times: 1000
SC-1  | CPSV-W Start: 1670420578370 End: 1670420581625 Spent: 3255 Times: 10000
SC-1  | CPSV-R Start: 1670420584626 End: 1670420588753 Spent: 4127 Times: 10000
```