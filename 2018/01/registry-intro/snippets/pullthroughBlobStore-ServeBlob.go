func (pbs *pullthroughBlobStore) ServeBlob(...) error {
	err := pbs.BlobStore.ServeBlob(ctx, w, req, dgst)
	switch {
	case err == distribution.ErrBlobUnknown:
	case err != nil:
		context.GetLogger(ctx).Errorf("Failed to find blob %q: %#v", dgst.String(), err)
		fallthrough
	default:
		return err
	}
	if pbs.mirror {
		if _, ok := inflight[dgst]; ok {
			context.GetLogger(ctx).Infof("Serving %q while mirroring in background", dgst)
			_, err := copyContent(ctx, pbs.remoteBlobGetter, dgst, w, req)
			return err
		}
		inflight[dgst] = struct{}{}

		storeLocalInBackground(ctx, pbs.imageStream, pbs.writeLimiter, pbs.BlobStore, dgst)
	}
	_, err = copyContent(ctx, pbs.remoteBlobGetter, dgst, w, req)
	return err
}
