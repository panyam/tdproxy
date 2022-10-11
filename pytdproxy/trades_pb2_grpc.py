# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

from protos import trades_pb2 as protos_dot_trades__pb2


class TradeServiceStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.GetTrades = channel.unary_unary(
                '/protos.TradeService/GetTrades',
                request_serializer=protos_dot_trades__pb2.GetTradesRequest.SerializeToString,
                response_deserializer=protos_dot_trades__pb2.GetTradesResponse.FromString,
                )
        self.SaveTrades = channel.unary_unary(
                '/protos.TradeService/SaveTrades',
                request_serializer=protos_dot_trades__pb2.SaveTradesRequest.SerializeToString,
                response_deserializer=protos_dot_trades__pb2.SaveTradesResponse.FromString,
                )
        self.ListTrades = channel.unary_unary(
                '/protos.TradeService/ListTrades',
                request_serializer=protos_dot_trades__pb2.ListTradesRequest.SerializeToString,
                response_deserializer=protos_dot_trades__pb2.ListTradesResponse.FromString,
                )
        self.RemoveTrades = channel.unary_unary(
                '/protos.TradeService/RemoveTrades',
                request_serializer=protos_dot_trades__pb2.RemoveTradesRequest.SerializeToString,
                response_deserializer=protos_dot_trades__pb2.RemoveTradesResponse.FromString,
                )


class TradeServiceServicer(object):
    """Missing associated documentation comment in .proto file."""

    def GetTrades(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def SaveTrades(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def ListTrades(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def RemoveTrades(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_TradeServiceServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'GetTrades': grpc.unary_unary_rpc_method_handler(
                    servicer.GetTrades,
                    request_deserializer=protos_dot_trades__pb2.GetTradesRequest.FromString,
                    response_serializer=protos_dot_trades__pb2.GetTradesResponse.SerializeToString,
            ),
            'SaveTrades': grpc.unary_unary_rpc_method_handler(
                    servicer.SaveTrades,
                    request_deserializer=protos_dot_trades__pb2.SaveTradesRequest.FromString,
                    response_serializer=protos_dot_trades__pb2.SaveTradesResponse.SerializeToString,
            ),
            'ListTrades': grpc.unary_unary_rpc_method_handler(
                    servicer.ListTrades,
                    request_deserializer=protos_dot_trades__pb2.ListTradesRequest.FromString,
                    response_serializer=protos_dot_trades__pb2.ListTradesResponse.SerializeToString,
            ),
            'RemoveTrades': grpc.unary_unary_rpc_method_handler(
                    servicer.RemoveTrades,
                    request_deserializer=protos_dot_trades__pb2.RemoveTradesRequest.FromString,
                    response_serializer=protos_dot_trades__pb2.RemoveTradesResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'protos.TradeService', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class TradeService(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def GetTrades(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/protos.TradeService/GetTrades',
            protos_dot_trades__pb2.GetTradesRequest.SerializeToString,
            protos_dot_trades__pb2.GetTradesResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def SaveTrades(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/protos.TradeService/SaveTrades',
            protos_dot_trades__pb2.SaveTradesRequest.SerializeToString,
            protos_dot_trades__pb2.SaveTradesResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def ListTrades(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/protos.TradeService/ListTrades',
            protos_dot_trades__pb2.ListTradesRequest.SerializeToString,
            protos_dot_trades__pb2.ListTradesResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def RemoveTrades(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/protos.TradeService/RemoveTrades',
            protos_dot_trades__pb2.RemoveTradesRequest.SerializeToString,
            protos_dot_trades__pb2.RemoveTradesResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
