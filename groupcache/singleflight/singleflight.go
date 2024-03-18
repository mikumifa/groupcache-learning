package singleflight

import "sync"

// call is an in-flight or completed Do call
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

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.keyToCall == nil {
		g.keyToCall = make(map[string]*call)
	}
	if myCall, ok := g.keyToCall[key]; ok {
		g.mu.Unlock()
		myCall.wg.Wait()
		return myCall.val, myCall.err
	}
	myCall := new(call)
	myCall.wg.Add(1)
	g.keyToCall[key] = myCall
	g.mu.Unlock()

	myCall.val, myCall.err = fn()
	myCall.wg.Done()

	g.mu.Lock()
	delete(g.keyToCall, key)
	g.mu.Unlock()

	return myCall.val, myCall.err
}
