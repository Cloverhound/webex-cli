package search

// DefaultFields maps each query type to its default GraphQL field selection.
var DefaultFields = map[QueryType]string{
	QueryTask: `tasks {
      id
      status
      channelType
      direction
      isActive
      createdTime
      endedTime
      origin
      destination
      totalDuration
      connectedDuration
      queueDuration
      holdDuration
      wrapupDuration
      ringingDuration
      selfserviceDuration
      owner { id name }
      lastEntryPoint { id name }
      lastQueue { id name }
      lastTeam { id name }
      lastWrapupCodeName
      aggregation { name value }
      intervalStartTime
    }
    pageInfo { hasNextPage endCursor }`,

	QueryTaskDetails: `tasks {
      id
      status
      channelType
      direction
      isActive
      createdTime
      endedTime
      origin
      destination
      totalDuration
      connectedDuration
      queueDuration
      holdDuration
      wrapupDuration
      ringingDuration
      selfserviceDuration
      owner { id name }
      lastEntryPoint { id name }
      lastQueue { id name }
      lastTeam { id name }
      lastWrapupCodeName
      csatScore
      terminationType
      contributors { id name }
      aggregation { name value }
      intervalStartTime
    }
    pageInfo { hasNextPage endCursor }`,

	QueryTaskLegDetails: `tasks {
      id
      taskId
      status
      channelType
      direction
      createdTime
      endedTime
      isActive
      origin
      destination
      entryPoint { id name }
      queue { id name }
      owner { id name }
      site { id name }
      team { id name }
      callLegType
      handleType
      isTaskLegHandled
      isWithinServiceLevel
      connectedDuration
      holdDuration
      queueDuration
      wrapupDuration
      ringingDuration
      aggregation { name value }
      intervalStartTime
    }
    pageInfo { hasNextPage endCursor }`,

	QueryAgentSession: `agentSessions {
      agentSessionId
      agentId
      agentName
      userLoginId
      isActive
      state
      startTime
      endTime
      siteId
      siteName
      teamId
      teamName
      channelInfo {
        channelType
        totalDuration
        idleDuration
        idleCount
        availableDuration
        availableCount
        connectedDuration
        connectedCount
        holdDuration
        holdCount
        wrapupDuration
        wrapupCount
        notRespondedCount
      }
      aggregation { name value }
      intervalStartTime
    }
    pageInfo { hasNextPage endCursor }`,

	QueryFlowInteractions: `flowInteractions {
      interactionId
      flowId
      flowVersionId
      entryPointName
      lastExecutedActivity
      flowStartTime
      flowEndTime
      status
    }
    pageInfo { hasNextPage totalPages totalHits }`,

	QueryFlowTraceEvents: `flowTraceEvents {
      interactionId
      flowId
      activityName
      activityProcessId
      outcome
      activityInputs { name value }
      activityOutput { name type value }
      modifiedFlowVariables { name type value }
      flowStartTime
      flowEndTime
    }
    pageInfo { hasNextPage totalPages totalHits }`,
}
