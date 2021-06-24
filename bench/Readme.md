# Bench test for remote dialer

Test case is to compare the preformance between direct access and tunneling.

```golang
//======================
//
//                      (direct)
//         .--------------->----------------.
//        /              tunnel              \
// request--->server:8125------>client---->fileserver:8124
//
//
//======================
```

The test client will request following sizes file from file server in each request.

- 4kb
- 256kb
- 1Mb
- 4Mb
- 16Mb
- 256Mb

My test result is following:

Here is the result of buffer size 4k

| mod    | size  | content transfer | total        |
| ------ | ----- | ---------------- | ------------ |
| direct | 4kb   | 37.456µs         | 704.588µs    |
| direct | 256kb | 527.383µs        | 835.892µs    |
| direct | 1Mb   | 1.457393ms       | 2.398064ms   |
| direct | 4Mb   | 6.466144ms       | 7.261271ms   |
| direct | 16Mb  | 26.824295ms      | 27.973309ms  |
| direct | 256Mb | 453.623477ms     | 457.843783ms |
| proxy  | 4kb   | 32.583µs         | 1.264073ms   |
| proxy  | 256kb | 361.282µs        | 1.956524ms   |
| proxy  | 1Mb   | 1.982524ms       | 3.381918ms   |
| proxy  | 4Mb   | 9.284765ms       | 10.635328ms  |
| proxy  | 16Mb  | 40.229588ms      | 43.377169ms  |
| proxy  | 256Mb | 701.352388ms     | 720.673702ms |

Here is the result of buffer size 64k

| mod    | size  | content transfer | total        |
| ------ | ----- | ---------------- | ------------ |
| direct | 4kb   | 38.202µs         | 698.372µs    |
| direct | 256kb | 489.048µs        | 813.301µs    |
| direct | 1Mb   | 1.84549ms        | 2.67673ms    |
| direct | 4Mb   | 6.035001ms       | 6.860431ms   |
| direct | 16Mb  | 24.608875ms      | 25.655674ms  |
| direct | 256Mb | 478.944781ms     | 483.328069ms |
| proxy  | 4kb   | 26.409µs         | 3.055611ms   |
| proxy  | 256kb | 229.48µs         | 2.276647ms   |
| proxy  | 1Mb   | 1.348357ms       | 7.583335ms   |
| proxy  | 4Mb   | 10.90011ms       | 15.405417ms  |
| proxy  | 16Mb  | 34.601686ms      | 38.962934ms  |
| proxy  | 256Mb | 574.315839ms     | 595.014586ms |

> PS:
> Here is the log to table regexp `INFO\[[0-9]+\] mod: ([a-z]+), +size: ([0-9]+) +content transfer: ([0-9]+.[0-9]+[µ|m]s) +total: ([0-9]+.[0-9]+[µ|m]s)`
