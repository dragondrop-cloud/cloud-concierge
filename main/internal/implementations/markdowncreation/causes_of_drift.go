package markdowncreation

import (
	"fmt"
	"sort"

	"github.com/atsushinee/go-markdown-generator/doc"
)

// CloudActor represents a cloud actor name
type CloudActor string

// CloudAction represents a cloud action (modified, creation)
type CloudAction string

// ActorActionCount represents the count of an actor action
type ActorActionCount struct {
	Actor  CloudActor
	Action CloudAction
	Count  int
}

// setRootCausesOfDriftData sets the root causes of drift data in the markdown report
func (m *MarkdownCreator) setRootCausesOfDriftData(report *doc.MarkDownDoc) {
	report.Write("# Root Causes of Drift").Writeln().Writeln()
	report.Write("## Cloud Actors Causing Changes").Writeln().Writeln()

	if len(m.resourcesToCloudActions) == 0 {
		report.Write("No identified Cloud Actor actions.").Writeln().Writeln()
		return
	}

	actorActions := map[CloudActor]map[CloudAction]int{}
	for _, cloudActions := range m.resourcesToCloudActions {
		key := "Create Resource"
		detail := cloudActions["creation"]

		if detail == (CloudActionDetail{}) {
			key = "Modify Resource"
			detail = cloudActions["modified"]
		}

		if actorActions[CloudActor(detail.Actor)] == nil {
			actorActions[CloudActor(detail.Actor)] = map[CloudAction]int{}
		}
		actorActions[CloudActor(detail.Actor)][CloudAction(key)]++
	}

	var actorActionCounts []ActorActionCount
	for actor, actions := range actorActions {
		for action, count := range actions {
			actorActionCounts = append(actorActionCounts, ActorActionCount{actor, action, count})
		}
	}

	sort.Slice(actorActionCounts, func(i, j int) bool {
		return actorActionCounts[i].Count > actorActionCounts[j].Count
	})

	report.Write("|Actor|Action|Count|\n| :---: | :---: | :---: |").Writeln()

	for _, actorActionCount := range actorActionCounts {
		report.Write(fmt.Sprintf("|%s", string(actorActionCount.Actor)))
		report.Write(fmt.Sprintf("|%s", string(actorActionCount.Action)))
		report.Write(fmt.Sprintf("|%d|", actorActionCount.Count)).Writeln()
	}

	report.Writeln()
}
