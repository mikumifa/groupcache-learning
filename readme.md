

```Go
// A Getter loads data for a key.
type Getter interface {
	Get(key string) ([]byte, error)
}

// A GetterFunc implements Getter with a function.
type GetterFunc func(key string) ([]byte, error)

// Get implements Getter interface function
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}
```

这样，既能够将普通的函数类型（需类型转换）作为参数，也可以将结构体作为参数，使用更为灵活，可读性也更好，这就是接口型函数的价值。

## 实现原理

有一个hash函数CH

对于一个节点A， 为这个节点提供peers（伙伴），A和peers具有同样的功能（注册方法，hostname->int）

对于一个索要请求的参数 `/<basepath>/<groupname>/<key>`  前缀都是统一的， 也就是要 `groupname`和`key`从而返回这个key的值

groupname的特点和计组里面的cache类似，没有删除，只有拿取之后的缓存

group的作用就是用来标记cache，一个机器可能会有多个分别是cache，给分布式cahche加一个namespace可以区分

### 处理流程

启用一个API server用于用户感知， 也可以理解成cache的一个调入的口

- 第一步，在本地查找 mainCache里面
- 本地找不到
  - 通过CH计算key的hash值，一致性hash，所有的分布式节点计算的hash都是一样的，然后找到这个key应该存的分布式节点的hostname
  - 然后远程去调用这个hostname来得到key，：没找到就说明真的没有， 然后就会自己去Getter

> 实现，如果key算出来的hostname就是自己则会Getter



> **缓存雪崩**：缓存在同一时刻全部失效，造成瞬时DB请求量大、压力骤增，引起雪崩。缓存雪崩通常因为缓存服务器宕机、缓存的 key 设置了相同的过期时间等引起。

> **缓存击穿**：一个存在的key，在缓存过期的一刻，同时有大量的请求，这些请求都会击穿到 DB ，造成瞬时DB请求量大、压力骤增。
>
> 缓存过期那一瞬间，DB查找并没有存到Cache，所以会一直去读取数据库	

//对于同一个读取请求，保证统一时间只有一次

//也就是互斥锁

> **缓存穿透**：查询一个不存在的数据，因为不存在则不会写到缓存中，所以每次都会去请求 DB，如果瞬间流量过大，穿透到 DB，导致宕机。

```
singleFlight实现

type call struct {
	wg sync.WaitGroup
	//存储返回值
	val interface{}
	err error
}

type Group struct {
	mu        sync.Mutex       // protects keyToCall
	keyToCall map[string]*call // lazily initialized
}
```

