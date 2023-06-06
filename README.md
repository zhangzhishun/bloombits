# Bloom Filter Test

## Build

```shell
# 生成测试地址
go build -o genaddress ./cmd/gen_address.go
# 测试布隆过滤器
go build -o bloomtest ./cmd/test_bloom.go
```

# Generate Address

```shell
./genaddress ${生成地址数量}
```

example:

```shell
./genaddress 1000000000
```

# Test Bloom

```shell
./bloomtest ${布隆过滤器长度} ${布隆过滤器哈希次数} ${导入布隆过滤器的地址文件路径} ${测试假阳率数据量}
```

example:

```shell
./bloomtest 1000000000 20 ./address_test/address_50w.txt 1000000
./bloomtest 1000000000 20 ./address_test/address_100w.txt 1000000
```