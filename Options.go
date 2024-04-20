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
		r.DirLength = length
	}
}

func WithResourceDir(dir string) Option {
	return func(r *Repository) {
		r.ResourceDir = dir
	}
}
