func (pbs *pullthroughBlobStore) Stat(...) (distribution.Descriptor, error) {
	desc, err := pbs.BlobStore.Stat(ctx, dgst)
	switch {
	case err == distribution.ErrBlobUnknown:
	case err != nil:
		context.GetLogger(ctx).Errorf("Failed to find blob %q: %#v", dgst.String(), err)
		fallthrough
	default:
		return desc, err
	}

	return pbs.remoteBlobGetter.Stat(ctx, dgst)
}
