/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "poktroll.poktroll";

export interface StakeInfo {
  /** address of the staker (bech32 prefixed) */
  address: string;
  coinsStaked: string;
}

/** Custom parameters for each actor */
export interface WatcherParams {
}

/** Placeholder for now */
export interface PortalParams {
}

/** Placeholder for now */
export interface ServicerParams {
}

/** Placeholder for now */
export interface ApplicationParams {
}

/** Actor definitions */
export interface Watcher {
  stakeInfo: StakeInfo | undefined;
  watcherParams: WatcherParams | undefined;
}

export interface Portal {
  stakeInfo: StakeInfo | undefined;
  portalParams: PortalParams | undefined;
}

export interface Servicer {
  stakeInfo: StakeInfo | undefined;
  servicerParams: ServicerParams | undefined;
}

export interface Application {
  stakeInfo: StakeInfo | undefined;
  applicationParams: ApplicationParams | undefined;
}

/** Wrapping all actors in one message for potential use cases */
export interface Actor {
  watcher: Watcher | undefined;
  portal: Portal | undefined;
  servicer: Servicer | undefined;
  application: Application | undefined;
}

function createBaseStakeInfo(): StakeInfo {
  return { address: "", coinsStaked: "" };
}

export const StakeInfo = {
  encode(message: StakeInfo, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    if (message.coinsStaked !== "") {
      writer.uint32(18).string(message.coinsStaked);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StakeInfo {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStakeInfo();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;
        case 2:
          message.coinsStaked = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): StakeInfo {
    return {
      address: isSet(object.address) ? String(object.address) : "",
      coinsStaked: isSet(object.coinsStaked) ? String(object.coinsStaked) : "",
    };
  },

  toJSON(message: StakeInfo): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = message.address);
    message.coinsStaked !== undefined && (obj.coinsStaked = message.coinsStaked);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<StakeInfo>, I>>(object: I): StakeInfo {
    const message = createBaseStakeInfo();
    message.address = object.address ?? "";
    message.coinsStaked = object.coinsStaked ?? "";
    return message;
  },
};

function createBaseWatcherParams(): WatcherParams {
  return {};
}

export const WatcherParams = {
  encode(_: WatcherParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): WatcherParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseWatcherParams();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): WatcherParams {
    return {};
  },

  toJSON(_: WatcherParams): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<WatcherParams>, I>>(_: I): WatcherParams {
    const message = createBaseWatcherParams();
    return message;
  },
};

function createBasePortalParams(): PortalParams {
  return {};
}

export const PortalParams = {
  encode(_: PortalParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PortalParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePortalParams();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): PortalParams {
    return {};
  },

  toJSON(_: PortalParams): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<PortalParams>, I>>(_: I): PortalParams {
    const message = createBasePortalParams();
    return message;
  },
};

function createBaseServicerParams(): ServicerParams {
  return {};
}

export const ServicerParams = {
  encode(_: ServicerParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ServicerParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseServicerParams();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): ServicerParams {
    return {};
  },

  toJSON(_: ServicerParams): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ServicerParams>, I>>(_: I): ServicerParams {
    const message = createBaseServicerParams();
    return message;
  },
};

function createBaseApplicationParams(): ApplicationParams {
  return {};
}

export const ApplicationParams = {
  encode(_: ApplicationParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ApplicationParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseApplicationParams();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): ApplicationParams {
    return {};
  },

  toJSON(_: ApplicationParams): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ApplicationParams>, I>>(_: I): ApplicationParams {
    const message = createBaseApplicationParams();
    return message;
  },
};

function createBaseWatcher(): Watcher {
  return { stakeInfo: undefined, watcherParams: undefined };
}

export const Watcher = {
  encode(message: Watcher, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.stakeInfo !== undefined) {
      StakeInfo.encode(message.stakeInfo, writer.uint32(10).fork()).ldelim();
    }
    if (message.watcherParams !== undefined) {
      WatcherParams.encode(message.watcherParams, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Watcher {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseWatcher();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.stakeInfo = StakeInfo.decode(reader, reader.uint32());
          break;
        case 2:
          message.watcherParams = WatcherParams.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Watcher {
    return {
      stakeInfo: isSet(object.stakeInfo) ? StakeInfo.fromJSON(object.stakeInfo) : undefined,
      watcherParams: isSet(object.watcherParams) ? WatcherParams.fromJSON(object.watcherParams) : undefined,
    };
  },

  toJSON(message: Watcher): unknown {
    const obj: any = {};
    message.stakeInfo !== undefined
      && (obj.stakeInfo = message.stakeInfo ? StakeInfo.toJSON(message.stakeInfo) : undefined);
    message.watcherParams !== undefined
      && (obj.watcherParams = message.watcherParams ? WatcherParams.toJSON(message.watcherParams) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Watcher>, I>>(object: I): Watcher {
    const message = createBaseWatcher();
    message.stakeInfo = (object.stakeInfo !== undefined && object.stakeInfo !== null)
      ? StakeInfo.fromPartial(object.stakeInfo)
      : undefined;
    message.watcherParams = (object.watcherParams !== undefined && object.watcherParams !== null)
      ? WatcherParams.fromPartial(object.watcherParams)
      : undefined;
    return message;
  },
};

function createBasePortal(): Portal {
  return { stakeInfo: undefined, portalParams: undefined };
}

export const Portal = {
  encode(message: Portal, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.stakeInfo !== undefined) {
      StakeInfo.encode(message.stakeInfo, writer.uint32(10).fork()).ldelim();
    }
    if (message.portalParams !== undefined) {
      PortalParams.encode(message.portalParams, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Portal {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePortal();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.stakeInfo = StakeInfo.decode(reader, reader.uint32());
          break;
        case 2:
          message.portalParams = PortalParams.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Portal {
    return {
      stakeInfo: isSet(object.stakeInfo) ? StakeInfo.fromJSON(object.stakeInfo) : undefined,
      portalParams: isSet(object.portalParams) ? PortalParams.fromJSON(object.portalParams) : undefined,
    };
  },

  toJSON(message: Portal): unknown {
    const obj: any = {};
    message.stakeInfo !== undefined
      && (obj.stakeInfo = message.stakeInfo ? StakeInfo.toJSON(message.stakeInfo) : undefined);
    message.portalParams !== undefined
      && (obj.portalParams = message.portalParams ? PortalParams.toJSON(message.portalParams) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Portal>, I>>(object: I): Portal {
    const message = createBasePortal();
    message.stakeInfo = (object.stakeInfo !== undefined && object.stakeInfo !== null)
      ? StakeInfo.fromPartial(object.stakeInfo)
      : undefined;
    message.portalParams = (object.portalParams !== undefined && object.portalParams !== null)
      ? PortalParams.fromPartial(object.portalParams)
      : undefined;
    return message;
  },
};

function createBaseServicer(): Servicer {
  return { stakeInfo: undefined, servicerParams: undefined };
}

export const Servicer = {
  encode(message: Servicer, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.stakeInfo !== undefined) {
      StakeInfo.encode(message.stakeInfo, writer.uint32(10).fork()).ldelim();
    }
    if (message.servicerParams !== undefined) {
      ServicerParams.encode(message.servicerParams, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Servicer {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseServicer();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.stakeInfo = StakeInfo.decode(reader, reader.uint32());
          break;
        case 2:
          message.servicerParams = ServicerParams.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Servicer {
    return {
      stakeInfo: isSet(object.stakeInfo) ? StakeInfo.fromJSON(object.stakeInfo) : undefined,
      servicerParams: isSet(object.servicerParams) ? ServicerParams.fromJSON(object.servicerParams) : undefined,
    };
  },

  toJSON(message: Servicer): unknown {
    const obj: any = {};
    message.stakeInfo !== undefined
      && (obj.stakeInfo = message.stakeInfo ? StakeInfo.toJSON(message.stakeInfo) : undefined);
    message.servicerParams !== undefined
      && (obj.servicerParams = message.servicerParams ? ServicerParams.toJSON(message.servicerParams) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Servicer>, I>>(object: I): Servicer {
    const message = createBaseServicer();
    message.stakeInfo = (object.stakeInfo !== undefined && object.stakeInfo !== null)
      ? StakeInfo.fromPartial(object.stakeInfo)
      : undefined;
    message.servicerParams = (object.servicerParams !== undefined && object.servicerParams !== null)
      ? ServicerParams.fromPartial(object.servicerParams)
      : undefined;
    return message;
  },
};

function createBaseApplication(): Application {
  return { stakeInfo: undefined, applicationParams: undefined };
}

export const Application = {
  encode(message: Application, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.stakeInfo !== undefined) {
      StakeInfo.encode(message.stakeInfo, writer.uint32(10).fork()).ldelim();
    }
    if (message.applicationParams !== undefined) {
      ApplicationParams.encode(message.applicationParams, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Application {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseApplication();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.stakeInfo = StakeInfo.decode(reader, reader.uint32());
          break;
        case 2:
          message.applicationParams = ApplicationParams.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Application {
    return {
      stakeInfo: isSet(object.stakeInfo) ? StakeInfo.fromJSON(object.stakeInfo) : undefined,
      applicationParams: isSet(object.applicationParams)
        ? ApplicationParams.fromJSON(object.applicationParams)
        : undefined,
    };
  },

  toJSON(message: Application): unknown {
    const obj: any = {};
    message.stakeInfo !== undefined
      && (obj.stakeInfo = message.stakeInfo ? StakeInfo.toJSON(message.stakeInfo) : undefined);
    message.applicationParams !== undefined && (obj.applicationParams = message.applicationParams
      ? ApplicationParams.toJSON(message.applicationParams)
      : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Application>, I>>(object: I): Application {
    const message = createBaseApplication();
    message.stakeInfo = (object.stakeInfo !== undefined && object.stakeInfo !== null)
      ? StakeInfo.fromPartial(object.stakeInfo)
      : undefined;
    message.applicationParams = (object.applicationParams !== undefined && object.applicationParams !== null)
      ? ApplicationParams.fromPartial(object.applicationParams)
      : undefined;
    return message;
  },
};

function createBaseActor(): Actor {
  return { watcher: undefined, portal: undefined, servicer: undefined, application: undefined };
}

export const Actor = {
  encode(message: Actor, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.watcher !== undefined) {
      Watcher.encode(message.watcher, writer.uint32(10).fork()).ldelim();
    }
    if (message.portal !== undefined) {
      Portal.encode(message.portal, writer.uint32(18).fork()).ldelim();
    }
    if (message.servicer !== undefined) {
      Servicer.encode(message.servicer, writer.uint32(26).fork()).ldelim();
    }
    if (message.application !== undefined) {
      Application.encode(message.application, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Actor {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseActor();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.watcher = Watcher.decode(reader, reader.uint32());
          break;
        case 2:
          message.portal = Portal.decode(reader, reader.uint32());
          break;
        case 3:
          message.servicer = Servicer.decode(reader, reader.uint32());
          break;
        case 4:
          message.application = Application.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Actor {
    return {
      watcher: isSet(object.watcher) ? Watcher.fromJSON(object.watcher) : undefined,
      portal: isSet(object.portal) ? Portal.fromJSON(object.portal) : undefined,
      servicer: isSet(object.servicer) ? Servicer.fromJSON(object.servicer) : undefined,
      application: isSet(object.application) ? Application.fromJSON(object.application) : undefined,
    };
  },

  toJSON(message: Actor): unknown {
    const obj: any = {};
    message.watcher !== undefined && (obj.watcher = message.watcher ? Watcher.toJSON(message.watcher) : undefined);
    message.portal !== undefined && (obj.portal = message.portal ? Portal.toJSON(message.portal) : undefined);
    message.servicer !== undefined && (obj.servicer = message.servicer ? Servicer.toJSON(message.servicer) : undefined);
    message.application !== undefined
      && (obj.application = message.application ? Application.toJSON(message.application) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Actor>, I>>(object: I): Actor {
    const message = createBaseActor();
    message.watcher = (object.watcher !== undefined && object.watcher !== null)
      ? Watcher.fromPartial(object.watcher)
      : undefined;
    message.portal = (object.portal !== undefined && object.portal !== null)
      ? Portal.fromPartial(object.portal)
      : undefined;
    message.servicer = (object.servicer !== undefined && object.servicer !== null)
      ? Servicer.fromPartial(object.servicer)
      : undefined;
    message.application = (object.application !== undefined && object.application !== null)
      ? Application.fromPartial(object.application)
      : undefined;
    return message;
  },
};

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
