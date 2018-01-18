func (bs *blobDescriptorService) Stat(ctx context.Context, dgst digest.Digest) (distribution.Descriptor, error) {
	desc, err := bs.BlobDescriptorService.Stat(ctx, dgst)
	if err == nil {
		return desc, nil
	}

	desc, err = bs.repo.app.BlobStatter().Stat(ctx, dgst)
	if err == nil {
		if !isEmptyDigest(dgst) {
			if !imageStreamHasBlob(ctx, bs.repo.imageStream, dgst, !bs.repo.app.config.Pullthrough.Enabled) {
				return distribution.Descriptor{}, distribution.ErrBlobUnknown
			}
		}
		return desc, nil
	}

	if err == distribution.ErrBlobUnknown && remoteBlobAccessCheckEnabledFrom(ctx) {
		desc, err = bs.repo.remoteBlobGetter.Stat(ctx, dgst)
	}
	return desc, err
}
