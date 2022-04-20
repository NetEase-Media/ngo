package httplib

var defaultHttpClient = newDefaultHttpClient()

func newDefaultHttpClient() *HttpClient {
	opt := NewDefaultOptions()
	client, _ := New(opt)
	return client
}

func DefaultHttpClient() *HttpClient {
	return defaultHttpClient
}

func SetDefaultHttpClient(client *HttpClient) {
	defaultHttpClient = client
}

// Get 调用默认http客户端的GET方法
func Get(url string) *DataFlow {
	return defaultHttpClient.Get(url)
}

// Post 调用默认http客户端的POST方法
func Post(url string) *DataFlow {
	return defaultHttpClient.Post(url)
}

// Put 调用默认http客户端的PUT方法
func Put(url string) *DataFlow {
	return defaultHttpClient.Put(url)
}

// Delete 调用默认http客户端的DELETE方法
func Delete(url string) *DataFlow {
	return defaultHttpClient.Delete(url)
}

// Patch 调用默认http客户端的PATCH方法
func Patch(url string) *DataFlow {
	return defaultHttpClient.Patch(url)
}

func Close() {
	defaultHttpClient.Close()
}
