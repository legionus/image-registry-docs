func (r *repository) Blobs(ctx context.Context) distribution.BlobStore {
	bs := r.Repository.Blobs(ctx)

	if r.app.quotaEnforcing.enforcementEnabled {
		bs = &quotaRestrictedBlobStore{BlobStore: bs}
	}

	if r.app.config.Pullthrough.Enabled {
		bs = &pullthroughBlobStore{BlobStore: bs}
	}

	bs = newPendingErrorsBlobStore(bs, r)

	if audit.LoggerExists(ctx) {
		bs = audit.NewBlobStore(ctx, bs)
	}

	if r.app.config.Metrics.Enabled {
		bs = metrics.NewBlobStore(bs, r.Named().Name())
	}

	return bs
}
