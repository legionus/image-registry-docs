func (r *repository) Manifests(ctx context.Context, options ...distribution.ManifestServiceOption) (distribution.ManifestService, error) {
	opts := append(options, registrystorage.SkipLayerVerification())
	ms, err := r.Repository.Manifests(ctx, opts...)
	if err != nil {
		return nil, err
	}

	ms = &manifestService{
		manifests:     ms,
		blobStore:     r.Blobs(ctx),
	}

	if r.app.config.Pullthrough.Enabled {
		ms = &pullthroughManifestService{ManifestService: ms}
	}

	ms = newPendingErrorsManifestService(ms, r)

	...

	return ms, nil
}
