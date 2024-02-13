// @generated by protobuf-ts 2.9.3 with parameter generate_dependencies
// @generated from protobuf file "rpc.proto" (syntax proto3)
// tslint:disable
import { Admin } from "./rpc";
import type { ListDeviceFeaturesResp } from "./rpc";
import type { ListLocationsResp } from "./rpc";
import type { UpdateDeviceReq } from "./rpc";
import type { SetDeviceDisableReq } from "./rpc";
import type { GetDeviceResp } from "./rpc";
import type { GetDeviceReq } from "./rpc";
import type { DeleteDeviceReq } from "./rpc";
import type { CreateDeviceResp } from "./rpc";
import type { CreateDeviceReq } from "./rpc";
import type { UpdateGroupReq } from "./rpc";
import type { SetGroupDisableReq } from "./rpc";
import type { GetGroupResp } from "./rpc";
import type { GetGroupReq } from "./rpc";
import type { DeleteGroupReq } from "./rpc";
import type { CreateGroupResp } from "./rpc";
import type { CreateGroupReq } from "./rpc";
import type { SetUserDisableReq } from "./rpc";
import type { SetUserAdminReq } from "./rpc";
import type { ResetUserPasswordReq } from "./rpc";
import type { DeleteUserReq } from "./rpc";
import type { UpdateUserReq } from "./rpc";
import type { GetUserResp } from "./rpc";
import type { GetUserReq } from "./rpc";
import type { CreateUserReq } from "./rpc";
import type { GetAdminUsersPageResp } from "./rpc";
import type { GetAdminUsersPageReq } from "./rpc";
import type { GetAdminGroupsPageResp } from "./rpc";
import type { GetAdminGroupsPageReq } from "./rpc";
import type { GetAdminGroupsIDPageResp } from "./rpc";
import type { GetAdminGroupsIDPageReq } from "./rpc";
import type { GetAdminDevicesPageResp } from "./rpc";
import type { GetAdminDevicesPageReq } from "./rpc";
import type { GetAdminDevicesIDPageResp } from "./rpc";
import type { GetAdminDevicesIDPageReq } from "./rpc";
import { User } from "./rpc";
import type { ListDeviceStorageResp } from "./rpc";
import type { ListDeviceStorageReq } from "./rpc";
import type { ListDeviceLicensesResp } from "./rpc";
import type { ListDeviceLicensesReq } from "./rpc";
import type { GetDeviceSoftwareVersionResp } from "./rpc";
import type { GetDeviceSoftwareVersionReq } from "./rpc";
import type { GetDeviceDetailResp } from "./rpc";
import type { GetDeviceDetailReq } from "./rpc";
import type { GetDeviceRPCStatusResp } from "./rpc";
import type { GetDeviceRPCStatusReq } from "./rpc";
import type { RevokeMySessionReq } from "./rpc";
import type { UpdateMyPasswordReq } from "./rpc";
import type { UpdateMyUsernameReq } from "./rpc";
import type { GetEventsPageResp } from "./rpc";
import type { GetEventsPageReq } from "./rpc";
import type { GetEmailsIDPageResp } from "./rpc";
import type { GetEmailsIDPageReq } from "./rpc";
import type { GetEmailsPageResp } from "./rpc";
import type { GetEmailsPageReq } from "./rpc";
import type { GetProfilePageResp } from "./rpc";
import type { GetDevicesPageResp } from "./rpc";
import type { GetDevicesPageReq } from "./rpc";
import type { GetHomePageResp } from "./rpc";
import { Public } from "./rpc";
import type { ForgotPasswordReq } from "./rpc";
import type { Empty } from "./google/protobuf/empty";
import type { SignUpReq } from "./rpc";
import type { RpcTransport } from "@protobuf-ts/runtime-rpc";
import type { ServiceInfo } from "@protobuf-ts/runtime-rpc";
import { HelloWorld } from "./rpc";
import { stackIntercept } from "@protobuf-ts/runtime-rpc";
import type { HelloResp } from "./rpc";
import type { HelloReq } from "./rpc";
import type { UnaryCall } from "@protobuf-ts/runtime-rpc";
import type { RpcOptions } from "@protobuf-ts/runtime-rpc";
// ---------- HelloWorld

/**
 * @generated from protobuf service HelloWorld
 */
export interface IHelloWorldClient {
    /**
     * @generated from protobuf rpc: Hello(HelloReq) returns (HelloResp);
     */
    hello(input: HelloReq, options?: RpcOptions): UnaryCall<HelloReq, HelloResp>;
}
// ---------- HelloWorld

/**
 * @generated from protobuf service HelloWorld
 */
export class HelloWorldClient implements IHelloWorldClient, ServiceInfo {
    typeName = HelloWorld.typeName;
    methods = HelloWorld.methods;
    options = HelloWorld.options;
    constructor(private readonly _transport: RpcTransport) {
    }
    /**
     * @generated from protobuf rpc: Hello(HelloReq) returns (HelloResp);
     */
    hello(input: HelloReq, options?: RpcOptions): UnaryCall<HelloReq, HelloResp> {
        const method = this.methods[0], opt = this._transport.mergeOptions(options);
        return stackIntercept<HelloReq, HelloResp>("unary", this._transport, method, opt, input);
    }
}
// ---------- Public

/**
 * @generated from protobuf service Public
 */
export interface IPublicClient {
    /**
     * @generated from protobuf rpc: SignUp(SignUpReq) returns (google.protobuf.Empty);
     */
    signUp(input: SignUpReq, options?: RpcOptions): UnaryCall<SignUpReq, Empty>;
    /**
     * @generated from protobuf rpc: ForgotPassword(ForgotPasswordReq) returns (google.protobuf.Empty);
     */
    forgotPassword(input: ForgotPasswordReq, options?: RpcOptions): UnaryCall<ForgotPasswordReq, Empty>;
}
// ---------- Public

/**
 * @generated from protobuf service Public
 */
export class PublicClient implements IPublicClient, ServiceInfo {
    typeName = Public.typeName;
    methods = Public.methods;
    options = Public.options;
    constructor(private readonly _transport: RpcTransport) {
    }
    /**
     * @generated from protobuf rpc: SignUp(SignUpReq) returns (google.protobuf.Empty);
     */
    signUp(input: SignUpReq, options?: RpcOptions): UnaryCall<SignUpReq, Empty> {
        const method = this.methods[0], opt = this._transport.mergeOptions(options);
        return stackIntercept<SignUpReq, Empty>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: ForgotPassword(ForgotPasswordReq) returns (google.protobuf.Empty);
     */
    forgotPassword(input: ForgotPasswordReq, options?: RpcOptions): UnaryCall<ForgotPasswordReq, Empty> {
        const method = this.methods[1], opt = this._transport.mergeOptions(options);
        return stackIntercept<ForgotPasswordReq, Empty>("unary", this._transport, method, opt, input);
    }
}
// ---------- User

/**
 * @generated from protobuf service User
 */
export interface IUserClient {
    /**
     * Pages
     *
     * @generated from protobuf rpc: GetHomePage(google.protobuf.Empty) returns (GetHomePageResp);
     */
    getHomePage(input: Empty, options?: RpcOptions): UnaryCall<Empty, GetHomePageResp>;
    /**
     * @generated from protobuf rpc: GetDevicesPage(GetDevicesPageReq) returns (GetDevicesPageResp);
     */
    getDevicesPage(input: GetDevicesPageReq, options?: RpcOptions): UnaryCall<GetDevicesPageReq, GetDevicesPageResp>;
    /**
     * @generated from protobuf rpc: GetProfilePage(google.protobuf.Empty) returns (GetProfilePageResp);
     */
    getProfilePage(input: Empty, options?: RpcOptions): UnaryCall<Empty, GetProfilePageResp>;
    /**
     * @generated from protobuf rpc: GetEmailsPage(GetEmailsPageReq) returns (GetEmailsPageResp);
     */
    getEmailsPage(input: GetEmailsPageReq, options?: RpcOptions): UnaryCall<GetEmailsPageReq, GetEmailsPageResp>;
    /**
     * @generated from protobuf rpc: GetEmailsIDPage(GetEmailsIDPageReq) returns (GetEmailsIDPageResp);
     */
    getEmailsIDPage(input: GetEmailsIDPageReq, options?: RpcOptions): UnaryCall<GetEmailsIDPageReq, GetEmailsIDPageResp>;
    /**
     * @generated from protobuf rpc: GetEventsPage(GetEventsPageReq) returns (GetEventsPageResp);
     */
    getEventsPage(input: GetEventsPageReq, options?: RpcOptions): UnaryCall<GetEventsPageReq, GetEventsPageResp>;
    /**
     * User
     *
     * @generated from protobuf rpc: UpdateMyUsername(UpdateMyUsernameReq) returns (google.protobuf.Empty);
     */
    updateMyUsername(input: UpdateMyUsernameReq, options?: RpcOptions): UnaryCall<UpdateMyUsernameReq, Empty>;
    /**
     * @generated from protobuf rpc: UpdateMyPassword(UpdateMyPasswordReq) returns (google.protobuf.Empty);
     */
    updateMyPassword(input: UpdateMyPasswordReq, options?: RpcOptions): UnaryCall<UpdateMyPasswordReq, Empty>;
    /**
     * @generated from protobuf rpc: RevokeMySession(RevokeMySessionReq) returns (google.protobuf.Empty);
     */
    revokeMySession(input: RevokeMySessionReq, options?: RpcOptions): UnaryCall<RevokeMySessionReq, Empty>;
    /**
     * @generated from protobuf rpc: RevokeAllMySessions(google.protobuf.Empty) returns (google.protobuf.Empty);
     */
    revokeAllMySessions(input: Empty, options?: RpcOptions): UnaryCall<Empty, Empty>;
    /**
     * Device
     *
     * @generated from protobuf rpc: GetDeviceRPCStatus(GetDeviceRPCStatusReq) returns (GetDeviceRPCStatusResp);
     */
    getDeviceRPCStatus(input: GetDeviceRPCStatusReq, options?: RpcOptions): UnaryCall<GetDeviceRPCStatusReq, GetDeviceRPCStatusResp>;
    /**
     * @generated from protobuf rpc: GetDeviceDetail(GetDeviceDetailReq) returns (GetDeviceDetailResp);
     */
    getDeviceDetail(input: GetDeviceDetailReq, options?: RpcOptions): UnaryCall<GetDeviceDetailReq, GetDeviceDetailResp>;
    /**
     * @generated from protobuf rpc: GetDeviceSoftwareVersion(GetDeviceSoftwareVersionReq) returns (GetDeviceSoftwareVersionResp);
     */
    getDeviceSoftwareVersion(input: GetDeviceSoftwareVersionReq, options?: RpcOptions): UnaryCall<GetDeviceSoftwareVersionReq, GetDeviceSoftwareVersionResp>;
    /**
     * @generated from protobuf rpc: ListDeviceLicenses(ListDeviceLicensesReq) returns (ListDeviceLicensesResp);
     */
    listDeviceLicenses(input: ListDeviceLicensesReq, options?: RpcOptions): UnaryCall<ListDeviceLicensesReq, ListDeviceLicensesResp>;
    /**
     * @generated from protobuf rpc: ListDeviceStorage(ListDeviceStorageReq) returns (ListDeviceStorageResp);
     */
    listDeviceStorage(input: ListDeviceStorageReq, options?: RpcOptions): UnaryCall<ListDeviceStorageReq, ListDeviceStorageResp>;
}
// ---------- User

/**
 * @generated from protobuf service User
 */
export class UserClient implements IUserClient, ServiceInfo {
    typeName = User.typeName;
    methods = User.methods;
    options = User.options;
    constructor(private readonly _transport: RpcTransport) {
    }
    /**
     * Pages
     *
     * @generated from protobuf rpc: GetHomePage(google.protobuf.Empty) returns (GetHomePageResp);
     */
    getHomePage(input: Empty, options?: RpcOptions): UnaryCall<Empty, GetHomePageResp> {
        const method = this.methods[0], opt = this._transport.mergeOptions(options);
        return stackIntercept<Empty, GetHomePageResp>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: GetDevicesPage(GetDevicesPageReq) returns (GetDevicesPageResp);
     */
    getDevicesPage(input: GetDevicesPageReq, options?: RpcOptions): UnaryCall<GetDevicesPageReq, GetDevicesPageResp> {
        const method = this.methods[1], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetDevicesPageReq, GetDevicesPageResp>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: GetProfilePage(google.protobuf.Empty) returns (GetProfilePageResp);
     */
    getProfilePage(input: Empty, options?: RpcOptions): UnaryCall<Empty, GetProfilePageResp> {
        const method = this.methods[2], opt = this._transport.mergeOptions(options);
        return stackIntercept<Empty, GetProfilePageResp>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: GetEmailsPage(GetEmailsPageReq) returns (GetEmailsPageResp);
     */
    getEmailsPage(input: GetEmailsPageReq, options?: RpcOptions): UnaryCall<GetEmailsPageReq, GetEmailsPageResp> {
        const method = this.methods[3], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetEmailsPageReq, GetEmailsPageResp>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: GetEmailsIDPage(GetEmailsIDPageReq) returns (GetEmailsIDPageResp);
     */
    getEmailsIDPage(input: GetEmailsIDPageReq, options?: RpcOptions): UnaryCall<GetEmailsIDPageReq, GetEmailsIDPageResp> {
        const method = this.methods[4], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetEmailsIDPageReq, GetEmailsIDPageResp>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: GetEventsPage(GetEventsPageReq) returns (GetEventsPageResp);
     */
    getEventsPage(input: GetEventsPageReq, options?: RpcOptions): UnaryCall<GetEventsPageReq, GetEventsPageResp> {
        const method = this.methods[5], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetEventsPageReq, GetEventsPageResp>("unary", this._transport, method, opt, input);
    }
    /**
     * User
     *
     * @generated from protobuf rpc: UpdateMyUsername(UpdateMyUsernameReq) returns (google.protobuf.Empty);
     */
    updateMyUsername(input: UpdateMyUsernameReq, options?: RpcOptions): UnaryCall<UpdateMyUsernameReq, Empty> {
        const method = this.methods[6], opt = this._transport.mergeOptions(options);
        return stackIntercept<UpdateMyUsernameReq, Empty>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: UpdateMyPassword(UpdateMyPasswordReq) returns (google.protobuf.Empty);
     */
    updateMyPassword(input: UpdateMyPasswordReq, options?: RpcOptions): UnaryCall<UpdateMyPasswordReq, Empty> {
        const method = this.methods[7], opt = this._transport.mergeOptions(options);
        return stackIntercept<UpdateMyPasswordReq, Empty>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: RevokeMySession(RevokeMySessionReq) returns (google.protobuf.Empty);
     */
    revokeMySession(input: RevokeMySessionReq, options?: RpcOptions): UnaryCall<RevokeMySessionReq, Empty> {
        const method = this.methods[8], opt = this._transport.mergeOptions(options);
        return stackIntercept<RevokeMySessionReq, Empty>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: RevokeAllMySessions(google.protobuf.Empty) returns (google.protobuf.Empty);
     */
    revokeAllMySessions(input: Empty, options?: RpcOptions): UnaryCall<Empty, Empty> {
        const method = this.methods[9], opt = this._transport.mergeOptions(options);
        return stackIntercept<Empty, Empty>("unary", this._transport, method, opt, input);
    }
    /**
     * Device
     *
     * @generated from protobuf rpc: GetDeviceRPCStatus(GetDeviceRPCStatusReq) returns (GetDeviceRPCStatusResp);
     */
    getDeviceRPCStatus(input: GetDeviceRPCStatusReq, options?: RpcOptions): UnaryCall<GetDeviceRPCStatusReq, GetDeviceRPCStatusResp> {
        const method = this.methods[10], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetDeviceRPCStatusReq, GetDeviceRPCStatusResp>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: GetDeviceDetail(GetDeviceDetailReq) returns (GetDeviceDetailResp);
     */
    getDeviceDetail(input: GetDeviceDetailReq, options?: RpcOptions): UnaryCall<GetDeviceDetailReq, GetDeviceDetailResp> {
        const method = this.methods[11], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetDeviceDetailReq, GetDeviceDetailResp>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: GetDeviceSoftwareVersion(GetDeviceSoftwareVersionReq) returns (GetDeviceSoftwareVersionResp);
     */
    getDeviceSoftwareVersion(input: GetDeviceSoftwareVersionReq, options?: RpcOptions): UnaryCall<GetDeviceSoftwareVersionReq, GetDeviceSoftwareVersionResp> {
        const method = this.methods[12], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetDeviceSoftwareVersionReq, GetDeviceSoftwareVersionResp>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: ListDeviceLicenses(ListDeviceLicensesReq) returns (ListDeviceLicensesResp);
     */
    listDeviceLicenses(input: ListDeviceLicensesReq, options?: RpcOptions): UnaryCall<ListDeviceLicensesReq, ListDeviceLicensesResp> {
        const method = this.methods[13], opt = this._transport.mergeOptions(options);
        return stackIntercept<ListDeviceLicensesReq, ListDeviceLicensesResp>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: ListDeviceStorage(ListDeviceStorageReq) returns (ListDeviceStorageResp);
     */
    listDeviceStorage(input: ListDeviceStorageReq, options?: RpcOptions): UnaryCall<ListDeviceStorageReq, ListDeviceStorageResp> {
        const method = this.methods[14], opt = this._transport.mergeOptions(options);
        return stackIntercept<ListDeviceStorageReq, ListDeviceStorageResp>("unary", this._transport, method, opt, input);
    }
}
// ---------- Admin

/**
 * @generated from protobuf service Admin
 */
export interface IAdminClient {
    /**
     * Pages
     *
     * @generated from protobuf rpc: GetAdminDevicesIDPage(GetAdminDevicesIDPageReq) returns (GetAdminDevicesIDPageResp);
     */
    getAdminDevicesIDPage(input: GetAdminDevicesIDPageReq, options?: RpcOptions): UnaryCall<GetAdminDevicesIDPageReq, GetAdminDevicesIDPageResp>;
    /**
     * @generated from protobuf rpc: GetAdminDevicesPage(GetAdminDevicesPageReq) returns (GetAdminDevicesPageResp);
     */
    getAdminDevicesPage(input: GetAdminDevicesPageReq, options?: RpcOptions): UnaryCall<GetAdminDevicesPageReq, GetAdminDevicesPageResp>;
    /**
     * @generated from protobuf rpc: GetAdminGroupsIDPage(GetAdminGroupsIDPageReq) returns (GetAdminGroupsIDPageResp);
     */
    getAdminGroupsIDPage(input: GetAdminGroupsIDPageReq, options?: RpcOptions): UnaryCall<GetAdminGroupsIDPageReq, GetAdminGroupsIDPageResp>;
    /**
     * @generated from protobuf rpc: GetAdminGroupsPage(GetAdminGroupsPageReq) returns (GetAdminGroupsPageResp);
     */
    getAdminGroupsPage(input: GetAdminGroupsPageReq, options?: RpcOptions): UnaryCall<GetAdminGroupsPageReq, GetAdminGroupsPageResp>;
    /**
     * @generated from protobuf rpc: GetAdminUsersPage(GetAdminUsersPageReq) returns (GetAdminUsersPageResp);
     */
    getAdminUsersPage(input: GetAdminUsersPageReq, options?: RpcOptions): UnaryCall<GetAdminUsersPageReq, GetAdminUsersPageResp>;
    /**
     * User
     *
     * @generated from protobuf rpc: CreateUser(CreateUserReq) returns (google.protobuf.Empty);
     */
    createUser(input: CreateUserReq, options?: RpcOptions): UnaryCall<CreateUserReq, Empty>;
    /**
     * @generated from protobuf rpc: GetUser(GetUserReq) returns (GetUserResp);
     */
    getUser(input: GetUserReq, options?: RpcOptions): UnaryCall<GetUserReq, GetUserResp>;
    /**
     * @generated from protobuf rpc: UpdateUser(UpdateUserReq) returns (google.protobuf.Empty);
     */
    updateUser(input: UpdateUserReq, options?: RpcOptions): UnaryCall<UpdateUserReq, Empty>;
    /**
     * @generated from protobuf rpc: DeleteUser(DeleteUserReq) returns (google.protobuf.Empty);
     */
    deleteUser(input: DeleteUserReq, options?: RpcOptions): UnaryCall<DeleteUserReq, Empty>;
    /**
     * @generated from protobuf rpc: ResetUserPassword(ResetUserPasswordReq) returns (google.protobuf.Empty);
     */
    resetUserPassword(input: ResetUserPasswordReq, options?: RpcOptions): UnaryCall<ResetUserPasswordReq, Empty>;
    /**
     * @generated from protobuf rpc: SetUserAdmin(SetUserAdminReq) returns (google.protobuf.Empty);
     */
    setUserAdmin(input: SetUserAdminReq, options?: RpcOptions): UnaryCall<SetUserAdminReq, Empty>;
    /**
     * @generated from protobuf rpc: SetUserDisable(SetUserDisableReq) returns (google.protobuf.Empty);
     */
    setUserDisable(input: SetUserDisableReq, options?: RpcOptions): UnaryCall<SetUserDisableReq, Empty>;
    /**
     * Group
     *
     * @generated from protobuf rpc: CreateGroup(CreateGroupReq) returns (CreateGroupResp);
     */
    createGroup(input: CreateGroupReq, options?: RpcOptions): UnaryCall<CreateGroupReq, CreateGroupResp>;
    /**
     * @generated from protobuf rpc: DeleteGroup(DeleteGroupReq) returns (google.protobuf.Empty);
     */
    deleteGroup(input: DeleteGroupReq, options?: RpcOptions): UnaryCall<DeleteGroupReq, Empty>;
    /**
     * @generated from protobuf rpc: GetGroup(GetGroupReq) returns (GetGroupResp);
     */
    getGroup(input: GetGroupReq, options?: RpcOptions): UnaryCall<GetGroupReq, GetGroupResp>;
    /**
     * @generated from protobuf rpc: SetGroupDisable(SetGroupDisableReq) returns (google.protobuf.Empty);
     */
    setGroupDisable(input: SetGroupDisableReq, options?: RpcOptions): UnaryCall<SetGroupDisableReq, Empty>;
    /**
     * @generated from protobuf rpc: UpdateGroup(UpdateGroupReq) returns (google.protobuf.Empty);
     */
    updateGroup(input: UpdateGroupReq, options?: RpcOptions): UnaryCall<UpdateGroupReq, Empty>;
    /**
     * Device
     *
     * @generated from protobuf rpc: CreateDevice(CreateDeviceReq) returns (CreateDeviceResp);
     */
    createDevice(input: CreateDeviceReq, options?: RpcOptions): UnaryCall<CreateDeviceReq, CreateDeviceResp>;
    /**
     * @generated from protobuf rpc: DeleteDevice(DeleteDeviceReq) returns (google.protobuf.Empty);
     */
    deleteDevice(input: DeleteDeviceReq, options?: RpcOptions): UnaryCall<DeleteDeviceReq, Empty>;
    /**
     * @generated from protobuf rpc: GetDevice(GetDeviceReq) returns (GetDeviceResp);
     */
    getDevice(input: GetDeviceReq, options?: RpcOptions): UnaryCall<GetDeviceReq, GetDeviceResp>;
    /**
     * @generated from protobuf rpc: SetDeviceDisable(SetDeviceDisableReq) returns (google.protobuf.Empty);
     */
    setDeviceDisable(input: SetDeviceDisableReq, options?: RpcOptions): UnaryCall<SetDeviceDisableReq, Empty>;
    /**
     * @generated from protobuf rpc: UpdateDevice(UpdateDeviceReq) returns (google.protobuf.Empty);
     */
    updateDevice(input: UpdateDeviceReq, options?: RpcOptions): UnaryCall<UpdateDeviceReq, Empty>;
    /**
     * Misc
     *
     * @generated from protobuf rpc: ListLocations(google.protobuf.Empty) returns (ListLocationsResp);
     */
    listLocations(input: Empty, options?: RpcOptions): UnaryCall<Empty, ListLocationsResp>;
    /**
     * @generated from protobuf rpc: ListDeviceFeatures(google.protobuf.Empty) returns (ListDeviceFeaturesResp);
     */
    listDeviceFeatures(input: Empty, options?: RpcOptions): UnaryCall<Empty, ListDeviceFeaturesResp>;
}
// ---------- Admin

/**
 * @generated from protobuf service Admin
 */
export class AdminClient implements IAdminClient, ServiceInfo {
    typeName = Admin.typeName;
    methods = Admin.methods;
    options = Admin.options;
    constructor(private readonly _transport: RpcTransport) {
    }
    /**
     * Pages
     *
     * @generated from protobuf rpc: GetAdminDevicesIDPage(GetAdminDevicesIDPageReq) returns (GetAdminDevicesIDPageResp);
     */
    getAdminDevicesIDPage(input: GetAdminDevicesIDPageReq, options?: RpcOptions): UnaryCall<GetAdminDevicesIDPageReq, GetAdminDevicesIDPageResp> {
        const method = this.methods[0], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetAdminDevicesIDPageReq, GetAdminDevicesIDPageResp>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: GetAdminDevicesPage(GetAdminDevicesPageReq) returns (GetAdminDevicesPageResp);
     */
    getAdminDevicesPage(input: GetAdminDevicesPageReq, options?: RpcOptions): UnaryCall<GetAdminDevicesPageReq, GetAdminDevicesPageResp> {
        const method = this.methods[1], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetAdminDevicesPageReq, GetAdminDevicesPageResp>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: GetAdminGroupsIDPage(GetAdminGroupsIDPageReq) returns (GetAdminGroupsIDPageResp);
     */
    getAdminGroupsIDPage(input: GetAdminGroupsIDPageReq, options?: RpcOptions): UnaryCall<GetAdminGroupsIDPageReq, GetAdminGroupsIDPageResp> {
        const method = this.methods[2], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetAdminGroupsIDPageReq, GetAdminGroupsIDPageResp>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: GetAdminGroupsPage(GetAdminGroupsPageReq) returns (GetAdminGroupsPageResp);
     */
    getAdminGroupsPage(input: GetAdminGroupsPageReq, options?: RpcOptions): UnaryCall<GetAdminGroupsPageReq, GetAdminGroupsPageResp> {
        const method = this.methods[3], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetAdminGroupsPageReq, GetAdminGroupsPageResp>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: GetAdminUsersPage(GetAdminUsersPageReq) returns (GetAdminUsersPageResp);
     */
    getAdminUsersPage(input: GetAdminUsersPageReq, options?: RpcOptions): UnaryCall<GetAdminUsersPageReq, GetAdminUsersPageResp> {
        const method = this.methods[4], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetAdminUsersPageReq, GetAdminUsersPageResp>("unary", this._transport, method, opt, input);
    }
    /**
     * User
     *
     * @generated from protobuf rpc: CreateUser(CreateUserReq) returns (google.protobuf.Empty);
     */
    createUser(input: CreateUserReq, options?: RpcOptions): UnaryCall<CreateUserReq, Empty> {
        const method = this.methods[5], opt = this._transport.mergeOptions(options);
        return stackIntercept<CreateUserReq, Empty>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: GetUser(GetUserReq) returns (GetUserResp);
     */
    getUser(input: GetUserReq, options?: RpcOptions): UnaryCall<GetUserReq, GetUserResp> {
        const method = this.methods[6], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetUserReq, GetUserResp>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: UpdateUser(UpdateUserReq) returns (google.protobuf.Empty);
     */
    updateUser(input: UpdateUserReq, options?: RpcOptions): UnaryCall<UpdateUserReq, Empty> {
        const method = this.methods[7], opt = this._transport.mergeOptions(options);
        return stackIntercept<UpdateUserReq, Empty>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: DeleteUser(DeleteUserReq) returns (google.protobuf.Empty);
     */
    deleteUser(input: DeleteUserReq, options?: RpcOptions): UnaryCall<DeleteUserReq, Empty> {
        const method = this.methods[8], opt = this._transport.mergeOptions(options);
        return stackIntercept<DeleteUserReq, Empty>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: ResetUserPassword(ResetUserPasswordReq) returns (google.protobuf.Empty);
     */
    resetUserPassword(input: ResetUserPasswordReq, options?: RpcOptions): UnaryCall<ResetUserPasswordReq, Empty> {
        const method = this.methods[9], opt = this._transport.mergeOptions(options);
        return stackIntercept<ResetUserPasswordReq, Empty>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: SetUserAdmin(SetUserAdminReq) returns (google.protobuf.Empty);
     */
    setUserAdmin(input: SetUserAdminReq, options?: RpcOptions): UnaryCall<SetUserAdminReq, Empty> {
        const method = this.methods[10], opt = this._transport.mergeOptions(options);
        return stackIntercept<SetUserAdminReq, Empty>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: SetUserDisable(SetUserDisableReq) returns (google.protobuf.Empty);
     */
    setUserDisable(input: SetUserDisableReq, options?: RpcOptions): UnaryCall<SetUserDisableReq, Empty> {
        const method = this.methods[11], opt = this._transport.mergeOptions(options);
        return stackIntercept<SetUserDisableReq, Empty>("unary", this._transport, method, opt, input);
    }
    /**
     * Group
     *
     * @generated from protobuf rpc: CreateGroup(CreateGroupReq) returns (CreateGroupResp);
     */
    createGroup(input: CreateGroupReq, options?: RpcOptions): UnaryCall<CreateGroupReq, CreateGroupResp> {
        const method = this.methods[12], opt = this._transport.mergeOptions(options);
        return stackIntercept<CreateGroupReq, CreateGroupResp>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: DeleteGroup(DeleteGroupReq) returns (google.protobuf.Empty);
     */
    deleteGroup(input: DeleteGroupReq, options?: RpcOptions): UnaryCall<DeleteGroupReq, Empty> {
        const method = this.methods[13], opt = this._transport.mergeOptions(options);
        return stackIntercept<DeleteGroupReq, Empty>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: GetGroup(GetGroupReq) returns (GetGroupResp);
     */
    getGroup(input: GetGroupReq, options?: RpcOptions): UnaryCall<GetGroupReq, GetGroupResp> {
        const method = this.methods[14], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetGroupReq, GetGroupResp>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: SetGroupDisable(SetGroupDisableReq) returns (google.protobuf.Empty);
     */
    setGroupDisable(input: SetGroupDisableReq, options?: RpcOptions): UnaryCall<SetGroupDisableReq, Empty> {
        const method = this.methods[15], opt = this._transport.mergeOptions(options);
        return stackIntercept<SetGroupDisableReq, Empty>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: UpdateGroup(UpdateGroupReq) returns (google.protobuf.Empty);
     */
    updateGroup(input: UpdateGroupReq, options?: RpcOptions): UnaryCall<UpdateGroupReq, Empty> {
        const method = this.methods[16], opt = this._transport.mergeOptions(options);
        return stackIntercept<UpdateGroupReq, Empty>("unary", this._transport, method, opt, input);
    }
    /**
     * Device
     *
     * @generated from protobuf rpc: CreateDevice(CreateDeviceReq) returns (CreateDeviceResp);
     */
    createDevice(input: CreateDeviceReq, options?: RpcOptions): UnaryCall<CreateDeviceReq, CreateDeviceResp> {
        const method = this.methods[17], opt = this._transport.mergeOptions(options);
        return stackIntercept<CreateDeviceReq, CreateDeviceResp>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: DeleteDevice(DeleteDeviceReq) returns (google.protobuf.Empty);
     */
    deleteDevice(input: DeleteDeviceReq, options?: RpcOptions): UnaryCall<DeleteDeviceReq, Empty> {
        const method = this.methods[18], opt = this._transport.mergeOptions(options);
        return stackIntercept<DeleteDeviceReq, Empty>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: GetDevice(GetDeviceReq) returns (GetDeviceResp);
     */
    getDevice(input: GetDeviceReq, options?: RpcOptions): UnaryCall<GetDeviceReq, GetDeviceResp> {
        const method = this.methods[19], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetDeviceReq, GetDeviceResp>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: SetDeviceDisable(SetDeviceDisableReq) returns (google.protobuf.Empty);
     */
    setDeviceDisable(input: SetDeviceDisableReq, options?: RpcOptions): UnaryCall<SetDeviceDisableReq, Empty> {
        const method = this.methods[20], opt = this._transport.mergeOptions(options);
        return stackIntercept<SetDeviceDisableReq, Empty>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: UpdateDevice(UpdateDeviceReq) returns (google.protobuf.Empty);
     */
    updateDevice(input: UpdateDeviceReq, options?: RpcOptions): UnaryCall<UpdateDeviceReq, Empty> {
        const method = this.methods[21], opt = this._transport.mergeOptions(options);
        return stackIntercept<UpdateDeviceReq, Empty>("unary", this._transport, method, opt, input);
    }
    /**
     * Misc
     *
     * @generated from protobuf rpc: ListLocations(google.protobuf.Empty) returns (ListLocationsResp);
     */
    listLocations(input: Empty, options?: RpcOptions): UnaryCall<Empty, ListLocationsResp> {
        const method = this.methods[22], opt = this._transport.mergeOptions(options);
        return stackIntercept<Empty, ListLocationsResp>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: ListDeviceFeatures(google.protobuf.Empty) returns (ListDeviceFeaturesResp);
     */
    listDeviceFeatures(input: Empty, options?: RpcOptions): UnaryCall<Empty, ListDeviceFeaturesResp> {
        const method = this.methods[23], opt = this._transport.mergeOptions(options);
        return stackIntercept<Empty, ListDeviceFeaturesResp>("unary", this._transport, method, opt, input);
    }
}
