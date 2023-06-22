package terraformValueObjects

// Design struct to handle data output. Should be defined globally, not just for gcp. Likely:
// "provider": {
//     "divisionName": {
//          "resourceType.resourceName": {
//             "creator":  {
//                 actor: "AUTH ACCOUNT",
//				   timestamp: "XYZ"
//             },
//             "modifier": {}
//         },
//          "resourceName2": {},
//          "resourceName3": {}
//      }
// }

// Timestamp is a string that represents the time of a cloud actor's action.
type Timestamp string

// CloudActor an entity that make changes to a cloud environment.
type CloudActor string

// CloudActorTimeStamp is a struct containing a cloud actor and a timestamp
// for one of their actions.
type CloudActorTimeStamp struct {

	// Actor is an entity that make changes to a cloud environment.
	Actor CloudActor

	// Timestamp is a string that represents the time of a cloud actor's action.
	Timestamp Timestamp
}

// ResourceActions is a struct containing the cloud actor and timestamp of the most recent (if any)
// resource modification.
type ResourceActions struct {

	// Creator is the cloud actor and timestamp of the resource creation.
	Creator CloudActorTimeStamp

	// Modifier is the cloud actor and timestamp of the most recent (if any) resource modification.
	Modifier CloudActorTimeStamp
}

// DivisionResourceActions is a mapping between a division name and a map between resource
// names and resource actions. Captures all resource actions for a particular division.
type DivisionResourceActions map[Division]map[ResourceName]ResourceActions

// ProviderResourceActions is a mapping between a provider name and DivisionResourceActions.
// Captures all resource actions for a particular provider.
type ProviderResourceActions map[Provider]DivisionResourceActions
