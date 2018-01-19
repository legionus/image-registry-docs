func (t tagService) Get(ctx context.Context, tag string) (distribution.Descriptor, error) {
	imageStream, err := t.imageStream.imageStreamGetter.get()
	if err != nil {
		return distribution.Descriptor{}, distribution.ErrRepositoryUnknown{Name: t.imageStream.Reference()}
	}
	te := imageapiv1.LatestTaggedImage(imageStream, tag)
	...
	dgst, err := digest.ParseDigest(te.Image)
	...

	if !t.pullthroughEnabled {
		image, err := t.imageStream.getImage(ctx, dgst)
		if err != nil {
			return distribution.Descriptor{}, err
		}

		if !isImageManaged(image) {
			return distribution.Descriptor{}, distribution.ErrTagUnknown{Tag: tag}
		}
	}
	return distribution.Descriptor{Digest: dgst}, nil
}
