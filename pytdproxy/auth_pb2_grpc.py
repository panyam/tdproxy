# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

from protos import auth_pb2 as protos_dot_auth__pb2


class AuthServiceStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.AddAuthToken = channel.unary_unary(
                '/protos.AuthService/AddAuthToken',
                request_serializer=protos_dot_auth__pb2.AddAuthTokenRequest.SerializeToString,
                response_deserializer=protos_dot_auth__pb2.AddAuthTokenResponse.FromString,
                )
        self.StartLogin = channel.unary_unary(
                '/protos.AuthService/StartLogin',
                request_serializer=protos_dot_auth__pb2.StartAuthRequest.SerializeToString,
                response_deserializer=protos_dot_auth__pb2.StartAuthResponse.FromString,
                )
        self.CompleteLogin = channel.unary_unary(
                '/protos.AuthService/CompleteLogin',
                request_serializer=protos_dot_auth__pb2.CompleteAuthRequest.SerializeToString,
                response_deserializer=protos_dot_auth__pb2.CompleteAuthResponse.FromString,
                )


class AuthServiceServicer(object):
    """Missing associated documentation comment in .proto file."""

    def AddAuthToken(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def StartLogin(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def CompleteLogin(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_AuthServiceServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'AddAuthToken': grpc.unary_unary_rpc_method_handler(
                    servicer.AddAuthToken,
                    request_deserializer=protos_dot_auth__pb2.AddAuthTokenRequest.FromString,
                    response_serializer=protos_dot_auth__pb2.AddAuthTokenResponse.SerializeToString,
            ),
            'StartLogin': grpc.unary_unary_rpc_method_handler(
                    servicer.StartLogin,
                    request_deserializer=protos_dot_auth__pb2.StartAuthRequest.FromString,
                    response_serializer=protos_dot_auth__pb2.StartAuthResponse.SerializeToString,
            ),
            'CompleteLogin': grpc.unary_unary_rpc_method_handler(
                    servicer.CompleteLogin,
                    request_deserializer=protos_dot_auth__pb2.CompleteAuthRequest.FromString,
                    response_serializer=protos_dot_auth__pb2.CompleteAuthResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'protos.AuthService', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class AuthService(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def AddAuthToken(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/protos.AuthService/AddAuthToken',
            protos_dot_auth__pb2.AddAuthTokenRequest.SerializeToString,
            protos_dot_auth__pb2.AddAuthTokenResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def StartLogin(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/protos.AuthService/StartLogin',
            protos_dot_auth__pb2.StartAuthRequest.SerializeToString,
            protos_dot_auth__pb2.StartAuthResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def CompleteLogin(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/protos.AuthService/CompleteLogin',
            protos_dot_auth__pb2.CompleteAuthRequest.SerializeToString,
            protos_dot_auth__pb2.CompleteAuthResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
