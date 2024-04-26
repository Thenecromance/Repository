package Repository

type Option func(*Repository)

func (o Option) apply(r *Repository) {
	o(r)
}

func WithAlgorithm(algId int) Option {
	return func(r *Repository) {
		//r.contentHashAlg = applyAlgorithm(algId)

		r.hash = NewHash(algId)
	}
}

func WithDirLength(length int) Option {
	return func(r *Repository) {
		r.dirLen = length
	}
}

func WithResourceDir(dir string) Option {
	return func(r *Repository) {
		r.rootDir = dir
	}
}

//func test() {
//	data, err := syscall.Mmap( /*int(f.Fd())*/ 0, 0 /*int(size)*/, 0, syscall.PROT_READ, syscall.MAP_SHARED)
//}
