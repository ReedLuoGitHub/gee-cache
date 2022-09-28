package gee_cache

type Getter interface {
	Get(string) ([]byte, error)
}
