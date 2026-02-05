package wpml

import (
	"fmt"
)

const DefaultAuthor = "DJI WPML SDK"

func ConvertWaylinesToWPMLMission(waylines *Waylines) (*WPMLMission, error) {
	waylines.ApplyDefaults()

	if err := waylines.Validate(); err != nil {
		return nil, fmt.Errorf(ErrWaylinesValidationFailed, err)
	}

	mission := NewWPMLMission()
	mission.SetAuthor(DefaultAuthor)
	mission.UpdateTimestamp()
	missionConfig, err := convertToMissionConfig(waylines)
	if err != nil {
		return nil, fmt.Errorf(ErrConvertMissionConfig, err)
	}
	mission.SetMissionConfig(*missionConfig)

	templateFolder, err := convertToTemplateFolder(waylines)
	if err != nil {
		return nil, fmt.Errorf(ErrConvertTemplateFolder, err)
	}
	mission.Template.Document.Folders = []TemplateFolder{*templateFolder}

	waylineFolder, err := convertToWaylineFolder(waylines)
	if err != nil {
		return nil, fmt.Errorf(ErrConvertWaylineFolder, err)
	}
	mission.Waylines.Document.Folders = []WaylineFolder{*waylineFolder}

	return mission, nil
}

func convertToMissionConfig(waylines *Waylines) (*MissionConfig, error) {
	flyToWaylineMode := FlightModeSafely
	finishAction := FinishActionGoHome
	if waylines.FinishAction != "" {
		finishAction = waylines.FinishAction
	}

	exitOnRCLost := RCLostActionGoContinue
	executeRCLostAction := ExecuteRCLostActionGoBack

	takeOffSecurityHeight := waylines.TakeOffSecurityHeight
	if takeOffSecurityHeight == 0 {
		takeOffSecurityHeight = 20
	}

	globalTransitionalSpeed := waylines.GlobalTransitionalSpeed
	if globalTransitionalSpeed == 0 {
		globalTransitionalSpeed = 6
	}

	globalRTHHeight := 100.0
	if waylines.GlobalRTHHeight > 0 {
		globalRTHHeight = waylines.GlobalRTHHeight
	} else if waylines.TakeOffSecurityHeight > 0 {
		globalRTHHeight = waylines.TakeOffSecurityHeight
	}

	droneInfo := DroneInfo{
		DroneEnumValue:    int(waylines.DroneModel),
		DroneSubEnumValue: 0,
	}

	payloadInfo := PayloadInfo{
		PayloadEnumValue:     int(waylines.PayloadModel),
		PayloadPositionIndex: int(waylines.PayloadPositionIndex),
	}

	var takeOffRefPoint *string
	var takeOffRefPointAGLHeight *float64

	if waylines.TakeOffRefPointLatitude != 0 && waylines.TakeOffRefPointLongitude != 0 {
		takeOffPointStr := fmt.Sprintf("%.6f,%.6f,%.1f",
			waylines.TakeOffRefPointLatitude,
			waylines.TakeOffRefPointLongitude,
			waylines.TakeOffRefPointHeight)
		takeOffRefPoint = &takeOffPointStr

		if waylines.TakeOffRefPointAGLHeight != nil {
			takeOffRefPointAGLHeight = waylines.TakeOffRefPointAGLHeight
		}
	} else if len(waylines.Waypoints) > 0 {
		firstWaypoint := waylines.Waypoints[0]
		takeOffPointStr := fmt.Sprintf("%.6f,%.6f,%.1f",
			firstWaypoint.Latitude,
			firstWaypoint.Longitude,
			firstWaypoint.Height)
		takeOffRefPoint = &takeOffPointStr

		aglHeight := 0.0
		takeOffRefPointAGLHeight = &aglHeight
	}

	return &MissionConfig{
		FlyToWaylineMode:         flyToWaylineMode,
		FinishAction:             finishAction,
		ExitOnRCLost:             exitOnRCLost,
		ExecuteRCLostAction:      &executeRCLostAction,
		TakeOffSecurityHeight:    takeOffSecurityHeight,
		TakeOffRefPoint:          takeOffRefPoint,
		TakeOffRefPointAGLHeight: takeOffRefPointAGLHeight,
		GlobalTransitionalSpeed:  globalTransitionalSpeed,
		GlobalRTHHeight:          &globalRTHHeight,
		DroneInfo:                droneInfo,
		PayloadInfo:              payloadInfo,
	}, nil
}

func convertToTemplateFolder(waylines *Waylines) (*TemplateFolder, error) {
	heightMode := HeightModeRelativeToStartPoint

	globalShootHeight := waylines.GlobalHeight
	surfaceRelativeHeight := waylines.GlobalHeight
	positioningType := PositioningTypeGPS

	waylineCoordSysParam := &WaylineCoordinateSysParam{
		CoordinateMode:          CoordinateModeWGS84,
		HeightMode:              heightMode,
		GlobalShootHeight:       &globalShootHeight,
		PositioningType:         &positioningType,
		SurfaceFollowModeEnable: intPtr(0),
		SurfaceRelativeHeight:   &surfaceRelativeHeight,
	}

	placemarks := make([]Placemark, len(waylines.Waypoints))
	for i, wp := range waylines.Waypoints {
		placemark, err := convertToTemplatePlacemark(wp, i)
		if err != nil {
			return nil, fmt.Errorf(ErrConvertWaypoint, i, err)
		}
		placemarks[i] = *placemark
	}

	return &TemplateFolder{
		TemplateType:               waylines.TemplateType,
		TemplateID:                 0,
		AutoFlightSpeed:            waylines.GlobalSpeed,
		GlobalHeight:               &waylines.GlobalHeight,
		WaylineCoordinateSysParam:  waylineCoordSysParam,
		GimbalPitchMode:            stringPtr(waylines.GimbalPitchMode),
		GlobalWaypointHeadingParam: convertGlobalHeadingParam(waylines),
		Placemarks:                 placemarks,
	}, nil
}

func convertToWaylineFolder(waylines *Waylines) (*WaylineFolder, error) {
	executeHeightMode := ExecuteHeightModeRelativeToStartPoint
	if waylines.HeightType == HeightModeRealTimeFollowSurface {
		executeHeightMode = ExecuteHeightModeRealTimeFollowSurface
	}
	placemarks := make([]Placemark, 0, len(waylines.Waypoints))
	for i, wp := range waylines.Waypoints {
		placemark, err := convertToWaylinePlacemark(wp, i, waylines)
		if err != nil {
			return nil, fmt.Errorf(ErrConvertWaypoint, i, err)
		}
		placemarks = append(placemarks, *placemark)
	}

	return &WaylineFolder{
		TemplateID:        0,
		WaylineID:         0,
		AutoFlightSpeed:   waylines.GlobalSpeed,
		ExecuteHeightMode: executeHeightMode,
		Distance:          waylines.Distance,
		Duration:          waylines.Duration,
		Placemarks:        placemarks,
	}, nil
}

func convertToTemplatePlacemark(waypoint WaylinesWaypoint, index int) (*Placemark, error) {
	point := &Point{
		Coordinates: formatCoordinates(waypoint.Longitude, waypoint.Latitude),
	}

	ellipsoidHeight := waypoint.Height
	height := waypoint.Height

	var waypointSpeed *float64
	if waypoint.Speed > 0 {
		waypointSpeed = &waypoint.Speed
	}

	useGlobalHeight := 0
	if waypoint.Height == 0 {
		useGlobalHeight = 1
	}

	useGlobalSpeed := 0
	if waypoint.Speed == 0 {
		useGlobalSpeed = 1
	}

	headingParam := &WaypointHeadingParam{
		WaypointHeadingMode:        HeadingModeFollowWayline,
		WaypointHeadingAngle:       float64Ptr(0),
		WaypointPoiPoint:           stringPtr("0.000000,0.000000,0.000000"),
		WaypointHeadingAngleEnable: intPtr(0),
		WaypointHeadingPathMode:    "followBadArc",
		WaypointHeadingPoiIndex:    intPtr(0),
	}

	turnParam := &WaypointTurnParam{
		WaypointTurnMode:        TurnModeToPointAndStopWithContinuityCurvature,
		WaypointTurnDampingDist: float64Ptr(0.2),
	}

	var actionGroups []ActionGroup
	if len(waypoint.Actions) > 0 {
		actionGroup := convertToActionGroup(waypoint.Actions, waypoint.TriggerType, index)
		if actionGroup != nil {
			actionGroups = append(actionGroups, *actionGroup)
		}
	}

	return &Placemark{
		Point:                 point,
		Index:                 index,
		EllipsoidHeight:       &ellipsoidHeight,
		Height:                &height,
		UseGlobalHeight:       &useGlobalHeight,
		UseGlobalSpeed:        &useGlobalSpeed,
		WaypointSpeed:         waypointSpeed,
		WaypointHeadingParam:  headingParam,
		WaypointTurnParam:     turnParam,
		UseGlobalHeadingParam: intPtr(1),
		UseGlobalTurnParam:    intPtr(1),
		UseStraightLine:       intPtr(1),
		ActionGroups:          actionGroups,
		IsRisky:               intPtr(0),
		WaypointWorkType:      intPtr(0),
	}, nil
}

func convertToWaylinePlacemark(waypoint WaylinesWaypoint, index int, waylines *Waylines) (*Placemark, error) {
	point := &Point{
		Coordinates: formatCoordinates(waypoint.Longitude, waypoint.Latitude),
	}

	executeHeight := waypoint.Height

	speed := waypoint.Speed
	if speed == 0 {
		speed = waylines.GlobalSpeed
	}

	headingParam := &WaypointHeadingParam{
		WaypointHeadingMode:        HeadingModeFollowWayline,
		WaypointHeadingAngle:       float64Ptr(0),
		WaypointPoiPoint:           stringPtr("0.000000,0.000000,0.000000"),
		WaypointHeadingAngleEnable: intPtr(0),
		WaypointHeadingPathMode:    "followBadArc",
		WaypointHeadingPoiIndex:    intPtr(0),
	}

	gimbalHeadingParam := &WaypointGimbalHeadingParam{
		WaypointGimbalPitchAngle: float64Ptr(0),
		WaypointGimbalYawAngle:   float64Ptr(0),
	}

	turnParam := createWaypointTurnParam(waypoint, waylines)

	useStraightLine := getUseStraightLine(waypoint, waylines, turnParam.WaypointTurnMode)

	var actionGroups []ActionGroup
	if len(waypoint.Actions) > 0 {
		actionGroup := convertToActionGroup(waypoint.Actions, waypoint.TriggerType, index)
		if actionGroup != nil {
			actionGroups = append(actionGroups, *actionGroup)
		}
	}

	return &Placemark{
		Point:                      point,
		Index:                      index,
		ExecuteHeight:              &executeHeight,
		WaypointSpeed:              &speed,
		WaypointHeadingParam:       headingParam,
		WaypointTurnParam:          turnParam,
		UseStraightLine:            useStraightLine,
		ActionGroups:               actionGroups,
		WaypointGimbalHeadingParam: gimbalHeadingParam,
		IsRisky:                    intPtr(0),
		WaypointWorkType:           intPtr(0),
	}, nil
}

func convertToActionGroup(actions []ActionRequest, triggerType string, waypointIndex int) *ActionGroup {
	if len(actions) == 0 {
		return nil
	}

	triggerTypeStr := TriggerTypeReachPoint
	switch triggerType {
	case "reachPoint":
		triggerTypeStr = TriggerTypeReachPoint
	case "passPoint":
		triggerTypeStr = TriggerTypePassPoint
	case "manual":
		triggerTypeStr = TriggerTypeManual
	case "betweenAdjacentPoints":
		triggerTypeStr = TriggerTypeBetweenAdjacentPoints
	case "multipleTiming":
		triggerTypeStr = TriggerTypeMultipleTiming
	case "multipleDistance":
		triggerTypeStr = TriggerTypeMultipleDistance
	}

	trigger := ActionTrigger{
		ActionTriggerType: triggerTypeStr,
	}

	actionList := make([]Action, len(actions))
	for i, actionReq := range actions {
		actionList[i] = Action{
			ActionID:                i,
			ActionActuatorFunc:      actionReq.Type,
			ActionActuatorFuncParam: convertActionParams(actionReq),
		}
	}

	return &ActionGroup{
		ActionGroupID:         waypointIndex,
		ActionGroupStartIndex: waypointIndex,
		ActionGroupEndIndex:   waypointIndex,
		ActionGroupMode:       ActionGroupModeSequence,
		ActionTrigger:         trigger,
		Actions:               actionList,
	}
}

func convertActionParams(actionReq ActionRequest) *ActionActuatorFuncParam {
	param := &ActionActuatorFuncParam{}

	switch actionReq.Type {
	case ActionTypeTakePhoto:
		if takePhotoAction, ok := actionReq.Action.(*TakePhotoAction); ok {
			posIndex := int(takePhotoAction.PayloadPositionIndex)
			param.PayloadPositionIndex = &posIndex
			param.PayloadLensIndex = takePhotoAction.PayloadLensIndex
			param.FileSuffix = &takePhotoAction.FileSuffix

			useGlobal := 0
			if takePhotoAction.UseGlobalPayloadLensIndex {
				useGlobal = 1
			}
			param.UseGlobalPayloadLensIndex = &useGlobal
		}
	case ActionTypeStartRecord:
		if startRecordAction, ok := actionReq.Action.(*StartRecordAction); ok {
			posIndex := int(startRecordAction.PayloadPositionIndex)
			param.PayloadPositionIndex = &posIndex
			param.FileSuffix = &startRecordAction.FileSuffix
			param.PayloadLensIndex = startRecordAction.PayloadLensIndex

			useGlobal := 0
			if startRecordAction.UseGlobalPayloadLensIndex {
				useGlobal = 1
			}
			param.UseGlobalPayloadLensIndex = &useGlobal
		}
	case ActionTypeStopRecord:
		if stopRecordAction, ok := actionReq.Action.(*StopRecordAction); ok {
			posIndex := int(stopRecordAction.PayloadPositionIndex)
			param.PayloadPositionIndex = &posIndex
			param.PayloadLensIndex = stopRecordAction.PayloadLensIndex
		}
	case ActionTypeHover:
		if hoverAction, ok := actionReq.Action.(*HoverAction); ok {
			param.HoverTime = &hoverAction.HoverTime
		}
	case ActionTypeZoom:
		if zoomAction, ok := actionReq.Action.(*ZoomAction); ok {
			posIndex := int(zoomAction.PayloadPositionIndex)
			param.PayloadPositionIndex = &posIndex
			param.FocalLength = &zoomAction.FocalLength
		}
	case ActionTypeFocus:
		if focusAction, ok := actionReq.Action.(*FocusAction); ok {
			posIndex := int(focusAction.PayloadPositionIndex)
			param.PayloadPositionIndex = &posIndex

			pointFocus := 0
			if focusAction.IsPointFocus {
				pointFocus = 1
			}
			param.IsPointFocus = &pointFocus

			param.FocusX = &focusAction.FocusX
			param.FocusY = &focusAction.FocusY

			infiniteFocus := 0
			if focusAction.IsInfiniteFocus {
				infiniteFocus = 1
			}
			param.IsInfiniteFocus = &infiniteFocus

			param.FocusRegionWidth = focusAction.FocusRegionWidth
			param.FocusRegionHeight = focusAction.FocusRegionHeight
		}
	case ActionTypeGimbalRotate:
		if gimbalRotateAction, ok := actionReq.Action.(*GimbalRotateAction); ok {
			posIndex := int(gimbalRotateAction.PayloadPositionIndex)
			param.PayloadPositionIndex = &posIndex
			param.GimbalHeadingYawBase = &gimbalRotateAction.GimbalHeadingYawBase
			param.GimbalRotateMode = &gimbalRotateAction.GimbalRotateMode

			pitchEnable := 0
			if gimbalRotateAction.GimbalPitchRotateEnable {
				pitchEnable = 1
			}
			param.GimbalPitchRotateEnable = &pitchEnable
			param.GimbalPitchRotateAngle = &gimbalRotateAction.GimbalPitchRotateAngle

			rollEnable := 0
			if gimbalRotateAction.GimbalRollRotateEnable {
				rollEnable = 1
			}
			param.GimbalRollRotateEnable = &rollEnable
			param.GimbalRollRotateAngle = &gimbalRotateAction.GimbalRollRotateAngle

			yawEnable := 0
			if gimbalRotateAction.GimbalYawRotateEnable {
				yawEnable = 1
			}
			param.GimbalYawRotateEnable = &yawEnable
			param.GimbalYawRotateAngle = &gimbalRotateAction.GimbalYawRotateAngle

			timeEnable := 0
			if gimbalRotateAction.GimbalRotateTimeEnable {
				timeEnable = 1
			}
			param.GimbalRotateTimeEnable = &timeEnable
			param.GimbalRotateTime = &gimbalRotateAction.GimbalRotateTime
		}
	case ActionTypeRotateYaw:
		if rotateYawAction, ok := actionReq.Action.(*RotateYawAction); ok {
			param.AircraftHeading = &rotateYawAction.AircraftHeading
			param.AircraftPathMode = rotateYawAction.AircraftPathMode
		}
	case ActionTypeCustomDirName:
		if customDirNameAction, ok := actionReq.Action.(*CustomDirNameAction); ok {
			posIndex := int(customDirNameAction.PayloadPositionIndex)
			param.PayloadPositionIndex = &posIndex
			param.DirectoryName = &customDirNameAction.DirectoryName
		}
	case ActionTypeGimbalEvenlyRotate:
		if gimbalEvenlyRotateAction, ok := actionReq.Action.(*GimbalEvenlyRotateAction); ok {
			posIndex := int(gimbalEvenlyRotateAction.PayloadPositionIndex)
			param.PayloadPositionIndex = &posIndex
			param.GimbalPitchRotateAngle = &gimbalEvenlyRotateAction.GimbalPitchRotateAngle
		}
	case ActionTypeAccurateShoot:
		if accurateShootAction, ok := actionReq.Action.(*AccurateShootAction); ok {
			posIndex := int(accurateShootAction.PayloadPositionIndex)
			param.PayloadPositionIndex = &posIndex
			param.GimbalPitchRotateAngle = &accurateShootAction.GimbalPitchRotateAngle
			param.GimbalYawRotateAngle = &accurateShootAction.GimbalYawRotateAngle
			param.FocusX = float64Ptr(float64(accurateShootAction.FocusX))
			param.FocusY = float64Ptr(float64(accurateShootAction.FocusY))
			param.FocusRegionWidth = float64Ptr(float64(accurateShootAction.FocusRegionWidth))
			param.FocusRegionHeight = float64Ptr(float64(accurateShootAction.FocusRegionHeight))
			param.FocalLength = &accurateShootAction.FocalLength
			param.AircraftHeading = &accurateShootAction.AircraftHeading
			param.TargetAngle = &accurateShootAction.TargetAngle
			param.ImageWidth = &accurateShootAction.ImageWidth
			param.ImageHeight = &accurateShootAction.ImageHeight
			param.AFPos = &accurateShootAction.AFPos
			param.GimbalPort = &accurateShootAction.GimbalPort
			param.PayloadLensIndex = accurateShootAction.PayloadLensIndex

			accurateFrameValid := 0
			if accurateShootAction.AccurateFrameValid {
				accurateFrameValid = 1
			}
			param.AccurateFrameValid = &accurateFrameValid

			useGlobal := 0
			if accurateShootAction.UseGlobalPayloadLensIndex {
				useGlobal = 1
			}
			param.UseGlobalPayloadLensIndex = &useGlobal
		}
	case ActionTypeOrientedShoot:
		if orientedShootAction, ok := actionReq.Action.(*OrientedShootAction); ok {
			posIndex := int(orientedShootAction.PayloadPositionIndex)
			param.PayloadPositionIndex = &posIndex
			param.GimbalPitchRotateAngle = &orientedShootAction.GimbalPitchRotateAngle
			param.GimbalYawRotateAngle = &orientedShootAction.GimbalYawRotateAngle
			param.FocusX = float64Ptr(float64(orientedShootAction.FocusX))
			param.FocusY = float64Ptr(float64(orientedShootAction.FocusY))
			param.FocusRegionWidth = float64Ptr(float64(orientedShootAction.FocusRegionWidth))
			param.FocusRegionHeight = float64Ptr(float64(orientedShootAction.FocusRegionHeight))
			param.FocalLength = &orientedShootAction.FocalLength
			param.AircraftHeading = &orientedShootAction.AircraftHeading
			param.TargetAngle = &orientedShootAction.TargetAngle
			param.ActionUUID = &orientedShootAction.ActionUUID
			param.ImageWidth = &orientedShootAction.ImageWidth
			param.ImageHeight = &orientedShootAction.ImageHeight
			param.AFPos = &orientedShootAction.AFPos
			param.GimbalPort = &orientedShootAction.GimbalPort
			param.OrientedCameraType = &orientedShootAction.OrientedCameraType
			param.OrientedFilePath = &orientedShootAction.OrientedFilePath
			param.OrientedFileMD5 = &orientedShootAction.OrientedFileMD5
			param.OrientedFileSize = &orientedShootAction.OrientedFileSize
			param.OrientedFileSuffix = &orientedShootAction.OrientedFileSuffix
			param.OrientedCameraApertue = &orientedShootAction.OrientedCameraApertue
			param.OrientedCameraLuminance = &orientedShootAction.OrientedCameraLuminance
			param.OrientedCameraShutterTime = &orientedShootAction.OrientedCameraShutterTime
			param.OrientedCameraISO = &orientedShootAction.OrientedCameraISO
			param.OrientedPhotoMode = &orientedShootAction.OrientedPhotoMode
			param.PayloadLensIndex = orientedShootAction.PayloadLensIndex

			accurateFrameValid := 0
			if orientedShootAction.AccurateFrameValid {
				accurateFrameValid = 1
			}
			param.AccurateFrameValid = &accurateFrameValid

			useGlobal := 0
			if orientedShootAction.UseGlobalPayloadLensIndex {
				useGlobal = 1
			}
			param.UseGlobalPayloadLensIndex = &useGlobal
		}
	case ActionTypePanoShot:
		if panoShotAction, ok := actionReq.Action.(*PanoShotAction); ok {
			posIndex := int(panoShotAction.PayloadPositionIndex)
			param.PayloadPositionIndex = &posIndex
			param.PanoShotSubMode = &panoShotAction.PanoShotSubMode
			param.PayloadLensIndex = panoShotAction.PayloadLensIndex

			useGlobal := 0
			if panoShotAction.UseGlobalPayloadLensIndex {
				useGlobal = 1
			}
			param.UseGlobalPayloadLensIndex = &useGlobal
		}
	case ActionTypeRecordPointCloud:
		if recordPointCloudAction, ok := actionReq.Action.(*RecordPointCloudAction); ok {
			posIndex := int(recordPointCloudAction.PayloadPositionIndex)
			param.PayloadPositionIndex = &posIndex
			param.RecordPointCloudOperate = &recordPointCloudAction.RecordPointCloudOperate
		}
	case ActionTypeGimbalAngleLock:
		if gimbalAngleLockAction, ok := actionReq.Action.(*GimbalAngleLockAction); ok {
			posIndex := int(gimbalAngleLockAction.PayloadPositionIndex)
			param.PayloadPositionIndex = &posIndex
		}
	case ActionTypeGimbalAngleUnlock:
		if gimbalAngleUnlockAction, ok := actionReq.Action.(*GimbalAngleUnlockAction); ok {
			posIndex := int(gimbalAngleUnlockAction.PayloadPositionIndex)
			param.PayloadPositionIndex = &posIndex
		}
	case ActionTypeStartSmartOblique:
		if startSmartObliqueAction, ok := actionReq.Action.(*StartSmartObliqueAction); ok {
			posIndex := int(startSmartObliqueAction.PayloadPositionIndex)
			param.PayloadPositionIndex = &posIndex
		}
	case ActionTypeStartTimeLapse:
		if startTimeLapseAction, ok := actionReq.Action.(*StartTimeLapseAction); ok {
			posIndex := int(startTimeLapseAction.PayloadPositionIndex)
			param.PayloadPositionIndex = &posIndex
		}
	case ActionTypeStopTimeLapse:
		if stopTimeLapseAction, ok := actionReq.Action.(*StopTimeLapseAction); ok {
			posIndex := int(stopTimeLapseAction.PayloadPositionIndex)
			param.PayloadPositionIndex = &posIndex
		}
	case ActionTypeSetFocusType:
		if setFocusTypeAction, ok := actionReq.Action.(*SetFocusTypeAction); ok {
			posIndex := int(setFocusTypeAction.PayloadPositionIndex)
			param.PayloadPositionIndex = &posIndex
		}
	case ActionTypeTargetDetection:
		if targetDetectionAction, ok := actionReq.Action.(*TargetDetectionAction); ok {
			posIndex := int(targetDetectionAction.PayloadPositionIndex)
			param.PayloadPositionIndex = &posIndex
		}
	}

	return param
}

func convertGlobalHeadingParam(waylines *Waylines) *GlobalWaypointHeadingParam {
	headingMode := HeadingModeFollowWayline
	if waylines.AircraftYawMode != "" {
		switch waylines.AircraftYawMode {
		case "followWayline":
			headingMode = HeadingModeFollowWayline
		case "followRoute":
			headingMode = HeadingModeFollowWayline
		case "manual":
			headingMode = HeadingModeManually
		case "free":
			headingMode = HeadingModeFree
		}
	}

	return &GlobalWaypointHeadingParam{
		WaypointHeadingMode: headingMode,
	}
}

func intPtr(v int) *int {
	return &v
}

func stringPtr(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

func float64Ptr(v float64) *float64 {
	return &v
}

func formatCoordinates(lon, lat float64) string {
	return fmt.Sprintf("%g,%g", lon, lat)
}

func createWaypointTurnParam(waypoint WaylinesWaypoint, waylines *Waylines) *WaypointTurnParam {
	turnMode := TurnModeToPointAndStopWithDiscontinuityCurvature

	if waypoint.WaypointTurnMode != "" {
		turnMode = waypoint.WaypointTurnMode
	} else if waylines.GlobalWaypointTurnMode != "" {
		turnMode = waylines.GlobalWaypointTurnMode
	}

	dampingDist := 0.0
	if waypoint.TurnDampingDist > 0 {
		dampingDist = waypoint.TurnDampingDist
	} else if waylines.GlobalTurnDampingDist > 0 {
		dampingDist = waylines.GlobalTurnDampingDist
	}

	return &WaypointTurnParam{
		WaypointTurnMode:        turnMode,
		WaypointTurnDampingDist: &dampingDist,
	}
}

func getUseStraightLine(waypoint WaylinesWaypoint, waylines *Waylines, turnMode string) *int {
	if turnMode != TurnModeCoordinateTurn &&
		turnMode != TurnModeToPointAndStopWithContinuityCurvature &&
		turnMode != TurnModeToPointAndPassWithContinuityCurvature &&
		turnMode != TurnModeToPointAndStopWithDiscontinuityCurvature {
		return nil
	}

	var useStraightLine *bool
	if waypoint.UseStraightLine != nil {
		useStraightLine = waypoint.UseStraightLine
	} else if waylines.GlobalUseStraightLine != nil {
		useStraightLine = waylines.GlobalUseStraightLine
	}

	if useStraightLine != nil {
		if *useStraightLine {
			return intPtr(1)
		} else {
			return intPtr(0)
		}
	}

	return intPtr(1)
}
