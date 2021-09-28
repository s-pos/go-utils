# Go Utils

Utils atau _Private Package_ ini digunakan untuk penggunaan middleware, standard logging output, atau kebutuhan-kebutuhan lainnya yg dipergunakan di setiap service.

## How To Use
Pertama-tama, Anda harus melakukan config pada git global anda terlebih dahulu dengan cara
```bash
$ git config -- url."https://${username_github}:${access_token}@github.com".insteadOf "https://github.com"
```

Setelah itu, anda harus melakukan config untuk mengijinkan golang mengakses library dari private repo dengan cara
```bash
 nano ~/.profile
```
lalu tambahkan config ini
```bash
 export GOSUMDB=off
```
simpan file, lalu reload profile dengan cara
```bash
 source ~/.profile
```

Setelah selesai melakukan config pada langkah sebelumnya, Anda bisa menggunakan _Private Package_ ini dengan cara menyimpan module kedalam project anda dengan cara
```bash
$ go get github.com/s-pos/go-utils@master
```

Anda akan mendapatkan versi stable dari _Private Package_ tersebut.

Setelah itu anda bisa melakukannya dengan cara import package-package yg ada pada _Private Package_ tersebut dengan sesuai kebutuhan.