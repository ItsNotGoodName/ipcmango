/* eslint-disable */
// api v1.0.0 d3caf539c7bed0d32b8d24b7b44793f47e618dd5
// --
// Code generated by webrpc-gen@v0.13.0 with github.com/ItsNotGoodName/gen-typescript-nuxt@heads/master generator. DO NOT EDIT.
//
// webrpc-gen -schema=./server/api.ridl -target=github.com/ItsNotGoodName/gen-typescript-nuxt@heads/master -client -out=./ui/core/client.gen.ts

// WebRPC description and code-gen version
export const WebRPCVersion = "v1"

// Schema version of your RIDL schema
export const WebRPCSchemaVersion = "v1.0.0"

// Schema hash generated from your RIDL schema
export const WebRPCSchemaHash = "d3caf539c7bed0d32b8d24b7b44793f47e618dd5"

//
// Types
//


export interface User {
  id: number
  email: string
  username: string
  createdAt: string
}

export interface AuthLoginResponse {
  session: string
  token: string
  expiredAt: string
}

export interface AuthService {
  register(args: RegisterArgs, headers?: object, signal?: AbortSignal): Promise<RegisterReturn>
  login(args: LoginArgs, headers?: object, signal?: AbortSignal): Promise<LoginReturn>
  refresh(args: RefreshArgs, headers?: object, signal?: AbortSignal): Promise<RefreshReturn>
  logout(args: LogoutArgs, headers?: object, signal?: AbortSignal): Promise<LogoutReturn>
}

export interface RegisterArgs {
  email: string
  username: string
  password: string
  passwordConfirm: string
}

export interface RegisterReturn {  
}
export interface LoginArgs {
  usernameOrEmail: string
  password: string
  clientId: string
  ipAddress: string
}

export interface LoginReturn {
  res: AuthLoginResponse  
}
export interface RefreshArgs {
  session: string
}

export interface RefreshReturn {
  token: string  
}
export interface LogoutArgs {
  session: string
}

export interface LogoutReturn {  
}

export interface UserService {
  me(headers?: object, signal?: AbortSignal): Promise<MeReturn>
}

export interface MeArgs {
}

export interface MeReturn {
  user: User  
}

export interface DahuaService {
  cameraCount(headers?: object, signal?: AbortSignal): Promise<CameraCountReturn>
  activeScannerCount(headers?: object, signal?: AbortSignal): Promise<ActiveScannerCountReturn>
}

export interface CameraCountArgs {
}

export interface CameraCountReturn {
  count: number  
}
export interface ActiveScannerCountArgs {
}

export interface ActiveScannerCountReturn {
  activeScanners: number
  totalScanners: number  
}


  
//
// Client
//
export class AuthService implements AuthService {
  protected hostname: string
  protected fetch: Fetch
  protected path = '/rpc/AuthService/'

  constructor(hostname: string, fetch: Fetch) {
    this.hostname = hostname
    this.fetch = (input, init) => fetch(input, init)
  }

  private url(name: string): string {
    return this.hostname + this.path + name
  }
  
  register = (args: RegisterArgs, headers?: object, signal?: AbortSignal): Promise<RegisterReturn> => {
    return this.fetch(
      this.url('Register'),
      createHTTPRequest(args, headers, signal)).then((res) => {
      return buildResponse(res).then(_data => {
        return {}
      })
    }, (error) => {
      throw WebrpcRequestFailedError.new({ cause: `fetch(): ${error.message || ''}` })
    })
  }
  
  login = (args: LoginArgs, headers?: object, signal?: AbortSignal): Promise<LoginReturn> => {
    return this.fetch(
      this.url('Login'),
      createHTTPRequest(args, headers, signal)).then((res) => {
      return buildResponse(res).then(_data => {
        return {
          res: <AuthLoginResponse>(_data.res),
        }
      })
    }, (error) => {
      throw WebrpcRequestFailedError.new({ cause: `fetch(): ${error.message || ''}` })
    })
  }
  
  refresh = (args: RefreshArgs, headers?: object, signal?: AbortSignal): Promise<RefreshReturn> => {
    return this.fetch(
      this.url('Refresh'),
      createHTTPRequest(args, headers, signal)).then((res) => {
      return buildResponse(res).then(_data => {
        return {
          token: <string>(_data.token),
        }
      })
    }, (error) => {
      throw WebrpcRequestFailedError.new({ cause: `fetch(): ${error.message || ''}` })
    })
  }
  
  logout = (args: LogoutArgs, headers?: object, signal?: AbortSignal): Promise<LogoutReturn> => {
    return this.fetch(
      this.url('Logout'),
      createHTTPRequest(args, headers, signal)).then((res) => {
      return buildResponse(res).then(_data => {
        return {}
      })
    }, (error) => {
      throw WebrpcRequestFailedError.new({ cause: `fetch(): ${error.message || ''}` })
    })
  }
  
}

export class UserService implements UserService {
  protected hostname: string
  protected fetch: Fetch
  protected path = '/rpc/UserService/'

  constructor(hostname: string, fetch: Fetch) {
    this.hostname = hostname
    this.fetch = (input, init) => fetch(input, init)
  }

  private url(name: string): string {
    return this.hostname + this.path + name
  }
  
  me = (headers?: object, signal?: AbortSignal): Promise<MeReturn> => {
    return this.fetch(
      this.url('Me'),
      createHTTPRequest({}, headers, signal)
      ).then((res) => {
      return buildResponse(res).then(_data => {
        return {
          user: <User>(_data.user),
        }
      })
    }, (error) => {
      throw WebrpcRequestFailedError.new({ cause: `fetch(): ${error.message || ''}` })
    })
  }
  
}

export class DahuaService implements DahuaService {
  protected hostname: string
  protected fetch: Fetch
  protected path = '/rpc/DahuaService/'

  constructor(hostname: string, fetch: Fetch) {
    this.hostname = hostname
    this.fetch = (input, init) => fetch(input, init)
  }

  private url(name: string): string {
    return this.hostname + this.path + name
  }
  
  cameraCount = (headers?: object, signal?: AbortSignal): Promise<CameraCountReturn> => {
    return this.fetch(
      this.url('CameraCount'),
      createHTTPRequest({}, headers, signal)
      ).then((res) => {
      return buildResponse(res).then(_data => {
        return {
          count: <number>(_data.count),
        }
      })
    }, (error) => {
      throw WebrpcRequestFailedError.new({ cause: `fetch(): ${error.message || ''}` })
    })
  }
  
  activeScannerCount = (headers?: object, signal?: AbortSignal): Promise<ActiveScannerCountReturn> => {
    return this.fetch(
      this.url('ActiveScannerCount'),
      createHTTPRequest({}, headers, signal)
      ).then((res) => {
      return buildResponse(res).then(_data => {
        return {
          activeScanners: <number>(_data.activeScanners),
          totalScanners: <number>(_data.totalScanners),
        }
      })
    }, (error) => {
      throw WebrpcRequestFailedError.new({ cause: `fetch(): ${error.message || ''}` })
    })
  }
  
}

  const createHTTPRequest = (body: object = {}, headers: object = {}, signal: AbortSignal | null = null): object => {
  return {
    method: 'POST',
    headers: { ...headers, 'Content-Type': 'application/json' },
    body: JSON.stringify(body || {}),
    signal,
    ignoreResponseError: true,
    parseResponse: (txt: any) => txt
  }
}

const buildResponse = (res: Response & { _data?: any }): Promise<any> => {
  return Promise.resolve().then(() => {
    let data
    try {
      data = JSON.parse(res._data)
    } catch(error) {
      let message = ''
      if (error instanceof Error)  {
        message = error.message
      }
      throw WebrpcBadResponseError.new({
        status: res.status,
        cause: `JSON.parse(): ${message}: response text: ${res._data}`},
      )
    }
    if (!res.ok) {
      const code: number = (typeof data.code === 'number') ? data.code : 0
      throw (webrpcErrorByCode[code] || WebrpcError).new(data)
    }
    return data
  })
}

//
// Errors
//

export class WebrpcError extends Error {
  name: string
  code: number
  message: string
  status: number
  cause?: string

  /** @deprecated Use message instead of msg. Deprecated in webrpc v0.11.0. */
  msg: string

  constructor(name: string, code: number, message: string, status: number, cause?: string) {
    super(message)
    this.name = name || 'WebrpcError'
    this.code = typeof code === 'number' ? code : 0
    this.message = message || `endpoint error ${this.code}`
    this.msg = this.message
    this.status = typeof status === 'number' ? status : 0
    this.cause = cause
    Object.setPrototypeOf(this, WebrpcError.prototype)
  }

  static new(payload: any): WebrpcError {
    return new this(payload.error, payload.code, payload.message || payload.msg, payload.status, payload.cause)
  }
}

// Webrpc errors

export class WebrpcEndpointError extends WebrpcError {
  constructor(
    name: string = 'WebrpcEndpoint',
    code: number = 0,
    message: string = 'endpoint error',
    status: number = 0,
    cause?: string
  ) {
    super(name, code, message, status, cause)
    Object.setPrototypeOf(this, WebrpcEndpointError.prototype)
  }
}

export class WebrpcRequestFailedError extends WebrpcError {
  constructor(
    name: string = 'WebrpcRequestFailed',
    code: number = -1,
    message: string = 'request failed',
    status: number = 0,
    cause?: string
  ) {
    super(name, code, message, status, cause)
    Object.setPrototypeOf(this, WebrpcRequestFailedError.prototype)
  }
}

export class WebrpcBadRouteError extends WebrpcError {
  constructor(
    name: string = 'WebrpcBadRoute',
    code: number = -2,
    message: string = 'bad route',
    status: number = 0,
    cause?: string
  ) {
    super(name, code, message, status, cause)
    Object.setPrototypeOf(this, WebrpcBadRouteError.prototype)
  }
}

export class WebrpcBadMethodError extends WebrpcError {
  constructor(
    name: string = 'WebrpcBadMethod',
    code: number = -3,
    message: string = 'bad method',
    status: number = 0,
    cause?: string
  ) {
    super(name, code, message, status, cause)
    Object.setPrototypeOf(this, WebrpcBadMethodError.prototype)
  }
}

export class WebrpcBadRequestError extends WebrpcError {
  constructor(
    name: string = 'WebrpcBadRequest',
    code: number = -4,
    message: string = 'bad request',
    status: number = 0,
    cause?: string
  ) {
    super(name, code, message, status, cause)
    Object.setPrototypeOf(this, WebrpcBadRequestError.prototype)
  }
}

export class WebrpcBadResponseError extends WebrpcError {
  constructor(
    name: string = 'WebrpcBadResponse',
    code: number = -5,
    message: string = 'bad response',
    status: number = 0,
    cause?: string
  ) {
    super(name, code, message, status, cause)
    Object.setPrototypeOf(this, WebrpcBadResponseError.prototype)
  }
}

export class WebrpcServerPanicError extends WebrpcError {
  constructor(
    name: string = 'WebrpcServerPanic',
    code: number = -6,
    message: string = 'server panic',
    status: number = 0,
    cause?: string
  ) {
    super(name, code, message, status, cause)
    Object.setPrototypeOf(this, WebrpcServerPanicError.prototype)
  }
}

export class WebrpcInternalErrorError extends WebrpcError {
  constructor(
    name: string = 'WebrpcInternalError',
    code: number = -7,
    message: string = 'internal error',
    status: number = 0,
    cause?: string
  ) {
    super(name, code, message, status, cause)
    Object.setPrototypeOf(this, WebrpcInternalErrorError.prototype)
  }
}


// Schema errors

export class InvalidSessionError extends WebrpcError {
  constructor(
    name: string = 'InvalidSession',
    code: number = 700,
    message: string = 'invalid session',
    status: number = 0,
    cause?: string
  ) {
    super(name, code, message, status, cause)
    Object.setPrototypeOf(this, InvalidSessionError.prototype)
  }
}

export class InvalidTokenError extends WebrpcError {
  constructor(
    name: string = 'InvalidToken',
    code: number = 701,
    message: string = 'invalid token',
    status: number = 0,
    cause?: string
  ) {
    super(name, code, message, status, cause)
    Object.setPrototypeOf(this, InvalidTokenError.prototype)
  }
}


export enum errors {
  WebrpcEndpoint = 'WebrpcEndpoint',
  WebrpcRequestFailed = 'WebrpcRequestFailed',
  WebrpcBadRoute = 'WebrpcBadRoute',
  WebrpcBadMethod = 'WebrpcBadMethod',
  WebrpcBadRequest = 'WebrpcBadRequest',
  WebrpcBadResponse = 'WebrpcBadResponse',
  WebrpcServerPanic = 'WebrpcServerPanic',
  WebrpcInternalError = 'WebrpcInternalError',
  InvalidSession = 'InvalidSession',
  InvalidToken = 'InvalidToken',
}

const webrpcErrorByCode: { [code: number]: any } = {
  [0]: WebrpcEndpointError,
  [-1]: WebrpcRequestFailedError,
  [-2]: WebrpcBadRouteError,
  [-3]: WebrpcBadMethodError,
  [-4]: WebrpcBadRequestError,
  [-5]: WebrpcBadResponseError,
  [-6]: WebrpcServerPanicError,
  [-7]: WebrpcInternalErrorError,
  [700]: InvalidSessionError,
  [701]: InvalidTokenError,
}

import { type NitroFetchRequest, NitroFetchOptions, TypedInternalResponse, ExtractedRouteMethod } from "nitropack"
import { type FetchResponse } from "ofetch"

export type Fetch<T = unknown, R extends NitroFetchRequest = NitroFetchRequest, O extends NitroFetchOptions<R> = NitroFetchOptions<R>> = (request: R, opts?: O) => Promise<FetchResponse<TypedInternalResponse<R, T, ExtractedRouteMethod<R, O>>>>;
