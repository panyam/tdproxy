# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: protos/streamer.proto
"""Generated protocol buffer code."""
from google.protobuf.internal import builder as _builder
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import struct_pb2 as google_dot_protobuf_dot_struct__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x15protos/streamer.proto\x12\x0etdproxy_protos\x1a\x1cgoogle/protobuf/struct.proto\"\x0e\n\x0c\x45mptyMessage\"\x1c\n\x0cSubscription\x12\x0c\n\x04name\x18\x01 \x01(\t\"4\n\x10SubscribeRequest\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\x12\n\ntopic_name\x18\x02 \x01(\t\"j\n\x0bSendRequest\x12\x10\n\x08sub_name\x18\x01 \x01(\t\x12\x0f\n\x07service\x18\x02 \x01(\t\x12\x0f\n\x07\x63ommand\x18\x03 \x01(\t\x12\'\n\x06params\x18\x04 \x01(\x0b\x32\x17.google.protobuf.Struct\"\x1e\n\x0cSendResponse\x12\x0e\n\x06status\x18\x01 \x01(\x03\"0\n\x07Message\x12%\n\x04info\x18\x01 \x01(\x0b\x32\x17.google.protobuf.Struct2\xe9\x01\n\x0fStreamerService\x12H\n\tSubscribe\x12 .tdproxy_protos.SubscribeRequest\x1a\x17.tdproxy_protos.Message0\x01\x12I\n\x0bUnsubscribe\x12\x1c.tdproxy_protos.Subscription\x1a\x1c.tdproxy_protos.EmptyMessage\x12\x41\n\x04Send\x12\x1b.tdproxy_protos.SendRequest\x1a\x1c.tdproxy_protos.SendResponseB\x10Z\x0etdproxy/protosb\x06proto3')

_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, globals())
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'protos.streamer_pb2', globals())
if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'Z\016tdproxy/protos'
  _EMPTYMESSAGE._serialized_start=71
  _EMPTYMESSAGE._serialized_end=85
  _SUBSCRIPTION._serialized_start=87
  _SUBSCRIPTION._serialized_end=115
  _SUBSCRIBEREQUEST._serialized_start=117
  _SUBSCRIBEREQUEST._serialized_end=169
  _SENDREQUEST._serialized_start=171
  _SENDREQUEST._serialized_end=277
  _SENDRESPONSE._serialized_start=279
  _SENDRESPONSE._serialized_end=309
  _MESSAGE._serialized_start=311
  _MESSAGE._serialized_end=359
  _STREAMERSERVICE._serialized_start=362
  _STREAMERSERVICE._serialized_end=595
# @@protoc_insertion_point(module_scope)