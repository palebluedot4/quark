import hashlib
import hmac
from typing import Protocol, final, override, runtime_checkable


class HashMismatchError(Exception):
    def __init__(self, message: str = "digest does not match the data") -> None:
        super().__init__(message)


@runtime_checkable
class Hasher(Protocol):
    def digest(self, data: bytes) -> bytes: ...
    def verify(self, digest: bytes, data: bytes) -> None: ...


@final
class SHA3_256Hasher(Hasher):
    __slots__ = ()

    @override
    def digest(self, data: bytes) -> bytes:
        return hashlib.sha3_256(data).digest()

    @override
    def verify(self, digest: bytes, data: bytes) -> None:
        computed = hashlib.sha3_256(data).digest()
        if not hmac.compare_digest(digest, computed):
            raise HashMismatchError()


@final
class SHA256Hasher(Hasher):
    __slots__ = ()

    @override
    def digest(self, data: bytes) -> bytes:
        return hashlib.sha256(data).digest()

    @override
    def verify(self, digest: bytes, data: bytes) -> None:
        computed = hashlib.sha256(data).digest()
        if not hmac.compare_digest(digest, computed):
            raise HashMismatchError()
