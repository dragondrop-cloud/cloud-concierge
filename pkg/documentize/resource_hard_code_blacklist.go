package documentize

// newResourceBlackList returns a set of blacklisted resources. A resource is only
// added to the list if it is known that Terraformer outputs broken configuration for that resource.
// This list presents an opportunity of fixes to be implemented for Terraformer. Ideally the length would
// be 0.
func newResourceBlackList() map[string]bool {
	return map[string]bool{
		"google_storage_bucket_iam_binding": true,
		"google_storage_bucket_object_acl":  true,
		"google_storage_default_object_acl": true,
		"google_storage_bucket_iam_policy":  true,
		"google_storage_bucket_acl":         true,
	}
}
