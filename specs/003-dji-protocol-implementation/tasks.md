# Tasks: DJI Protocol Implementation

**Input**: Design documents from `/specs/003-dji-protocol-implementation/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/
**Depends On**: 001-project-setup (✅), 002-protocol-adapter-design (✅)

**Tests**: TDD approach required - tests written first, must FAIL before implementation. Coverage ≥ 80%.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1-US9)
- Include exact file paths in descriptions

## Path Conventions

- **Go project**: `pkg/`, `cmd/`, `internal/` at repository root
- **Tests**: `*_test.go` files alongside implementation

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization - already completed in 001/002

- [x] T001 Project structure exists from 001-project-setup
- [x] T002 Go 1.22 project with dependencies initialized
- [x] T003 [P] Linting and formatting tools configured
- [x] T004 [P] DJI adapter skeleton exists from 002-protocol-adapter-design

**Note**: Setup tasks completed in previous features.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### Tests for Foundational

- [x] T005 [P] Unit test for Handler interface in pkg/adapter/dji/handler/handler_test.go
- [x] T006 [P] Unit test for Router interface in pkg/adapter/dji/router/router_test.go
- [x] T007 [P] Unit test for config constants in pkg/adapter/dji/config/config_test.go

### Implementation for Foundational

- [x] T008 [P] Create Handler interface in pkg/adapter/dji/handler/handler.go
  - `Handle(ctx context.Context, msg *Message, topic *TopicInfo) (*StandardMessage, error)`
  - `GetTopicType() TopicType`

- [x] T009 [P] Create Router types in pkg/adapter/dji/router/router.go
  - `ErrMethodNotFound` error
  - `ErrMethodAlreadyRegistered` error
  - ServiceRouter and EventRouter in separate files

- [x] T010 [P] Create config constants in pkg/adapter/dji/config/config.go
  - `SERVICE_CALL_TIMEOUT = 30s`
  - `DRC_HEARTBEAT_TIMEOUT = 3s`
  - `UNKNOWN_DEVICE_POLICY = "discard"`

- [x] T011 Create HandlerRegistry in pkg/adapter/dji/handler/registry.go
  - `Register(handler Handler)`
  - `Get(topicType TopicType) Handler`
  - Auto-dispatch based on topic type

**Checkpoint**: Foundation ready - user story implementation can now begin

---

## Phase 3: User Story 1 - 设备遥测数据上报 (Priority: P1) 🎯 MVP

**Goal**: 完整解析 OSD 遥测数据，集成到现有 adapter

**Independent Test**: 发送 OSD 消息，验证能解析出完整的 AircraftOSD/DockOSD/RCOSD 结构

### Tests for User Story 1

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T012 [P] [US1] Unit test for OSDHandler in pkg/adapter/dji/handler/osd_handler_test.go
  - Test AircraftOSD parsing (80+ fields)
  - Test DockOSD parsing
  - Test RCOSD parsing
  - Test nested device data extraction

- [x] T013 [P] [US1] Unit test for OSDParser in pkg/adapter/dji/integration/osd_parser_test.go
  - Test ParseAircraftOSD with full data
  - Test ParseAircraftOSD with partial data (pointer fields)
  - Test ParseDockOSD
  - Test ParseRCOSD

- [ ] T014 [P] [US1] Integration test for OSD flow in tests/integration/dji_osd_test.go
  - Mock RabbitMQ message with OSD payload
  - Verify StandardMessage output

### Implementation for User Story 1

- [x] T015 [US1] Create OSDParser in pkg/adapter/dji/integration/osd_parser.go
  - `ParseAircraftOSD(data json.RawMessage) (*aircraft.AircraftOSD, error)`
  - `ParseDockOSD(data json.RawMessage) (*aircraft.DockOSD, error)`
  - `ParseRCOSD(data json.RawMessage) (*aircraft.RCOSD, error)`
  - Handle partial updates (all pointer fields)

- [x] T016 [US1] Create OSDHandler in pkg/adapter/dji/handler/osd_handler.go
  - Implement Handler interface
  - Use OSDParser for data extraction
  - Extract device topology from OSD
  - Convert to StandardMessage with full data

- [x] T017 [US1] Update adapter.go to use OSDHandler in pkg/adapter/dji/adapter.go
  - Add MessageHandler and HandlerRegistry interfaces
  - Add SetHandlerRegistry method
  - Add HandleMessage method for handler-based processing
  - Create pkg/adapter/dji/init package for handler initialization

- [x] T018 [P] [US1] Add OSD mock messages in tests/mocks/dji_osd_messages.go
  - Full AircraftOSD sample
  - Full DockOSD sample
  - Full RCOSD sample
  - Partial update samples

**Checkpoint**: User Story 1 complete - OSD 完整解析可用

---

## Phase 4: User Story 2 & 3 - 属性变化与设备状态 (Priority: P1)

**Goal**: 完整处理 State 和 Status 消息

**Independent Test**: 发送 State/Status 消息，验证属性变化和设备拓扑解析

### Tests for User Story 2 & 3

- [x] T019 [P] [US2] Unit test for StateHandler in pkg/adapter/dji/handler/state_handler_test.go
  - Test property change parsing
  - Test multiple properties change
  - Test device topology in state

- [x] T020 [P] [US3] Unit test for StatusHandler in pkg/adapter/dji/handler/status_handler_test.go
  - Test online status parsing
  - Test offline status parsing
  - Test device topology extraction

### Implementation for User Story 2 & 3

- [x] T021 [US2] Create StateHandler in pkg/adapter/dji/handler/state_handler.go
  - Implement Handler interface
  - Parse property changes
  - Convert to StandardMessage

- [x] T022 [US3] Create StatusHandler in pkg/adapter/dji/handler/status_handler.go
  - Implement Handler interface
  - Parse online/offline status
  - Extract device topology (gateway → aircraft → payloads)
  - Convert to StandardMessage

- [x] T023 [P] [US2] Add State/Status mock messages in tests/mocks/dji_state_messages.go

**Checkpoint**: User Stories 2 & 3 complete - State/Status 处理完成

---

## Phase 5: User Story 4 - 服务调用下发 (Priority: P1)

**Goal**: 实现服务调用的双向转换和路由

**Independent Test**: 发送服务调用请求，验证转换为 DJI 格式并能接收响应

### Tests for User Story 4

- [x] T024 [P] [US4] Unit test for ServiceRouter in pkg/adapter/dji/router/service_router_test.go
  - Test method registration
  - Test routing to correct handler
  - Test unknown method handling

- [x] T025 [P] [US4] Unit test for ServiceHandler in pkg/adapter/dji/handler/service_handler_test.go
  - Test services request parsing
  - Test services_reply parsing
  - Test timeout handling (30s)

- [x] T026 [P] [US4] Unit test for device commands in pkg/adapter/dji/router/device_commands_test.go
  - Test cover_open/close
  - Test drone_open/close
  - Test device_reboot

### Implementation for User Story 4

- [x] T027 [US4] Create ServiceRouter in pkg/adapter/dji/router/service_router.go
  - Implement Router interface
  - Method → Handler mapping
  - Support 50+ service methods

- [x] T028 [US4] Create ServiceHandler in pkg/adapter/dji/handler/service_handler.go
  - Implement Handler interface
  - Use ServiceRouter for method dispatch
  - Handle services and services_reply topics
  - Support timeout (30s default)

- [x] T029 [US4] Register device commands in pkg/adapter/dji/router/device_commands.go
  - `cover_open`, `cover_close`, `cover_force_close`
  - `drone_open`, `drone_close`
  - `charge_open`, `charge_close`
  - `device_reboot`, `device_format`, `drone_format`
  - `debug_mode_open`, `debug_mode_close`
  - `battery_maintenance_switch`
  - `air_conditioner_mode_switch`
  - `alarm_state_switch`
  - `sdr_workmode_switch`

- [x] T030 [P] [US4] Add Services mock messages in tests/mocks/dji_service_messages.go

**Checkpoint**: User Story 4 complete - 基础服务调用可用

---

## Phase 6: User Story 5 - 事件上报与确认 (Priority: P1)

**Goal**: 实现事件处理的双向转换

**Independent Test**: 发送事件消息，验证解析和回复生成

### Tests for User Story 5

- [x] T031 [P] [US5] Unit test for EventRouter in pkg/adapter/dji/router/event_router_test.go
  - Test event registration
  - Test routing to correct handler
  - Test need_reply handling

- [x] T032 [P] [US5] Unit test for EventHandler in pkg/adapter/dji/handler/event_handler_test.go
  - Test events parsing
  - Test events_reply generation
  - Test HMS event parsing

### Implementation for User Story 5

- [x] T033 [US5] Create EventRouter in pkg/adapter/dji/router/event_router.go
  - Implement Router interface
  - Method → Handler mapping
  - Support need_reply flag

- [x] T034 [US5] Create EventHandler in pkg/adapter/dji/handler/event_handler.go
  - Implement Handler interface
  - Use EventRouter for method dispatch
  - Handle events and events_reply topics
  - Generate reply for need_reply events

- [x] T035 [US5] Register core events in pkg/adapter/dji/router/core_events.go
  - `device_exit_homing_notify`
  - `device_temp_ntfy_need_clear`
  - `file_upload_callback`
  - `hms` (HMS 健康管理)

- [x] T036 [P] [US5] Add Events mock messages in tests/mocks/dji_event_messages.go

**Checkpoint**: User Story 5 complete - P1 核心协议完成

---

## Phase 7: User Story 6 - 航线任务管理 (Priority: P2)

**Goal**: 实现航线任务管理

**Independent Test**: 发送航线任务命令，验证任务创建/执行/进度上报

### Tests for User Story 6

- [x] T037 [P] [US6] Unit test for wayline commands in pkg/adapter/dji/router/wayline_commands_test.go
  - Test flighttask_create
  - Test flighttask_prepare
  - Test flighttask_execute
  - Test flighttask_pause/recovery/undo
  - Test return_home/return_home_cancel

- [x] T038 [P] [US6] Unit test for wayline events in pkg/adapter/dji/router/wayline_events_test.go
  - Test flighttask_progress
  - Test flighttask_ready
  - Test return_home_info

### Implementation for User Story 6

- [x] T039 [US6] Register wayline commands in pkg/adapter/dji/router/wayline_commands.go
  - `flighttask_create` with CreateData
  - `flighttask_prepare` with PrepareData
  - `flighttask_execute` with ExecuteData
  - `flighttask_pause`, `flighttask_recovery`
  - `flighttask_undo` with UndoData
  - `return_home`, `return_home_cancel`

- [x] T040 [US6] Register wayline events in pkg/adapter/dji/router/wayline_events.go
  - `flighttask_progress` with FlightTaskProgressData
  - `flighttask_ready` with FlightTaskReadyData (need_reply)
  - `return_home_info` with ReturnHomeInfoData

- [x] T041 [P] [US6] Add Wayline mock messages in tests/mocks/dji_wayline_messages.go

**Checkpoint**: User Story 6 complete - 航线任务管理可用

---

## Phase 8: User Story 7 - 相机控制 (Priority: P2)

**Goal**: 实现相机控制

**Independent Test**: 发送相机控制命令，验证拍照/录像/云台控制

### Tests for User Story 7

- [x] T042 [P] [US7] Unit test for camera commands in pkg/adapter/dji/router/camera_commands_test.go
  - Test camera_mode_switch
  - Test camera_photo_take
  - Test camera_recording_start/stop
  - Test gimbal_reset

### Implementation for User Story 7

- [x] T043 [US7] Register camera commands in pkg/adapter/dji/router/camera_commands.go
  - `camera_mode_switch` with CameraModeSwitchData
  - `camera_photo_take` with CameraPhotoTakeData
  - `camera_recording_start`, `camera_recording_stop` with CameraRecordingData
  - `camera_aim` with CameraAimData
  - `camera_focal_length_set` with CameraFocalLengthData
  - `gimbal_reset` with GimbalResetData
  - IR metering commands

- [x] T044 [P] [US7] Add Camera mock messages in tests/mocks/dji_camera_messages.go

**Checkpoint**: User Story 7 complete - 相机控制可用

---

## Phase 9: User Story 4 Extension - 配置管理 (Priority: P2)

**Goal**: 实现设备配置管理和设备请求处理

**Independent Test**: 发送配置读取/设置命令，验证配置管理功能

### Tests for User Story 4 Extension

- [x] T045 [P] [US4] Unit test for config commands in pkg/adapter/dji/router/config_commands_test.go
- [x] T046 [P] [US4] Unit test for RequestHandler in pkg/adapter/dji/handler/request_handler_test.go

### Implementation for User Story 4 Extension

- [x] T047 [US4] Register config commands in pkg/adapter/dji/router/config_commands.go
  - `config_get` with ConfigGetData
  - `config_set` with ConfigSetData

- [x] T048 [US4] Create RequestHandler in pkg/adapter/dji/handler/request_handler.go
  - Handle device-initiated requests
  - Support requests and requests_reply topics

**Checkpoint**: P2 业务功能完成

---

## Phase 10: User Story 8 - 实时控制 DRC (Priority: P3)

**Goal**: 实现 DRC 实时控制

**Independent Test**: 建立 DRC 连接，发送虚拟摇杆指令，验证心跳超时

### Tests for User Story 8

- [x] T049 [P] [US8] Unit test for DRC commands in pkg/adapter/dji/router/drc_commands_test.go
  - Test drc_mode_enter/exit
  - Test drone_control
  - Test heart (heartbeat)
  - Test drone_emergency_stop

- [x] T050 [P] [US8] Unit test for DRC events in pkg/adapter/dji/router/drc_events_test.go
  - Test joystick_invalid_notify
  - Test drc_status_notify

- [x] T051 [P] [US8] Unit test for DRCHandler in pkg/adapter/dji/handler/drc_handler_test.go
  - Test heartbeat timeout (3s)
  - Test connection state management

### Implementation for User Story 8

- [x] T052 [US8] Create DRCHandler in pkg/adapter/dji/handler/drc_handler.go
  - Handle drc/up and drc/down topics
  - Heartbeat management (3s timeout)
  - Connection state tracking

- [x] T053 [US8] Register DRC commands in pkg/adapter/dji/router/drc_commands.go
  - `drc_mode_enter` with DRCModeEnterData
  - `drc_mode_exit`
  - `drone_control` with DroneControlData
  - `drone_emergency_stop`
  - `heart` with HeartData

- [x] T054 [US8] Register DRC events in pkg/adapter/dji/router/drc_events.go
  - `joystick_invalid_notify`
  - `drc_status_notify`

- [x] T055 [P] [US8] Add DRC mock messages in tests/mocks/dji_drc_messages.go

**Checkpoint**: User Story 8 complete - DRC 实时控制可用

---

## Phase 11: User Story 9 - 媒体文件管理 (Priority: P3)

**Goal**: 实现媒体文件管理

**Independent Test**: 获取文件列表，验证文件上传进度事件

### Tests for User Story 9

- [x] T056 [P] [US9] Unit test for file commands in pkg/adapter/dji/router/file_commands_test.go
- [x] T057 [P] [US9] Unit test for file events in pkg/adapter/dji/router/file_events_test.go

### Implementation for User Story 9

- [x] T058 [US9] Register file commands in pkg/adapter/dji/router/file_commands.go
  - `file_upload_start` with FileUploadStartData
  - `file_upload_finish` with FileUploadFinishData
  - `file_upload_list` with FileUploadListData

- [x] T059 [US9] Register file events in pkg/adapter/dji/router/file_events.go
  - `highest_priority_upload_flighttask_media` (need_reply)
  - `file_upload_progress`
  - `fileupload_progress` (alias)

**Checkpoint**: User Story 9 complete - 文件管理可用

---

## Phase 12: Extended Features - Firmware & Live (Priority: P3)

**Goal**: 实现固件升级和实时视频

**Independent Test**: 创建 OTA 任务，验证推流控制

### Tests for Extended Features

- [x] T060 [P] Unit test for firmware commands in pkg/adapter/dji/router/firmware_commands_test.go
- [x] T061 [P] Unit test for live commands in pkg/adapter/dji/router/live_commands_test.go

### Implementation for Extended Features

- [x] T062 Register firmware commands in pkg/adapter/dji/router/firmware_commands.go
  - `ota_create` with OTACreateData

- [x] T063 Register firmware events in pkg/adapter/dji/router/firmware_events.go
  - `ota_progress` with OTAProgressData

- [x] T064 Register live commands in pkg/adapter/dji/router/live_commands.go
  - `live_start_push` with LiveStartPushData
  - `live_stop_push` with LiveStopPushData
  - `live_set_quality` with LiveSetQualityData
  - `live_lens_change` with LiveLensChangeData

**Checkpoint**: P3 高级功能完成

---

## Phase 13: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

### Integration Tests

- [ ] T065 [P] End-to-end test for OSD flow in tests/integration/dji_e2e_osd_test.go
- [ ] T066 [P] End-to-end test for Service flow in tests/integration/dji_e2e_service_test.go
- [ ] T067 [P] End-to-end test for Event flow in tests/integration/dji_e2e_event_test.go
- [ ] T068 [P] Performance test for 1000 devices in tests/integration/dji_performance_test.go

### Validation

- [x] T069 Run all tests and verify ≥80% coverage
- [ ] T070 Run golangci-lint and fix any issues
- [ ] T071 Verify message processing latency < 50ms (P95)

### Documentation

- [ ] T072 [P] Update pkg/adapter/dji/README.md with handler/router documentation
- [ ] T073 [P] Update quickstart.md with testing examples

**Checkpoint**: 003 Feature 完成

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: ✅ Already complete from 001/002
- **Foundational (Phase 2)**: Depends on Setup - BLOCKS all user stories
- **User Stories (Phase 3-12)**: All depend on Foundational phase completion
  - P1 stories (US1-US5) should complete first
  - P2 stories (US6-US7) can start after P1
  - P3 stories (US8-US9) can start after P1
- **Polish (Phase 13)**: Depends on all user stories being complete

### User Story Dependencies

```
Phase 2 (Foundational)
         │
         ├──────────────────────────────────────────────────┐
         │                                                  │
         ▼                                                  │
    US1 (OSD) ─────► US2/US3 (State/Status) ─────► US4 (Services)
                                                      │
                                                      ▼
                                                 US5 (Events)
                                                      │
         ┌────────────────────────────────────────────┼────────────────┐
         │                                            │                │
         ▼                                            ▼                ▼
    US6 (Wayline)                               US7 (Camera)     US4-Ext (Config)
         │                                            │                │
         └────────────────────────────────────────────┴────────────────┘
                                                      │
         ┌────────────────────────────────────────────┼────────────────┐
         │                                            │                │
         ▼                                            ▼                ▼
    US8 (DRC)                                   US9 (File)      Extended (FW/Live)
         │                                            │                │
         └────────────────────────────────────────────┴────────────────┘
                                                      │
                                                      ▼
                                              Phase 13 (Polish)
```

### Within Each User Story (TDD Order)

1. **Tests First**: Write tests, ensure they FAIL
2. **Implementation**: Implement to make tests pass
3. **Refactor**: Clean up code while keeping tests green

### Parallel Opportunities

- **Phase 2**: All tasks marked [P] can run in parallel
- **Phase 3-6**: P1 stories should be sequential (core protocol)
- **Phase 7-9**: P2 stories can run in parallel after P1 completes
- **Phase 10-12**: P3 stories can run in parallel after P1 completes
- **Phase 13**: All tasks marked [P] can run in parallel

---

## Parallel Example: User Story 1

```bash
# Launch all tests first (TDD):
Task: "Unit test for OSDHandler in pkg/adapter/dji/handler/osd_handler_test.go"
Task: "Unit test for OSDParser in pkg/adapter/dji/integration/osd_parser_test.go"
Task: "Integration test for OSD flow in tests/integration/dji_osd_test.go"

# Then implement (after tests fail):
Task: "Create OSDParser in pkg/adapter/dji/integration/osd_parser.go"
Task: "Create OSDHandler in pkg/adapter/dji/handler/osd_handler.go"
Task: "Update adapter.go to use OSDHandler"
```

---

## Implementation Strategy

### MVP First (User Stories 1-5 Only)

1. Complete Phase 1: Setup (✅ already done)
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3-6: User Stories 1-5 (P1 Core Protocol)
4. **STOP and VALIDATE**: Test all P1 stories independently
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 (OSD) → Test independently → Deploy/Demo
3. Add User Stories 2-3 (State/Status) → Test independently → Deploy/Demo
4. Add User Story 4 (Services) → Test independently → Deploy/Demo
5. Add User Story 5 (Events) → Test independently → Deploy/Demo (MVP!)
6. Add User Stories 6-7 (Wayline/Camera) → Test independently → Deploy/Demo
7. Add User Stories 8-9 (DRC/File) → Test independently → Deploy/Demo
8. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (OSD)
   - Developer B: User Story 2-3 (State/Status)
3. After P1 core completes:
   - Developer A: User Story 6 (Wayline)
   - Developer B: User Story 7 (Camera)
   - Developer C: User Story 8 (DRC)
4. Stories complete and integrate independently

---

## Task Summary

**Total Tasks**: 73

**Tasks per User Story**:

| User Story | Phase | Tasks | Priority |
|------------|-------|-------|----------|
| Foundational | 2 | 7 | P0 |
| US1 (OSD) | 3 | 7 | P1 |
| US2/US3 (State/Status) | 4 | 5 | P1 |
| US4 (Services) | 5 | 7 | P1 |
| US5 (Events) | 6 | 6 | P1 |
| US6 (Wayline) | 7 | 5 | P2 |
| US7 (Camera) | 8 | 3 | P2 |
| US4-Ext (Config) | 9 | 4 | P2 |
| US8 (DRC) | 10 | 7 | P3 |
| US9 (File) | 11 | 4 | P3 |
| Extended (FW/Live) | 12 | 5 | P3 |
| Polish | 13 | 9 | - |

**Priority Summary**:

| Priority | User Stories | Tasks | Phases |
|----------|--------------|-------|--------|
| P0 (Blocking) | Foundational | 7 | 2 |
| P1 (Core) | US1-US5 | 25 | 3-6 |
| P2 (Business) | US6-US7, US4-Ext | 12 | 7-9 |
| P3 (Advanced) | US8-US9, Extended | 16 | 10-12 |
| Polish | - | 9 | 13 |

**Parallel Opportunities**:

- Phase 2: 6 parallel tasks
- Phase 3: 4 parallel tasks (tests + mock)
- Phase 7-9: Can run in parallel (3 phases)
- Phase 10-12: Can run in parallel (3 phases)
- Phase 13: 6 parallel tasks

**MVP Scope**: Phase 1-6 (Foundational + US1-US5) = 36 tasks

---

## Incremental Delivery Milestones

### Milestone 1: OSD 可用 (Phase 2-3)
- 完整解析 AircraftOSD/DockOSD/RCOSD
- 验证: 发送 OSD 消息，查看解析结果

### Milestone 2: 核心协议完成 (Phase 4-6)
- State/Status/Services/Events 全部可用
- 验证: 端到端服务调用测试

### Milestone 3: 业务功能完成 (Phase 7-9)
- 航线/相机/配置管理可用
- 验证: 执行航线任务测试

### Milestone 4: 高级功能完成 (Phase 10-12)
- DRC/文件/固件/视频可用
- 验证: DRC 实时控制测试

### Milestone 5: 生产就绪 (Phase 13)
- 测试覆盖率 ≥ 80%
- 性能达标 (< 50ms P95)
- 文档完整

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- **TDD Required**: Write tests FIRST, ensure they FAIL, then implement
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- All data types reference `pkg/adapter/dji/protocol/*` existing structures
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
