/**
 * @fileoverview
 * @enhanceable
 * @suppress {messageConventions} JS Compiler reports an error if a variable or
 *     field starts with 'MSG_' and isn't a translatable message.
 * @public
 */
// GENERATED CODE -- DO NOT EDIT!

var jspb = require('google-protobuf');
var goog = jspb;
var global = Function('return this')();

goog.exportSymbol('proto.CompanyResponse', null, global);
goog.exportSymbol('proto.Currencies', null, global);
goog.exportSymbol('proto.PriceResponse', null, global);
goog.exportSymbol('proto.TickerRequest', null, global);
goog.exportSymbol('proto.TickerResponse', null, global);

/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.TickerRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.TickerRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.TickerRequest.displayName = 'proto.TickerRequest';
}


if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.TickerRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.TickerRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.TickerRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.TickerRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    ticker: jspb.Message.getFieldWithDefault(msg, 1, ""),
    destination: jspb.Message.getFieldWithDefault(msg, 2, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.TickerRequest}
 */
proto.TickerRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.TickerRequest;
  return proto.TickerRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.TickerRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.TickerRequest}
 */
proto.TickerRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setTicker(value);
      break;
    case 2:
      var value = /** @type {!proto.Currencies} */ (reader.readEnum());
      msg.setDestination(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.TickerRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.TickerRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.TickerRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.TickerRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getTicker();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getDestination();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
};


/**
 * optional string Ticker = 1;
 * @return {string}
 */
proto.TickerRequest.prototype.getTicker = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.TickerRequest.prototype.setTicker = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional Currencies Destination = 2;
 * @return {!proto.Currencies}
 */
proto.TickerRequest.prototype.getDestination = function() {
  return /** @type {!proto.Currencies} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/** @param {!proto.Currencies} value */
proto.TickerRequest.prototype.setDestination = function(value) {
  jspb.Message.setProto3EnumField(this, 2, value);
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.TickerResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.TickerResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.TickerResponse.displayName = 'proto.TickerResponse';
}


if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.TickerResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.TickerResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.TickerResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.TickerResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    symbol: jspb.Message.getFieldWithDefault(msg, 1, ""),
    open: jspb.Message.getFieldWithDefault(msg, 2, ""),
    high: jspb.Message.getFieldWithDefault(msg, 3, ""),
    low: jspb.Message.getFieldWithDefault(msg, 4, ""),
    price: jspb.Message.getFieldWithDefault(msg, 5, ""),
    volume: jspb.Message.getFieldWithDefault(msg, 6, ""),
    ltd: jspb.Message.getFieldWithDefault(msg, 7, ""),
    prevclose: jspb.Message.getFieldWithDefault(msg, 8, ""),
    change: jspb.Message.getFieldWithDefault(msg, 9, ""),
    percentchange: jspb.Message.getFieldWithDefault(msg, 10, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.TickerResponse}
 */
proto.TickerResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.TickerResponse;
  return proto.TickerResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.TickerResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.TickerResponse}
 */
proto.TickerResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setSymbol(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setOpen(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setHigh(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setLow(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setPrice(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.setVolume(value);
      break;
    case 7:
      var value = /** @type {string} */ (reader.readString());
      msg.setLtd(value);
      break;
    case 8:
      var value = /** @type {string} */ (reader.readString());
      msg.setPrevclose(value);
      break;
    case 9:
      var value = /** @type {string} */ (reader.readString());
      msg.setChange(value);
      break;
    case 10:
      var value = /** @type {string} */ (reader.readString());
      msg.setPercentchange(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.TickerResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.TickerResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.TickerResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.TickerResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSymbol();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getOpen();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getHigh();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getLow();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getPrice();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getVolume();
  if (f.length > 0) {
    writer.writeString(
      6,
      f
    );
  }
  f = message.getLtd();
  if (f.length > 0) {
    writer.writeString(
      7,
      f
    );
  }
  f = message.getPrevclose();
  if (f.length > 0) {
    writer.writeString(
      8,
      f
    );
  }
  f = message.getChange();
  if (f.length > 0) {
    writer.writeString(
      9,
      f
    );
  }
  f = message.getPercentchange();
  if (f.length > 0) {
    writer.writeString(
      10,
      f
    );
  }
};


/**
 * optional string Symbol = 1;
 * @return {string}
 */
proto.TickerResponse.prototype.getSymbol = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.TickerResponse.prototype.setSymbol = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string Open = 2;
 * @return {string}
 */
proto.TickerResponse.prototype.getOpen = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.TickerResponse.prototype.setOpen = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string High = 3;
 * @return {string}
 */
proto.TickerResponse.prototype.getHigh = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.TickerResponse.prototype.setHigh = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string Low = 4;
 * @return {string}
 */
proto.TickerResponse.prototype.getLow = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/** @param {string} value */
proto.TickerResponse.prototype.setLow = function(value) {
  jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * optional string Price = 5;
 * @return {string}
 */
proto.TickerResponse.prototype.getPrice = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/** @param {string} value */
proto.TickerResponse.prototype.setPrice = function(value) {
  jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * optional string Volume = 6;
 * @return {string}
 */
proto.TickerResponse.prototype.getVolume = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 6, ""));
};


/** @param {string} value */
proto.TickerResponse.prototype.setVolume = function(value) {
  jspb.Message.setProto3StringField(this, 6, value);
};


/**
 * optional string LTD = 7;
 * @return {string}
 */
proto.TickerResponse.prototype.getLtd = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 7, ""));
};


/** @param {string} value */
proto.TickerResponse.prototype.setLtd = function(value) {
  jspb.Message.setProto3StringField(this, 7, value);
};


/**
 * optional string PrevClose = 8;
 * @return {string}
 */
proto.TickerResponse.prototype.getPrevclose = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 8, ""));
};


/** @param {string} value */
proto.TickerResponse.prototype.setPrevclose = function(value) {
  jspb.Message.setProto3StringField(this, 8, value);
};


/**
 * optional string Change = 9;
 * @return {string}
 */
proto.TickerResponse.prototype.getChange = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 9, ""));
};


/** @param {string} value */
proto.TickerResponse.prototype.setChange = function(value) {
  jspb.Message.setProto3StringField(this, 9, value);
};


/**
 * optional string PercentChange = 10;
 * @return {string}
 */
proto.TickerResponse.prototype.getPercentchange = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 10, ""));
};


/** @param {string} value */
proto.TickerResponse.prototype.setPercentchange = function(value) {
  jspb.Message.setProto3StringField(this, 10, value);
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.CompanyResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.CompanyResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.CompanyResponse.displayName = 'proto.CompanyResponse';
}


if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.CompanyResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.CompanyResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.CompanyResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.CompanyResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    ticker: jspb.Message.getFieldWithDefault(msg, 1, ""),
    name: jspb.Message.getFieldWithDefault(msg, 2, ""),
    exchange: jspb.Message.getFieldWithDefault(msg, 3, ""),
    sector: jspb.Message.getFieldWithDefault(msg, 4, ""),
    marketcap: jspb.Message.getFieldWithDefault(msg, 5, ""),
    peratio: jspb.Message.getFieldWithDefault(msg, 6, ""),
    pegratio: jspb.Message.getFieldWithDefault(msg, 7, ""),
    divpershare: jspb.Message.getFieldWithDefault(msg, 8, ""),
    eps: jspb.Message.getFieldWithDefault(msg, 9, ""),
    revpershare: jspb.Message.getFieldWithDefault(msg, 10, ""),
    profitmargin: jspb.Message.getFieldWithDefault(msg, 11, ""),
    yearhigh: jspb.Message.getFieldWithDefault(msg, 12, ""),
    yearlow: jspb.Message.getFieldWithDefault(msg, 13, ""),
    sharesoutstanding: jspb.Message.getFieldWithDefault(msg, 14, ""),
    pricetobookratio: jspb.Message.getFieldWithDefault(msg, 15, ""),
    beta: jspb.Message.getFieldWithDefault(msg, 16, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.CompanyResponse}
 */
proto.CompanyResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.CompanyResponse;
  return proto.CompanyResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.CompanyResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.CompanyResponse}
 */
proto.CompanyResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setTicker(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setExchange(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setSector(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setMarketcap(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.setPeratio(value);
      break;
    case 7:
      var value = /** @type {string} */ (reader.readString());
      msg.setPegratio(value);
      break;
    case 8:
      var value = /** @type {string} */ (reader.readString());
      msg.setDivpershare(value);
      break;
    case 9:
      var value = /** @type {string} */ (reader.readString());
      msg.setEps(value);
      break;
    case 10:
      var value = /** @type {string} */ (reader.readString());
      msg.setRevpershare(value);
      break;
    case 11:
      var value = /** @type {string} */ (reader.readString());
      msg.setProfitmargin(value);
      break;
    case 12:
      var value = /** @type {string} */ (reader.readString());
      msg.setYearhigh(value);
      break;
    case 13:
      var value = /** @type {string} */ (reader.readString());
      msg.setYearlow(value);
      break;
    case 14:
      var value = /** @type {string} */ (reader.readString());
      msg.setSharesoutstanding(value);
      break;
    case 15:
      var value = /** @type {string} */ (reader.readString());
      msg.setPricetobookratio(value);
      break;
    case 16:
      var value = /** @type {string} */ (reader.readString());
      msg.setBeta(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.CompanyResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.CompanyResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.CompanyResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.CompanyResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getTicker();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getExchange();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getSector();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getMarketcap();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getPeratio();
  if (f.length > 0) {
    writer.writeString(
      6,
      f
    );
  }
  f = message.getPegratio();
  if (f.length > 0) {
    writer.writeString(
      7,
      f
    );
  }
  f = message.getDivpershare();
  if (f.length > 0) {
    writer.writeString(
      8,
      f
    );
  }
  f = message.getEps();
  if (f.length > 0) {
    writer.writeString(
      9,
      f
    );
  }
  f = message.getRevpershare();
  if (f.length > 0) {
    writer.writeString(
      10,
      f
    );
  }
  f = message.getProfitmargin();
  if (f.length > 0) {
    writer.writeString(
      11,
      f
    );
  }
  f = message.getYearhigh();
  if (f.length > 0) {
    writer.writeString(
      12,
      f
    );
  }
  f = message.getYearlow();
  if (f.length > 0) {
    writer.writeString(
      13,
      f
    );
  }
  f = message.getSharesoutstanding();
  if (f.length > 0) {
    writer.writeString(
      14,
      f
    );
  }
  f = message.getPricetobookratio();
  if (f.length > 0) {
    writer.writeString(
      15,
      f
    );
  }
  f = message.getBeta();
  if (f.length > 0) {
    writer.writeString(
      16,
      f
    );
  }
};


/**
 * optional string Ticker = 1;
 * @return {string}
 */
proto.CompanyResponse.prototype.getTicker = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.CompanyResponse.prototype.setTicker = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string Name = 2;
 * @return {string}
 */
proto.CompanyResponse.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.CompanyResponse.prototype.setName = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string Exchange = 3;
 * @return {string}
 */
proto.CompanyResponse.prototype.getExchange = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.CompanyResponse.prototype.setExchange = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string Sector = 4;
 * @return {string}
 */
proto.CompanyResponse.prototype.getSector = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/** @param {string} value */
proto.CompanyResponse.prototype.setSector = function(value) {
  jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * optional string MarketCap = 5;
 * @return {string}
 */
proto.CompanyResponse.prototype.getMarketcap = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/** @param {string} value */
proto.CompanyResponse.prototype.setMarketcap = function(value) {
  jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * optional string PERatio = 6;
 * @return {string}
 */
proto.CompanyResponse.prototype.getPeratio = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 6, ""));
};


/** @param {string} value */
proto.CompanyResponse.prototype.setPeratio = function(value) {
  jspb.Message.setProto3StringField(this, 6, value);
};


/**
 * optional string PEGRatio = 7;
 * @return {string}
 */
proto.CompanyResponse.prototype.getPegratio = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 7, ""));
};


/** @param {string} value */
proto.CompanyResponse.prototype.setPegratio = function(value) {
  jspb.Message.setProto3StringField(this, 7, value);
};


/**
 * optional string DivPerShare = 8;
 * @return {string}
 */
proto.CompanyResponse.prototype.getDivpershare = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 8, ""));
};


/** @param {string} value */
proto.CompanyResponse.prototype.setDivpershare = function(value) {
  jspb.Message.setProto3StringField(this, 8, value);
};


/**
 * optional string EPS = 9;
 * @return {string}
 */
proto.CompanyResponse.prototype.getEps = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 9, ""));
};


/** @param {string} value */
proto.CompanyResponse.prototype.setEps = function(value) {
  jspb.Message.setProto3StringField(this, 9, value);
};


/**
 * optional string RevPerShare = 10;
 * @return {string}
 */
proto.CompanyResponse.prototype.getRevpershare = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 10, ""));
};


/** @param {string} value */
proto.CompanyResponse.prototype.setRevpershare = function(value) {
  jspb.Message.setProto3StringField(this, 10, value);
};


/**
 * optional string ProfitMargin = 11;
 * @return {string}
 */
proto.CompanyResponse.prototype.getProfitmargin = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 11, ""));
};


/** @param {string} value */
proto.CompanyResponse.prototype.setProfitmargin = function(value) {
  jspb.Message.setProto3StringField(this, 11, value);
};


/**
 * optional string YearHigh = 12;
 * @return {string}
 */
proto.CompanyResponse.prototype.getYearhigh = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 12, ""));
};


/** @param {string} value */
proto.CompanyResponse.prototype.setYearhigh = function(value) {
  jspb.Message.setProto3StringField(this, 12, value);
};


/**
 * optional string YearLow = 13;
 * @return {string}
 */
proto.CompanyResponse.prototype.getYearlow = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 13, ""));
};


/** @param {string} value */
proto.CompanyResponse.prototype.setYearlow = function(value) {
  jspb.Message.setProto3StringField(this, 13, value);
};


/**
 * optional string SharesOutstanding = 14;
 * @return {string}
 */
proto.CompanyResponse.prototype.getSharesoutstanding = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 14, ""));
};


/** @param {string} value */
proto.CompanyResponse.prototype.setSharesoutstanding = function(value) {
  jspb.Message.setProto3StringField(this, 14, value);
};


/**
 * optional string PriceToBookRatio = 15;
 * @return {string}
 */
proto.CompanyResponse.prototype.getPricetobookratio = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 15, ""));
};


/** @param {string} value */
proto.CompanyResponse.prototype.setPricetobookratio = function(value) {
  jspb.Message.setProto3StringField(this, 15, value);
};


/**
 * optional string Beta = 16;
 * @return {string}
 */
proto.CompanyResponse.prototype.getBeta = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 16, ""));
};


/** @param {string} value */
proto.CompanyResponse.prototype.setBeta = function(value) {
  jspb.Message.setProto3StringField(this, 16, value);
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.PriceResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.PriceResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.PriceResponse.displayName = 'proto.PriceResponse';
}


if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.PriceResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.PriceResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.PriceResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.PriceResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    ticker: jspb.Message.getFieldWithDefault(msg, 1, ""),
    stockprice: +jspb.Message.getFieldWithDefault(msg, 2, 0.0),
    currency: jspb.Message.getFieldWithDefault(msg, 3, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.PriceResponse}
 */
proto.PriceResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.PriceResponse;
  return proto.PriceResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.PriceResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.PriceResponse}
 */
proto.PriceResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setTicker(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readDouble());
      msg.setStockprice(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setCurrency(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.PriceResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.PriceResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.PriceResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.PriceResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getTicker();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getStockprice();
  if (f !== 0.0) {
    writer.writeDouble(
      2,
      f
    );
  }
  f = message.getCurrency();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
};


/**
 * optional string Ticker = 1;
 * @return {string}
 */
proto.PriceResponse.prototype.getTicker = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.PriceResponse.prototype.setTicker = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional double StockPrice = 2;
 * @return {number}
 */
proto.PriceResponse.prototype.getStockprice = function() {
  return /** @type {number} */ (+jspb.Message.getFieldWithDefault(this, 2, 0.0));
};


/** @param {number} value */
proto.PriceResponse.prototype.setStockprice = function(value) {
  jspb.Message.setProto3FloatField(this, 2, value);
};


/**
 * optional string Currency = 3;
 * @return {string}
 */
proto.PriceResponse.prototype.getCurrency = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.PriceResponse.prototype.setCurrency = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * @enum {number}
 */
proto.Currencies = {
  USD: 0,
  EUR: 1,
  JPY: 2,
  BGN: 3,
  CZK: 4,
  DKK: 5,
  GBP: 6,
  HUF: 7,
  PLN: 8,
  RON: 9,
  SEK: 10,
  CHF: 11,
  ISK: 12,
  NOK: 13,
  HRK: 14,
  RUB: 15,
  TRY: 16,
  AUD: 17,
  BRL: 18,
  CAD: 19,
  CNY: 20,
  HKD: 21,
  IDR: 22,
  ILS: 23,
  INR: 24,
  KRW: 25,
  MXN: 26,
  MYR: 27,
  NZD: 28,
  PHP: 29,
  SGD: 30,
  THB: 31,
  ZAR: 32
};

goog.object.extend(exports, proto);
