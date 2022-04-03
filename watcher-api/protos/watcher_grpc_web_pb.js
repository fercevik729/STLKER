/**
 * @fileoverview gRPC-Web generated client stub for 
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');

const proto = require('./watcher_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.WatcherClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options.format = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.WatcherPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options.format = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.TickerRequest,
 *   !proto.TickerResponse>}
 */
const methodDescriptor_Watcher_GetInfo = new grpc.web.MethodDescriptor(
  '/Watcher/GetInfo',
  grpc.web.MethodType.UNARY,
  proto.TickerRequest,
  proto.TickerResponse,
  /**
   * @param {!proto.TickerRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.TickerResponse.deserializeBinary
);


/**
 * @param {!proto.TickerRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.TickerResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.TickerResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.WatcherClient.prototype.getInfo =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/Watcher/GetInfo',
      request,
      metadata || {},
      methodDescriptor_Watcher_GetInfo,
      callback);
};


/**
 * @param {!proto.TickerRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.TickerResponse>}
 *     Promise that resolves to the response
 */
proto.WatcherPromiseClient.prototype.getInfo =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/Watcher/GetInfo',
      request,
      metadata || {},
      methodDescriptor_Watcher_GetInfo);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.TickerRequest,
 *   !proto.PriceResponse>}
 */
const methodDescriptor_Watcher_SubscribeTicker = new grpc.web.MethodDescriptor(
  '/Watcher/SubscribeTicker',
  grpc.web.MethodType.SERVER_STREAMING,
  proto.TickerRequest,
  proto.PriceResponse,
  /**
   * @param {!proto.TickerRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.PriceResponse.deserializeBinary
);


/**
 * @param {!proto.TickerRequest} request The request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.PriceResponse>}
 *     The XHR Node Readable Stream
 */
proto.WatcherClient.prototype.subscribeTicker =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/Watcher/SubscribeTicker',
      request,
      metadata || {},
      methodDescriptor_Watcher_SubscribeTicker);
};


/**
 * @param {!proto.TickerRequest} request The request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.PriceResponse>}
 *     The XHR Node Readable Stream
 */
proto.WatcherPromiseClient.prototype.subscribeTicker =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/Watcher/SubscribeTicker',
      request,
      metadata || {},
      methodDescriptor_Watcher_SubscribeTicker);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.TickerRequest,
 *   !proto.CompanyResponse>}
 */
const methodDescriptor_Watcher_MoreInfo = new grpc.web.MethodDescriptor(
  '/Watcher/MoreInfo',
  grpc.web.MethodType.UNARY,
  proto.TickerRequest,
  proto.CompanyResponse,
  /**
   * @param {!proto.TickerRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.CompanyResponse.deserializeBinary
);


/**
 * @param {!proto.TickerRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.CompanyResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.CompanyResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.WatcherClient.prototype.moreInfo =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/Watcher/MoreInfo',
      request,
      metadata || {},
      methodDescriptor_Watcher_MoreInfo,
      callback);
};


/**
 * @param {!proto.TickerRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.CompanyResponse>}
 *     Promise that resolves to the response
 */
proto.WatcherPromiseClient.prototype.moreInfo =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/Watcher/MoreInfo',
      request,
      metadata || {},
      methodDescriptor_Watcher_MoreInfo);
};


module.exports = proto;

