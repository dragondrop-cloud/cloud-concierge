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
	Actor CloudActor `json:"actor"`

	// Timestamp is a string that represents the time of a cloud actor's action.
	Timestamp Timestamp `json:"timestamp"`
}

// ResourceActions is a struct containing the cloud actor and timestamp of the most recent (if any)
// resource modification.
type ResourceActions struct {

	// Creator is the cloud actor and timestamp of the resource creation.
	Creator *CloudActorTimeStamp `json:"creation,omitempty"`

	// Modifier is the cloud actor and timestamp of the most recent (if any) resource modification.
	Modifier *CloudActorTimeStamp `json:"modified,omitempty"`
}

// ResourceActionMap is a mapping between a resource name and resource actions.
type ResourceActionMap map[ResourceName]*ResourceActions
